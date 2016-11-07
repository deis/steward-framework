package claim

import (
	"context"
	"strconv"
	"time"

	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/k8s"
	"github.com/deis/steward-framework/k8s/claim/state"
	"github.com/deis/steward-framework/lib"
)

const (
	asyncProvisionRespOperationKey = "provision-resp-operation"
	asyncProvisionPollStateKey     = "provision-poll-state"
	asyncProvisionPollCountKey     = "provision-poll-count"
)

func pollProvisionState(
	ctx context.Context,
	serviceID,
	planID,
	operation,
	instanceID string,
	opStatusRetriever framework.OperationStatusRetriever,
	claimCh chan<- state.Update,
) framework.OperationState {

	pollNum := 0
	pollState := framework.OperationStateInProgress
	pollErrCount := 0

	startTime := time.Now()
	for {
		if pollState == framework.OperationStateSucceeded || pollState == framework.OperationStateFailed {
			// if the polling went into success or failed state, just return that
			return pollState
		}
		if pollState == framework.OperationStateGone {
			// When provisioning, treat "gone" as a failure
			return framework.OperationStateFailed
		}

		// If maxAsyncDuration has been exceeded
		if time.Since(startTime) > maxAsyncDuration {
			select {
			case claimCh <- state.FullUpdate(
				k8s.StatusFailed,
				"asynchronous provisionining has exceeded the one hour allotted; service state is unknown",
				instanceID,
				"",
				lib.EmptyJSONObject(),
			):
			case <-ctx.Done():
			}
			return framework.OperationStateFailed
		}

		// otherwise continue provisioning state
		update := state.FullUpdate(
			k8s.StatusProvisioningAsync,
			"polling for asynchronous provisionining",
			instanceID,
			"",
			lib.JSONObject(map[string]interface{}{
				asyncProvisionRespOperationKey: operation,
				asyncProvisionPollStateKey:     pollState.String(),
				asyncProvisionPollCountKey:     strconv.Itoa(pollNum),
			}))
		select {
		case claimCh <- update:
		case <-ctx.Done():
		}
		resp, err := opStatusRetriever.GetOperationStatus(ctx, &framework.OperationStatusRequest{
			InstanceID: instanceID,
			ServiceID:  serviceID,
			PlanID:     planID,
			Operation:  operation,
		})
		if err != nil {
			if pollErrCount < 3 {
				pollErrCount++
			} else {
				// After threee consecutive polling errors, we'll consider provisioning failed
				select {
				case claimCh <- state.FullUpdate(
					k8s.StatusFailed,
					"polling for asynchronous provisionining has failed (repeatedly); service state is unknown",
					instanceID,
					"",
					lib.EmptyJSONObject(),
				):
				case <-ctx.Done():
				}
				return framework.OperationStateFailed
			}
		} else {
			// Reset error count to zero
			pollErrCount = 0
		}
		pollNum++
		newState, err := resp.GetState()
		if err != nil {
			select {
			case claimCh <- state.ErrUpdate(err):
			case <-ctx.Done():
			}
			return framework.OperationStateFailed
		}
		pollState = newState
		time.Sleep(30 * time.Second)
	}
}
