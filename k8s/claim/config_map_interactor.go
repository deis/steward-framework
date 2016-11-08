package claim

import (
	"context"

	"github.com/deis/steward-framework/k8s"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	v1types "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/watch"
)

type cmInterface struct {
	cm v1.ConfigMapInterface
}

func (c cmInterface) Get(name string) (*k8s.ServicePlanClaimWrapper, error) {
	cm, err := c.cm.Get(name)
	if err != nil {
		return nil, err
	}
	return k8s.ServicePlanClaimWrapperFromConfigMap(cm)
}

func (c cmInterface) List(opts v1types.ListOptions) (*k8s.ServicePlanClaimsListWrapper, error) {
	cms, err := c.cm.List(opts)
	if err != nil {
		return nil, err
	}
	claims := make([]*k8s.ServicePlanClaimWrapper, len(cms.Items))
	for i, cm := range cms.Items {
		wr, err := k8s.ServicePlanClaimWrapperFromConfigMap(&cm)
		if err != nil {
			return nil, err
		}
		claims[i] = wr
	}
	return &k8s.ServicePlanClaimsListWrapper{
		ResourceVersion: cms.ResourceVersion,
		Claims:          claims,
	}, nil
}

func (c cmInterface) Update(spc *k8s.ServicePlanClaimWrapper) (*k8s.ServicePlanClaimWrapper, error) {
	cm := &v1types.ConfigMap{
		Data:       spc.Claim.ToMap(),
		ObjectMeta: spc.ObjectMeta,
	}
	logger.Debugf("updating ConfigMap %s with data %s", cm.Name, cm.Data)
	newCM, err := c.cm.Update(cm)
	if err != nil {
		return nil, err
	}
	return k8s.ServicePlanClaimWrapperFromConfigMap(newCM)
}

func (c cmInterface) Watch(ctx context.Context, opts v1types.ListOptions) Watcher {
	return newConfigMapWatcher(ctx, func() (watch.Interface, error) {
		return c.cm.Watch(opts)
	})
}
