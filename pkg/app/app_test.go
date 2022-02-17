package app

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/k8smetadata/pkg/label"
	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_NewCR(t *testing.T) {
	testCases := []struct {
		name     string
		config   Config
		expected *v1alpha1.App
	}{
		{
			name: "flawlesss",
			config: Config{
				AppCatalog:   "giantswarm",
				AppName:      "hello-world-app",
				AppNamespace: "hello-world",
				AppVersion:   "0.2.0",
				Name:         "hello-world",
				Namespace:    "giantswarm",
			},
			expected: &v1alpha1.App{
				TypeMeta: v1alpha1.NewAppTypeMeta(),
				ObjectMeta: metav1.ObjectMeta{
					Name:      "hello-world",
					Namespace: "giantswarm",
					Annotations: map[string]string{
						"chart-operator.giantswarm.io/force-helm-upgrade": "true",
					},
					Labels: map[string]string{
						"app-operator.giantswarm.io/version": "0.0.0",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: false,
					},
					Name:      "hello-world-app",
					Namespace: "hello-world",
					Version:   "0.2.0",
				},
			},
		},
		{
			name: "flawless with extra labels and annotations",
			config: Config{
				Annotations: map[string]string{
					"test-key": "test-value",
				},
				AppCatalog:   "giantswarm",
				AppName:      "hello-world-app",
				AppNamespace: "hello-world",
				AppVersion:   "0.2.0",
				Labels: map[string]string{
					label.ManagedBy: "flux",
				},
				Name:      "hello-world",
				Namespace: "giantswarm",
			},
			expected: &v1alpha1.App{
				TypeMeta: v1alpha1.NewAppTypeMeta(),
				ObjectMeta: metav1.ObjectMeta{
					Name:      "hello-world",
					Namespace: "giantswarm",
					Annotations: map[string]string{
						"chart-operator.giantswarm.io/force-helm-upgrade": "true",
						"test-key": "test-value",
					},
					Labels: map[string]string{
						"app-operator.giantswarm.io/version": "0.0.0",
						"giantswarm.io/managed-by":           "flux",
					},
				},
				Spec: v1alpha1.AppSpec{
					Catalog: "giantswarm",
					KubeConfig: v1alpha1.AppSpecKubeConfig{
						InCluster: false,
					},
					Name:      "hello-world-app",
					Namespace: "hello-world",
					Version:   "0.2.0",
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d: %s", i, tc.name), func(t *testing.T) {
			app := NewCR(tc.config)

			if !reflect.DeepEqual(app, tc.expected) {
				t.Fatalf("Want matching object \n %s", cmp.Diff(app, tc.expected))
			}
		})
	}
}
