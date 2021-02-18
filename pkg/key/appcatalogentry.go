package key

import (
	"fmt"
)

func AppCatalogEntryManagedBy(projectName string) string {
	return fmt.Sprintf("%s-unique", projectName)
}

func AppCatalogEntryName(catalogName, appName, appVersion string) string {
	return fmt.Sprintf("%s-%s-%s", catalogName, appName, appVersion)
}
