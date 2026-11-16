package xsql

import (
	"context"
	"database/sql"

	"github.com/gogozs/zlib/tools"
	"github.com/jmoiron/sqlx"
)

type DB interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type TransactionDB interface {
	DB
	Beginx() (*sqlx.Tx, error)
}

func DoTx(ctx context.Context, db TransactionDB, transactions ...func(db DB) error) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
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
	for _, f := range transactions {
		if err = f(tx); err != nil {
			return err
		}
	}
	return
}
