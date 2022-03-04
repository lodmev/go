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

var errsmap = map[ErrorType]bool{}

type customError struct {
	eType         ErrorType
	originalError error
	context       map[string]interface{}
}

// New creates a new customError
func (et ErrorType) New(msg string) error {
	errsmap[et] = true
	return customError{eType: et, originalError: errors.New(msg)}
}

// New creates a new customError with formatted message
func (et ErrorType) Newf(msg string, args ...interface{}) error {
	errsmap[et] = true
	return customError{eType: et, originalError: fmt.Errorf(msg, args...)}
}

// Wrap creates a new wrapped error
func (et ErrorType) Wrap(err error, msg string) error {
	return et.Wrapf(err, msg)
}

// Wrap creates a new wrapped error with formatted message
func (et ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	errsmap[et] = true
	return customError{eType: et, originalError: errors.Wrapf(err, msg, args...)}
}

// GotSuches return true if errors of this type was got
func (et ErrorType) GotSuches() bool {
	_, ok := errsmap[et]
	return ok
}

// Error returns the message of a customError
func (error customError) Error() string {
	return error.originalError.Error()
}

// GetContext returns the customError.context
func (error customError) GetContext() map[string]interface{} {
	return error.context
}

// New creates a no type error
func New(msg string) error {
	return customError{eType: NoType, originalError: errors.New(msg)}
}

// Newf creates a no type error with formatted message
func Newf(msg string, args ...interface{}) error {
	return customError{eType: NoType, originalError: errors.New(fmt.Sprintf(msg, args...))}
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
		errsmap[customErr.eType] = true
		return customError{
			eType:         customErr.eType,
			originalError: wrappedError,
			context:       customErr.context,
		}
	}
	errsmap[NoType] = true
	return customError{eType: NoType, originalError: wrappedError}
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
		errsmap[customErr.eType] = true
		return customError{eType: customErr.eType, originalError: customErr.originalError, context: customErr.context}
	}

	errsmap[NoType] = true
	return customError{eType: NoType, originalError: err, context: context}
}

// GetContext returns the error context
func GetContext(err error) map[string]interface{} {
	if customErr, ok := err.(customError); ok {
		return customErr.GetContext()
	}

	return map[string]interface{}{}
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(customError); ok {
		return customErr.eType
	}

	return NoType
}
