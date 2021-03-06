// Package annotation contains common Kubernetes metadata. These are defined in
// https://github.com/giantswarm/fmt/blob/master/kubernetes/annotations_and_labels.md.
package annotation

const (
	// AppOperatorPrefix the prefix for annotations that control logic inside app-operator.
	AppOperatorPrefix = "app-operator.giantswarm.io"

	// ChartOperatorPrefix is the prefix for annotations that control logic inside chart-operator.
	ChartOperatorPrefix = "chart-operator.giantswarm.io"

	// CordonReason is the name of the annotation that indicates
	// the reason of why operators should not apply any update to this app CR.
	CordonReason = "cordon-reason"

	// CordonUntil is the name of the annotation that indicates
	// the expiration date for this cordon rule.
	CordonUntil = "cordon-until"

	// LatestConfigMapVersion is the highest resource version among the configmaps
	// app CRs depends on.
	LatestConfigMapVersion = "latest-configmap-version"

	// LatestSecretVersion is the highest resource version among the secret
	// app CRs depends on.
	LatestSecretVersion = "latest-secret-version"

	// Owners annotation is defined in Chart.yaml and added to AppCatalogEntry CRs.
	// It is used when an app is owned by multiple teams.
	Owners = "application.giantswarm.io/owners"

	// Team annotation is defined in Chart.yaml and added to AppCatalogEntry CRs.
	// It is used when an app is owned by a single team.
	Team = "application.giantswarm.io/team"

	// WebhookURL is the URL that chart-operator reports chart updates.
	WebhookURL = "webhook-url"
)
