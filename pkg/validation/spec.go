package validation

import (
	"context"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/application/v1alpha1"
)

type Interface interface {
	// ValidateApp checks if given App CR is valid. It returns nil if it is.
	//
	// It returns error matched by IsValidation if the App CR is invalid.
	//
	// It returns error matched by IsAppConfigMapNotFound if the ConfigMap
	// containing values to template is not found in the cluster. Usually
	// it should be generated by app-operator after a short delay.
	//
	// It returns error matched by IsKubeConfigNotFound if the Secret
	// containing kubeconfig is not generated yet. Usually it should be
	// generated by app-operator after a short delay.
	ValidateApp(ctx context.Context, app v1alpha1.App) error
}
