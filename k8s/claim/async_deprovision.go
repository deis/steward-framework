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
	asyncDeprovisionRespOperationKey = "deprovision-resp-operation"
	asyncDeprovisionPollStateKey     = "deprovision-poll-state"
	asyncDeprovisionPollCountKey     = "deprovision-poll-count"
)

func pollDeprovisionState(
	ctx context.Context,
	serviceID,
	planID,
	operation,
	instanceID string,
	lastOpGetter framework.LastOperationGetter,
	claimCh chan<- state.Update,
) framework.LastOperationState {
	pollNum := 0
	pollState := framework.LastOperationStateInProgress
	pollErrCount := 0

	startTime := time.Now()
	for {
		if pollState == framework.LastOperationStateSucceeded || pollState == framework.LastOperationStateFailed {
			// if the polling went into success or failed state, just return that
			return pollState
		}
		if pollState == framework.LastOperationStateGone {
			// When deprovisioning, treat "gone" as success
			return framework.LastOperationStateSucceeded
		}

		// If maxAsyncDuration has been exceeded
		if time.Since(startTime) > maxAsyncDuration {
			select {
			case claimCh <- state.FullUpdate(
				k8s.StatusFailed,
				"asynchronous deprovisionining has exceeded the one hour allotted; service state is unknown",
				instanceID,
				"",
				lib.EmptyJSONObject(),
			):
			case <-ctx.Done():
			}
			return framework.LastOperationStateFailed
		}

		// otherwise continue provisioning state
		update := state.FullUpdate(
			k8s.StatusDeprovisioningAsync,
			"polling for asynchronous deprovisionining",
			instanceID,
			"",
			lib.JSONObject(map[string]interface{}{
				asyncDeprovisionRespOperationKey: operation,
				asyncDeprovisionPollStateKey:     pollState.String(),
				asyncDeprovisionPollCountKey:     strconv.Itoa(pollNum),
			}))
		select {
		case claimCh <- update:
		case <-ctx.Done():
		}
		resp, err := lastOpGetter.GetLastOperation(ctx, &framework.GetLastOperationRequest{
			InstanceID: instanceID,
			ServiceID:  serviceID,
			PlanID:     planID,
			Operation:  operation,
		})
		if err != nil {
			if pollErrCount < 3 {
				pollErrCount++
			} else {
				// After threee consecutive polling errors, we'll consider deprovisioning failed
				select {
				case claimCh <- state.FullUpdate(
					k8s.StatusFailed,
					"polling for asynchronous depprovisionining has failed (repeatedly); service state is unknown",
					instanceID,
					"",
					lib.EmptyJSONObject(),
				):
				case <-ctx.Done():
				}
				return framework.LastOperationStateFailed
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
			return framework.LastOperationStateFailed
		}
		pollState = newState
		time.Sleep(30 * time.Second)
	}
}
