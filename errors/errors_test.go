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
func TestGotSuches(t *testing.T) {
	const (
		NoType errors.ErrorType = iota
		NewTypeErr
		NewFTypeErr
		WrapfTypeErr
		NotSuchErrors
		ContextErr
	)
	NoType.New("no type err")
	er := ContextErr.New("cont err")
	errors.AddContext(er, "cont", 0)
	err := NewTypeErr.New("some error")
	NewFTypeErr.Newf("any %s", "error")
	WrapfTypeErr.Wrapf(err, "and %s", "another")
	switch {
	case NotSuchErrors.GotSuches():
		t.Fatalf("expected for NotSuchError %t, got %t", false, true)
	case !ContextErr.GotSuches():
		t.Fatalf("expected for ContextErr %t, got %t", true, false)
	case !NoType.GotSuches():
		t.Fatalf("expected for NoType %t, got %t", true, false)
	case !NewTypeErr.GotSuches():
		t.Fatalf("expected for NewTypeErr %t, got %t", true, false)
	case !NewFTypeErr.GotSuches():
		t.Fatalf("expected for NewFTypeErr %t, got %t", true, false)
	case !WrapfTypeErr.GotSuches():
		t.Fatalf("expected for WrapfTypeErr %t, got %t", true, false)

	}
}
