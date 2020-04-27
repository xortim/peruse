package k8sclient

import (
	"reflect"

	"go.uber.org/zap"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	// ExternalDNSHostnameAnnotation is the external-dns service's annotaiton denoting the hostname
	ExternalDNSHostnameAnnotation = "external-dns.alpha.kunernetes.io/hostname"

	// KubernetesIngressTLSRedirect annotation indicating enforcement of https redirect
	// Traefik and some other ingresses read this
	KubernetesIngressTLSRedirect = "ingress.kubernetes.io/ssl-redirect"

	// NginxIngressTLSRedirect annotation indicating enforcement of https redirect
	NginxIngressTLSRedirect = "nginx.ingress.kubernetes.io/ssl-redirect"

	// AppGWIngressTLSRedirect annotation indicating enforcement of https redirect
	AppGWIngressTLSRedirect = "appgw.ingress.kubernetes.io/ssl-redirect"

	// IngressClassAnnotation determines what ingress controller is responsible for the Ingress
	IngressClassAnnotation = "kubernetes.io/ingress.class"
)

var (
	// IngressTLSRedirectAnnotations collection of annotations whose boolean value enforces https redirect at the Ingress
	IngressTLSRedirectAnnotations = []string{
		KubernetesIngressTLSRedirect,
		NginxIngressTLSRedirect,
		AppGWIngressTLSRedirect,
	}
)

// DeploymentIngressPath represents the deployment -> ingress path.
type DeploymentIngressPath struct {
	Deployment  v1.Deployment
	StatefulSet v1.StatefulSet
	Pods        []apiv1.Pod
	Services    []apiv1.Service
	Ingresses   []v1beta1.Ingress
}

// DeploymentIngressPaths represents a slice of DeploymentIngressPath structs
type DeploymentIngressPaths []DeploymentIngressPath

// NewClient returns a new kubernetes.clientset
func NewClient(masterURL, kubeconfig string) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	config, err = rest.InClusterConfig()
	if err != nil {
		zap.S().Info("could not perform incluster config. falling back to KUBECONFIG")
		config, err = clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	}
	if err != nil {
		zap.S().Error("could not authenticate to cluster\n")
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

// GetDeploymentIngressPaths ...
func GetDeploymentIngressPaths(clientset *kubernetes.Clientset, namespace string) (DeploymentIngressPaths, error) {
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	zap.S().Debugf("Listing deployments in namespace %q\n", namespace)
	dList, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		zap.S().Fatalf(err.Error())
	}

	serviceDeployments, _ := GetServiceDeployments(clientset, namespace)
	if err != nil {
		return nil, err
	}

	dips := DeploymentIngressPaths{}
	for _, deployment := range dList.Items {
		dip := DeploymentIngressPath{}
		dip.Deployment = deployment

		zap.S().Debugf("Getting pods associated with deployment %q\n", dip.Deployment.Name)
		pods, err := DeploymentPods(clientset, dip.Deployment)
		if err != nil {
			return nil, err
		}
		dip.Pods = pods.Items

		zap.S().Debugf("Finding the services that select deployment %q", dip.Deployment.Name)
		for _, sd := range serviceDeployments {
			if sd.SelectsDeployment(deployment) {
				dip.Services = append(dip.Services, sd.Service)
			}
		}

		for _, s := range dip.Services {
			zap.S().Debugf("Finding ingresses that select service %q", s.Name)
			ingresses, _ := GetServiceIngresses(clientset, s)
			dip.Ingresses = append(dip.Ingresses, ingresses.Items...)
		}
		dips = append(dips, dip)
	}
	return dips, nil
}

// ListContains is a helper for determining if a deployment pointer exists in a <T>List
func ListContains(haystack interface{}, needle interface{}) bool {
	ValueIface := reflect.ValueOf(haystack)
	// see if it's a pointer
	if ValueIface.Type().Kind() != reflect.Ptr {
		// convert it to a pointer
		ValueIface = reflect.New(reflect.TypeOf(haystack))
	}

	// Ensure that the passed interface has field Items
	items := ValueIface.Elem().FieldByName("Items")
	if !items.IsValid() {
		return false
	}

	// Ensure that the Items field is a slice and loop over it
	switch reflect.TypeOf(items).Kind() {
	case reflect.Slice:
		for i := 0; i < items.Len(); i++ {
			if reflect.DeepEqual(needle, items.Index(i)) {
				return true
			}
		}
	}

	return false
}
