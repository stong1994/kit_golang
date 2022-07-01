// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a simplified abstraction of the zap.Logger
type Logger interface {
	Debug(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
	With(fields ...zapcore.Field) Logger
	WithOptions(options ...Option) Logger
}

// logger delegates all calls to the underlying zap.Logger
type logger struct {
	logger Logger
}

// Info logs an info msg with fields
func (l logger) Debug(msg string, fields ...zapcore.Field) {
	l.logger.Debug(msg, fields...)
}

// Info logs an info msg with fields
func (l logger) Info(msg string, fields ...zapcore.Field) {
	l.logger.Info(msg, fields...)
}

// Error logs an warn msg with fields
func (l logger) Warn(msg string, fields ...zapcore.Field) {
	l.logger.Warn(msg, fields...)
}

// Error logs an error msg with fields
func (l logger) Error(msg string, fields ...zapcore.Field) {
	l.logger.Error(msg, fields...)
}

// Fatal logs a fatal error msg with fields
func (l logger) Fatal(msg string, fields ...zapcore.Field) {
	l.logger.Fatal(msg, fields...)
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (l logger) With(fields ...zapcore.Field) Logger {
	return logger{logger: l.logger.With(fields...)}
}

func (l logger) WithOptions(options ...Option) Logger {
	return l.logger.WithOptions(options...)
}

func Debug(msg string, args ...zap.Field) {
	zap.L().Debug(msg, args...)
}

func Info(msg string, args ...zap.Field) {
	zap.L().Info(msg, args...)
}

func Warn(msg string, args ...zap.Field) {
	zap.L().Warn(msg, args...)
}

func Error(msg string, args ...zap.Field) {
	zap.L().Error(msg, args...)
}

func Panic(msg string, args ...zap.Field) {
	zap.L().Panic(msg, args...)
}

func Fatal(msg string, args ...zap.Field) {
	zap.L().Fatal(msg, args...)
}
