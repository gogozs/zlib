package log

import (
	"context"
	"testing"
)

func TestDebug(t *testing.T) {
	Debug("test: %s", "content")
	Info("test: %s", "content")
	Warn("test: %s", "content")
	Error("test: %s", "content")
	DebugContext(context.Background(), "test: %s", "content")
	InfoContext(context.Background(), "test: %s", "content")
	WarnContext(context.Background(), "test: %s", "content")
	ErrorContext(context.Background(), "test: %s", "content")
}
