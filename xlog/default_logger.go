package xlog

import (
	"context"

	"go.uber.org/zap"
)

type (
	Logger struct {
		logger     *zap.SugaredLogger
		callerSkip int
	}
	traceKey struct {
	}
)

const traceLogItem = "_TRACE_ID_"

func NewLogger(logger *zap.SugaredLogger) ILogger {
	return Logger{logger: logger}
}

func (l Logger) Debug(ctx context.Context, msg string, args ...interface{}) {
	l.parse(ctx, 1).Debugf(msg, args...)
}

func (l Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.parse(ctx, 1).Infof(msg, args...)
}

func (l Logger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.parse(ctx, 1).Warnf(msg, args...)
}

func (l Logger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.parse(ctx, 1).Errorf(msg, args...)
}

func (l Logger) Fatal(ctx context.Context, msg string, args ...interface{}) {
	l.parse(ctx, 1).Fatalf(msg, args...)
}

func (l Logger) With(fields ...interface{}) ILogger {
	return Logger{logger: l.logger.With(fields...), callerSkip: l.callerSkip}
}

func (l Logger) MsgItem(msg string, value interface{}) ILogger {
	return Logger{logger: l.logger.With(msg, value), callerSkip: l.callerSkip}
}

func (l Logger) Sync() {
	_ = l.logger.Sync()
}

func (l Logger) parse(ctx context.Context, callerSkip int) *zap.SugaredLogger {
	logger := l.logger
	callerSkip += l.callerSkip
	if callerSkip > 0 {
		logger = logger.Desugar().WithOptions(zap.AddCallerSkip(callerSkip)).Sugar()
	}
	if exist, v := TraceID(ctx); exist {
		logger = logger.With(traceLogItem, v)
	}
	return logger
}
