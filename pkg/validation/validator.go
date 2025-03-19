package validation

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Config struct {
	G8sClient client.Client
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger

	IsAdmissionController bool
	Provider              string
}

type Validator struct {
	g8sClient client.Client
	k8sClient kubernetes.Interface
	logger    micrologger.Logger

	isAdmissionController bool
	provider              string
}

func NewValidator(config Config) (*Validator, error) {
	if config.G8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.G8sClient must not be empty", config)
	}
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Provider == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Provider must not be empty", config)
	}

	validator := &Validator{
		g8sClient: config.G8sClient,
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		isAdmissionController: config.IsAdmissionController,
		provider:              config.Provider,
	}

	return validator, nil
}
