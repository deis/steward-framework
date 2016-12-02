package k8s

import (
	"context"

	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/k8s/refs"
	"github.com/deis/steward-framework/k8s/servicebinding"
	"github.com/deis/steward-framework/k8s/servicebroker"
	"github.com/deis/steward-framework/k8s/serviceinstance"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// StartControlLoops is a convenience function that initiates all of Steward's individual control
// loops, of which there is one per relevant resource type-- ServiceBroker, ServiceInstance, and
// ServiceBinding.
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

	// Start service broker loop
	go func() {
		updateServiceBrokerFn := servicebroker.NewK8sUpdateServiceBrokerFunc(dynamicCl)
		watchServiceBrokerFn := servicebroker.NewK8sWatchServiceBrokerFunc(dynamicCl)
		createSvcClassFn := servicebroker.NewK8sCreateServiceClassFunc(dynamicCl)
		if err := servicebroker.RunLoop(
			ctx,
			globalNamespace,
			watchServiceBrokerFn,
			updateServiceBrokerFn,
			cataloger,
			createSvcClassFn,
		); err != nil {
			logger.Errorf("service broker loop (%s)", err)
			errCh <- err
		}
	}()

	// Start service instance loop
	go func() {
		if err := serviceinstance.RunLoop(
			ctx,
			serviceinstance.NewK8sWatchServiceInstanceFunc(dynamicCl),
			serviceinstance.NewK8sUpdateServiceInstanceFunc(dynamicCl),
			refs.NewK8sServiceClassGetterFunc(dynamicCl),
			refs.NewK8sServiceBrokerGetterFunc(dynamicCl),
			lifecycler,
		); err != nil {
			logger.Errorf("service instance loop (%s)", err)
			errCh <- err
		}
	}()

	// start service binding loop
	go func() {
		watchServiceBindingFn := servicebinding.NewK8sWatchServiceBindingFunc(dynamicCl)
		updateServiceBindingFn := servicebinding.NewK8sUpdateServiceBindingFunc(dynamicCl)

		serviceBrokerGetterFn := refs.NewK8sServiceBrokerGetterFunc(dynamicCl)
		svcClassGetterFn := refs.NewK8sServiceClassGetterFunc(dynamicCl)
		svcInstanceGetterFn := refs.NewK8sServiceInstanceGetterFunc(dynamicCl)

		// TODO: remove the hard-coded "default" namespace and instead watch all namespaces and be able
		// to create secrets in any given namespace.
		//
		// See https://github.com/deis/steward-framework/issues/30 and
		// https://github.com/deis/steward-framework/issues/29 for details on how to resolve
		secretWriterFunc := servicebinding.NewK8sSecretWriterFunc(k8sClient.Core().Secrets("default"))
		if err := servicebinding.RunLoop(
			ctx,
			"default",
			lifecycler,
			secretWriterFunc,
			watchServiceBindingFn,
			updateServiceBindingFn,
			serviceBrokerGetterFn,
			svcClassGetterFn,
			svcInstanceGetterFn,
		); err != nil {
			logger.Errorf("service binding loop (%s)", err)
			errCh <- err
		}
	}()

	<-ctx.Done()
}
