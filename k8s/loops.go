package k8s

import (
	"context"

	"github.com/deis/steward-framework"
	"k8s.io/client-go/kubernetes"
)

// StartControlLoops is a convenience function that initiates all of Steward's individual control
// loops, of which there is one per relevant resource type-- Broker, Instance, and Binding.
//
// In order to stop all loops, pass a cancellable context to this function and call its cancel
// function.
func StartControlLoops(
	ctx context.Context,
	k8sClient *kubernetes.Clientset,
	cataloger framework.Cataloger,
	lifecycler framework.Lifecycler,
	errCh chan<- error,
) {
	// TODO: Start broker / catalog publish loop
	// TODO: Start instance loop
	// TODO: Start binding loop
}
