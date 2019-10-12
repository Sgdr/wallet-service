package logger

import (
	"context"
	"os"

	"github.com/go-kit/kit/log"
)

const (
	loggerKey      string = "loggerCtxKey"
	timestampField string = "ts"
)

var defaultLogger log.Logger = nil

func Init() log.Logger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, timestampField, log.DefaultTimestampUTC)
	defaultLogger = logger
	return logger
}

func ToContext(ctx context.Context, l log.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// FromContext returns logger from context with previously added fields.
// If context has not logger returns default logger
func FromContext(ctx context.Context) log.Logger {
	if l, ok := ctx.Value(loggerKey).(log.Logger); ok {
		return l
	} else {
		return defaultLogger
	}
}
