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

// Package zephyrus provides commonly used logging functions by wrapping the uber-go/zap package.
package zephyrus

import (
	"context"
	"reflect"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var reflectTypeToZapFieldType = map[reflect.Kind]zapcore.FieldType{
	reflect.String: zapcore.StringType,
	reflect.Int64:  zapcore.Int64Type,
}

type loggerKeyType string

var loggerKey loggerKeyType = "logger"

// Logger wraps zap.Logger.
type Logger struct {
	ctx    context.Context
	logger *zap.Logger
	fields []zap.Field
	err    error
}

// Options is the logger options.
// TODO: Adapt to zap.Options
type Options struct {
}

// Field represents a logger field.
type Field struct {
	Key   string
	Value interface{}
}

// New creates a new logger. Make sure to call defer logger.Clean() after calling this.
// TODO: Benchmarking
func New(ctx context.Context, debug bool) (*Logger, error) {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	if !debug {
		config.DisableStacktrace = true
		config.EncoderConfig.CallerKey = ""
		config.EncoderConfig.TimeKey = ""
	}
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{
		ctx:    ctx,
		logger: logger,
		fields: make([]zap.Field, 0),
	}, nil
}

// Clean cleans up the log states. Make sure to call this after creating a new logger.
func (l *Logger) Clean() {
	l.logger.Sync()
}

// NewContextWithLogger returns a new context with a logger and panics if it doesn't match the interface.
func NewContextWithLogger(parentCtx context.Context, logger *Logger) context.Context {
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
	l.logger.With(l.fields...).Debug(msg)
}

// Info logs a Info level message.
func (l *Logger) Info(msg string) {
	l.logger.With(l.fields...).Info(msg)
}

// Warn logs a Warning level message.
func (l *Logger) Warn(msg string) {
	l.logger.With(l.fields...).Warn(msg)
}

// Error logs a Error level message.
func (l *Logger) Error(msg string) {
	l.logger.With(l.fields...).Error(msg)
}

// Fatal logs a Fatal message.
func (l *Logger) Fatal(msg string) {
	l.logger.With(l.fields...).Fatal(msg)
}

// WithField returns a logger with the provided field.
func (l *Logger) WithField(key string, val interface{}) *Logger {
	fieldType := reflectTypeToZapFieldType[reflect.TypeOf(val).Kind()]
	switch fieldType {
	case zapcore.StringType:
		l.fields = append(l.fields, zap.Field{Key: key, String: val.(string), Type: fieldType})
	case zapcore.Int64Type:
		l.fields = append(l.fields, zap.Field{Key: key, Integer: val.(int64), Type: fieldType})
	default:
		// Skip this since we don't know the type
	}
	return l
}

// WithFields returns a logger with the providedfields.
func (l *Logger) WithFields(fields []Field) *Logger {
	for _, field := range fields {
		l.fields = append(l.fields, zap.Field{Key: field.Key, Interface: field.Value})
	}
	return l
}

// WithError returns a logger with the current error padded.
// TODO: How to append this error?
func (l *Logger) WithError(err error) *Logger {
	l.err = err
	return l
}
