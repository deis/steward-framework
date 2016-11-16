package binding

import (
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func newFakeSecretWriterFunc(retErr error) (SecretWriterFunc, *[]*apiv1.Secret) {
	secrets := new([]*apiv1.Secret)
	return func(sec *apiv1.Secret) (*apiv1.Secret, error) {
		if retErr != nil {
			return nil, retErr
		}
		*secrets = append(*secrets, sec)
		return sec, nil
	}, secrets
}
