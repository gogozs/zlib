package log

import (
	"context"

	"go.uber.org/zap"
)

type (
	Logger struct {
		logger *zap.SugaredLogger
	}
)

func NewLogger(logger *zap.SugaredLogger) ILogger {
	return Logger{
		logger: logger,
	}
}

func (l Logger) Debug(ctx context.Context, msg string, args ...interface{}) {
	l.logger.Debugf(msg, args...)
}

func (l Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.logger.Infof(msg, args...)
}

func (l Logger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.logger.Warnf(msg, args...)
}

func (l Logger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.logger.Errorf(msg, args...)
}

func (l Logger) Fatal(ctx context.Context, msg string, args ...interface{}) {
	l.logger.Fatalf(msg, args...)
}

func (l Logger) Sync() {
	_ = l.logger.Sync()
}
