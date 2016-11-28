package data

import (
	"encoding/json"
	"fmt"

	"k8s.io/client-go/pkg/runtime"
)

type ErrMismatchedKinds struct {
	RawKind  string
	Expected string
}

func (e ErrMismatchedKinds) Error() string {
	return fmt.Sprintf("unstructured kind %s doesn't match TPR kind %s", e.RawKind, e.Expected)
}

// TranslateToTPR translates a raw runtime.Object into the given tpr. In many cases, this will be
// used to convert a *runtime.Unstructured into a third party resource.
//
// If the translation fails, a non-nil error will be returned. In ese cases, tpr may have been
// written to, but its contents are not guaranteed.
//
// This code was inspired by the TPRObjectToSCObject function introduced in
// https://github.com/kubernetes-incubator/service-catalog/pull/37/files
func TranslateToTPR(raw runtime.Object, tpr runtime.Object, expectedKind string) error {
	jsonEncoded, err := json.Marshal(raw)
	if err != nil {
		logger.Errorf("encoding to JSON (%s)", err)
		return err
	}

	if err := json.Unmarshal(jsonEncoded, &tpr); err != nil {
		logger.Errorf("decoding object to third party resource (%s)", err)
		return err
	}
	rawKind := raw.GetObjectKind().GroupVersionKind().Kind
	if rawKind != expectedKind {
		return ErrMismatchedKinds{RawKind: rawKind, Expected: expectedKind}
	}
	return nil
}

// TranslateToUnstructured converts a runtime.Object into a *runtime.Unstructured.
//
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
