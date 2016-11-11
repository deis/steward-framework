package restutil

import (
	"strings"
)

const (
	// APIVersionBase is the default base API version to be used in the AbsPath function & similar
	APIVersionBase = "steward.deis.io"
	// APIVersion is the default API version to be used in the AbsPath function & similar
	APIVersion = "v1"
)

// AbsPath returns a slice of strings representing the path to use in a REST call to the Kubernetes
// API
func AbsPath(apiVersionBase, apiVersion string, watch bool, namespace, pluralName string) []string {
	ret := []string{
		"apis",
		apiVersionBase,
		apiVersion,
	}
	if watch {
		ret = append(ret, "watch")
	}
	ret = append(ret, "namespaces", namespace, strings.ToLower(pluralName))
	return ret
}
