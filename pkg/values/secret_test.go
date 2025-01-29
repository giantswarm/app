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

func Test_MergeSecretData(t *testing.T) {
	tests := []struct {
		name         string
		app          v1alpha1.App
		catalog      v1alpha1.Catalog
		secrets      []*corev1.Secret
		expectedData map[string]interface{}
		errorMatcher func(error) bool
	}{
		{
			name: "case 0: secret is nil when there are no secrets",
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
			name: "case 1: basic match with app secrets",
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
						Secret: v1alpha1.AppSpecConfigSecret{
							Name:      "test-cluster-secrets",
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
			secrets: []*corev1.Secret{
				getSecretDefinition("test-cluster-secrets", "giantswarm", map[string][]byte{
					"secrets": []byte("cluster: yaml\n"),
				}),
			},
			expectedData: map[string]interface{}{
				"cluster": "yaml",
			},
		},
		{
			name: "case 2: basic match with catalog secrets",
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
			catalog: getSimpleTestCatalogDefinitionWithSecret(),
			secrets: []*corev1.Secret{
				getTestCatalogSecretDefinition(map[string][]byte{
					"secrets": []byte("catalog: yaml\n"),
				}),
			},
			expectedData: map[string]interface{}{
				"catalog": "yaml",
			},
		},
		{
			name: "case 3: non-intersecting catalog and app secrets are merged",
			app: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-test-app",
					Namespace: "giantswarm",
				},
				Spec: v1alpha1.AppSpec{
					Catalog: "test-catalog",
					Config: v1alpha1.AppSpecConfig{
						Secret: v1alpha1.AppSpecConfigSecret{
							Name:      "test-cluster-secrets",
							Namespace: "giantswarm",
						},
					},
					Name:      "test-app",
					Namespace: "giantswarm",
				},
			},
			catalog: getSimpleTestCatalogDefinitionWithSecret(),
			secrets: []*corev1.Secret{
				getTestCatalogSecretDefinition(map[string][]byte{
					"values": []byte("catalog: yaml\n"),
				}),
				getSecretDefinition("test-cluster-secrets", "giantswarm", map[string][]byte{
					"values": []byte("cluster: yaml\n"),
				}),
			},
			expectedData: map[string]interface{}{
				"catalog": "yaml",
				"cluster": "yaml",
			},
		},
		{
			name: "case 4: intersecting catalog and app secrets, app overwrites catalog",
			app: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-test-app",
					Namespace: "giantswarm",
				},
				Spec: v1alpha1.AppSpec{
					Catalog: "test-catalog",
					Config: v1alpha1.AppSpecConfig{
						Secret: v1alpha1.AppSpecConfigSecret{
							Name:      "test-cluster-secrets",
							Namespace: "giantswarm",
						},
					},
					Name:      "test-app",
					Namespace: "giantswarm",
				},
			},
			catalog: getSimpleTestCatalogDefinitionWithSecret(),
			secrets: []*corev1.Secret{
				getTestCatalogSecretDefinition(map[string][]byte{
					"values": []byte("catalog: yaml\ntest: catalog\n"),
				}),
				getSecretDefinition("test-cluster-secrets", "giantswarm", map[string][]byte{
					"values": []byte("cluster: yaml\ntest: app\n"),
				}),
			},
			expectedData: map[string]interface{}{
				// "test: app" overrides "test: catalog".
				"catalog": "yaml",
				"cluster": "yaml",
				"test":    "app",
			},
		},
		{
			name: "case 5: intersecting catalog, app and user secrets are merged, user is preferred",
			app: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-test-app",
					Namespace: "giantswarm",
				},
				Spec: v1alpha1.AppSpec{
					Catalog: "test-catalog",
					Config: v1alpha1.AppSpecConfig{
						Secret: v1alpha1.AppSpecConfigSecret{
							Name:      "test-cluster-secrets",
							Namespace: "giantswarm",
						},
					},
					Name:      "test-app",
					Namespace: "giantswarm",
					UserConfig: v1alpha1.AppSpecUserConfig{
						Secret: v1alpha1.AppSpecUserConfigSecret{
							Name:      "test-user-secrets",
							Namespace: "giantswarm",
						},
					},
				},
			},
			catalog: getSimpleTestCatalogDefinitionWithSecret(),
			secrets: []*corev1.Secret{
				getTestCatalogSecretDefinition(map[string][]byte{
					"values": []byte("catalog: test\ntest: catalog\n"),
				}),
				getSecretDefinition("test-cluster-secrets", "giantswarm", map[string][]byte{
					"values": []byte("cluster: test\ntest: app\n"),
				}),
				getSecretDefinition("test-user-secrets", "giantswarm", map[string][]byte{
					"values": []byte("user: test\ntest: user\n"),
				}),
			},
			expectedData: map[string]interface{}{
				// "test: user" overrides "test: catalog" and "test: app".
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
					Name:      "my-prometheus",
					Namespace: "giantswarm",
				},
				Spec: v1alpha1.AppSpec{
					Catalog:   "app-catalog",
					Name:      "prometheus",
					Namespace: "monitoring",
					UserConfig: v1alpha1.AppSpecUserConfig{
						Secret: v1alpha1.AppSpecUserConfigSecret{
							Name:      "test-cluster-user-secrets",
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
			secrets: []*corev1.Secret{
				getSecretDefinition("test-cluster-user-secrets", "giantswarm", map[string][]byte{
					"secrets": []byte("cluster --\n"),
				}),
			},
			errorMatcher: IsParsingError,
		},
		{
			name: "case multi layer 1: pre cluster overrides",
			app: v1alpha1.App{
				Spec: v1alpha1.AppSpec{
					Catalog: "test-catalog",
					Config: v1alpha1.AppSpecConfig{
						Secret: v1alpha1.AppSpecConfigSecret{
							Name:      "test-cluster-secrets",
							Namespace: "giantswarm",
						},
					},
					ExtraConfigs: []v1alpha1.AppExtraConfig{
						{
							Kind:      "secret",
							Name:      "pre-cluster-secret-overrides-1",
							Namespace: "giantswarm",
							Priority:  v1alpha1.ConfigPriorityCluster - 1,
						},
						{
							Name:      "pre-cluster-config-map-overrides",
							Namespace: "giantswarm",
						},
						{
							Kind:      "secret",
							Name:      "pre-cluster-secret-overrides-2",
							Namespace: "giantswarm",
						},
					},
					Name:      "test-app",
					Namespace: "giantswarm",
				},
			},
			catalog: getSimpleTestCatalogDefinitionWithSecret(),
			secrets: []*corev1.Secret{
				getTestCatalogSecretDefinition(map[string][]byte{
					"values": []byte("catalog: test\nfoo: bar\n"),
				}),
				getSecretDefinition("test-cluster-secrets", "giantswarm", map[string][]byte{
					"values": []byte("foo: cluster\ncluster: fallthrough\n"),
				}),
				getSecretDefinition("pre-cluster-secret-overrides-1", "giantswarm", map[string][]byte{
					"values": []byte("color: red\nfoo: hello\napple: pear\n"),
				}),
				getSecretDefinition("pre-cluster-secret-overrides-2", "giantswarm", map[string][]byte{
					"values": []byte("color: green\nfoo: baz\ntop: nope\n"),
				}),
			},
			expectedData: map[string]interface{}{
				"catalog": "test",
				"foo":     "cluster",
				"color":   "red",
				"apple":   "pear",
				"top":     "nope",
				"cluster": "fallthrough",
			},
		},
		{
			name: "case multi layer 2: post cluster overrides",
			app: v1alpha1.App{
				Spec: v1alpha1.AppSpec{
					Catalog: "test-catalog",
					Config: v1alpha1.AppSpecConfig{
						Secret: v1alpha1.AppSpecConfigSecret{
							Name:      "test-cluster-secrets",
							Namespace: "giantswarm",
						},
					},
					ExtraConfigs: []v1alpha1.AppExtraConfig{
						{
							Kind:      "secret",
							Name:      "post-cluster-secret-overrides-1",
							Namespace: "giantswarm",
							Priority:  v1alpha1.ConfigPriorityUser,
						},
						{
							Kind:      "secret",
							Name:      "pre-cluster-secret-overrides",
							Namespace: "giantswarm",
							Priority:  v1alpha1.ConfigPriorityCluster - 1,
						},
						{
							Kind:      "secret",
							Name:      "post-cluster-secret-overrides-2",
							Namespace: "giantswarm",
							Priority:  v1alpha1.ConfigPriorityCluster + v1alpha1.ConfigPriorityDistance/2,
						},
					},
					Name:      "test-app",
					Namespace: "giantswarm",
					UserConfig: v1alpha1.AppSpecUserConfig{
						Secret: v1alpha1.AppSpecUserConfigSecret{
							Name:      "test-cluster-user-secrets",
							Namespace: "giantswarm",
						},
					},
				},
			},
			catalog: getSimpleTestCatalogDefinitionWithSecret(),
			secrets: []*corev1.Secret{
				getTestCatalogSecretDefinition(map[string][]byte{
					"values": []byte("catalog: test\nfoo: bar\n"),
				}),
				getSecretDefinition("pre-cluster-secret-overrides", "giantswarm", map[string][]byte{
					"values": []byte("pre-cluster: giantswarm\ncluster: nope"),
				}),
				getSecretDefinition("test-cluster-secrets", "giantswarm", map[string][]byte{
					"values": []byte("foo: cluster\ncluster: fallthrough\n"),
				}),
				getSecretDefinition("post-cluster-secret-overrides-1", "giantswarm", map[string][]byte{
					"values": []byte("color: red\nfoo: hello\napple: pear\n"),
				}),
				getSecretDefinition("post-cluster-secret-overrides-2", "giantswarm", map[string][]byte{
					"values": []byte("color: green\nfoo: baz\ntop: nope\n"),
				}),
				getSecretDefinition("test-cluster-user-secrets", "giantswarm", map[string][]byte{
					"values": []byte("user: test\ntop: max\nfoo: user\n"),
				}),
			},
			expectedData: map[string]interface{}{
				"catalog":     "test",
				"pre-cluster": "giantswarm",
				"foo":         "user",
				"color":       "red",
				"apple":       "pear",
				"top":         "max",
				"cluster":     "fallthrough",
				"user":        "test",
			},
		},
		{
			name: "case multi layer 3: post user overrides",
			app: v1alpha1.App{
				Spec: v1alpha1.AppSpec{
					Catalog: "test-catalog",
					ExtraConfigs: []v1alpha1.AppExtraConfig{
						{
							Kind:      "secret",
							Name:      "post-user-secret-overrides-1",
							Namespace: "giantswarm",
							Priority:  v1alpha1.ConfigPriorityMaximum,
						},
						{
							Kind:      "secret",
							Name:      "pre-user-secret-overrides",
							Namespace: "giantswarm",
							Priority:  v1alpha1.ConfigPriorityUser,
						},
						{
							Kind:      "secret",
							Name:      "post-user-secret-overrides-2",
							Namespace: "giantswarm",
							Priority:  v1alpha1.ConfigPriorityUser + 1,
						},
					},
					Name:      "test-app",
					Namespace: "giantswarm",
					UserConfig: v1alpha1.AppSpecUserConfig{
						Secret: v1alpha1.AppSpecUserConfigSecret{
							Name:      "test-cluster-user-secrets",
							Namespace: "giantswarm",
						},
					},
				},
			},
			catalog: getSimpleTestCatalogDefinitionWithSecret(),
			secrets: []*corev1.Secret{
				getTestCatalogSecretDefinition(map[string][]byte{
					"values": []byte("catalog: test\nfoo: bar\n"),
				}),
				getSecretDefinition("pre-user-secret-overrides", "giantswarm", map[string][]byte{
					"values": []byte("pre-cluster: giantswarm\nuser: nope"),
				}),
				getSecretDefinition("test-cluster-user-secrets", "giantswarm", map[string][]byte{
					"values": []byte("user: test\ntop: nope\nfoo: user\n"),
				}),
				getSecretDefinition("post-user-secret-overrides-2", "giantswarm", map[string][]byte{
					"values": []byte("color: green\nfoo: baz\ntop: almost\n"),
				}),
				getSecretDefinition("post-user-secret-overrides-1", "giantswarm", map[string][]byte{
					"values": []byte("color: red\nfoo: hello\napple: pear\n"),
				}),
			},
			expectedData: map[string]interface{}{
				"catalog":     "test",
				"pre-cluster": "giantswarm",
				"foo":         "hello",
				"color":       "red",
				"apple":       "pear",
				"top":         "almost",
				"user":        "test",
			},
		},
		{
			name: "case multi layer 4: data should be nil when multi layer secret not found",
			app: v1alpha1.App{
				Spec: v1alpha1.AppSpec{
					Catalog: "test-catalog",
					ExtraConfigs: []v1alpha1.AppExtraConfig{
						{
							Kind:      "secret",
							Name:      "pre-user-secret-overrides",
							Namespace: "giantswarm",
							Priority:  v1alpha1.ConfigPriorityUser,
						},
					},
					Name:      "test-app",
					Namespace: "giantswarm",
					UserConfig: v1alpha1.AppSpecUserConfig{
						Secret: v1alpha1.AppSpecUserConfigSecret{
							Name:      "test-cluster-user-secrets",
							Namespace: "giantswarm",
						},
					},
				},
			},
			catalog: getSimpleTestCatalogDefinitionWithSecret(),
			secrets: []*corev1.Secret{
				getTestCatalogSecretDefinition(map[string][]byte{
					"values": []byte("catalog: test\nfoo: bar\n"),
				}),
				getSecretDefinition("test-cluster-user-secrets", "giantswarm", map[string][]byte{
					"values": []byte("user: test\ntop: nope\nfoo: user\n"),
				}),
			},
			errorMatcher: IsNotFound,
		},
		{
			name: "case multi layer 5: data should be nil when multi layer secret is invalid yaml",
			app: v1alpha1.App{
				Spec: v1alpha1.AppSpec{
					Catalog: "test-catalog",
					ExtraConfigs: []v1alpha1.AppExtraConfig{
						{
							Kind:      "secret",
							Name:      "post-user-secret-overrides",
							Namespace: "giantswarm",
							Priority:  v1alpha1.ConfigPriorityMaximum,
						},
					},
					Name:      "test-app",
					Namespace: "giantswarm",
					UserConfig: v1alpha1.AppSpecUserConfig{
						Secret: v1alpha1.AppSpecUserConfigSecret{
							Name:      "test-cluster-user-secrets",
							Namespace: "giantswarm",
						},
					},
				},
			},
			catalog: getSimpleTestCatalogDefinitionWithSecret(),
			secrets: []*corev1.Secret{
				getTestCatalogSecretDefinition(map[string][]byte{
					"values": []byte("catalog: test\nfoo: bar\n"),
				}),
				getSecretDefinition("test-cluster-user-secrets", "giantswarm", map[string][]byte{
					"values": []byte("user: test\ntop: nope\nfoo: user\n"),
				}),
				getSecretDefinition("post-user-secret-overrides", "giantswarm", map[string][]byte{
					"values": []byte("this-will-not-end-well"),
				}),
			},
			errorMatcher: IsParsingError,
		},
		{
			name: "case multi layer 7: no catalog config map",
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
							Kind:      "secret",
							Name:      "pre-cluster-overrides-1",
							Namespace: "giantswarm",
						},
						{
							Kind:      "secret",
							Name:      "pre-cluster-overrides-2",
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
			secrets: []*corev1.Secret{
				getSecretDefinition("pre-cluster-overrides-1", "giantswarm", map[string][]byte{
					"values": []byte("foo: bar\nhello: world"),
				}),
				getSecretDefinition("pre-cluster-overrides-2", "giantswarm", map[string][]byte{
					"values": []byte("foo: baz"),
				}),
			},
			expectedData: map[string]interface{}{
				"foo":   "baz",
				"hello": "world",
			},
		},
		{
			name: "case multi layer 8: no extra secrets",
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
							Kind:      "configMap",
							Name:      "test-config-map",
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
			expectedData: nil, // we do not want the map[string]interface{}{} here
		},
	}
	ctx := context.Background()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			objs := make([]runtime.Object, 0)
			for _, cm := range tc.secrets {
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

			result, err := v.MergeSecretData(ctx, tc.app, tc.catalog)
			switch {
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case err != nil && !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}

			if result != nil && tc.expectedData == nil {
				t.Fatalf("expected nil secret got %#v", result)
			}
			if result == nil && tc.expectedData != nil {
				t.Fatal("expected non-nil secret got nil")
			}

			if tc.expectedData != nil {
				if !reflect.DeepEqual(result, tc.expectedData) {
					t.Fatalf("want matching data \n %s", cmp.Diff(result, tc.expectedData))
				}
			}
		})
	}
}

func getSimpleTestCatalogDefinitionWithSecret() v1alpha1.Catalog {
	return v1alpha1.Catalog{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-catalog",
		},
		Spec: v1alpha1.CatalogSpec{
			Title: "test-catalog",
			Config: &v1alpha1.CatalogSpecConfig{
				Secret: &v1alpha1.CatalogSpecConfigSecret{
					Name:      "test-catalog-secrets",
					Namespace: "giantswarm",
				},
			},
		},
	}
}

func getTestCatalogSecretDefinition(data map[string][]byte) *corev1.Secret {
	return getSecretDefinition("test-catalog-secrets", "giantswarm", data)
}

func getSecretDefinition(name string, namespace string, data map[string][]byte) *corev1.Secret {
	return &corev1.Secret{
		Data: data,
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}
