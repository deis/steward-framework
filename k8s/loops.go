package k8s

import (
	"context"

	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/k8s/binding"
	"github.com/deis/steward-framework/k8s/broker"
	"github.com/deis/steward-framework/k8s/instance"
	"github.com/deis/steward-framework/k8s/refs"
	"k8s.io/client-go/dynamic"
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
	dynamicCl *dynamic.Client,
	cataloger framework.Cataloger,
	lifecycler framework.Lifecycler,
	globalNamespace string,
	errCh chan<- error,
) {

	// Start broker loop
	go func() {
		updateBrokerFn := broker.NewK8sUpdateBrokerFunc(dynamicCl)
		watchBrokerFn := broker.NewK8sWatchBrokerFunc(dynamicCl)
		createSvcClassFn := broker.NewK8sCreateServiceClassFunc(dynamicCl)
		if err := broker.RunLoop(
			ctx,
			globalNamespace,
			watchBrokerFn,
			updateBrokerFn,
			cataloger,
			createSvcClassFn,
		); err != nil {
			logger.Errorf("broker loop (%s)", err)
			errCh <- err
		}
	}()

	// Start instance loop
	go func() {
		if err := instance.RunLoop(
			ctx,
			instance.NewK8sWatchInstanceFunc(dynamicCl),
			instance.NewK8sUpdateInstanceFunc(dynamicCl),
			refs.NewK8sServiceClassGetterFunc(dynamicCl),
			refs.NewK8sBrokerGetterFunc(dynamicCl),
			lifecycler,
		); err != nil {
			logger.Errorf("instance loop (%s)", err)
			errCh <- err
		}
	}()

	// start binding loop
	go func() {
		watchBindingFn := binding.NewK8sWatchBindingFunc(dynamicCl)
		updateBindingFn := binding.NewK8sUpdateBindingFunc(dynamicCl)

		brokerGetterFn := refs.NewK8sBrokerGetterFunc(dynamicCl)
		svcClassGetterFn := refs.NewK8sServiceClassGetterFunc(dynamicCl)
		instanceGetterFn := refs.NewK8sInstanceGetterFunc(dynamicCl)

		// TODO: remove the hard-coded "default" namespace and instead watch all namespaces and be able
		// to create secrets in any given namespace.
		//
		// See https://github.com/deis/steward-framework/issues/30 and
		// https://github.com/deis/steward-framework/issues/29 for details on how to resolve
		secretWriterFunc := binding.NewK8sSecretWriterFunc(k8sClient.Core().Secrets("default"))
		if err := binding.RunLoop(
			ctx,
			"default",
			lifecycler,
			secretWriterFunc,
			watchBindingFn,
			updateBindingFn,
			brokerGetterFn,
			svcClassGetterFn,
			instanceGetterFn,
		); err != nil {
			logger.Errorf("binding loop (%s)", err)
			errCh <- err
		}
	}()

	<-ctx.Done()
}
