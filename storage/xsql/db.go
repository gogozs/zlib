package xsql

import (
	"context"
	"database/sql"
	"github.com/gogozs/zlib/tools"
	"github.com/jmoiron/sqlx"
)

type (
	DB interface {
		GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	}
	TransactionDB interface {
		DB
		DoTx(ctx context.Context, transactions ...func(ctx context.Context, db DB) error) error
	}
	txDB struct {
		*sqlx.DB
	}
)

func (db txDB) DoTx(ctx context.Context, transactions ...func(ctx context.Context, db DB) error) (err error) {
	var tx *sqlx.Tx
	tx, err = db.Beginx()
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
		if err = f(ctx, tx); err != nil {
			return err
		}
	}
	return nil
}
