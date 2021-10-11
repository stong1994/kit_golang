package slog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

var level zap.AtomicLevel

func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

func InitLog(addCaller bool, skip int, isLocal bool) {
	zap.ReplaceGlobals(initLog(addCaller, skip, isLocal))
}

func initLog(addCaller bool, skip int, isLocal bool) *zap.Logger {
	level = zap.NewAtomicLevel()
	var allCore []zapcore.Core

	localWriter := zapcore.Lock(os.Stdout)
	if isLocal {
		cfg := zap.NewDevelopmentEncoderConfig()
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(cfg)
		allCore = append(allCore, zapcore.NewCore(consoleEncoder, localWriter, level))
	} else {
		cfg := zap.NewProductionEncoderConfig()
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder
		jsonEncoder := zapcore.NewJSONEncoder(cfg)
		allCore = append(allCore, zapcore.NewCore(jsonEncoder, localWriter, level))
	}

	core := zapcore.NewTee(allCore...)

	var options []zap.Option
	if addCaller {
		options = append(options, zap.AddCaller(), zap.AddCallerSkip(skip))
	}

	return zap.New(core).WithOptions(options...)
}

func SetLevel(lvl string) {
	level.SetLevel(getLoggerLevel(lvl))
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

func DPanic(msg string, args ...zap.Field) {
	zap.L().DPanic(msg, args...)
}

func Panic(msg string, args ...zap.Field) {
	zap.L().Panic(msg, args...)
}

func Fatal(msg string, args ...zap.Field) {
	zap.L().Fatal(msg, args...)
}
