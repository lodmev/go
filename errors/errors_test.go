package errors_test

import (
	"testing"

	"github.com/lodmev/go/errors"
	"github.com/stretchr/testify/assert"
)

const (
	Ntype = errors.NoType
	BadRequest
)

func TestContext(t *testing.T) {

	err := BadRequest.New("an_error")
	errWithContext := errors.AddContext(err, "a_field", "the field is empty")

	expectedContext := map[string]interface{}{"a_field": "the field is empty"}

	assert.Equal(t, BadRequest, errors.GetType(errWithContext))
	assert.Equal(t, expectedContext, errors.GetContext(errWithContext))
	assert.Equal(t, err.Error(), errWithContext.Error())
}

func TestContextInNoTypeError(t *testing.T) {
	err := errors.New("a custom error")

	errWithContext := errors.AddContext(err, "a_field", "the field is empty")

	expectedContext := map[string]interface{}{"a_field": "the field is empty"}

	assert.Equal(t, errors.NoType, errors.GetType(errWithContext))
	assert.Equal(t, expectedContext, errors.GetContext(errWithContext))
	assert.Equal(t, err.Error(), errWithContext.Error())
}

func TestWrapf(t *testing.T) {
	err := errors.New("an_error")
	wrappedError := BadRequest.Wrapf(err, "error %s", "1")

	assert.Equal(t, BadRequest, errors.GetType(wrappedError))
	assert.EqualError(t, wrappedError, "error 1: an_error")
}

func TestWrapfInNoTypeError(t *testing.T) {
	err := errors.Newf("an_error %s", "2")
	wrappedError := errors.Wrapf(err, "error %s", "1")

	assert.Equal(t, errors.NoType, errors.GetType(wrappedError))
	assert.EqualError(t, wrappedError, "error 1: an_error 2")
}
