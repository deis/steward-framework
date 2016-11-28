package k8s

import (
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/errors"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	kubeConfigFile = "/go/src/github.com/deis/steward-framework/kubeconfig.yaml"
)

var (
	restCfg     *rest.Config
	clientset   *kubernetes.Clientset
	restCfgOnce sync.Once
	clientOnce  sync.Once
)

func GetRESTConfig() (*rest.Config, error) {
	var err error
	restCfgOnce.Do(func() {
		restCfg, err = clientcmd.BuildConfigFromFlags("", kubeConfigFile)
	})
	return restCfg, err
}

// GetClientset returns a *kubernetes.Clientset suitable for use with integration tests against
// a leased k8s cluster
func GetClientset() (*kubernetes.Clientset, error) {
	var err error
	clientOnce.Do(func() {
		var cfg *rest.Config
		cfg, err = GetRESTConfig()
		if err != nil {
			return
		}
		clientset, err = kubernetes.NewForConfig(cfg)
	})
	return clientset, err
}

// EnsureNamespace ensures the existence of the specified namespace
func EnsureNamespace(namespaceStr string) error {
	clientset, err := GetClientset()
	if err != nil {
		return err
	}
	nsClient := clientset.Namespaces()
	// Just try to create the namespace. If it already exists, that's fine.
	_, err = nsClient.Create(&v1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: namespaceStr,
		},
	})
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// DeleteNamespace deletes the specified namespace
func DeleteNamespace(namespaceStr string) error {
	clientset, err := GetClientset()
	if err != nil {
		return err
	}
	nsClient := clientset.Namespaces()
	// If the problem is just that the namespace doesn't exist, ignore it
	if err := nsClient.Delete(namespaceStr, &v1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}
