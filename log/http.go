package log

import (
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"time"
)

type StructuredLoggerEntry struct {
	Logger *zap.Logger
}

func (l *StructuredLoggerEntry) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: zap.L()}
	var logFields []zapcore.Field

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields = append(logFields, zapcore.Field{
			Key:    "req_id",
			Type:   zapcore.StringType,
			String: reqID,
		})
	}

	logFields = append(logFields, zapcore.Field{
		Key:    "http_method",
		Type:   zapcore.StringType,
		String: r.Method,
	}, zapcore.Field{
		Key:    "remote_addr",
		Type:   zapcore.StringType,
		String: r.RemoteAddr,
	}, zapcore.Field{
		Key:    "uri",
		Type:   zapcore.StringType,
		String: r.RequestURI,
	})

	entry.Logger = entry.Logger.With(logFields...)

	entry.Logger.Info("Request started")

	return entry
}

func NewStructuredLogger() func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLoggerEntry{zap.L()})
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger.Info("Request completed", zap.Field{
		Key:     "resp_status",
		Type:    zapcore.Int64Type,
		Integer: int64(status),
	}, zap.Field{
		Key:     "resp_bytes_length",
		Type:    zapcore.Int64Type,
		Integer: int64(bytes),
	}, zap.Field{
		Key:       "resp_elapsed",
		Type:      zapcore.DurationType,
		Interface: elapsed,
	})
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger.Info("panic", zap.Field{
		Key:    "stack",
		String: string(stack),
		Type:   zapcore.StringerType,
	}, zap.Field{
		Key:       "panic",
		Type:      zapcore.StringerType,
		Interface: v,
	})
}
