package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Log struct {
	level      zap.AtomicLevel
	logger     *zap.Logger
	addCaller  bool
	localWrite bool
}

func NewLogger(options ...Option) *Log {
	logger := new(Log)
	for _, op := range options {
		op.apply(logger)
	}
	if logger.logger == nil {
		logger.initLog()
	}
	return logger
}

func (l *Log) initLog() {
	var allCore []zapcore.Core

	localWriter := zapcore.Lock(os.Stdout)
	if l.localWrite {
		cfg := zap.NewDevelopmentEncoderConfig()
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(cfg)
		allCore = append(allCore, zapcore.NewCore(consoleEncoder, localWriter, l.level))
	} else {
		cfg := zap.NewProductionEncoderConfig()
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder
		jsonEncoder := zapcore.NewJSONEncoder(cfg)
		allCore = append(allCore, zapcore.NewCore(jsonEncoder, localWriter, l.level))
	}

	core := zapcore.NewTee(allCore...)

	var options []zap.Option
	if l.addCaller {
		options = append(options, zap.AddCaller())
	}

	l.logger = zap.New(core).WithOptions(options...)
}

func (l *Log) clone() *Log {
	copy := *l
	return &copy
}

func (l *Log) SetLevel(lvl string) {
	l.level.SetLevel(getLoggerLevel(lvl))
}

func (l *Log) Debug(msg string, args ...zap.Field) {
	l.logger.Debug(msg, args...)
}

func (l *Log) Info(msg string, args ...zap.Field) {
	l.logger.Info(msg, args...)
}

func (l *Log) Warn(msg string, args ...zap.Field) {
	l.logger.Warn(msg, args...)
}

func (l *Log) Error(msg string, args ...zap.Field) {
	l.logger.Error(msg, args...)
}

func (l *Log) Panic(msg string, args ...zap.Field) {
	l.logger.Panic(msg, args...)
}

func (l *Log) Fatal(msg string, args ...zap.Field) {
	l.logger.Fatal(msg, args...)
}

func (l *Log) With(fields ...zapcore.Field) Logger {
	c := l.clone()
	l.logger = l.logger.With(fields...)
	return c
}

func (l *Log) WithOptions(options ...Option) Logger {
	c := l.clone()
	for _, op := range options {
		op.apply(c)
	}
	return c
}
