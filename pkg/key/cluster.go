package key

import (
	"fmt"

	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
)

func ClusterConfigMapName(customResource v1alpha1.App) string {
	if IsInOrgNamespace(customResource) {
		return fmt.Sprintf("%s-cluster-values", ClusterLabel(customResource))
	}

	return fmt.Sprintf("%s-cluster-values", customResource.Namespace)
}

func ClusterKubeConfigSecretName(customResource v1alpha1.App) string {
	if IsInOrgNamespace(customResource) {
		return fmt.Sprintf("%s-kubeconfig", ClusterLabel(customResource))
	}

	return fmt.Sprintf("%s-kubeconfig", customResource.Namespace)
}
