package k8sclient

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetDeploymentServices returns a ServiceList whose selectors match the labels on passed deployment
func GetDeploymentServices(clientset *kubernetes.Clientset, deployment v1.Deployment) (*apiv1.ServiceList, error) {
	// get all services
	servicesClient := clientset.CoreV1().Services(deployment.Namespace)
	deploymentsClient := clientset.AppsV1().Deployments(deployment.Namespace)
	svcList, err := servicesClient.List(metav1.ListOptions{FieldSelector: ""})
	if err != nil {
		return nil, err
	}

	result := &apiv1.ServiceList{Items: []apiv1.Service{}}
	for _, svc := range svcList.Items {

		selector := metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: svc.Spec.Selector})
		if selector == "<none>" {
			continue
		}
		list, err := deploymentsClient.List(metav1.ListOptions{
			LabelSelector: selector,
		})
		if err != nil {
			return result, err
		}

		if DeploymentListContains(list, deployment) {
			result.Items = append(result.Items, svc)
		}
	}

	return result, nil
}

// DeploymentListContains is a helper for determining if a deployment pointer exists in a DeploymentList
func DeploymentListContains(deploymentList *v1.DeploymentList, deployment v1.Deployment) bool {
	for _, d := range deploymentList.Items {
		if reflect.DeepEqual(deployment, d) {
			return true
		}
	}
	return false
}

// DeploymentPods uses the Deployment's spec.Selector.matchLabels field to select pods
func DeploymentPods(clientset *kubernetes.Clientset, deployment v1.Deployment) (*apiv1.PodList, error) {
	return clientset.CoreV1().Pods(deployment.Namespace).List(
		metav1.ListOptions{
			LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: deployment.Spec.Selector.MatchLabels}),
		},
	)
}

// NewTable creates a populated table writer
func (dips DeploymentIngressPaths) NewTable() table.Writer {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Deployment", "Version", "Service", "Ingress"})
	for _, dip := range dips {
		row := table.Row{}

		depStr := []string{"Name: " + dip.Deployment.Name}
		for _, p := range dip.Pods {
			depStr = append(depStr, p.Status.PodIP)
		}
		row = append(row, strings.Join(depStr, "\n"))

		imageStr := []string{}
		for _, container := range dip.Deployment.Spec.Template.Spec.Containers {
			imageStr = append(imageStr, container.Image)
		}
		row = append(row, strings.Join(imageStr, "\n"))

		svcStr := []string{}
		for _, s := range dip.Services {
			svcStr = append(svcStr, fmt.Sprintf("%s", s.ObjectMeta.Name))
		}
		row = append(row, strings.Join(svcStr, "\n"))

		ingStr := []string{}
		for _, ing := range dip.Ingresses {
			ingClass := ing.ObjectMeta.Annotations[IngressClassAnnotation]
			ingStr = append(ingStr, fmt.Sprintf("%s: %s", ing.Name, ingClass))
			for _, url := range IngressURLs(ing) {
				ingStr = append(ingStr, url.String())
			}
		}
		row = append(row, strings.Join(ingStr, "\n"))

		t.AppendRow(row)
	}
	return t
}

// FPrintTable prints the DeploymentIngressPath as an ascii table
func (dips DeploymentIngressPaths) FPrintTable(w io.Writer) {
	t := dips.NewTable()
	t.SetOutputMirror(w)
	t.SetStyle(table.StyleLight)
	t.Render()
}
