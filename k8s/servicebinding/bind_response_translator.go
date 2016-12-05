package servicebinding

import (
	"fmt"
)

func bindResponseCredsToSecretData(creds map[string]interface{}) map[string][]byte {
	// TODO: make sure the values of the secret are valid DNS subdomains.
	// https://tools.ietf.org/html/rfc4648#section-4 for the spec
	// (steal logic from the k8s client or steward here)
	secretData := make(map[string][]byte, len(creds))
	for k, v := range creds {
		vStr := fmt.Sprintf("%s", v)
		secretData[k] = []byte(vStr)
	}
	return secretData
}
