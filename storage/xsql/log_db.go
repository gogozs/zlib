package xsql

import (
	"context"
	"database/sql"
	"github.com/gogozs/zlib/xlog"
)

type LogDB struct {
	db DB
}

func WrapLog(db DB) DB {
	return &LogDB{db: db}
}

func (d LogDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	xlog.MsgItem("sql", query).
		MsgItem("args", args).
		Info(ctx, "[SQL] GET")
	return d.db.GetContext(ctx, dest, query, args...)
}

func (d LogDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	xlog.MsgItem("sql", query).
		MsgItem("args", args).
		Info(ctx, "[SQL] SELECT")
	return d.db.SelectContext(ctx, dest, query, args...)
}

func (d LogDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	xlog.MsgItem("sql", query).
		MsgItem("args", args).
		Info(ctx, "[SQL] EXEC")
	return d.db.ExecContext(ctx, query, args...)
}
