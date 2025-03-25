package key

import "github.com/giantswarm/microerror"

// EmptyValueError is returned when a value is empty.
var EmptyValueError = &microerror.Error{
	Kind: "emptyValueError",
}

// IsEmptyValueError asserts emptyValueError.
func IsEmptyValueError(err error) bool {
	return microerror.Cause(err) == EmptyValueError
}

// WrongTypeError is returned when a value has the wrong type.
var WrongTypeError = &microerror.Error{
	Kind: "wrongTypeError",
}

// IsWrongTypeError asserts WrongTypeError.
func IsWrongTypeError(err error) bool {
	return microerror.Cause(err) == WrongTypeError
}
