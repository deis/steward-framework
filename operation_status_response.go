package framework

import (
	"fmt"
)

// ErrUnknownOperationState is an error indicating that a string is not a known operation state.
type ErrUnknownOperationState string

// Error is the error interface implementation.
func (e ErrUnknownOperationState) Error() string {
	return fmt.Sprintf("unknown operation state '%s'", string(e))
}

// OperationState represents the state returned in a "get operation status" call. This type
// implements fmt.Stringer.
type OperationState string

// String is the fmt.Stringer interface implementation.
func (l OperationState) String() string {
	return string(l)
}

const (
	// OperationStateSucceeded is the OperationState indicating that the operation has succeeded
	OperationStateSucceeded OperationState = "succeeded"
	// OperationStateFailed is the OperationState indicating that the operation has failed.
	OperationStateFailed OperationState = "failed"
	// OperationStateInProgress is the OperationState indicating that the operation is still in
	// progress.
	OperationStateInProgress OperationState = "in progress"
	// OperationStateGone is the OperationState indicating that the service broker has deleted the
	// service instance in question. In the case of async deprovisioning, this is an indicator of
	// success.
	OperationStateGone OperationState = "gone"
)

// OperationStatusResponse represents a response to a OperationStatusRequest.
type OperationStatusResponse struct {
	State string
}

// GetState returns the OperationState corresponding to a OperationStatusResponse instance's
// State attribute, or an error if that State attribute does not correspond to a valid
// OperationState.
func (o *OperationStatusResponse) GetState() (OperationState, error) {
	switch o.State {
	case OperationStateSucceeded.String():
		return OperationStateSucceeded, nil
	case OperationStateFailed.String():
		return OperationStateFailed, nil
	case OperationStateInProgress.String():
		return OperationStateInProgress, nil
	case OperationStateGone.String():
		return OperationStateGone, nil
	default:
		return OperationState(""), ErrUnknownOperationState(o.State)
	}
}
