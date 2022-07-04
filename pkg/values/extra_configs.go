package values

import "github.com/giantswarm/apiextensions-application/api/v1alpha1"

func getExtraConfigs(appExtraConfig []v1alpha1.AppExtraConfig, kind string, minPriority, maxPriority int) []v1alpha1.AppExtraConfig {
	extraConfigs := []v1alpha1.AppExtraConfig{}

	for _, extraConfig := range appExtraConfig {
		var extraConfigKind string
		if extraConfig.Kind != "" {
			extraConfigKind = extraConfig.Kind
		} else {
			extraConfigKind = "configMap"
		}

		var extraConfigPriority int
		if extraConfig.Priority != 0 {
			extraConfigPriority = extraConfig.Priority
		} else {
			extraConfigPriority = v1alpha1.ConfigPriorityDefault
		}

		if extraConfigKind == kind && extraConfigPriority > minPriority && extraConfigPriority <= maxPriority {
			extraConfigs = append(extraConfigs, extraConfig)
		}
	}

	return extraConfigs
}
