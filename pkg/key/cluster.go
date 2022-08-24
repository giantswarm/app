package key

import (
	"fmt"

	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
)

const (
	ingressControllerConfigMapName = "ingress-controller-values"
	nginxIngressControllerAppName  = "nginx-ingress-controller-app"
)

func ClusterConfigMapName(customResource v1alpha1.App) string {
	// For the org-namespaced apps DO NOT distinguish two cases, namely
	// the NGINX vs other apps, but instead use only one and the same
	// `%s-cluster-values`
	if IsInOrgNamespace(customResource) {
		return fmt.Sprintf("%s-cluster-values", ClusterLabel(customResource))
	}

	// A separate config map is used for Nginx Ingress Controller.
	if AppName(customResource) == nginxIngressControllerAppName {
		return ingressControllerConfigMapName
	}

	return fmt.Sprintf("%s-cluster-values", customResource.Namespace)
}

func ClusterKubeConfigSecretName(customResource v1alpha1.App) string {
	if IsInOrgNamespace(customResource) {
		return fmt.Sprintf("%s-kubeconfig", ClusterLabel(customResource))
	}

	return fmt.Sprintf("%s-kubeconfig", customResource.Namespace)
}
