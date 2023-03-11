package xlog

import (
	"context"

	"go.uber.org/zap"
)

type (
	Logger struct {
		logger *zap.SugaredLogger
	}
	traceKey struct {
	}
)

const traceLogItem = "_TRACE_ID_"

func NewLogger(logger *zap.SugaredLogger) ILogger {
	return Logger{
		logger: logger,
	}
}

func WrapTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceKey{}, traceID)
}

func (l Logger) Debug(ctx context.Context, msg string, args ...interface{}) {
	l.parse(ctx).Debugf(msg, args...)
}

func (l Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.parse(ctx).Infof(msg, args...)
}

func (l Logger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.parse(ctx).Warnf(msg, args...)
}

func (l Logger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.parse(ctx).Errorf(msg, args...)
}

func (l Logger) Fatal(ctx context.Context, msg string, args ...interface{}) {
	l.parse(ctx).Fatalf(msg, args...)
}

func (l Logger) With(fields ...interface{}) ILogger {
	return NewLogger(l.logger.With(fields...))
}

func (l Logger) MsgItem(msg string, value interface{}) ILogger {
	return NewLogger(l.logger.With(msg, value))
}

func (l Logger) Sync() {
	_ = l.logger.Sync()
}

func (l Logger) parse(ctx context.Context) *zap.SugaredLogger {
	logger := l.logger
	if v := ctx.Value(traceKey{}); v != nil {
		logger = logger.With(traceLogItem, v)
	}
	return logger
}
