package key

import (
	"fmt"
	"testing"

	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/k8smetadata/pkg/label"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_ClusterConfigMapName(t *testing.T) {
	testCases := []struct {
		name          string
		obj           v1alpha1.App
		expectedValue string
	}{
		{
			name: "Vintage",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "hello-world",
					Namespace: "demo",
				},
				Spec: v1alpha1.AppSpec{
					Name: "hello-world",
				},
			},
			expectedValue: "demo-cluster-values",
		},
		{
			name: "CAPI",
			obj: v1alpha1.App{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "hello-world",
					Namespace: "org-demo",
					Labels: map[string]string{
						label.Cluster: "demo",
					},
				},
				Spec: v1alpha1.AppSpec{
					Name: "hello-world",
				},
			},
			expectedValue: "demo-cluster-values",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d: %s", i, tc.name), func(t *testing.T) {
			t.Log(tc.name)

			name := ClusterConfigMapName(tc.obj)

			if name != tc.expectedValue {
				t.Fatalf("AppConfigMapName %#q, want %#q", name, tc.expectedValue)
			}
		})
	}
}
