package claim

import (
	"context"
	"testing"

	"github.com/deis/steward-framework/fake"
	"github.com/deis/steward-framework/k8s"
)

func TestStartLoop(t *testing.T) {
	t.Skip("TODO")
}

func TestReceiveEvent(t *testing.T) {
	ctx := context.Background()
	evt := getEvent(getClaim(k8s.ActionProvision))
	iface := &FakeInteractor{}
	secretsNamespacer := &k8s.FakeSecretsNamespacer{}
	lookup := k8s.NewServiceCatalogLookup(nil) // TODO: add service/plan to the catalog
	lifecycler := &fake.Lifecycler{}
	receiveEvent(ctx, evt, iface, secretsNamespacer, lookup, lifecycler)
}

func TestStopLoop(t *testing.T) {
	t.Skip("TODO")
}

func TestWatchChanClosed(t *testing.T) {
	t.Skip("TODO")
}
