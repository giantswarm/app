package key

import (
	"fmt"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/k8smetadata/pkg/annotation"
	"github.com/giantswarm/k8smetadata/pkg/label"
	"github.com/giantswarm/microerror"
)

const (
	ChartOperatorAppName = "chart-operator"
	// LegacyAppVersionLabel was used for app CRs deployed with Helm 2.
	// We now always default the value for this label.
	LegacyAppVersionLabel = "1.0.0"
	UniqueAppVersionLabel = "0.0.0"
)

func AppConfigMapName(customResource v1alpha1.App) string {
	return customResource.Spec.Config.ConfigMap.Name
}

func AppConfigMapNamespace(customResource v1alpha1.App) string {
	return customResource.Spec.Config.ConfigMap.Namespace
}

func AppName(customResource v1alpha1.App) string {
	return customResource.Spec.Name
}

func AppNamespace(customResource v1alpha1.App) string {
	return customResource.Spec.Namespace
}

func AppNamespaceAnnotation(customResource v1alpha1.App) string {
	return customResource.GetAnnotations()[annotation.AppNamespace]
}

func AppKubernetesNameLabel(customResource v1alpha1.App) string {
	if val, ok := customResource.Labels[label.AppKubernetesName]; ok {
		return val
	}

	return ""
}

func AppLabel(customResource v1alpha1.App) string {
	if val, ok := customResource.Labels[label.App]; ok {
		return val
	}

	return ""
}

func AppNamespaceAnnotations(customResource v1alpha1.App) map[string]string {
	return customResource.Spec.NamespaceConfig.Annotations
}

func AppNamespaceLabels(customResource v1alpha1.App) map[string]string {
	return customResource.Spec.NamespaceConfig.Labels
}

func AppSecretName(customResource v1alpha1.App) string {
	return customResource.Spec.Config.Secret.Name
}

func AppSecretNamespace(customResource v1alpha1.App) string {
	return customResource.Spec.Config.Secret.Namespace
}

func AppStatus(customResource v1alpha1.App) v1alpha1.AppStatus {
	return customResource.Status
}

func AppTeam(customResource v1alpha1.App) string {
	return customResource.Annotations[annotation.AppTeam]
}

func CatalogName(customResource v1alpha1.App) string {
	return customResource.Spec.Catalog
}

func CatalogNamespace(customResource v1alpha1.App) string {
	return customResource.Spec.CatalogNamespace
}

func ClusterID(customResource v1alpha1.App) string {
	return customResource.GetLabels()[label.Cluster]
}

func ClusterLabel(customResource v1alpha1.App) string {
	if val, ok := customResource.Labels[label.Cluster]; ok {
		return val
	}

	return ""
}

func ClusterValuesConfigMapName(customResource v1alpha1.App) string {
	return fmt.Sprintf("%s-cluster-values", customResource.GetNamespace())
}

func ConfigMapExtraConfigs(customResource v1alpha1.App) []v1alpha1.AppExtraConfig {
	extraConfigs := []v1alpha1.AppExtraConfig{}

	for _, v := range customResource.Spec.ExtraConfigs {
		if v.Kind != "secret" {
			extraConfigs = append(extraConfigs, v)
		}
	}

	return extraConfigs
}

func CordonReason(customResource v1alpha1.App) string {
	return customResource.GetAnnotations()[annotation.AppOperatorCordonReason]
}

func CordonUntil(customResource v1alpha1.App) string {
	return customResource.GetAnnotations()[annotation.AppOperatorCordonUntil]
}

func CordonUntilDate() string {
	return time.Now().Add(1 * time.Hour).Format("2006-01-02T15:04:05")
}

func DefaultCatalogStorageURL() string {
	return "https://giantswarm.github.io/default-catalog"
}

func InCluster(customResource v1alpha1.App) bool {
	return customResource.Spec.KubeConfig.InCluster
}

func InstallSkipCRDs(customResource v1alpha1.App) bool {
	return customResource.Spec.Install.SkipCRDs
}

func InstallTimeout(customResource v1alpha1.App) *metav1.Duration {
	return customResource.Spec.Install.Timeout
}

func IsAppCordoned(customResource v1alpha1.App) (bool, error) {
	cordoned := false

	cordonedUntilStr := CordonUntil(customResource)
	if cordonedUntilStr != "" {
		cordonedUntil, err := time.Parse(time.RFC3339, cordonedUntilStr)
		if err != nil {
			return true, err
		}

		cordoned = time.Now().Before(cordonedUntil)
	}

	return cordoned, nil
}

func IsDeleted(customResource v1alpha1.App) bool {
	return customResource.DeletionTimestamp != nil
}

func IsInOrgNamespace(customResource v1alpha1.App) bool {
	return strings.HasPrefix(customResource.Namespace, "org-")
}

func KubeConfigContextName(customResource v1alpha1.App) string {
	return customResource.Spec.KubeConfig.Context.Name
}

func KubeConfigFinalizer(customResource v1alpha1.App) string {
	return fmt.Sprintf("app-operator.giantswarm.io/app-%s", customResource.GetName())
}

func KubeConfigSecretName(customResource v1alpha1.App) string {
	return customResource.Spec.KubeConfig.Secret.Name
}

func KubeConfigSecretNamespace(customResource v1alpha1.App) string {
	return customResource.Spec.KubeConfig.Secret.Namespace
}

func ManagedByLabel(customResource v1alpha1.App) string {
	return customResource.Labels[label.ManagedBy]
}

func Namespace(customResource v1alpha1.App) string {
	return customResource.Spec.Namespace
}

func OrganizationID(customResource v1alpha1.App) string {
	return customResource.GetLabels()[label.Organization]
}

func ReleaseName(customResource v1alpha1.App) string {
	return customResource.Spec.Name
}

func RollbackTimeout(customResource v1alpha1.App) *metav1.Duration {
	return customResource.Spec.Rollback.Timeout
}

func SecretExtraConfigs(customResource v1alpha1.App) []v1alpha1.AppExtraConfig {
	extraConfigs := []v1alpha1.AppExtraConfig{}

	for _, v := range customResource.Spec.ExtraConfigs {
		if v.Kind == "secret" {
			extraConfigs = append(extraConfigs, v)
		}
	}

	return extraConfigs
}

func ToApp(v interface{}) (v1alpha1.App, error) {
	customResource, ok := v.(*v1alpha1.App)
	if !ok {
		return v1alpha1.App{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &v1alpha1.App{}, v)
	}

	if customResource == nil {
		return v1alpha1.App{}, microerror.Maskf(emptyValueError, "empty value cannot be converted to customResource")
	}

	return *customResource, nil
}

func UninstallTimeout(customResource v1alpha1.App) *metav1.Duration {
	return customResource.Spec.Uninstall.Timeout
}

func UpgradeTimeout(customResource v1alpha1.App) *metav1.Duration {
	return customResource.Spec.Upgrade.Timeout
}

func UserConfigMapName(customResource v1alpha1.App) string {
	return customResource.Spec.UserConfig.ConfigMap.Name
}

func ExtraConfigs(customResource v1alpha1.App) []v1alpha1.AppExtraConfig {
	var extraConfigs []v1alpha1.AppExtraConfig
	if customResource.Spec.ExtraConfigs != nil {
		extraConfigs = customResource.Spec.ExtraConfigs
	} else {
		extraConfigs = []v1alpha1.AppExtraConfig{}
	}
	return extraConfigs
}

func UserConfigMapNamespace(customResource v1alpha1.App) string {
	return customResource.Spec.UserConfig.ConfigMap.Namespace
}

func UserSecretName(customResource v1alpha1.App) string {
	return customResource.Spec.UserConfig.Secret.Name
}

func UserSecretNamespace(customResource v1alpha1.App) string {
	return customResource.Spec.UserConfig.Secret.Namespace
}

func Version(customResource v1alpha1.App) string {
	// This enables Flux to use `v`-prefixed image tags as the App
	// CR version to automatically update them.
	return strings.TrimPrefix(customResource.Spec.Version, "v")
}

func VersionLabel(customResource v1alpha1.App) string {
	if val, ok := customResource.Labels[label.AppOperatorVersion]; ok {
		return val
	}

	return ""
}
