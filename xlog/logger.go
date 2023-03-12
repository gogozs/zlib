package xlog

import (
	"context"
	"os"

	"github.com/gogozs/zlib/xtrace"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	ILogger interface {
		Debug(ctx context.Context, msg string, args ...interface{})
		Info(ctx context.Context, msg string, args ...interface{})
		Warn(ctx context.Context, msg string, args ...interface{})
		Error(ctx context.Context, msg string, args ...interface{})
		Fatal(ctx context.Context, msg string, args ...interface{})
		With(...interface{}) ILogger
		MsgItem(string, interface{}) ILogger
		// Sync flushes any buffered log entries.
		Sync()
	}
)

var (
	globalLogger ILogger
)

func init() {
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2)) // skip calling stack
	globalLogger = NewLogger(logger.Sugar())
}

func SetLoggerLevel(level string) {
	zapLevel, _ := zapcore.ParseLevel(level)
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zapLevel,
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2)) // skip calling stack
	globalLogger = NewLogger(logger.Sugar())
}

// SetGlobalLogger set customized logger
func SetGlobalLogger(logger ILogger) {
	globalLogger = logger
}

// Sync flushing any buffered log
func Sync() {
	globalLogger.Sync()
}

func Debug(msg string, args ...interface{}) {
	globalLogger.Debug(context.Background(), msg, args...)
}

func Info(msg string, args ...interface{}) {
	globalLogger.Info(context.Background(), msg, args...)
}

func Warn(msg string, args ...interface{}) {
	globalLogger.Warn(context.Background(), msg, args...)
}

func Error(msg string, args ...interface{}) {
	globalLogger.Error(context.Background(), msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	globalLogger.Fatal(context.Background(), msg, args...)
}

func DebugContext(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Debug(ctx, msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Info(ctx, msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Warn(ctx, msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Error(ctx, msg, args...)
}

func FatalContext(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Fatal(ctx, msg, args...)
}

func MsgItem(msg string, value interface{}) ILogger {
	return globalLogger.MsgItem(msg, value)
}

func WrapTraceWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceKey{}, traceID)
}

func TraceID(ctx context.Context) (bool, string) {
	v := ctx.Value(traceKey{})
	if v != nil {
		return true, v.(string)
	}
	return false, ""
}

func WrapTrace(ctx context.Context) context.Context {
	if exist, _ := TraceID(ctx); exist {
		return ctx
	}
	return context.WithValue(ctx, traceKey{}, xtrace.ParseOrGenTraceID(ctx))
}
