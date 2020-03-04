// Copyright © 2020 Krishna Iyer Easwaran
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package log provides commonly used logging functions by wrapping the uber-go/zap package.
package log

import (
	"context"

	"go.uber.org/zap"
)

type loggerKeyType string

var loggerKey loggerKeyType = "logger"

// Logger wraps zap.Logger.
type Logger struct {
	ctx    context.Context
	logger *zap.Logger
	fields map[string]Field
}

//Options is the logger options.
type Options struct {
}

// Field represents a logger field.
type Field struct {
}

// New creates a new logger. Make sure to call defer logger.Clean() after calling this.
func New(ctx context.Context) (*Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &Logger{
		ctx:    ctx,
		logger: logger,
	}, nil
}

// Clean cleans up the log states. Make sure to call this after creating a new logger.
func (l *Logger) Clean() {
	l.logger.Sync()
}

// NewLoggerWithContext returns a new context with a logger and panics if it doesn't match the interface.
func NewLoggerWithContext(parentCtx context.Context, logger *Logger) context.Context {
	if logger == nil {
		panic("Nil Logger")
	}
	return context.WithValue(parentCtx, loggerKey, logger)
}

// NewLoggerFromContext retrieves a logger from a context and panics if there isn't one.
func NewLoggerFromContext(ctx context.Context) *Logger {
	val := ctx.Value(loggerKey)
	logger, ok := val.(*Logger)
	if !ok {
		panic("No logger in context")
	}
	return logger
}

// Debug logs a Debug level message.
func (l *Logger) Debug(msg string) {

}

// Info logs a Info level message.
func (l *Logger) Info(msg string) {

}

// Warn logs a Warning level message.
func (l *Logger) Warn(msg string) {

}

// Error logs a Error level message.
func (l *Logger) Error(msg string) {

}

// Fatal logs a Fatal message.
func (l *Logger) Fatal(msg string) {

}

// Debugf logs a Debug level formatted message.
func (l *Logger) Debugf(format string, v ...interface{}) {

}

// Infof logs a Info level formatted message.
func (l *Logger) Infof(format string, v ...interface{}) {

}

// Warnf logs a Warning level formatted message.
func (l *Logger) Warnf(format string, v ...interface{}) {

}

// Errorf logs a Error formatted message.
func (l *Logger) Errorf(format string, v ...interface{}) {

}

// Fatalf logs a Fatal formatted message.
func (l *Logger) Fatalf(format string, v ...interface{}) {

}

// WithField returns a logger with the fields.
func (l *Logger) WithField(string, interface{}) Logger {

}

func (l *Logger) WithFields(Fields) Logger {

}

// WithError returns a logger with the current error padded.
func (l *Logger) WithError(error) Logger {

}
