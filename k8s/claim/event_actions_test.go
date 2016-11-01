package claim

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/fake"
	"github.com/deis/steward-framework/k8s"
	"github.com/deis/steward-framework/k8s/claim/state"
	"github.com/deis/steward-framework/lib"
	"github.com/pborman/uuid"
)

const (
	waitDur = 100 * time.Millisecond
)

var (
	ctx = context.Background()
)

func TestUpdateIsTerminal(t *testing.T) {
	// claims that should be marked terminal
	terminalUpdates := []state.Update{
		state.StatusUpdate(k8s.StatusFailed),
		state.StatusUpdate(k8s.StatusProvisioned),
		state.StatusUpdate(k8s.StatusBound),
		state.StatusUpdate(k8s.StatusUnbound),
		state.StatusUpdate(k8s.StatusDeprovisioned),
	}
	for i, update := range terminalUpdates {
		if !state.UpdateIsTerminal(update) {
			t.Fatalf("update %d (%s) was not marked terminal", i, update)
		}
	}

	// normal updates that should not be marked terminal
	normalUpdates := []state.Update{
		state.StatusUpdate(k8s.StatusProvisioning),
		state.StatusUpdate(k8s.StatusBinding),
		state.StatusUpdate(k8s.StatusUnbinding),
		state.StatusUpdate(k8s.StatusDeprovisioning),
	}
	for i, update := range normalUpdates {
		if state.UpdateIsTerminal(update) {
			t.Fatalf("update %d (%s) was marked terminal", i, update)
		}
	}
}

func TestNewErrClaimUpdate(t *testing.T) {
	err := errors.New("test error")
	update := state.ErrUpdate(err)
	assert.True(t, state.UpdateIsTerminal(update), "new claim wasn't marked stop")
	assert.Equal(t, update.Status(), k8s.StatusFailed, "resulting status")
	assert.Equal(t, update.Description(), err.Error(), "resulting status description")
}

func getCatalogFromEvents(evts ...*Event) k8s.ServiceCatalogLookup {
	ret := k8s.NewServiceCatalogLookup(nil)
	for _, evt := range evts {
		ret.Set(&k8s.ServiceCatalogEntry{
			Info: framework.ServiceInfo{ID: evt.claim.Claim.ServiceID},
			Plan: framework.ServicePlan{ID: evt.claim.Claim.PlanID},
		})
	}
	return ret
}

func TestGetService(t *testing.T) {
	claim := getClaim(k8s.ActionProvision)
	catalog := k8s.NewServiceCatalogLookup(nil)
	svc, err := getService(claim, catalog)
	assert.True(t, isNoSuchServiceAndPlanErr(err), "returned error was not a errNoSuchServiceAndPlan")
	assert.Nil(t, svc, "returned service")
	catalog = k8s.NewServiceCatalogLookup([]*k8s.ServiceCatalogEntry{
		{
			Info: framework.ServiceInfo{ID: claim.ServiceID},
			Plan: framework.ServicePlan{ID: claim.PlanID},
		},
	})
	svc, err = getService(claim, catalog)
	assert.NoErr(t, err)
	assert.NotNil(t, svc, "returned service")
	claim.ServiceID = "doesnt-exist"
	svc, err = getService(claim, catalog)
	assert.True(t, isNoSuchServiceAndPlanErr(err), "returned error was not a errNoSuchServiceAndPlan")
	assert.Nil(t, svc, "returned service")
}

func TestProcessProvisionServiceNotFound(t *testing.T) {
	evt := getEvent(getClaim(k8s.ActionProvision))
	catalogLookup := k8s.NewServiceCatalogLookup(nil)
	ch := make(chan state.Update)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processProvision(cancelCtx, evt, nil, catalogLookup, nil, ch)
	select {
	case claimUpdate := <-ch:
		assert.True(t, state.UpdateIsTerminal(claimUpdate), "claim update was not returned as terminal")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}
}

func TestProcessProvisionServiceFound(t *testing.T) {
	evt := getEvent(getClaim(k8s.ActionProvision))
	catalogLookup := getCatalogFromEvents(evt)
	ch := make(chan state.Update)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	provisioner := &fake.Provisioner{
		Res: &framework.ProvisionResponse{
			Extra: lib.JSONObject(map[string]interface{}{
				uuid.New(): uuid.New(),
			}),
		},
	}
	lifecycler := &fake.Lifecycler{
		Provisioner: provisioner,
	}
	go processProvision(cancelCtx, evt, nil, catalogLookup, lifecycler, ch)

	// provisioning status
	select {
	case claimUpdate := <-ch:
		assert.False(t, state.UpdateIsTerminal(claimUpdate), "update was marked terminal when it shouldn't have been")
		assert.Equal(t, claimUpdate.Status(), k8s.StatusProvisioning, "new status")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}

	// provisioned status
	select {
	case claimUpdate := <-ch:
		assert.True(t, state.UpdateIsTerminal(claimUpdate), "update was not marked terminal when it should have been")
		assert.Equal(t, claimUpdate.Status(), k8s.StatusProvisioned, "new status")
		assert.True(t, len(claimUpdate.InstanceID()) > 0, "no instance ID written")
		assert.Equal(t, len(claimUpdate.BindingID()), 0, "bind ID written when it shouldn't have been")
		assert.Equal(t, claimUpdate.Extra(), lib.JSONObject(provisioner.Res.Extra), "extra data")
		assert.Equal(t, len(provisioner.Reqs), 1, "number of provision calls")
		assert.True(t, state.UpdateIsTerminal(claimUpdate), "provisioned update was not marked terminal")
		req := provisioner.Reqs[0]
		assert.Equal(t, req.ServiceID, evt.claim.Claim.ServiceID, "service ID")
		assert.Equal(t, req.PlanID, evt.claim.Claim.PlanID, "plan ID")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}
}

func TestProcessBindServiceNotFound(t *testing.T) {
	evt := getEvent(getClaim(k8s.ActionBind))
	catalogLookup := k8s.NewServiceCatalogLookup(nil)
	ch := make(chan state.Update)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processBind(cancelCtx, evt, nil, catalogLookup, nil, ch)
	select {
	case claimUpdate := <-ch:
		assert.True(t, state.UpdateIsTerminal(claimUpdate), "claim update was not returned terminal when it should have been")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}
}

func TestProcessBindServiceFound(t *testing.T) {
	evt := getEvent(getClaim(k8s.ActionBind))
	catalogLookup := getCatalogFromEvents(evt)
	ch := make(chan state.Update)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processBind(cancelCtx, evt, nil, catalogLookup, nil, ch)

	// binding status
	select {
	case claimUpdate := <-ch:
		assert.False(t, state.UpdateIsTerminal(claimUpdate), "claim update was returned terminal when it shouldn't have been")
		assert.Equal(t, claimUpdate.Status(), k8s.StatusBinding, "new status")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}

	// missing instance ID status
	select {
	case claimUpdate := <-ch:
		assert.True(t, state.UpdateIsTerminal(claimUpdate), "claim update was not returned terminal when it should have been")
		assert.Equal(t, claimUpdate.Status(), k8s.StatusFailed, "new status")
		assert.Equal(t, claimUpdate.Description(), errMissingInstanceID.Error(), "new status description")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}
}

func TestProcessBindInstanceIDFound(t *testing.T) {
	evt := getEvent(getClaim(k8s.ActionBind))
	evt.claim.Claim.InstanceID = uuid.New()
	catalogLookup := getCatalogFromEvents(evt)
	binder := &fake.Binder{
		Res: &framework.BindResponse{
			Creds: lib.JSONObject(map[string]interface{}{
				"cred1": uuid.New(),
				"cred2": uuid.New(),
			}),
		},
	}
	lifecycler := &fake.Lifecycler{
		Binder: binder,
	}
	secretsNamespacer := k8s.NewFakeSecretsNamespacer()
	ch := make(chan state.Update)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processBind(cancelCtx, evt, secretsNamespacer, catalogLookup, lifecycler, ch)

	// binding status
	select {
	case claimUpdate := <-ch:
		assert.False(t, state.UpdateIsTerminal(claimUpdate), "claim update was marked terminal when it shouldn't have been")
		assert.Equal(t, claimUpdate.Status(), k8s.StatusBinding, "new status")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}

	// bound status
	select {
	case claimUpdate := <-ch:
		assert.True(t, state.UpdateIsTerminal(claimUpdate), "claim update was not marked terminal when it should have been")
		assert.Equal(t, claimUpdate.Status(), k8s.StatusBound, "new status")
		assert.True(t, len(claimUpdate.InstanceID()) > 0, "instance ID not found")
		assert.True(t, len(claimUpdate.BindingID()) > 0, "bind ID not found")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}

	// check the lifecycler
	assert.Equal(t, len(binder.Reqs), 1, "number of bind calls")
	req := binder.Reqs[0]
	assert.Equal(t, req.InstanceID, evt.claim.Claim.InstanceID, "instance ID")
	assert.Equal(t, req.ServiceID, evt.claim.Claim.ServiceID, "service ID")
	assert.Equal(t, req.PlanID, evt.claim.Claim.PlanID, "plan ID")

	// check the secrets namespacer
	assert.Equal(t, len(secretsNamespacer.Returned), 1, "number of returned secrets interfaces")
	secretsIface := secretsNamespacer.Returned["testns"]
	assert.Equal(t, len(secretsIface.Created), 1, "number of created secrets")
	secret := secretsIface.Created[0]
	assert.Equal(t, len(secret.Data), len(binder.Res.Creds), "number of creds in the secret")
	for k, v := range secret.Data {
		assert.Equal(t, string(v), binder.Res.Creds[k], fmt.Sprintf("value for key %s", k))
	}
}

func TestProcessUnbindServiceNotFound(t *testing.T) {
	evt := getEvent(getClaim(k8s.ActionUnbind))
	catalogLookup := k8s.NewServiceCatalogLookup(nil)
	ch := make(chan state.Update)
	go processUnbind(ctx, evt, nil, catalogLookup, nil, ch)
	select {
	case claimUpdate := <-ch:
		assert.True(t, state.UpdateIsTerminal(claimUpdate), "claim update was not marked terminal when it should have been")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}
}

func TestProcessUnbindServiceFound(t *testing.T) {
	evt := getEvent(getClaim(k8s.ActionBind))
	catalogLookup := getCatalogFromEvents(evt)
	ch := make(chan state.Update)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processUnbind(cancelCtx, evt, nil, catalogLookup, nil, ch)

	// unbinding status
	select {
	case claimUpdate := <-ch:
		assert.False(t, state.UpdateIsTerminal(claimUpdate), "claim update was marked terminal when it shouldn't have been")
		assert.Equal(t, claimUpdate.Status(), k8s.StatusUnbinding, "new status")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}

	// missing instance ID status
	select {
	case claimUpdate := <-ch:
		assert.True(t, state.UpdateIsTerminal(claimUpdate), "claim update was not marked terminal when it should have been")
		assert.Equal(t, claimUpdate.Description(), errMissingInstanceID.Error(), "new status description")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}
}

func TestProcessUnbindInstanceIDFound(t *testing.T) {
	claim := getClaim(k8s.ActionUnbind)
	claim.InstanceID = uuid.New()
	claim.BindingID = uuid.New()
	evt := getEvent(claim)
	catalogLookup := getCatalogFromEvents(evt)
	unbinder := &fake.Unbinder{}
	lifecycler := &fake.Lifecycler{Unbinder: unbinder}
	secretsNamespacer := k8s.NewFakeSecretsNamespacer()
	ch := make(chan state.Update)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processUnbind(cancelCtx, evt, secretsNamespacer, catalogLookup, lifecycler, ch)

	// unbinding status
	select {
	case claimUpdate := <-ch:
		assert.False(t, state.UpdateIsTerminal(claimUpdate), "claim update was marked terminal when it shouldn't have been")
		assert.Equal(t, claimUpdate.Status(), k8s.StatusUnbinding, "new status")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}

	// unbound status
	select {
	case claimUpdate := <-ch:
		assert.True(t, state.UpdateIsTerminal(claimUpdate), "claim update was not marked terminal when it should have been")
		assert.Equal(t, claimUpdate.Status(), k8s.StatusUnbound, "new status")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}

	// check the lifecycler
	assert.Equal(t, len(unbinder.Reqs), 1, "number of bind calls")
	req := unbinder.Reqs[0]
	assert.Equal(t, req.InstanceID, claim.InstanceID, "instance ID")
	assert.Equal(t, req.BindingID, claim.BindingID, "bind ID")

	// check the secrets namespacer
	assert.Equal(t, len(secretsNamespacer.Returned), 1, "number of returned secrets interfaces")
	secretsIface := secretsNamespacer.Returned["testns"]
	assert.Equal(t, len(secretsIface.Deleted), 1, "number of deleted secrets")
}

func TestProcessDeprovisionServiceNotFound(t *testing.T) {
	evt := getEvent(getClaim(k8s.ActionDeprovision))
	catalogLookup := k8s.NewServiceCatalogLookup(nil)
	ch := make(chan state.Update)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processDeprovision(cancelCtx, evt, nil, catalogLookup, nil, ch)
	select {
	case claimUpdate := <-ch:
		assert.True(t, state.UpdateIsTerminal(claimUpdate), "claim update was not marked terminal when it should have been")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}
}

func TestProcessDeprovisionServiceFound(t *testing.T) {
	evt := getEvent(getClaim(k8s.ActionDeprovision))
	catalogLookup := getCatalogFromEvents(evt)
	deprovisioner := &fake.Deprovisioner{}
	lifecycler := &fake.Lifecycler{
		Deprovisioner: deprovisioner,
	}
	ch := make(chan state.Update)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processDeprovision(cancelCtx, evt, nil, catalogLookup, lifecycler, ch)

	// deprovisioning status
	select {
	case claimUpdate := <-ch:
		assert.False(t, state.UpdateIsTerminal(claimUpdate), "claim update was marked terminal when it should have been")
		assert.Equal(t, claimUpdate.Status(), k8s.StatusDeprovisioning, "new status")
		assert.Equal(t, len(deprovisioner.Reqs), 0, "number of deprovision calls")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}

	// missing instance ID status
	select {
	case claimUpdate := <-ch:
		assert.True(t, state.UpdateIsTerminal(claimUpdate), "claim update was not marked terminal when it should have been")
		assert.Equal(t, claimUpdate.Status(), k8s.StatusFailed, "new status")
		assert.Equal(t, claimUpdate.Description(), errMissingInstanceID.Error(), "new status description")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}
}

func TestDeprovisionInstanceIDFound(t *testing.T) {
	claim := getClaim(k8s.ActionDeprovision)
	claim.InstanceID = uuid.New()
	claim.Extra = lib.JSONObject(map[string]interface{}{
		uuid.New(): uuid.New(),
	})
	evt := getEvent(claim)
	catalogLookup := getCatalogFromEvents(evt)
	ch := make(chan state.Update)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	deprovisioner := &fake.Deprovisioner{
		Res: &framework.DeprovisionResponse{Operation: "testop"},
	}
	lifecycler := &fake.Lifecycler{
		Deprovisioner: deprovisioner,
	}
	go processDeprovision(cancelCtx, evt, nil, catalogLookup, lifecycler, ch)

	// deprovisioning status
	select {
	case claimUpdate := <-ch:
		assert.False(t, state.UpdateIsTerminal(claimUpdate), "update was marked terminal when it shouldn't have been")
		assert.Equal(t, claimUpdate.Status(), k8s.StatusDeprovisioning, "new status")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}

	// deprovisioned status
	select {
	case claimUpdate := <-ch:
		assert.True(t, state.UpdateIsTerminal(claimUpdate), "update was not marked terminal when it should have been")
		assert.Equal(t, claimUpdate.Status(), k8s.StatusDeprovisioned, "new status")
		// assert.Equal(t, claimUpdate.Extra(), deprivi.Resp.Extra, "extra data")
		assert.Equal(t, len(deprovisioner.Reqs), 1, "number of provision calls")
		assert.True(t, state.UpdateIsTerminal(claimUpdate), "provisioned update was not marked terminal")
		req := deprovisioner.Reqs[0]
		assert.Equal(t, req.InstanceID, claim.InstanceID, "instance ID")
		assert.Equal(t, req.ServiceID, claim.ServiceID, "service ID")
		assert.Equal(t, req.PlanID, claim.PlanID, "plan ID")
		assert.Equal(t, lib.JSONObject(req.Parameters), claim.Extra, "extra")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}
}
