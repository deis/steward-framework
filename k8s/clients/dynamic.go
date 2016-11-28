package clients

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/rest"
)

func NewDynamic(cfg rest.Config) (*dynamic.Client, error) {
	cfg.ContentConfig.GroupVersion = &unversioned.GroupVersion{
		Group:   "steward.deis.io",
		Version: "v1",
	}
	cfg.APIPath = "apis"
	return dynamic.NewClient(&cfg)
}
