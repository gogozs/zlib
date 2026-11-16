package xsql

import (
	"github.com/gogozs/zlib/xerr"
	"github.com/jmoiron/sqlx"
)

type TransactionWrapperDB struct {
	DB
}

func (t TransactionWrapperDB) Beginx() (*sqlx.Tx, error) {
	db, ok := t.DB.(TransactionDB)
	if !ok {
		return nil, xerr.Errorf("transaction not supported")
	}
	return db.Beginx()
}

func WrapTransaction(db DB) TransactionWrapperDB {
	return TransactionWrapperDB{
		DB: db,
	}
}
