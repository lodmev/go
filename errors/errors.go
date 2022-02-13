package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// ErrorType is the type of an error
type ErrorType uint

const (
	// NoType error
	NoType ErrorType = iota
)

type customError struct {
	errorType     ErrorType
	originalError error
	context       map[string]interface{}
}

// New creates a new customError
func (errorType ErrorType) New(msg string) error {
	return customError{errorType: errorType, originalError: errors.New(msg)}
}

// New creates a new customError with formatted message
func (errorType ErrorType) Newf(msg string, args ...interface{}) error {
	return customError{errorType: errorType, originalError: fmt.Errorf(msg, args...)}
}

// Wrap creates a new wrapped error
func (errorType ErrorType) Wrap(err error, msg string) error {
	return errorType.Wrapf(err, msg)
}

// Wrap creates a new wrapped error with formatted message
func (errorType ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	return customError{errorType: errorType, originalError: errors.Wrapf(err, msg, args...)}
}

// Error returns the message of a customError
func (error customError) Error() string {
	return error.originalError.Error()
}

// New creates a no type error
func New(msg string) error {
	return customError{errorType: NoType, originalError: errors.New(msg)}
}

// Newf creates a no type error with formatted message
func Newf(msg string, args ...interface{}) error {
	return customError{errorType: NoType, originalError: errors.New(fmt.Sprintf(msg, args...))}
}

// Wrap an error with a string
func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

// Wrapf an error with format string
func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(customError); ok {
		return customError{
			errorType:     customErr.errorType,
			originalError: wrappedError,
			context:       customErr.context,
		}
	}

	return customError{errorType: NoType, originalError: wrappedError}
}

// AddContext adds a context to an error
func AddContext(err error, field string, message interface{}) error {
	context := map[string]interface{}{field: message}
	if customErr, ok := err.(customError); ok {
		if customErr.context != nil {
			customErr.context[field] = message
		} else {
			customErr.context = context
		}
		return customError{errorType: customErr.errorType, originalError: customErr.originalError, context: customErr.context}
	}

	return customError{errorType: NoType, originalError: err, context: context}
}

// GetContext returns the error context
func GetContext(err error) map[string]interface{} {
	if customErr, ok := err.(customError); ok {
		return customErr.context
	}

	return map[string]interface{}{}
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(customError); ok {
		return customErr.errorType
	}

	return NoType
}
