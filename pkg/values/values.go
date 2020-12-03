package values

import (
	"context"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/application/v1alpha1"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/imdario/mergo"
	"k8s.io/client-go/kubernetes"
)

// Config represents the configuration used to create a new values service.
type Config struct {
	// Dependencies.
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger
}

// Values implements the values service.
type Values struct {
	// Dependencies.
	k8sClient kubernetes.Interface
	logger    micrologger.Logger
}

// New creates a new configured values service.
func New(config Config) (*Values, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	r := &Values{
		// Dependencies.
		k8sClient: config.K8sClient,
		logger:    config.Logger,
	}

	return r, nil
}

// MergeAll merges both configmap and secret values to produce a single set of
// values that can be passed to Helm.
func (v *Values) MergeAll(ctx context.Context, app v1alpha1.App, appCatalog v1alpha1.AppCatalog) (map[string]interface{}, error) {
	configMapData, err := v.MergeConfigMapData(ctx, app, appCatalog)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	secretData, err := v.MergeSecretData(ctx, app, appCatalog)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	err = mergo.Merge(&configMapData, secretData, mergo.WithOverride)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return configMapData, nil
}

// toStringMap converts from a byte slice map to a string map.
func toStringMap(input map[string][]byte) map[string]string {
	if input == nil {
		return nil
	}

	result := map[string]string{}

	for k, v := range input {
		result[k] = string(v)
	}

	return result
}
