package tprs

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/errors"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
	ext "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

var serviceBroker3PR = &ext.ThirdPartyResource{
	TypeMeta: unversioned.TypeMeta{
		Kind:       "ThirdPartyResource",
		APIVersion: "extensions/v1beta1",
	},
	ObjectMeta: v1.ObjectMeta{
		Name: "service-broker.steward.deis.io",
		Labels: map[string]string{
			"heritage": "steward",
		},
	},
	Description: "A service broker which Steward may use to provision and bind to services",
	Versions: []ext.APIVersion{
		{Name: "v1"},
	},
}

var serviceClass3PR = &ext.ThirdPartyResource{
	TypeMeta: unversioned.TypeMeta{
		Kind:       "ThirdPartyResource",
		APIVersion: "extensions/v1beta1",
	},
	ObjectMeta: v1.ObjectMeta{
		Name: "service-class.steward.deis.io",
		Labels: map[string]string{
			"heritage": "steward",
		},
	},
	Description: "A type of service that Steward can provision through an underlying service broker",
	Versions: []ext.APIVersion{
		{Name: "v1"},
	},
}

var serviceInstance3PR = &ext.ThirdPartyResource{
	TypeMeta: unversioned.TypeMeta{
		Kind:       "ThirdPartyResource",
		APIVersion: "extensions/v1beta1",
	},
	ObjectMeta: v1.ObjectMeta{
		Name: "service-instance.steward.deis.io",
		Labels: map[string]string{
			"heritage": "steward",
		},
	},
	Description: "A service instance that Steward will provision through an underlying service broker",
	Versions: []ext.APIVersion{
		{Name: "v1"},
	},
}

var serviceBinding3PR = &ext.ThirdPartyResource{
	TypeMeta: unversioned.TypeMeta{
		Kind:       "ThirdPartyResource",
		APIVersion: "extensions/v1beta1",
	},
	ObjectMeta: v1.ObjectMeta{
		Name: "service-binding.steward.deis.io",
		Labels: map[string]string{
			"heritage": "steward",
		},
	},
	Description: "A binding to a service instance that was provisioned through Steward",
	Versions: []ext.APIVersion{
		{Name: "v1"},
	},
}

func Ensure3PRs(k8sClient *kubernetes.Clientset) error {
	tprs := k8sClient.Extensions().ThirdPartyResources()

	if _, err := tprs.Create(serviceBroker3PR); err != nil && !errors.IsAlreadyExists(err) {
		logger.Errorf("error creating the servicebroker third party resource (%s)", err)
		return errCreatingThirdPartyResource{Original: err}
	}

	if _, err := tprs.Create(serviceClass3PR); err != nil && !errors.IsAlreadyExists(err) {
		logger.Errorf("error creating the service class third party resource (%s)", err)
		return errCreatingThirdPartyResource{Original: err}
	}

	if _, err := tprs.Create(serviceInstance3PR); err != nil && !errors.IsAlreadyExists(err) {
		logger.Errorf("error creating the service instance third party resource (%s)", err)
		return errCreatingThirdPartyResource{Original: err}
	}

	if _, err := tprs.Create(serviceBinding3PR); err != nil && !errors.IsAlreadyExists(err) {
		logger.Errorf("error creating the service binding third party resource (%s)", err)
		return errCreatingThirdPartyResource{Original: err}
	}

	return nil
}
