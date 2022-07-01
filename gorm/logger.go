package gorm

import (
	"context"
	"errors"
	"fmt"
	"github.com/stong1994/kit_golang/log"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Logger struct {
	loggerFactory             *log.Factory
	LogLevel                  gormlogger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
	ignorePaths               []string
}

func (l Logger) Name() string {
	return "zapTraceLogger"
}

func (l Logger) Initialize(db *gorm.DB) error {
	if l.loggerFactory == nil {
		return errors.New("not set logger factory yet")
	}
	return nil
}

func NewLogger(loggerFactory *log.Factory) *Logger {
	return &Logger{
		loggerFactory:             loggerFactory,
		LogLevel:                  gormlogger.Warn,
		SlowThreshold:             100 * time.Millisecond,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: false,
		ignorePaths:               []string{gormPackage, logPackage},
	}
}

func (l Logger) SetAsDefault() {
	gormlogger.Default = l
}

func (l *Logger) AddIgnorePaths(paths ...string) {
	for _, path := range paths {
		l.ignorePaths = append(l.ignorePaths, path)
	}
	gormlogger.Default = l
}

func (l Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return Logger{
		loggerFactory:             l.loggerFactory,
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		SkipCallerLookup:          l.SkipCallerLookup,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
	}
}

func (l Logger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Info {
		return
	}
	l.logger(ctx).Debug(fmt.Sprintf(str, args...))
}

func (l Logger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Warn {
		return
	}
	l.logger(ctx).Warn(fmt.Sprintf(str, args...))
}

func (l Logger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Error {
		return
	}
	l.logger(ctx).Error(fmt.Sprintf(str, args...))
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		l.logger(ctx).Error("gorm", zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		l.logger(ctx).Warn("gorm", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.LogLevel >= gormlogger.Info:
		sql, rows := fc()
		l.logger(ctx).Debug("gorm", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	}
}

var (
	gormPackage = "kit_golang/gorm"
	logPackage  = "kit_golang/log"
)

func (l Logger) logger(ctx context.Context) log.Logger {
	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		if !ok {
			continue
		}
		contain := false
		for _, v := range l.ignorePaths {
			if strings.Contains(file, v) {
				contain = true
				break
			}
		}
		if !contain {
			return l.loggerFactory.For(ctx).WithOptions(log.AddCallerSkip(i))
		}
	}
	return l.loggerFactory.For(ctx)
}
