package k8sclient

import (
	apiv1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

// GetServiceIngresses returns a ServiceList whose selectors match the labels on passed deployment
func GetServiceIngresses(clientset *kubernetes.Clientset, service apiv1.Service) (*v1beta1.IngressList, error) {
	// get all services
	ingressClient := clientset.ExtensionsV1beta1().Ingresses(service.Namespace)
	ingList, err := ingressClient.List(metav1.ListOptions{FieldSelector: ""})
	if err != nil {
		return nil, err
	}

	result := &v1beta1.IngressList{Items: []v1beta1.Ingress{}}
	for _, ing := range ingList.Items {
		for _, rule := range ing.Spec.Rules {
			for _, path := range rule.HTTP.Paths {
				if service.Name == path.Backend.ServiceName && ServicePortsContains(service.Spec.Ports, path.Backend.ServicePort) {
					result.Items = append(result.Items, ing)
				}
			}
		}
	}
	return result, nil
}

// ServicePortsContains is a helper for determining if a Port exists in a slice of ServicePorts
func ServicePortsContains(servicePorts []apiv1.ServicePort, port intstr.IntOrString) bool {
	for _, p := range servicePorts {
		if port.IntVal == p.Port {
			return true
		}
	}
	return false
}
