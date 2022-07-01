package log

import "go.uber.org/zap"

type Option interface {
	apply(*Log)
}

type optionFunc func(*Log)

func (f optionFunc) apply(logger *Log) {
	f(logger)
}

func AddCaller() Option {
	return optionFunc(func(logger *Log) {
		logger.addCaller = true
	})
}

func SetLevel(level string) Option {
	return optionFunc(func(logger *Log) {
		logger.level = zap.NewAtomicLevelAt(getLoggerLevel(level))
	})
}

func SetLocalWrite(localWrite bool) Option {
	return optionFunc(func(logger *Log) {
		logger.localWrite = localWrite
	})
}

func SetLogger(logger *zap.Logger) Option {
	return optionFunc(func(l *Log) {
		l.logger = logger
	})
}

func AddCallerSkip(skip int) Option {
	return optionFunc(func(logger *Log) {
		logger.logger = logger.logger.WithOptions(zap.AddCallerSkip(skip))
	})
}
