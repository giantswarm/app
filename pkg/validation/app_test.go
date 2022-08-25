package validation

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/k8smetadata/pkg/annotation"
	"github.com/giantswarm/k8smetadata/pkg/label"
	"github.com/giantswarm/micrologger/microloggertest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgofake "k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/fake" //nolint:staticcheck
)

func Test_ValidateApp(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name                 string
		obj                  v1alpha1.App
		catalogs             []*v1alpha1.Catalog
		configMaps           []*corev1.ConfigMap
		secrets              []*corev1.Secret
		enableManagedByLabel bool
		expectedErr          string
	}{
		{
			name: "flawless flow",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
					Labels: map[string]string{
						label.AppOperatorVersion: "2.6.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					Config: v1alpha1.AppSpecConfig{
						ConfigMap: v1alpha1.AppSpecConfigConfigMap{
							Name:      "eggs2-cluster-values",
							Namespace: "eggs2",
						},
					},
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						Context: v1alpha1.AppSpecKubeConfigContext{
							Name: "eggs2-kubeconfig",
						},
						InCluster: false,
						Secret: v1alpha1.AppSpecKubeConfigSecret{
							Name:      "eggs2-kubeconfig",
							Namespace: "eggs2",
						},
					},
					UserConfig: v1alpha1.AppSpecUserConfig{
						ConfigMap: v1alpha1.AppSpecUserConfigConfigMap{
							Name:      "kiam-user-values",
							Namespace: "eggs2",
						},
					},
					Version: "1.4.0",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("giantswarm", "default"),
			},
			configMaps: []*corev1.ConfigMap{
				newTestConfigMap("eggs2-cluster-values", "eggs2"),
				newTestConfigMap("kiam-user-values", "eggs2"),
			},
			secrets: []*corev1.Secret{
				newTestSecret("eggs2-kubeconfig", "eggs2"),
			},
		},
		{
			name: "flawless org-managed flow",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "org-eggs2",
					Labels: map[string]string{
						label.Cluster: "eggs2",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					Config: v1alpha1.AppSpecConfig{
						ConfigMap: v1alpha1.AppSpecConfigConfigMap{
							Name:      "eggs2-cluster-values",
							Namespace: "org-eggs2",
						},
					},
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						Context: v1alpha1.AppSpecKubeConfigContext{
							Name: "eggs2-kubeconfig",
						},
						InCluster: false,
						Secret: v1alpha1.AppSpecKubeConfigSecret{
							Name:      "eggs2-kubeconfig",
							Namespace: "org-eggs2",
						},
					},
					UserConfig: v1alpha1.AppSpecUserConfig{
						ConfigMap: v1alpha1.AppSpecUserConfigConfigMap{
							Name:      "kiam-user-values",
							Namespace: "org-eggs2",
						},
					},
					Version: "1.4.0",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("giantswarm", "default"),
			},
			configMaps: []*corev1.ConfigMap{
				newTestConfigMap("eggs2-cluster-values", "org-eggs2"),
				newTestConfigMap("kiam-user-values", "org-eggs2"),
			},
			secrets: []*corev1.Secret{
				newTestSecret("eggs2-kubeconfig", "org-eggs2"),
			},
		},
		{
			name: "flawless in-cluster",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("control-plane-catalog", "giantswarm"),
			},
		},
		{
			name: "missing version label",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("control-plane-catalog", "giantswarm"),
			},
			expectedErr: "validation error: label `app-operator.giantswarm.io/version` not found",
		},
		{
			name: "missing cluster label",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "org-eggs2",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						Context: v1alpha1.AppSpecKubeConfigContext{
							Name: "eggs2-kubeconfig",
						},
						InCluster: false,
						Secret: v1alpha1.AppSpecKubeConfigSecret{
							Name:      "eggs2-kubeconfig",
							Namespace: "org-eggs2",
						},
					},
					Version: "1.4.0",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("giantswarm", "default"),
			},
			secrets: []*corev1.Secret{
				newTestSecret("eggs2-kubeconfig", "org-eggs2"),
			},
			expectedErr: "validation error: label `giantswarm.io/cluster` not found",
		},
		{
			name: "spec.catalog not found",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "missing",
					Name:      "dex-app",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("control-plane-catalog", "giantswarm"),
			},
			expectedErr: "validation error: catalog `missing` not found",
		},
		{
			name: "spec.config.configMap not found",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					Config: v1alpha1.AppSpecConfig{
						ConfigMap: v1alpha1.AppSpecConfigConfigMap{
							Name:      "dex-app-values",
							Namespace: "giantswarm",
						},
					},
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("control-plane-catalog", "giantswarm"),
			},
			expectedErr: "app config map not found error: configmap `dex-app-values` in namespace `giantswarm` not found",
		},
		{
			name: "spec.config.configMap no namespace specified",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					Config: v1alpha1.AppSpecConfig{
						ConfigMap: v1alpha1.AppSpecConfigConfigMap{
							Name:      "dex-app-values",
							Namespace: "",
						},
					},
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("control-plane-catalog", "giantswarm"),
			},
			expectedErr: "validation error: namespace is not specified for configmap `dex-app-values`",
		},
		{
			name: "spec.config.secret not found",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					Config: v1alpha1.AppSpecConfig{
						Secret: v1alpha1.AppSpecConfigSecret{
							Name:      "dex-app-secrets",
							Namespace: "giantswarm",
						},
					},
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("control-plane-catalog", "giantswarm"),
			},
			expectedErr: "validation error: secret `dex-app-secrets` in namespace `giantswarm` not found",
		},
		{
			name: "spec.config.secret no namespace specified",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					Config: v1alpha1.AppSpecConfig{
						Secret: v1alpha1.AppSpecConfigSecret{
							Name:      "dex-app-secrets",
							Namespace: "",
						},
					},
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("control-plane-catalog", "giantswarm"),
			},
			expectedErr: "validation error: namespace is not specified for secret `dex-app-secrets`",
		},
		{
			name: "spec.kubeConfig.secret not found",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
					Labels: map[string]string{
						label.AppOperatorVersion: "2.6.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						Context: v1alpha1.AppSpecKubeConfigContext{
							Name: "eggs2-kubeconfig",
						},
						InCluster: false,
						Secret: v1alpha1.AppSpecKubeConfigSecret{
							Name:      "eggs2-kubeconfig",
							Namespace: "eggs2",
						},
					},
					Version: "1.4.0",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("giantswarm", "default"),
			},
			expectedErr: "kube config not found error: kubeconfig secret `eggs2-kubeconfig` in namespace `eggs2` not found",
		},
		{
			name: "missing spec.kubeConfig.secret allowed when managed by Flux on conditional validation",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
					Labels: map[string]string{
						label.AppOperatorVersion: "2.6.0",
						label.ManagedBy:          "flux",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						Context: v1alpha1.AppSpecKubeConfigContext{
							Name: "eggs2-kubeconfig",
						},
						InCluster: false,
						Secret: v1alpha1.AppSpecKubeConfigSecret{
							Name:      "eggs2-kubeconfig",
							Namespace: "eggs2",
						},
					},
					Version: "1.4.0",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("giantswarm", "default"),
			},
			enableManagedByLabel: true,
		},
		{
			name: "missing spec.kubeConfig.secret not allowed when managed by Flux on uncoditional validation",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
					Labels: map[string]string{
						label.AppOperatorVersion: "2.6.0",
						label.ManagedBy:          "flux",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						Context: v1alpha1.AppSpecKubeConfigContext{
							Name: "eggs2-kubeconfig",
						},
						InCluster: false,
						Secret: v1alpha1.AppSpecKubeConfigSecret{
							Name:      "eggs2-kubeconfig",
							Namespace: "eggs2",
						},
					},
					Version: "1.4.0",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("giantswarm", "default"),
			},
			expectedErr: "kube config not found error: kubeconfig secret `eggs2-kubeconfig` in namespace `eggs2` not found",
		},
		{
			name: "spec.kubeConfig.secret no namespace specified",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
					Labels: map[string]string{
						label.AppOperatorVersion: "2.6.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						Context: v1alpha1.AppSpecKubeConfigContext{
							Name: "eggs2-kubeconfig",
						},
						InCluster: false,
						Secret: v1alpha1.AppSpecKubeConfigSecret{
							Name:      "eggs2-kubeconfig",
							Namespace: "",
						},
					},
					Version: "1.4.0",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("giantswarm", "default"),
			},
			expectedErr: "validation error: namespace is not specified for kubeconfig secret `eggs2-kubeconfig`",
		},
		{
			name: "spec.userConfig.configMap not found",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					UserConfig: v1alpha1.AppSpecUserConfig{
						ConfigMap: v1alpha1.AppSpecUserConfigConfigMap{
							Name:      "dex-app-user-values",
							Namespace: "giantswarm",
						},
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("control-plane-catalog", "default"),
			},
			expectedErr: "validation error: configmap `dex-app-user-values` in namespace `giantswarm` not found",
		},
		{
			name: "missing spec.userConfig.configMap allowed when managed by Flux on conditional validation",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
						label.ManagedBy:          "flux",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					UserConfig: v1alpha1.AppSpecUserConfig{
						ConfigMap: v1alpha1.AppSpecUserConfigConfigMap{
							Name:      "dex-app-user-values",
							Namespace: "giantswarm",
						},
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("control-plane-catalog", "default"),
			},
			enableManagedByLabel: true,
		},
		{
			name: "missing spec.userConfig.configMap not allowed when managed by Flux on unconditional validation",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
						label.ManagedBy:          "flux",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					UserConfig: v1alpha1.AppSpecUserConfig{
						ConfigMap: v1alpha1.AppSpecUserConfigConfigMap{
							Name:      "dex-app-user-values",
							Namespace: "giantswarm",
						},
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("control-plane-catalog", "default"),
			},
			expectedErr: "validation error: configmap `dex-app-user-values` in namespace `giantswarm` not found",
		},
		{
			name: "spec.userConfig.configMap no namespace specified",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					UserConfig: v1alpha1.AppSpecUserConfig{
						ConfigMap: v1alpha1.AppSpecUserConfigConfigMap{
							Name:      "dex-app-user-values",
							Namespace: "",
						},
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("control-plane-catalog", "giantswarm"),
			},
			expectedErr: "validation error: namespace is not specified for configmap `dex-app-user-values`",
		},
		{
			name: "spec.userConfig.secret not found",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					UserConfig: v1alpha1.AppSpecUserConfig{
						Secret: v1alpha1.AppSpecUserConfigSecret{
							Name:      "dex-app-user-secrets",
							Namespace: "giantswarm",
						},
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("control-plane-catalog", "giantswarm"),
			},
			expectedErr: "validation error: secret `dex-app-user-secrets` in namespace `giantswarm` not found",
		},
		{
			name: "missing spec.userConfig.secret allowed when managed by Flux on conditional validation",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
						label.ManagedBy:          "flux",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					UserConfig: v1alpha1.AppSpecUserConfig{
						Secret: v1alpha1.AppSpecUserConfigSecret{
							Name:      "dex-app-user-secrets",
							Namespace: "giantswarm",
						},
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("control-plane-catalog", "giantswarm"),
			},
			enableManagedByLabel: true,
		},
		{
			name: "missing spec.userConfig.secret not allowed when managed by Flux on unconditional validation",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
						label.ManagedBy:          "flux",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					UserConfig: v1alpha1.AppSpecUserConfig{
						Secret: v1alpha1.AppSpecUserConfigSecret{
							Name:      "dex-app-user-secrets",
							Namespace: "giantswarm",
						},
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("control-plane-catalog", "giantswarm"),
			},
			expectedErr: "validation error: secret `dex-app-user-secrets` in namespace `giantswarm` not found",
		},
		{
			name: "spec.userConfig.secret no namespace specified",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dex-app-unique",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "control-plane-catalog",
					Name:      "dex-app",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					UserConfig: v1alpha1.AppSpecUserConfig{
						Secret: v1alpha1.AppSpecUserConfigSecret{
							Name:      "dex-app-user-secrets",
							Namespace: "",
						},
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("control-plane-catalog", "giantswarm"),
			},
			expectedErr: "validation error: namespace is not specified for secret `dex-app-user-secrets`",
		},
		{
			name: "spec.userConfig.configMap.name incorrect for default catalog app",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "default",
					Name:      "kiam-app",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					UserConfig: v1alpha1.AppSpecUserConfig{
						ConfigMap: v1alpha1.AppSpecUserConfigConfigMap{
							Name:      "user-values",
							Namespace: "",
						},
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("default", "giantswarm"),
			},
			expectedErr: "validation error: user configmap must be named `kiam-user-values` for app in default catalog",
		},
		{
			name: "spec.userConfig.secret.name incorrect for default catalog app",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "default",
					Name:      "kiam-app",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					UserConfig: v1alpha1.AppSpecUserConfig{
						Secret: v1alpha1.AppSpecUserConfigSecret{
							Name:      "user-secrets",
							Namespace: "",
						},
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("default", "giantswarm"),
			},
			expectedErr: "validation error: user secret must be named `kiam-user-secrets` for app in default catalog",
		},
		{
			name: "metadata.name exceeds max length",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster-autoscaler-1.2.2-2b060b8bda545a7b6aeff1b8cb13951181ae30d3",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "default",
					Name:      "cluster-autoscaler",
					Namespace: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					Version: "1.2.2",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("default", "giantswarm"),
			},
			expectedErr: "validation error: name `cluster-autoscaler-1.2.2-2b060b8bda545a7b6aeff1b8cb13951181ae30d3` is 65 chars and exceeds max length of 53 chars",
		},
		{
			name: "legacy version label is rejected",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
					Labels: map[string]string{
						label.AppOperatorVersion: "1.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					Version: "1.4.0",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("giantswarm", "default"),
			},
			expectedErr: "validation error: label `app-operator.giantswarm.io/version` has invalid value `1.0.0`",
		},
		{
			name: "nginx user values configmap name is not restricted",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "giantswarm",
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "default",
					Name:      "nginx-ingress-controller-app",
					Namespace: "kube-system",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					UserConfig: v1alpha1.AppSpecUserConfig{
						ConfigMap: v1alpha1.AppSpecUserConfigConfigMap{
							Name:      "nginx-ingress-user-values",
							Namespace: "eggs2",
						},
					},
					Version: "1.4.0",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("default", "default"),
			},
			configMaps: []*corev1.ConfigMap{
				newTestConfigMap("nginx-ingress-user-values", "eggs2"),
			},
		},
		{
			name: "spec.kubeConfig.secret no name specified",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
					Labels: map[string]string{
						label.AppOperatorVersion: "2.6.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						Context: v1alpha1.AppSpecKubeConfigContext{
							Name: "eggs2-kubeconfig",
						},
						InCluster: false,
						Secret: v1alpha1.AppSpecKubeConfigSecret{
							Name:      "",
							Namespace: "default",
						},
					},
					Version: "1.4.0",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("giantswarm", "default"),
			},
			expectedErr: "validation error: name is not specified for kubeconfig secret",
		},
		{
			name: ".spec.namespace for in-cluster app not allowed outside org namespace",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "org-eggs2",
					Labels: map[string]string{
						label.AppOperatorVersion: "2.6.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					Version: "1.4.0",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("giantswarm", "default"),
			},
			expectedErr: "validation error: target namespace kube-system is not allowed",
		},
		{
			name: ".spec.namespace for in-cluster app not allowed outside WC namespace",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
					Labels: map[string]string{
						label.AppOperatorVersion: "2.6.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					Version: "1.4.0",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("giantswarm", "default"),
			},
			expectedErr: "validation error: target namespace kube-system is not allowed",
		},
		{
			name: ".spec.namespace for in-cluster app allowed when it matches org namespace",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "org-eggs2",
					Labels: map[string]string{
						label.AppOperatorVersion: "2.6.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "org-eggs2",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					Version: "1.4.0",
				},
			},
			catalogs: []*v1alpha1.Catalog{
				newTestCatalog("giantswarm", "default"),
			},
		},
		{
			name: "mismatch in annotation namespace is not allowed",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "hello-world",
					Namespace: "demo0",
					Annotations: map[string]string{
						annotation.AppNamespace: "giantswarm",
					},
					Labels: map[string]string{
						label.AppOperatorVersion: "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "hello-world",
					Namespace: "demo0",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
					Version: "0.3.0",
				},
			},
			expectedErr: "validation error: wrong `giantswarm` namespace for the `chart-operator.giantswarm.io/app-namespace` annotation",
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("case %d: %s", i, tc.name), func(t *testing.T) {
			g8sObjs := make([]runtime.Object, 0)
			for _, cat := range tc.catalogs {
				g8sObjs = append(g8sObjs, cat)
			}

			k8sObjs := make([]runtime.Object, 0)
			for _, cm := range tc.configMaps {
				k8sObjs = append(k8sObjs, cm)
			}

			for _, secret := range tc.secrets {
				k8sObjs = append(k8sObjs, secret)
			}

			scheme := runtime.NewScheme()
			_ = v1alpha1.AddToScheme(scheme)

			fakeCtrlClient := fake.NewFakeClientWithScheme(scheme, g8sObjs...)

			c := Config{
				G8sClient: fakeCtrlClient,
				K8sClient: clientgofake.NewSimpleClientset(k8sObjs...),
				Logger:    microloggertest.New(),

				ProjectName:          "app-admission-controller",
				Provider:             "aws",
				EnableManagedByLabel: tc.enableManagedByLabel,
			}
			r, err := NewValidator(c)
			if err != nil {
				t.Fatalf("error == %#v, want nil", err)
			}

			_, err = r.ValidateApp(ctx, tc.obj)
			switch {
			case err != nil && tc.expectedErr == "":
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.expectedErr != "":
				t.Fatalf("error == nil, want non-nil")
			}

			if err != nil && tc.expectedErr != "" {
				if !strings.Contains(err.Error(), tc.expectedErr) {
					t.Fatalf("error == %#v, want %#v ", err.Error(), tc.expectedErr)
				}

			}
		})
	}
}

func Test_ValidateAppUpdate(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		obj         v1alpha1.App
		currentApp  v1alpha1.App
		expectedErr string
	}{
		{
			name: "case 0: flawless",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					Version:   "1.4.0",
				},
			},
			currentApp: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					Version:   "1.4.0",
				},
			},
		},
		{
			name: "case 1: changed namespace is rejected",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "default",
					Version:   "1.4.0",
				},
			},
			currentApp: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					Version:   "1.4.0",
				},
			},
			expectedErr: "validation error: target namespace for app `kiam` cannot be changed from `kube-system` to `default`",
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			scheme := runtime.NewScheme()
			_ = v1alpha1.AddToScheme(scheme)

			fakeCtrlClient := fake.NewFakeClientWithScheme(scheme)

			c := Config{
				G8sClient: fakeCtrlClient,
				K8sClient: clientgofake.NewSimpleClientset(),
				Logger:    microloggertest.New(),

				ProjectName: "app-admission-controller",
				Provider:    "aws",
			}
			r, err := NewValidator(c)
			if err != nil {
				t.Fatalf("error == %#v, want nil", err)
			}

			_, err = r.ValidateAppUpdate(ctx, tc.obj, tc.currentApp)
			switch {
			case err != nil && tc.expectedErr == "":
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.expectedErr != "":
				t.Fatalf("error == nil, want non-nil")
			}

			if err != nil && tc.expectedErr != "" {
				if !strings.Contains(err.Error(), tc.expectedErr) {
					t.Fatalf("error == %#v, want %#v ", err.Error(), tc.expectedErr)
				}

			}
		})
	}
}

func Test_ValidateMetadataConstraints(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		obj          v1alpha1.App
		catalogEntry *v1alpha1.AppCatalogEntry
		apps         []*v1alpha1.App
		expectedErr  string
	}{
		{
			name: "case 0: flawless flow",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
					Labels: map[string]string{
						label.AppOperatorVersion: "2.6.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: metav1.NamespaceDefault,
					Version:   "1.4.0",
				},
			},
			catalogEntry: &v1alpha1.AppCatalogEntry{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "giantswarm-kiam-1.4.0",
					Namespace: metav1.NamespaceDefault,
				},
				Spec: v1alpha1.AppCatalogEntrySpec{
					Restrictions: &v1alpha1.AppCatalogEntrySpecRestrictions{
						FixedNamespace: metav1.NamespaceDefault,
					},
				},
			},
		},
		{
			name: "case 1: fixed namespace constraint",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					Version:   "1.4.0",
				},
			},
			catalogEntry: &v1alpha1.AppCatalogEntry{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "giantswarm-kiam-1.4.0",
					Namespace: metav1.NamespaceDefault,
				},
				Spec: v1alpha1.AppCatalogEntrySpec{
					Restrictions: &v1alpha1.AppCatalogEntrySpecRestrictions{
						FixedNamespace: "eggs1",
					},
				},
			},
			expectedErr: "validation error: app `kiam` can only be installed in namespace `eggs1` only, not `kube-system`",
		},
		{
			name: "case 2: cluster singleton constraint",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					Version:   "1.4.0",
				},
			},
			apps: []*v1alpha1.App{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "another-kiam",
						Namespace: "eggs2",
					},
					Spec: v1alpha1.AppSpec{
						Catalog:   "giantswarm",
						Name:      "kiam",
						Namespace: "giantswarm",
						Version:   "1.3.0-rc1",
					},
				},
			},
			catalogEntry: &v1alpha1.AppCatalogEntry{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "giantswarm-kiam-1.4.0",
					Namespace: metav1.NamespaceDefault,
				},
				Spec: v1alpha1.AppCatalogEntrySpec{
					Restrictions: &v1alpha1.AppCatalogEntrySpecRestrictions{
						ClusterSingleton: true,
					},
				},
			},
			expectedErr: "validation error: app `kiam` can only be installed once in cluster `eggs2`",
		},
		{
			name: "case 3: namespace singleton constraint",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					Version:   "1.4.0",
				},
			},
			apps: []*v1alpha1.App{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "another-kiam",
						Namespace: "eggs2",
					},
					Spec: v1alpha1.AppSpec{
						Catalog:   "giantswarm",
						Name:      "kiam",
						Namespace: "giantswarm",
						Version:   "1.3.0-rc1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "another-kiam-1",
						Namespace: "eggs2",
					},
					Spec: v1alpha1.AppSpec{
						Catalog:   "giantswarm",
						Name:      "kiam",
						Namespace: "kube-system",
						Version:   "1.3.0-rc1",
					},
				},
			},
			catalogEntry: &v1alpha1.AppCatalogEntry{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "giantswarm-kiam-1.4.0",
					Namespace: metav1.NamespaceDefault,
				},
				Spec: v1alpha1.AppCatalogEntrySpec{
					Restrictions: &v1alpha1.AppCatalogEntrySpecRestrictions{
						NamespaceSingleton: true,
					},
				},
			},
			expectedErr: "validation error: app `kiam` can only be installed only once in namespace `kube-system`",
		},
		{
			name: "case 4: compatible providers constraint",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					Version:   "1.4.0",
				},
			},
			catalogEntry: &v1alpha1.AppCatalogEntry{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "giantswarm-kiam-1.4.0",
					Namespace: metav1.NamespaceDefault,
				},
				Spec: v1alpha1.AppCatalogEntrySpec{
					Restrictions: &v1alpha1.AppCatalogEntrySpecRestrictions{
						CompatibleProviders: []string{"azure"},
					},
				},
			},
			expectedErr: "validation error: app `kiam` can only be installed for providers [`azure`] not `aws`",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g8sObjs := make([]runtime.Object, 0)

			if tc.catalogEntry != nil {
				g8sObjs = append(g8sObjs, tc.catalogEntry)
			}

			for _, app := range tc.apps {
				g8sObjs = append(g8sObjs, app)
			}

			scheme := runtime.NewScheme()
			_ = v1alpha1.AddToScheme(scheme)

			fakeCtrlClient := fake.NewFakeClientWithScheme(scheme, g8sObjs...)

			c := Config{
				G8sClient: fakeCtrlClient,
				K8sClient: clientgofake.NewSimpleClientset(),
				Logger:    microloggertest.New(),

				ProjectName: "app-admission-controller",
				Provider:    "aws",
			}
			r, err := NewValidator(c)
			if err != nil {
				t.Fatalf("error == %#v, want nil", err)
			}

			err = r.validateMetadataConstraints(ctx, tc.obj)
			switch {
			case err != nil && tc.expectedErr == "":
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.expectedErr != "":
				t.Fatalf("error == nil, want non-nil")
			}

			if err != nil && tc.expectedErr != "" {
				if !strings.Contains(err.Error(), tc.expectedErr) {
					t.Fatalf("error == %#v, want %#v ", err.Error(), tc.expectedErr)
				}

			}
		})
	}
}

func Test_ValidateNamespace(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		obj         v1alpha1.App
		apps        []*v1alpha1.App
		expectedErr string
	}{
		{
			name: "case 0: flawless",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					NamespaceConfig: v1alpha1.AppSpecNamespaceConfig{
						Annotations: map[string]string{
							"linkerd.io/inject": "enabled",
						},
					},
					Version: "1.4.0",
				},
			},
			apps: []*v1alpha1.App{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "another-kiam-1",
						Namespace: "eggs2",
					},
					Spec: v1alpha1.AppSpec{
						Catalog:   "giantswarm",
						Name:      "kiam",
						Namespace: "kube-system",
						NamespaceConfig: v1alpha1.AppSpecNamespaceConfig{
							Annotations: map[string]string{
								"linkerd.io/inject": "enabled",
							},
						},
						Version: "1.3.0-rc1",
					},
				},
			},
		},
		{
			name: "case 1: namespace annotation collision",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					NamespaceConfig: v1alpha1.AppSpecNamespaceConfig{
						Annotations: map[string]string{
							"linkerd.io/inject": "enabled",
						},
					},
					Version: "1.4.0",
				},
			},
			apps: []*v1alpha1.App{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "another-kiam-1",
						Namespace: "eggs2",
					},
					Spec: v1alpha1.AppSpec{
						Catalog:   "giantswarm",
						Name:      "kiam",
						Namespace: "kube-system",
						NamespaceConfig: v1alpha1.AppSpecNamespaceConfig{
							Annotations: map[string]string{
								"linkerd.io/inject": "disabled",
							},
						},
						Version: "1.3.0-rc1",
					},
				},
			},
			expectedErr: "app `kiam` annotation `linkerd.io/inject` for target namespace `kube-system` collides with value `disabled` for app `another-kiam-1`",
		},
		{
			name: "case 2: namespace label collision",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kiam",
					Namespace: "eggs2",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "kiam",
					Namespace: "kube-system",
					NamespaceConfig: v1alpha1.AppSpecNamespaceConfig{
						Labels: map[string]string{
							"monitoring": "enabled",
						},
					},
					Version: "1.4.0",
				},
			},
			apps: []*v1alpha1.App{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "another-kiam-1",
						Namespace: "eggs2",
					},
					Spec: v1alpha1.AppSpec{
						Catalog:   "giantswarm",
						Name:      "kiam",
						Namespace: "kube-system",
						NamespaceConfig: v1alpha1.AppSpecNamespaceConfig{
							Labels: map[string]string{
								"monitoring": "disabled",
							},
						},
						Version: "1.3.0-rc1",
					},
				},
			},
			expectedErr: "app `kiam` label `monitoring` for target namespace `kube-system` collides with value `disabled` for app `another-kiam-1`",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g8sObjs := make([]runtime.Object, 0)

			for _, app := range tc.apps {
				g8sObjs = append(g8sObjs, app)
			}

			scheme := runtime.NewScheme()
			_ = v1alpha1.AddToScheme(scheme)

			fakeCtrlClient := fake.NewFakeClientWithScheme(scheme, g8sObjs...)

			c := Config{
				G8sClient: fakeCtrlClient,
				K8sClient: clientgofake.NewSimpleClientset(),
				Logger:    microloggertest.New(),

				ProjectName: "app-admission-controller",
				Provider:    "aws",
			}
			r, err := NewValidator(c)
			if err != nil {
				t.Fatalf("error == %#v, want nil", err)
			}

			err = r.validateNamespaceConfig(ctx, tc.obj)
			switch {
			case err != nil && tc.expectedErr == "":
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.expectedErr != "":
				t.Fatalf("error == nil, want non-nil")
			}

			if err != nil && tc.expectedErr != "" {
				if !strings.Contains(err.Error(), tc.expectedErr) {
					t.Fatalf("error == %#v, want %#v ", err.Error(), tc.expectedErr)
				}

			}
		})
	}
}

func Test_ValidateUniqueInClusterAppName(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		obj         v1alpha1.App
		apps        []*v1alpha1.App
		expectedErr string
	}{
		{
			name: "case 0: not an in-cluster app",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "not-security-pack",
					Namespace: "abc01",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "not-security-pack",
					Namespace: "abc01",
					Version:   "1.2.0",
				},
			},
			apps: []*v1alpha1.App{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "not-security-pack",
						Namespace: "another-namespace",
					},
					Spec: v1alpha1.AppSpec{
						Catalog:   "giantswarm",
						Name:      "not-security-pack",
						Namespace: "another-namespace",
						Version:   "1.2.0",
					},
				},
			},
		},
		{
			name: "case 1: in-cluster app with a non-unique name",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "security-pack",
					Namespace: "abc01",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "security-pack",
					Namespace: "abc01",
					Version:   "1.2.0",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
				},
			},
			apps: []*v1alpha1.App{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "security-pack",
						Namespace: "another-namespace",
					},
					Spec: v1alpha1.AppSpec{
						Catalog:   "giantswarm",
						Name:      "security-pack",
						Namespace: "another-namespace",
						Version:   "1.2.0",
						KubeConfig: v1alpha1.AppSpecKubeConfig{
							InCluster: true,
						},
					},
				},
			},
			expectedErr: "in-cluster apps must be given a unique name, found an app named `security-pack` as well in the `another-namespace` namespace",
		},
		{
			name: "case 3: in-cluster app with a name is updated",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "security-pack",
					Namespace: "abc01",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "security-pack",
					Namespace: "abc01",
					// The version is updated, the app already exists
					Version: "1.3.0",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
				},
			},
			apps: []*v1alpha1.App{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "security-pack",
						Namespace: "abc01",
					},
					Spec: v1alpha1.AppSpec{
						Catalog:   "giantswarm",
						Name:      "security-pack",
						Namespace: "abc01",
						Version:   "1.2.0",
						KubeConfig: v1alpha1.AppSpecKubeConfig{
							InCluster: true,
						},
					},
				},
			},
		},
		{
			name: "case 4: in-cluster app with no name collision",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "abc01-security-pack",
					Namespace: "abc01",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "security-pack",
					Namespace: "abc01",
					Version:   "1.2.0",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
				},
			},
			apps: []*v1alpha1.App{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "qwe99-security-pack",
						Namespace: "qwe99",
					},
					Spec: v1alpha1.AppSpec{
						Catalog:   "giantswarm",
						Name:      "security-pack",
						Namespace: "qwe99",
						Version:   "1.2.0",
						KubeConfig: v1alpha1.AppSpecKubeConfig{
							InCluster: true,
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "totally-unrelated-app",
						Namespace: "xyz",
					},
					Spec: v1alpha1.AppSpec{
						Catalog:   "giantswarm-test",
						Name:      "random-other-app",
						Namespace: "xyz",
						Version:   "12.3.8",
					},
				},
			},
		},
		{
			name: "case 5: there is another app with the same name in a different namespace but it is not an in-cluster app",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "security-pack",
					Namespace: "abc01",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "security-pack",
					Namespace: "abc01",
					Version:   "1.2.0",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
				},
			},
			apps: []*v1alpha1.App{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "security-pack",
						Namespace: "qwe99",
					},
					Spec: v1alpha1.AppSpec{
						Catalog:   "giantswarm",
						Name:      "another-app-name",
						Namespace: "qwe99",
						Version:   "42.0.1",
					},
				},
			},
		},
		{
			name: "case 6: cover the edge case of installing another app with the same name in the special giantswarm namespace",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "security-pack",
					Namespace: "abc01",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "security-pack",
					Namespace: "abc01",
					Version:   "1.2.0",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: true,
					},
				},
			},
			apps: []*v1alpha1.App{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "security-pack",
						Namespace: "giantswarm",
					},
					Spec: v1alpha1.AppSpec{
						Catalog:   "giantswarm",
						Name:      "another-app",
						Namespace: "giantswarm",
						Version:   "2.4.0",
					},
				},
			},
			expectedErr: "found another app named `security-pack` installed into the `giantswarm` namespace",
		},
		{
			name: "case 7: there is an in-cluster app installed and a non-in-cluster app is applied to the giantswarm namespace with the same name",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "security-pack",
					Namespace: "giantswarm",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "giantswarm",
					Name:      "another-app",
					Namespace: "giantswarm",
					Version:   "2.4.0",
				},
			},
			apps: []*v1alpha1.App{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "security-pack",
						Namespace: "abc01",
					},
					Spec: v1alpha1.AppSpec{
						Catalog:   "giantswarm",
						Name:      "security-pack",
						Namespace: "abc01",
						Version:   "1.2.0",
						KubeConfig: v1alpha1.AppSpecKubeConfig{
							InCluster: true,
						},
					},
				},
			},
			expectedErr: "there is in-cluster app named `security-pack` already installed in the `abc01` namespace that would cause name collision with the currently applied app named `security-pack` in the `giantswarm` namespace",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g8sObjs := make([]runtime.Object, 0)

			for _, app := range tc.apps {
				g8sObjs = append(g8sObjs, app)
			}

			scheme := runtime.NewScheme()
			_ = v1alpha1.AddToScheme(scheme)

			fakeCtrlClient := fake.NewFakeClientWithScheme(scheme, g8sObjs...)

			c := Config{
				G8sClient: fakeCtrlClient,
				K8sClient: clientgofake.NewSimpleClientset(),
				Logger:    microloggertest.New(),

				ProjectName: "app-admission-controller",
				Provider:    "aws",
			}
			r, err := NewValidator(c)
			if err != nil {
				t.Fatalf("error == %#v, want nil", err)
			}

			err = r.validateUniqueInClusterAppName(ctx, tc.obj)
			switch {
			case err != nil && tc.expectedErr == "":
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.expectedErr != "":
				t.Fatalf("error == nil, want non-nil")
			}

			if err != nil && tc.expectedErr != "" {
				if !strings.Contains(err.Error(), tc.expectedErr) {
					t.Fatalf("error == %#v, want %#v ", err.Error(), tc.expectedErr)
				}

			}
		})
	}
}

func newTestCatalog(name, namespace string) *v1alpha1.Catalog {
	return &v1alpha1.Catalog{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.CatalogSpec{
			Description: name,
			Title:       name,
		},
	}
}

func newTestConfigMap(name, namespace string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		Data: map[string]string{
			"values": "cluster: yaml\n",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func newTestSecret(name, namespace string) *corev1.Secret {
	return &corev1.Secret{
		Data: map[string][]byte{
			"values": []byte("cluster: yaml\n"),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}
