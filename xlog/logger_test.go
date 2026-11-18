package xlog

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggerReportsCallSite(t *testing.T) {
	core, logs := observer.New(zapcore.InfoLevel)
	previous := globalLogger
	SetGlobalLogger(NewLogger(zap.New(core, zap.AddCaller()).Sugar()))
	t.Cleanup(func() { SetGlobalLogger(previous) })

	wantDirectLine := logDirectly()
	_, _, wantWrappedLine, _ := runtime.Caller(0)
	logThroughHelper()
	wantWrappedLine++
	entries := logs.All()
	if len(entries) != 2 {
		t.Fatalf("got %d log entries, want 2", len(entries))
	}
	for i, wantLine := range []int{wantDirectLine, wantWrappedLine} {
		entry := entries[i].Entry
		if filepath.Base(entry.Caller.File) != "logger_test.go" || entry.Caller.Line != wantLine {
			t.Errorf("entry %d caller = %s:%d, want logger_test.go:%d", i, entry.Caller.File, entry.Caller.Line, wantLine)
		}
	}
}

func logDirectly() int {
	_, _, line, _ := runtime.Caller(0)
	InfoContext(context.Background(), "direct call")
	return line + 1
}

func logThroughHelper() {
	logger := MsgItem("key", "value")
	WithCallerSkip(logger, 1).Info(context.Background(), "wrapped call")
}

func TestLogger(t *testing.T) {
	Debug("test: %s", "content")
	Info("test: %s", "content")
	Warn("test: %s", "content")
	Error("test: %s", "content")
	DebugContext(context.Background(), "test: %s", "content")
	InfoContext(context.Background(), "test: %s", "content")
	WarnContext(context.Background(), "test: %s", "content")
	ErrorContext(context.Background(), "test: %s", "content")
}

func TestWrapTraceID(t *testing.T) {
	ctx := WrapTraceWithTraceID(context.Background(), "testID")
	InfoContext(ctx, "test: %s", "content")
}

func TestWrapTrace(t *testing.T) {
	ctx := WrapTrace(context.Background())
	InfoContext(ctx, "test: %s", "content")
}
