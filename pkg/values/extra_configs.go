package values

import (
	"sort"

	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
)

// Generic

func isWithinPriorityLevel(minExclusive, maxInclusive int) func(priority int) bool {
	return func(priority int) bool {
		return minExclusive < priority && priority <= maxInclusive
	}
}

func getExtraConfigs(appExtraConfigs []v1alpha1.AppExtraConfig, kind string, priorityCondition func(int) bool) []v1alpha1.AppExtraConfig {
	extraConfigs := []v1alpha1.AppExtraConfig{}

	for _, extraConfig := range appExtraConfigs {
		var extraConfigKind string
		if extraConfig.Kind != "" {
			extraConfigKind = extraConfig.Kind
		} else {
			extraConfigKind = "configMap"
		}

		var extraConfigPriority = considerDefaultPriority(extraConfig.Priority)

		if extraConfigKind == kind && priorityCondition(extraConfigPriority) {
			extraConfigs = append(extraConfigs, extraConfig)
		}
	}

	// Sort from based on priority, keep original order on equal elements
	sort.SliceStable(extraConfigs, func(i, j int) bool {
		left := considerDefaultPriority(extraConfigs[i].Priority)
		right := considerDefaultPriority(extraConfigs[j].Priority)

		return left < right
	})

	return extraConfigs
}
func considerDefaultPriority(value int) int {
	if value != 0 {
		return value
	} else {
		return v1alpha1.ConfigPriorityDefault
	}
}

// Pre Cluster

var isPreClusterPriority = isWithinPriorityLevel(v1alpha1.ConfigPriorityCatalog, v1alpha1.ConfigPriorityCluster)

func getPreClusterExtraConfigs(appExtraConfigs []v1alpha1.AppExtraConfig, kind string) []v1alpha1.AppExtraConfig {
	return getExtraConfigs(appExtraConfigs, kind, isPreClusterPriority)
}

func getPreClusterExtraConfigMapEntries(appExtraConfigs []v1alpha1.AppExtraConfig) []v1alpha1.AppExtraConfig {
	return getPreClusterExtraConfigs(appExtraConfigs, "configMap")
}

func getPreClusterExtraSecretEntries(appExtraConfigs []v1alpha1.AppExtraConfig) []v1alpha1.AppExtraConfig {
	return getPreClusterExtraConfigs(appExtraConfigs, "secret")
}

// Post Cluster + Pre User

var isPostClusterPreUserPriority = isWithinPriorityLevel(v1alpha1.ConfigPriorityCluster, v1alpha1.ConfigPriorityUser)

func getPostClusterPreUserExtraConfigs(appExtraConfigs []v1alpha1.AppExtraConfig, kind string) []v1alpha1.AppExtraConfig {
	return getExtraConfigs(appExtraConfigs, kind, isPostClusterPreUserPriority)
}

func getPostClusterPreUserExtraConfigMapEntries(appExtraConfigs []v1alpha1.AppExtraConfig) []v1alpha1.AppExtraConfig {
	return getPostClusterPreUserExtraConfigs(appExtraConfigs, "configMap")
}

func getPostClusterPreUserExtraSecretEntries(appExtraConfigs []v1alpha1.AppExtraConfig) []v1alpha1.AppExtraConfig {
	return getPostClusterPreUserExtraConfigs(appExtraConfigs, "secret")
}

// Post User

var isPostUserPriority = isWithinPriorityLevel(v1alpha1.ConfigPriorityUser, v1alpha1.ConfigPriorityMaximum)

func getPostUserExtraConfigs(appExtraConfigs []v1alpha1.AppExtraConfig, kind string) []v1alpha1.AppExtraConfig {
	return getExtraConfigs(appExtraConfigs, kind, isPostUserPriority)
}

func getPostUserExtraConfigMapEntries(appExtraConfigs []v1alpha1.AppExtraConfig) []v1alpha1.AppExtraConfig {
	return getPostUserExtraConfigs(appExtraConfigs, "configMap")
}

func getPostUserExtraSecretEntries(appExtraConfigs []v1alpha1.AppExtraConfig) []v1alpha1.AppExtraConfig {
	return getPostUserExtraConfigs(appExtraConfigs, "secret")
}
