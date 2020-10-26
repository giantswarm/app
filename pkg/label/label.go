// Package label contains common Kubernetes object labels. These are defined in
// https://github.com/giantswarm/fmt/blob/master/kubernetes/annotations_and_labels.md.
package label

const (
	// AppKubernetesName label is used to identify Kubernetes resources.
	AppKubernetesName = "app.kubernetes.io/name"

	// AppKubernetesVersion label is used to identify the version of Kubernetes
	// resources.
	AppKubernetesVersion = "app.kubernetes.io/version"

	// CatalogName label is used to identify resources belonging to a Giant Swarm
	// app catalog.
	CatalogName = "application.giantswarm.io/catalog"

	// CatalogType label is used to identify the type of Giant Swarm app catalog
	// e.g. stable or test.
	CatalogType = "application.giantswarm.io/catalog-type"

	// CatalogVisibility label is used to determine how Giant Swarm app catalogs
	// are displayed in the UX. e.g. public or internal.
	CatalogVisibility = "application.giantswarm.io/catalog-visibility"

	// Latest label is added to appcatalogentry CRs to filter for the most
	// recent release.
	Latest = "latest"

	// Watching is the label added to configmaps watched by the app value controller.
	Watching = "app-operator.giantswarm.io/watching"
)
