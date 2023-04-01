package xsql

import (
	"context"
	"database/sql"
	"github.com/gogozs/zlib/tools"
	"github.com/jmoiron/sqlx"
)

type SQLConfig struct {
	Host     string
	Username string
	Password string
	Port     int
	Dbname   string
}

type DB interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func DoTx(ctx context.Context, db *sqlx.DB, transactions ...func(db DB) error) (err error) {
	tx, err := db.Beginx()
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			err = tools.ToPanicError(r)
			return
		}
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	if err != nil {
		return err
	}
	for _, f := range transactions {
		if err = f(tx); err != nil {
			return err
		}
	}
	return
}
