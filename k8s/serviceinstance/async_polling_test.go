package serviceinstance

import (
	"context"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/fake"
	"github.com/deis/steward-framework/k8s/data"
	"github.com/pborman/uuid"
)

func TestPollProvisionStateSuccess(t *testing.T) {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	serviceID := uuid.New()
	planID := uuid.New()
	instanceID := uuid.New()

	lastOpGetter := &fake.OperationStatusRetriever{
		Res: func() *framework.OperationStatusResponse {
			return &framework.OperationStatusResponse{
				State: framework.OperationStateSucceeded.String(),
			}
		},
	}

	finalState, _, err := pollProvisionState(
		ctx,
		framework.ServiceBrokerSpec{},
		lastOpGetter,
		instanceID,
		serviceID,
		planID,
	)
	assert.NoErr(t, err)
	assert.Equal(t, finalState, data.ServiceInstanceStateProvisioned, "final state")
}

func TestPollProvisionStateFailure(t *testing.T) {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	serviceID := uuid.New()
	planID := uuid.New()
	instanceID := uuid.New()

	lastOpGetter := &fake.OperationStatusRetriever{
		Res: func() *framework.OperationStatusResponse {
			return &framework.OperationStatusResponse{
				State: framework.OperationStateFailed.String(),
			}
		},
	}

	finalState, _, err := pollProvisionState(
		ctx,
		framework.ServiceBrokerSpec{},
		lastOpGetter,
		instanceID,
		serviceID,
		planID,
	)
	assert.NoErr(t, err)
	assert.Equal(t, finalState, data.ServiceInstanceStateFailed, "final state")
}

func TestPollProvisionStateTimeout(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFn()

	serviceID := uuid.New()
	planID := uuid.New()
	instanceID := uuid.New()

	lastOpGetter := &fake.OperationStatusRetriever{
		Res: func() *framework.OperationStatusResponse {
			return &framework.OperationStatusResponse{
				State: framework.OperationStateSucceeded.String(),
			}
		},
	}

	finalState, _, err := pollProvisionState(
		ctx,
		framework.ServiceBrokerSpec{},
		lastOpGetter,
		instanceID,
		serviceID,
		planID,
	)
	assert.NoErr(t, err)
	assert.Equal(t, finalState, data.ServiceInstanceStateUnknown, "final state")
}

func TestPollDeprovisionStateSuccess(t *testing.T) {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	serviceID := uuid.New()
	planID := uuid.New()
	instanceID := uuid.New()

	lastOpGetter := &fake.OperationStatusRetriever{
		Res: func() *framework.OperationStatusResponse {
			return &framework.OperationStatusResponse{
				State: framework.OperationStateSucceeded.String(),
			}
		},
	}

	err := pollDeprovisionState(
		ctx,
		framework.ServiceBrokerSpec{},
		lastOpGetter,
		instanceID,
		serviceID,
		planID,
	)
	assert.NoErr(t, err)
}

func TestPollDeprovisionStateFailed(t *testing.T) {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	serviceID := uuid.New()
	planID := uuid.New()
	instanceID := uuid.New()

	lastOpGetter := &fake.OperationStatusRetriever{
		Res: func() *framework.OperationStatusResponse {
			return &framework.OperationStatusResponse{
				State: framework.OperationStateFailed.String(),
			}
		},
	}

	err := pollDeprovisionState(
		ctx,
		framework.ServiceBrokerSpec{},
		lastOpGetter,
		instanceID,
		serviceID,
		planID,
	)
	assert.NoErr(t, err)
}

func TestPollDeprovisionStateTimeout(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFn()

	serviceID := uuid.New()
	planID := uuid.New()
	instanceID := uuid.New()

	lastOpGetter := &fake.OperationStatusRetriever{
		Res: func() *framework.OperationStatusResponse {
			return &framework.OperationStatusResponse{
				State: framework.OperationStateSucceeded.String(),
			}
		},
	}

	err := pollDeprovisionState(
		ctx,
		framework.ServiceBrokerSpec{},
		lastOpGetter,
		instanceID,
		serviceID,
		planID,
	)
	assert.NoErr(t, err)
}
