package serviceinstance

import (
	"context"
	"time"

	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/k8s/data"
)

func pollProvisionState(
	ctx context.Context,
	serviceBrokerSpec framework.ServiceBrokerSpec,
	lastOpStatGetter framework.OperationStatusRetriever,
	serviceInstanceID string,
	serviceClassID string,
	servicePlanID string,
) (data.ServiceInstanceState, string, error) {

	// Start by assuming provisioning is in-progress
	pollState := framework.OperationStateInProgress
	var pollErrCount int

	ticker := time.NewTicker(30 * time.Second)
	// At the start of each loop, we'll check if we've successfully determined service instance
	// status. If not, then we'll proceed into a select where we'll either poll at the scheduled
	// interval, conext will be cancelled, or the whole thing will time out.
	for {
		if pollState == framework.OperationStateSucceeded {
			return data.ServiceInstanceStateProvisioned, "", nil
		}
		if pollState == framework.OperationStateFailed {
			return data.ServiceInstanceStateFailed, "Asynchronous provisioning failed", nil
		}
		if pollState == framework.OperationStateGone {
			// When provisioning, treat "gone" as a failure
			return data.ServiceInstanceStateFailed, "Asynchronous provisioning failed", nil
		}
		select {
		case <-ticker.C:
			resp, err := lastOpStatGetter.GetOperationStatus(
				ctx,
				serviceBrokerSpec,
				&framework.OperationStatusRequest{
					InstanceID: serviceInstanceID,
					ServiceID:  serviceClassID,
					PlanID:     servicePlanID,
					Operation:  "provision",
				},
			)
			if err != nil {
				if pollErrCount < 3 {
					pollErrCount++
				} else {
					// After threee consecutive polling errors, we'll consider provisioning failed
					return data.ServiceInstanceStateUnknown, "Polling for asynchronous provisioning status has failed", nil
				}
			} else {
				// Reset error count to zero
				pollErrCount = 0
			}
			newState, err := resp.GetState()
			if err != nil {
				// After threee consecutive polling errors, we'll consider provisioning failed
				return data.ServiceInstanceStateUnknown, "Polling for asynchronous provisioning status has failed", nil
			}
			pollState = newState
		case <-ctx.Done():
			return data.ServiceInstanceStateUnknown, "Polling for asynchronous provisioning status was cancelled or has timed out", nil
		}
	}
}

func pollDeprovisionState(
	ctx context.Context,
	serviceBrokerSpec framework.ServiceBrokerSpec,
	lastOpStatGetter framework.OperationStatusRetriever,
	serviceInstanceID string,
	serviceClassID string,
	servicePlanID string,
) error {

	// Start by assuming provisioning is in-progress
	pollState := framework.OperationStateInProgress
	var pollErrCount int

	ticker := time.NewTicker(30 * time.Second)
	// At the start of each loop, we'll check if we've successfully determined service instance
	// status. If not, then we'll proceed into a select where we'll either poll at the scheduled
	// interval, conext will be cancelled, or the whole thing will time out.
	for {
		// When deprovisioning, treat "gone" as success
		if pollState == framework.OperationStateSucceeded || pollState == framework.OperationStateGone {
			return nil
		}
		if pollState == framework.OperationStateFailed {
			logger.Errorf("Asynchronous deprovisioning failed")
			return nil
		}
		select {
		case <-ticker.C:
			resp, err := lastOpStatGetter.GetOperationStatus(
				ctx,
				serviceBrokerSpec,
				&framework.OperationStatusRequest{
					InstanceID: serviceInstanceID,
					ServiceID:  serviceClassID,
					PlanID:     servicePlanID,
					Operation:  "deprovision",
				},
			)
			if err != nil {
				if pollErrCount < 3 {
					pollErrCount++
				} else {
					logger.Errorf("Encountered 3 consecutive polling errors; service instance status is unknown")
					return nil
				}
			} else {
				// Reset error count to zero
				pollErrCount = 0
			}
			newState, err := resp.GetState()
			if err != nil {
				logger.Errorf("Error determining operation state; service instance status is unknown")
				return nil
			}
			pollState = newState
		case <-ctx.Done():
			logger.Errorf("Polling for asynchronous provisioning status was cancelled or has timed out")
		}
	}
}
