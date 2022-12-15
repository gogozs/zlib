package log

import (
	"context"
	"testing"
)

func TestDebug(t *testing.T) {
	Info(context.Background(), "test: %s", "content")
}
