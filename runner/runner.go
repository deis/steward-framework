package runner

import (
	"context"
	"errors"
	"time"

	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/k8s"
	"github.com/deis/steward-framework/k8s/clients"
	"github.com/deis/steward-framework/k8s/tprs"
	"github.com/deis/steward-framework/web/api"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Run starts all event and control loops. Steward Framework implementations should invoke this
// function LAST in their main() function and can rely upon this function to block program
// their program from exiting until a fatal error is encountered.
func Run(
	cataloger framework.Cataloger,
	lifecycler framework.Lifecycler,
	maxAsyncDuration time.Duration,
	apiPort int,
) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return errGettingK8sConfig{Original: err}
	}
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errGettingK8sClient{Original: err}
	}
	dynClient, err := clients.NewDynamic(*config)
	if err != nil {
		return errGettingDynamicClient{Original: err}
	}

	if err := tprs.Ensure3PRs(k8sClient); err != nil {
		return err
	}

	rootCtx := context.Background()
	ctx, cancelFn := context.WithCancel(rootCtx)
	defer cancelFn()

	errCh := make(chan error)

	go k8s.StartControlLoops(
		ctx,
		k8sClient,
		dynClient,
		cataloger,
		lifecycler,
		k8s.GlobalNamespace,
		errCh,
	)

	go api.Serve(apiPort, errCh)

	select {
	case err := <-errCh:
		if err != nil {
			logger.Criticalf(err.Error())
			return err
		}
		msg := "unknown error, crashing"
		logger.Criticalf(msg)
		return errors.New(msg)
	}
}
