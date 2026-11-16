package xsql

import (
	"context"
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
)

type transactionStub struct {
	DB
	calls int
	tx    *sqlx.Tx
	err   error
}

func (s *transactionStub) Beginx() (*sqlx.Tx, error) {
	s.calls++
	return s.tx, s.err
}

func TestWrapTransaction(t *testing.T) {
	for _, depth := range []int{0, 1, 2} {
		stub := &transactionStub{tx: &sqlx.Tx{}}
		var db DB = stub
		for i := 0; i < depth; i++ {
			db = WrapLog(db)
		}
		tx, err := WrapTransaction(db).Beginx()
		if err != nil || tx != stub.tx || stub.calls != 1 {
			t.Fatalf("logging depth %d: tx=%p, err=%v, calls=%d", depth, tx, err, stub.calls)
		}
	}
}

func TestDoTxBeginError(t *testing.T) {
	expected := errors.New("begin failed")
	stub := &transactionStub{err: expected}
	called := false
	err := DoTx(context.Background(), WrapTransaction(WrapLog(stub)), func(DB) error {
		called = true
		return nil
	})
	if !errors.Is(err, expected) || called {
		t.Fatalf("err=%v, callback called=%v", err, called)
	}
}

func TestWrapTransactionUnsupported(t *testing.T) {
	db := WrapTransaction(WrapLog(&LogDB{}))
	if _, err := db.Beginx(); err == nil {
		t.Fatal("expected error for unsupported database")
	}
}
