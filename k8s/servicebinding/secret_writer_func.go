package servicebinding

import (
	"k8s.io/client-go/kubernetes/typed/core/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

// SecretWriterFunc is the function that writes a secret to storage
type SecretWriterFunc func(*apiv1.Secret) (*apiv1.Secret, error)

// NewK8sSecretWriterFunc returns a SecretWriterFunc implemented with a Kubernetes client
func NewK8sSecretWriterFunc(secretsIface v1.SecretInterface) SecretWriterFunc {
	return func(sec *apiv1.Secret) (*apiv1.Secret, error) {
		return secretsIface.Create(sec)
	}
}
