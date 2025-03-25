package key

import (
	"fmt"
	"strings"

	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/k8smetadata/pkg/annotation"
	"github.com/giantswarm/microerror"
)

func ChartConfigMapName(customResource v1alpha1.App) string {
	return fmt.Sprintf("%s-chart-values", customResource.GetName())
}

func ChartName(app v1alpha1.App, clusterID string) string {
	// Chart CR name should match the app CR name when installed in the
	// same cluster.
	if InCluster(app) {
		return app.Name
	}

	// If the app CR has the cluster ID as a prefix or suffix we remove it
	// as its redundant in the remote cluster.
	chartName := strings.TrimPrefix(app.Name, fmt.Sprintf("%s-", clusterID))
	return strings.TrimSuffix(chartName, fmt.Sprintf("-%s", clusterID))
}

func ChartSecretName(customResource v1alpha1.App) string {
	return fmt.Sprintf("%s-chart-secrets", customResource.GetName())
}

func ChartStatus(customResource v1alpha1.Chart) v1alpha1.ChartStatus {
	return customResource.Status
}

func IsChartCordoned(customResource v1alpha1.Chart) bool {
	_, reasonOk := customResource.Annotations[annotation.ChartOperatorCordonReason]
	_, untilOk := customResource.Annotations[annotation.ChartOperatorCordonUntil]

	if reasonOk && untilOk {
		return true
	}

	return false
}

func ToChart(v interface{}) (v1alpha1.Chart, error) {
	if v == nil {
		return v1alpha1.Chart{}, microerror.Maskf(EmptyValueError, "empty value cannot be converted to customResource")
	}

	customResource, ok := v.(*v1alpha1.Chart)
	if !ok {
		return v1alpha1.Chart{}, microerror.Maskf(WrongTypeError, "expected '%T', got '%T'", &v1alpha1.Chart{}, v)
	}

	return *customResource, nil
}
