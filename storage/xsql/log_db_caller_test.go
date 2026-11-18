package xsql

import (
	"context"
	"database/sql"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gogozs/zlib/xlog"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type callerTestDB struct{}

func (callerTestDB) GetContext(context.Context, interface{}, string, ...interface{}) error {
	return nil
}

func (callerTestDB) SelectContext(context.Context, interface{}, string, ...interface{}) error {
	return nil
}

func (callerTestDB) ExecContext(context.Context, string, ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (callerTestDB) Beginx() (*sqlx.Tx, error) { return nil, nil }

func TestLogDBReportsDAOCallSite(t *testing.T) {
	core, logs := observer.New(zapcore.InfoLevel)
	xlog.SetGlobalLogger(xlog.NewLogger(zap.New(core, zap.AddCaller()).Sugar()))
	t.Cleanup(func() { xlog.SetLoggerLevel("info") })

	wantLine := selectFromDAO()
	entries := logs.All()
	if len(entries) != 1 {
		t.Fatalf("got %d log entries, want 1", len(entries))
	}
	caller := entries[0].Caller
	if filepath.Base(caller.File) != "log_db_caller_test.go" || caller.Line != wantLine {
		t.Errorf("caller = %s:%d, want log_db_caller_test.go:%d", caller.File, caller.Line, wantLine)
	}
}

func selectFromDAO() int {
	_, _, line, _ := runtime.Caller(0)
	_ = WrapLog(callerTestDB{}).SelectContext(context.Background(), nil, "SELECT 1")
	return line + 1
}
