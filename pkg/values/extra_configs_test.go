package values

import (
	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"strconv"
	"testing"
)

func Test_GetExtraConfigs(t *testing.T) {
	tests := []struct {
		name           string
		appExtraConfig []v1alpha1.AppExtraConfig
		kind           string
		minPriority    int
		maxPriority    int
		expected       []v1alpha1.AppExtraConfig
	}{
		{
			"Empty list",
			[]v1alpha1.AppExtraConfig{},
			"configMap",
			0,
			50,
			[]v1alpha1.AppExtraConfig{},
		},
		{
			"List of a single config map and a single secret",
			[]v1alpha1.AppExtraConfig{
				{Name: "test-config-map-1", Namespace: "default"},
				{Kind: "secret", Name: "test-secret-1", Namespace: "default"},
			},
			"configMap",
			v1alpha1.ConfigPriorityCatalog,
			v1alpha1.ConfigPriorityCluster,
			[]v1alpha1.AppExtraConfig{
				{Name: "test-config-map-1", Namespace: "default"},
			},
		},
		{
			"List of a multiple config maps, pre-clusterlevel",
			[]v1alpha1.AppExtraConfig{
				{Name: "test-config-map-1", Namespace: "default", Priority: v1alpha1.ConfigPriorityCluster},
				{Name: "test-config-map-2", Namespace: "default", Priority: v1alpha1.ConfigPriorityCluster - 1},
				{Kind: "secret", Name: "test-secret-1", Namespace: "default"},
				{Name: "test-config-map-3", Namespace: "default", Priority: v1alpha1.ConfigPriorityCluster + 1},
				{Name: "test-config-map-4", Namespace: "default", Priority: v1alpha1.ConfigPriorityDefault},
				{Name: "test-config-map-5", Namespace: "default", Priority: v1alpha1.ConfigPriorityUser},
				{Name: "test-config-map-7", Namespace: "default"},
			},
			"configMap",
			v1alpha1.ConfigPriorityCatalog,
			v1alpha1.ConfigPriorityCluster,
			[]v1alpha1.AppExtraConfig{
				{Name: "test-config-map-1", Namespace: "default", Priority: v1alpha1.ConfigPriorityCluster},
				{Name: "test-config-map-2", Namespace: "default", Priority: v1alpha1.ConfigPriorityCluster - 1},
				{Name: "test-config-map-4", Namespace: "default", Priority: v1alpha1.ConfigPriorityDefault},
				{Name: "test-config-map-7", Namespace: "default"},
			},
		},
		{
			"List of a multiple config maps, post-cluster / pre-user level",
			[]v1alpha1.AppExtraConfig{
				{Name: "test-config-map-1", Namespace: "default"},
				{Name: "test-config-map-2", Namespace: "default", Priority: v1alpha1.ConfigPriorityCluster},
				{Name: "test-config-map-3", Namespace: "default", Priority: v1alpha1.ConfigPriorityCluster + 1},
				{Name: "test-config-map-4", Namespace: "default", Priority: v1alpha1.ConfigPriorityCluster + v1alpha1.ConfigPriorityDistance/2},
				{Name: "test-config-map-5", Namespace: "default", Priority: v1alpha1.ConfigPriorityUser},
				{Name: "test-config-map-6", Namespace: "default"},
				{Kind: "secret", Name: "test-secret-1", Namespace: "default", Priority: v1alpha1.ConfigPriorityCluster + v1alpha1.ConfigPriorityDistance/2},
				{Name: "test-config-map-7", Namespace: "default", Priority: v1alpha1.ConfigPriorityUser - 1},
				{Name: "test-config-map-8", Namespace: "default", Priority: v1alpha1.ConfigPriorityUser + 1},
			},
			"configMap",
			v1alpha1.ConfigPriorityCluster,
			v1alpha1.ConfigPriorityUser,
			[]v1alpha1.AppExtraConfig{
				{Name: "test-config-map-3", Namespace: "default", Priority: v1alpha1.ConfigPriorityCluster + 1},
				{Name: "test-config-map-4", Namespace: "default", Priority: v1alpha1.ConfigPriorityCluster + v1alpha1.ConfigPriorityDistance/2},
				{Name: "test-config-map-5", Namespace: "default", Priority: v1alpha1.ConfigPriorityUser},
				{Name: "test-config-map-7", Namespace: "default", Priority: v1alpha1.ConfigPriorityUser - 1},
			},
		},
		{
			"List of a multiple config maps, post-user level",
			[]v1alpha1.AppExtraConfig{
				{Name: "test-config-map-1", Namespace: "default", Priority: v1alpha1.ConfigPriorityUser},
				{Name: "test-config-map-2", Namespace: "default"},
				{Kind: "secret", Name: "test-secret-1", Namespace: "default"},
				{Name: "test-config-map-3", Namespace: "default", Priority: v1alpha1.ConfigPriorityUser + v1alpha1.ConfigPriorityDistance/2},
				{Name: "test-config-map-4", Namespace: "default", Priority: v1alpha1.ConfigPriorityMaximum},
				{Kind: "secret", Name: "test-secret-2", Namespace: "default", Priority: v1alpha1.ConfigPriorityCluster},
				{Name: "test-config-map-5", Namespace: "default", Priority: v1alpha1.ConfigPriorityUser - 1},
				{Name: "test-config-map-6", Namespace: "default", Priority: v1alpha1.ConfigPriorityUser + 1},
			},
			"configMap",
			v1alpha1.ConfigPriorityUser,
			v1alpha1.ConfigPriorityMaximum,
			[]v1alpha1.AppExtraConfig{
				{Name: "test-config-map-3", Namespace: "default", Priority: v1alpha1.ConfigPriorityUser + v1alpha1.ConfigPriorityDistance/2},
				{Name: "test-config-map-4", Namespace: "default", Priority: v1alpha1.ConfigPriorityMaximum},
				{Name: "test-config-map-6", Namespace: "default", Priority: v1alpha1.ConfigPriorityUser + 1},
			},
		},
	}

	print(v1alpha1.NewAppCR())

	for i, tc := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			result := getExtraConfigs(tc.appExtraConfig, tc.kind, tc.minPriority, tc.maxPriority)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Fatalf("want matching data \n %s", cmp.Diff(result, tc.expected))
			}
		})
	}
}
