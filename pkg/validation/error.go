package validation

import (
	"strings"

	"github.com/giantswarm/microerror"
)

const (
	appConfigMapNotFoundPattern string = "admission webhook \"apps.app-admission-controller-unique.giantswarm.io\" denied the request: app config map not found error"
	kubeConfigNotFoundPattern   string = "admission webhook \"apps.app-admission-controller-unique.giantswarm.io\" denied the request: kube config not found error"
)

var appConfigMapNotFoundError = &microerror.Error{
	Kind: "appConfigMapNotFoundError",
}

// IsAppConfigMapNotFound asserts appConfigMapNotFoundError.
func IsAppConfigMapNotFound(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)

	if strings.Contains(c.Error(), appConfigMapNotFoundPattern) {
		return true
	}

	if c == appConfigMapNotFoundError { //nolint:gosimple
		return true
	}

	return false
}

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var kubeConfigNotFoundError = &microerror.Error{
	Kind: "kubeConfigNotFoundError",
}

// IsKubeConfigNotFound asserts kubeConfigNotFoundError.
func IsKubeConfigNotFound(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)

	if strings.Contains(c.Error(), kubeConfigNotFoundPattern) {
		return true
	}

	if c == appConfigMapNotFoundError { //nolint:gosimple
		return true
	}

	return false
}

var notAllowedError = &microerror.Error{
	Kind: "notAllowedError",
}

// IsNotAllowed asserts notAllowedError.
func IsNotAllowed(err error) bool {
	return microerror.Cause(err) == notAllowedError
}

var notFoundError = &microerror.Error{
	Kind: "notFoundError",
}

// IsNotFound asserts notFoundError.
func IsNotFound(err error) bool {
	return microerror.Cause(err) == notFoundError
}

var parsingFailedError = &microerror.Error{
	Kind: "parsingFailedError",
}

// IsParsingFailed asserts parsingFailedError.
func IsParsingFailed(err error) bool {
	return microerror.Cause(err) == parsingFailedError
}

var validationError = &microerror.Error{
	Kind: "validationError",
}

// IsValidationError asserts validationError.
func IsValidationError(err error) bool {
	return microerror.Cause(err) == validationError
}
