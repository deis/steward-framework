package runner

import (
	"context"
	"errors"
	"time"

	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/k8s"
	"github.com/deis/steward-framework/k8s/claim"
	apiserver "github.com/deis/steward-framework/web/api"
	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/pkg/api"
	errs "k8s.io/client-go/1.4/pkg/api/errors"
	"k8s.io/client-go/1.4/rest"
)

func Run(
	brokerName string,
	namespaces []string,
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

	// Create Service Catalog 3PR without bombing out if it already exists.
	extensions := k8sClient.Extensions()
	tpr := extensions.ThirdPartyResources()

	_, err = tpr.Create(k8s.ServiceCatalog3PR)
	if err != nil && !errs.IsAlreadyExists(err) {
		return errCreatingThirdPartyResource{Original: err}
	}

	catalogInteractor := k8s.NewK8sServiceCatalogInteractor(k8sClient.CoreClient.RESTClient)

	rootCtx := context.Background()
	ctx, cancelFn := context.WithCancel(rootCtx)
	defer cancelFn()

	published, err := publishCatalog(ctx, brokerName, cataloger, catalogInteractor)
	if err != nil {
		return errPublishingServiceCatalog{Original: err}
	}
	logger.Infof("published %d entries into the service catalog", len(published))

	evtNamespacer := claim.NewConfigMapsInteractorNamespacer(k8sClient)
	lookup, err := k8s.FetchServiceCatalogLookup(catalogInteractor)
	if err != nil {
		return errGettingServiceCatalogLookupTable{Original: err}
	}
	logger.Infof("created service catalog lookup with %d items", lookup.Len())

	errCh := make(chan error)

	claim.StartControlLoops(
		ctx,
		evtNamespacer,
		k8sClient,
		*lookup,
		lifecycler,
		namespaces,
		maxAsyncDuration,
		errCh,
	)

	go apiserver.Serve(apiPort, errCh)

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

// Does the following:
//
//	1. fetches the service catalog from the backing broker
//	2. checks the 3pr for already-existing entries, and errors if one already exists
//	3. if none error-ed in #2, publishes 3prs for all of the catalog entries
//
// returns all of the entries it wrote into the catalog, or an error
func publishCatalog(ctx context.Context, brokerName string, cataloger framework.Cataloger, catalogEntries k8s.ServiceCatalogInteractor) ([]*k8s.ServiceCatalogEntry, error) {
	services, err := cataloger.List(ctx)
	if err != nil {
		return nil, errGettingServiceCatalog{Original: err}
	}

	published := []*k8s.ServiceCatalogEntry{}
	// Write all entries from cf catalog to 3prs
	for _, service := range services {
		for _, plan := range service.Plans {
			entry := k8s.NewServiceCatalogEntry(brokerName, api.ObjectMeta{}, service.ServiceInfo, plan)
			if _, err := catalogEntries.Create(entry); err != nil {
				logger.Errorf(
					"error publishing catalog entry (svc_name, plan_name) = (%s, %s) (%s), continuing",
					entry.Info.Name,
					entry.Plan.Name,
					err,
				)
				continue
			}
			published = append(published, entry)
		}
	}

	return published, nil
}
