package key

import (
	"fmt"

	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/k8smetadata/pkg/annotation"
)

func AppCatalogEntryCompatibleProviders(customResource v1alpha1.AppCatalogEntry) []string {
	if customResource.Spec.Restrictions == nil {
		return []string{}
	}
	return customResource.Spec.Restrictions.CompatibleProviders
}

func AppCatalogEntryManagedBy(projectName string) string {
	return fmt.Sprintf("%s-unique", projectName)
}

func AppCatalogEntryName(catalogName, appName, appVersion string) string {
	return fmt.Sprintf("%s-%s-%s", catalogName, appName, appVersion)
}

func AppCatalogEntryOwners(customResource v1alpha1.AppCatalogEntry) string {
	return customResource.Annotations[annotation.AppOwners]
}

func AppCatalogEntryTeam(customResource v1alpha1.AppCatalogEntry) string {
	return customResource.Annotations[annotation.AppTeam]
}
