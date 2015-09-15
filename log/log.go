/*
Package log defines logging for grpc-gateway.
*/
package log

import (
	"log"
	"os"
)

var globalLogger Logger = NewLogger(log.New(os.Stderr, "", log.LstdFlags))

// Logger provides an interface for logging implementations for grpc-gateway.
type Logger interface {
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
}

// SetLogger sets the logger that is used in grpc-gateway.
func SetLogger(l Logger) {
	globalLogger = l
}

// NewLogger returns a Logger that wraps a log.Logger from Go's standard log package.
func NewLogger(l *log.Logger) Logger {
	return newLogger(l)
}

// Infof logs an info message in the manner of fmt.Printf.
func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

// Infoln logs an info message in the manner of fmt.Println.
func Infoln(args ...interface{}) {
	globalLogger.Infoln(args...)
}

// Warnf logs an info message in the manner of fmt.Printf.
func Warnf(format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}

// Warnln logs an info message in the manner of fmt.Println.
func Warnln(args ...interface{}) {
	globalLogger.Warnln(args...)
}

// Errorf logs an info message in the manner of fmt.Printf.
func Errorf(format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}

// Errorln logs an info message in the manner of fmt.Println.
func Errorln(args ...interface{}) {
	globalLogger.Errorln(args...)
}
