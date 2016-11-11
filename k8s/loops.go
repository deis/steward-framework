package k8s

import (
	"context"

	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/k8s/broker"
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
	globalNamespace string,
	errCh chan<- error,
) {
	restClient := k8sClient.CoreClient.RESTClient()
	// start broker loop
	go func() {
		watchBrokerFn := broker.NewK8sWatchBrokerFunc(restClient)
		createSvcClassFn := broker.NewK8sCreateServiceClassFunc(restClient)
		if err := broker.RunLoop(
			ctx,
			globalNamespace,
			watchBrokerFn,
			cataloger,
			createSvcClassFn,
		); err != nil {
			errCh <- err
		}
	}()

	// TODO: Start instance loop
	// TODO: Start binding loop

	<-ctx.Done()
}
