package data

import (
	"encoding/json"

	"k8s.io/client-go/pkg/runtime"
)

// TranslateToTPR translates a runtime.Object into a third party resource given in tpr.
// This code was inspired by the TPRObjectToSCObject introduced in
// https://github.com/kubernetes-incubator/service-catalog/pull/37/files
func TranslateToTPR(obj runtime.Object, tpr interface{}) error {
	jsonEncoded, err := json.Marshal(obj)
	if err != nil {
		logger.Errorf("failed to encode to JSON (%s)", err)
		return err
	}

	if err := json.Unmarshal(jsonEncoded, &tpr); err != nil {
		logger.Errorf("failed to decode object to third party resource (%s)", err)
		return err
	}
	return nil
}
