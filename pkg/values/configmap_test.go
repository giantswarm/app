package values

import (
	"context"
	"reflect"
	"testing"

	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgofake "k8s.io/client-go/kubernetes/fake"
)

func Test_MergeConfigMapData(t *testing.T) {
	tests := []struct {
		name         string
		app          v1alpha1.App
		catalog      v1alpha1.Catalog
		configMaps   []*corev1.ConfigMap
		expectedData map[string]interface{}
		errorMatcher func(error) bool
	}{
		{
			name: "case 0: configmap is nil when there is no config",
			app: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-test-app",
					Namespace: "giantswarm",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "app-catalog",
					Name:      "test-app",
					Namespace: "kube-system",
				},
			},
			catalog: v1alpha1.Catalog{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-catalog",
				},
			},
			expectedData: nil,
		},
		{
			name: "case 1: basic match with app config",
			app: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-prometheus",
					Namespace: "giantswarm",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "app-catalog",
					Name:      "prometheus",
					Namespace: "monitoring",
					Config: v1alpha1.AppSpecConfig{
						ConfigMap: v1alpha1.AppSpecConfigConfigMap{
							Name:      "test-cluster-values",
							Namespace: "giantswarm",
						},
					},
				},
			},
			catalog: v1alpha1.Catalog{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-catalog",
				},
			},
			configMaps: []*corev1.ConfigMap{
				getConfigMapDefinition("test-cluster-values", "giantswarm", map[string]string{
					"values": "cluster: yaml\n",
				}),
			},
			expectedData: map[string]interface{}{
				"cluster": "yaml",
			},
		},
		{
			name: "case 2: basic match with catalog config",
			app: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-test-app",
					Namespace: "giantswarm",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "test-catalog",
					Name:      "test-app",
					Namespace: "giantswarm",
				},
			},
			catalog: getSimpleTestCatalogDefinition(),
			configMaps: []*corev1.ConfigMap{
				getTestCatalogConfigMapDefinition(map[string]string{
					"values": "catalog: yaml\n",
				}),
			},
			expectedData: map[string]interface{}{
				"catalog": "yaml",
			},
		},
		{
			name: "case 3: non-intersecting catalog and app config are merged",
			app: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-test-app",
					Namespace: "giantswarm",
				},
				Spec: v1alpha1.AppSpec{
					Name:      "test-app",
					Namespace: "giantswarm",
					Catalog:   "test-catalog",
					Config: v1alpha1.AppSpecConfig{
						ConfigMap: v1alpha1.AppSpecConfigConfigMap{
							Name:      "test-cluster-values",
							Namespace: "giantswarm",
						},
					},
				},
			},
			catalog: getSimpleTestCatalogDefinition(),
			configMaps: []*corev1.ConfigMap{
				getTestCatalogConfigMapDefinition(map[string]string{
					"values": "catalog: yaml\n",
				}),
				getConfigMapDefinition("test-cluster-values", "giantswarm", map[string]string{
					"values": "cluster: yaml\n",
				}),
			},
			expectedData: map[string]interface{}{
				"catalog": "yaml",
				"cluster": "yaml",
			},
		},
		{
			name: "case 4: intersecting catalog and app config are merged, app is preferred",
			app: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-test-app",
					Namespace: "giantswarm",
				},
				Spec: v1alpha1.AppSpec{
					Name:      "test-app",
					Namespace: "giantswarm",
					Catalog:   "test-catalog",
					Config: v1alpha1.AppSpecConfig{
						ConfigMap: v1alpha1.AppSpecConfigConfigMap{
							Name:      "test-cluster-values",
							Namespace: "giantswarm",
						},
					},
				},
			},
			catalog: getSimpleTestCatalogDefinition(),
			configMaps: []*corev1.ConfigMap{
				getTestCatalogConfigMapDefinition(map[string]string{
					"values": "test: catalog\n",
				}),
				getConfigMapDefinition("test-cluster-values", "giantswarm", map[string]string{
					"values": "test: app\n",
				}),
			},
			expectedData: map[string]interface{}{
				"test": "app",
			},
		},
		{
			name: "case 5: intersecting catalog, app and user config is merged, user is preferred",
			app: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-test-app",
					Namespace: "giantswarm",
				},
				Spec: v1alpha1.AppSpec{
					Name:      "test-app",
					Namespace: "giantswarm",
					Catalog:   "test-catalog",
					Config: v1alpha1.AppSpecConfig{
						ConfigMap: v1alpha1.AppSpecConfigConfigMap{
							Name:      "test-cluster-values",
							Namespace: "giantswarm",
						},
					},
					UserConfig: v1alpha1.AppSpecUserConfig{
						ConfigMap: v1alpha1.AppSpecUserConfigConfigMap{
							Name:      "test-user-values",
							Namespace: "giantswarm",
						},
					},
				},
			},
			catalog: getSimpleTestCatalogDefinition(),
			configMaps: []*corev1.ConfigMap{
				getTestCatalogConfigMapDefinition(map[string]string{
					"values": "catalog: test\ntest: catalog\n",
				}),
				getConfigMapDefinition("test-cluster-values", "giantswarm", map[string]string{
					"values": "cluster: test\ntest: app\n",
				}),
				getConfigMapDefinition("test-user-values", "giantswarm", map[string]string{
					"values": "user: test\ntest: user\n",
				}),
			},
			expectedData: map[string]interface{}{
				"catalog": "test",
				"cluster": "test",
				"test":    "user",
				"user":    "test",
			},
		},
		{
			name: "case 6: parsing error from wrong user values",
			app: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-test-app",
					Namespace: "giantswarm",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "test-catalog",
					Name:      "test-app",
					Namespace: "giantswarm",
					UserConfig: v1alpha1.AppSpecUserConfig{
						ConfigMap: v1alpha1.AppSpecUserConfigConfigMap{
							Name:      "user-values",
							Namespace: "giantswarm",
						},
					},
				},
			},
			catalog: getSimpleTestCatalogDefinition(),
			configMaps: []*corev1.ConfigMap{
				getTestCatalogConfigMapDefinition(map[string]string{
					"values": `values: val`,
				}),
				getConfigMapDefinition("user-values", "giantswarm", map[string]string{
					"values": `values: -`,
				}),
			},
			errorMatcher: IsParsingError,
		},
		{
			name: "case multi layer 1: pre cluster overrides catalog",
			app: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-test-app",
					Namespace: "giantswarm",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "test-catalog",
					Name:      "test-app",
					Namespace: "giantswarm",
					ExtraConfigs: []v1alpha1.AppExtraConfig{
						{
							Name:      "pre-cluster-overrides",
							Namespace: "giantswarm",
						},
					},
				},
			},
			catalog: getSimpleTestCatalogDefinition(),
			configMaps: []*corev1.ConfigMap{
				getTestCatalogConfigMapDefinition(map[string]string{
					"values": "foo: bar\ntest: catalog\n",
				}),
				getConfigMapDefinition("pre-cluster-overrides", "giantswarm", map[string]string{
					"values": "foo: baz\n",
				}),
			},
			expectedData: map[string]interface{}{
				"foo":  "baz",
				"test": "catalog",
			},
		},
		{
			name: "case multi layer 2: post cluster overrides pre cluster",
			app: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-test-app",
					Namespace: "giantswarm",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "test-catalog",
					Name:      "test-app",
					Namespace: "giantswarm",
					Config: v1alpha1.AppSpecConfig{
						ConfigMap: v1alpha1.AppSpecConfigConfigMap{
							Name:      "cluster-overrides",
							Namespace: "giantswarm",
						},
					},
					ExtraConfigs: []v1alpha1.AppExtraConfig{
						{
							Name:      "pre-cluster-overrides",
							Namespace: "giantswarm",
						},
						{
							Name:      "post-cluster-overrides",
							Namespace: "giantswarm",
							Priority:  v1alpha1.ConfigPriorityCluster + 1,
						},
					},
				},
			},
			catalog: getSimpleTestCatalogDefinition(),
			configMaps: []*corev1.ConfigMap{
				getTestCatalogConfigMapDefinition(map[string]string{
					"values": "foo: bar\ntest: catalog\n",
				}),
				getConfigMapDefinition("pre-cluster-overrides", "giantswarm", map[string]string{
					"values": "foo: baz\npre-cluster: test",
				}),
				getConfigMapDefinition("cluster-overrides", "giantswarm", map[string]string{
					"values": "cluster: something",
				}),
				getConfigMapDefinition("post-cluster-overrides", "giantswarm", map[string]string{
					"values": "foo: hello\npost-cluster: world",
				}),
			},
			expectedData: map[string]interface{}{
				"foo":          "hello",
				"test":         "catalog",
				"cluster":      "something",
				"pre-cluster":  "test",
				"post-cluster": "world",
			},
		},
		{
			name: "case multi layer 3: post user overrides all",
			app: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-test-app",
					Namespace: "giantswarm",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "test-catalog",
					Name:      "test-app",
					Namespace: "giantswarm",
					Config: v1alpha1.AppSpecConfig{
						ConfigMap: v1alpha1.AppSpecConfigConfigMap{
							Name:      "cluster-overrides",
							Namespace: "giantswarm",
						},
					},
					UserConfig: v1alpha1.AppSpecUserConfig{
						ConfigMap: v1alpha1.AppSpecUserConfigConfigMap{
							Name:      "user-overrides",
							Namespace: "giantswarm",
						},
					},
					ExtraConfigs: []v1alpha1.AppExtraConfig{
						{
							Name:      "pre-cluster-overrides",
							Namespace: "giantswarm",
						},
						{
							Name:      "post-cluster-overrides",
							Namespace: "giantswarm",
							Priority:  v1alpha1.ConfigPriorityCluster + 1,
						},
						{
							Name:      "post-user-overrides-1",
							Namespace: "giantswarm",
							Priority:  v1alpha1.ConfigPriorityUser + 1,
						},
						{
							Name:      "post-user-overrides-2",
							Namespace: "giantswarm",
							Priority:  v1alpha1.ConfigPriorityMaximum,
						},
					},
				},
			},
			catalog: getSimpleTestCatalogDefinition(),
			configMaps: []*corev1.ConfigMap{
				getTestCatalogConfigMapDefinition(map[string]string{
					"values": "foo: bar\ntest: catalog\n",
				}),
				getConfigMapDefinition("pre-cluster-overrides", "giantswarm", map[string]string{
					"values": "foo: baz\npre-cluster: test",
				}),
				getConfigMapDefinition("cluster-overrides", "giantswarm", map[string]string{
					"values": "cluster: something",
				}),
				getConfigMapDefinition("post-cluster-overrides", "giantswarm", map[string]string{
					"values": "foo: hello\npost-cluster: world",
				}),
				getConfigMapDefinition("user-overrides", "giantswarm", map[string]string{
					"values": "ping: pong\napple: pear",
				}),
				getConfigMapDefinition("post-user-overrides-1", "giantswarm", map[string]string{
					"values": "foo: post-user\napple: banana\ncolor: blue",
				}),
				getConfigMapDefinition("post-user-overrides-2", "giantswarm", map[string]string{
					"values": "cluster: max-priority\ncolor: yellow\ntop: max",
				}),
			},
			expectedData: map[string]interface{}{
				"foo":          "post-user",
				"test":         "catalog",
				"cluster":      "max-priority",
				"pre-cluster":  "test",
				"post-cluster": "world",
				"ping":         "pong",
				"apple":        "banana",
				"color":        "yellow",
				"top":          "max",
			},
		},
	}

	ctx := context.Background()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			objs := make([]runtime.Object, 0)
			for _, cm := range tc.configMaps {
				objs = append(objs, cm)
			}

			c := Config{
				K8sClient: clientgofake.NewSimpleClientset(objs...),
				Logger:    microloggertest.New(),
			}
			v, err := New(c)
			if err != nil {
				t.Fatalf("error == %#v, want nil", err)
			}

			result, err := v.MergeConfigMapData(ctx, tc.app, tc.catalog)
			switch {
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case err != nil && !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}

			if result != nil && tc.expectedData == nil {
				t.Fatalf("expected nil map got %#v", result)
			}
			if result == nil && tc.expectedData != nil {
				t.Fatal("expected non-nil gmap got nil")
			}

			if tc.expectedData != nil {
				if !reflect.DeepEqual(result, tc.expectedData) {
					t.Fatalf("want matching data \n %s", cmp.Diff(result, tc.expectedData))
				}
			}
		})
	}
}

func getSimpleTestCatalogDefinition() v1alpha1.Catalog {
	return v1alpha1.Catalog{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-catalog",
		},
		Spec: v1alpha1.CatalogSpec{
			Title: "test-catalog",
			Config: &v1alpha1.CatalogSpecConfig{
				ConfigMap: &v1alpha1.CatalogSpecConfigConfigMap{
					Name:      "test-catalog-values",
					Namespace: "giantswarm",
				},
			},
		},
	}
}

func getTestCatalogConfigMapDefinition(data map[string]string) *corev1.ConfigMap {
	return getConfigMapDefinition("test-catalog-values", "giantswarm", data)
}

func getConfigMapDefinition(name, namespace string, data map[string]string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		Data: data,
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}
