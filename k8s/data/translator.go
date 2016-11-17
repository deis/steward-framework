package data

import (
	"encoding/json"

	"k8s.io/client-go/pkg/runtime"
)

// TranslateToTPR translates a runtime.Object into a third party resource given in tpr.
// This code was inspired by the TPRObjectToSCObject function introduced in
// https://github.com/kubernetes-incubator/service-catalog/pull/37/files
func TranslateToTPR(obj runtime.Object, tpr interface{}) error {
	jsonEncoded, err := json.Marshal(obj)
	if err != nil {
		logger.Errorf("encoding to JSON (%s)", err)
		return err
	}

	if err := json.Unmarshal(jsonEncoded, &tpr); err != nil {
		logger.Errorf("decoding object to third party resource (%s)", err)
		return err
	}
	return nil
}

// TranslateToUnstructured converts a runtime.Object into a *runtime.Unstructured.
// This code was inspired by the SCObjectToTPRObject function introduced in
// https://github.com/kubernetes-incubator/service-catalog/pull/37
func TranslateToUnstructured(obj runtime.Object) (*runtime.Unstructured, error) {
	jsonEncoded, err := json.Marshal(obj)
	if err != nil {
		logger.Errorf("encoding to JSON (%s)", err)
		return nil, err
	}
	ret := new(runtime.Unstructured)
	if err := json.Unmarshal(jsonEncoded, ret); err != nil {
		logger.Errorf("decoding object to runtime.Unstructured (%s)", err)
		return nil, err
	}
	return ret, nil
}
