package errors

import (
	"database/sql"
	"errors"

	"getsturdy.com/api/pkg/auth"
)

type ResolverError interface {
	error
	Extensions() map[string]interface{}
}

var ErrNotFound = errors.New("NotFoundError")
var ErrBadRequest = errors.New("BadRequestError")
var ErrInternalServer = errors.New("InternalServerError")
var ErrForbidden = errors.New("ForbiddenError")
var ErrUnauthenticated = errors.New("UnauthenticatedError")
var ErrNotImplemented = errors.New("NotImplementedError")

var clientSideErrors = []error{
	ErrNotFound,
	ErrBadRequest,
	ErrForbidden,
	ErrUnauthenticated,
	ErrNotImplemented,
}

func IsClientSideError(err error) bool {
	for _, ce := range clientSideErrors {
		if errors.Is(err, ce) {
			return true
		}
	}
	return false
}

func Error(err error, kv ...string) ResolverError {
	data := make(map[string]interface{})
	for i := 0; i < len(kv); i += 2 {
		data[kv[i]] = kv[i+1]
	}

	// Stack trace
	// data["trace"] = string(debug.Stack())
	// data["error"] = err.Error()

	switch {
	case err == nil:
		return nil
	case errors.Is(err, sql.ErrNoRows):
		return &SturdyGraphqlError{err: ErrNotFound, data: data, originalError: err}
	case errors.Is(err, auth.ErrUnauthenticated):
		return &SturdyGraphqlError{err: ErrUnauthenticated, data: data, originalError: err}
	case errors.Is(err, auth.ErrForbidden):
		// if resource is forbidden - users shouldn't know if exists. thus, return not found
		return &SturdyGraphqlError{err: ErrNotFound, data: data, originalError: err}
	case errors.Is(err, ErrNotFound),
		errors.Is(err, ErrBadRequest),
		errors.Is(err, ErrForbidden),
		errors.Is(err, ErrInternalServer),
		errors.Is(err, ErrNotImplemented):
		return &SturdyGraphqlError{err: err, data: data, originalError: err}
	default:
		return &SturdyGraphqlError{err: ErrInternalServer, data: data, originalError: err}
	}
}

type SturdyGraphqlError struct {
	err  error // This error is exposed on the API
	data map[string]interface{}

	originalError error // not exposed by GraphQL, used for internal logging
}

func (e *SturdyGraphqlError) Error() string {
	return e.err.Error()
}

func (e *SturdyGraphqlError) Extensions() map[string]interface{} {
	return e.data
}

func (e *SturdyGraphqlError) OriginalError() error {
	return e.originalError
}

func (e *SturdyGraphqlError) Is(target error) bool {
	return target == e.err
}

func (e *SturdyGraphqlError) Unwrap() error {
	return e.originalError
}
