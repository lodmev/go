// Package log provides a global logger for zerolog.
package log

import (
	"fmt"
	"context"

	"github.com/rs/zerolog"
)

type GetContext interface {
	GetContext() map[string]interface{}
}

func Errf(format string, a ...any) {
	wrapedErr := fmt.Errorf(format, a...)
	Error().Err(wrapedErr).Send()
}

// Errt starts a new message with error and error type  
func Errt(err error) {
	Err(err).Str("error type", fmt.Sprintf("%T",err)).Send()
}

// Fatal simalar same method of standart library.
func Fatal(err error) {
	Logger.Fatal().Err(err).Send()
}

// Fatalf simalar same method of standart library.
func Fatalf(format string, a ...any) {
	wrapedErr := fmt.Errorf(format, a...)
	Logger.Fatal().Err(wrapedErr).Send()
}

func ErrContext(err error) *zerolog.Event {
	return Err(err).Dict("error context", DictErrFields(err))
}

func DictErrFields(err error) *zerolog.Event {
	var context map[string]interface{}
	if contextErr, ok := err.(GetContext); ok {
		context = contextErr.GetContext()
	}
	return zerolog.Dict().Fields(context)
}

// Output duplicates the global logger and sets w as its output.
var Output = Logger.Output

// With creates a child logger with the field added to its context.
var With = Logger.With

// Level creates a child logger with the minimum accepted level set to level.
var Level = Logger.Level

// Sample returns a logger with the s sampler.
var Sample = Logger.Sample

// Hook returns a logger with the h Hook.
var Hook = Logger.Hook

// Err starts a new message with error level with err as a field if not nil or
// with info level if err is nil.
//
// You must call Msg on the returned event in order to send the event.
var Err = Logger.Err

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
var Error = Logger.Error

// Trace starts a new message with trace level.
//
// You must call Msg on the returned event in order to send the event.
var Trace = Logger.Trace

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
var Debug = Logger.Debug

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
var Info = Logger.Info

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
var Warn = Logger.Warn

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call Msg on the returned event in order to send the event.
var Panic = Logger.Panic

// WithLevel starts a new message with level.
//
// You must call Msg on the returned event in order to send the event.
var WithLevel = Logger.WithLevel

// Log starts a new message with no level. Setting zerolog.GlobalLevel to
// zerolog.Disabled will still disable events produced by this method.
//
// You must call Msg on the returned event in order to send the event.
var Log = Logger.Log

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
var Print = Logger.Print

//Println is alias for Print
var Println = Logger.Print

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
var Printf = Logger.Printf

// Ctx returns the Logger associated with the ctx. If no logger
// is associated, a disabled logger is returned.
func Ctx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
