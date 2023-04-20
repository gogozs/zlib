package dbgenerator

import (
	"fmt"
	"os"
	"testing"

	"github.com/gogozs/zlib/storage/xsql"
	"github.com/gogozs/zlib/xlog"
	"github.com/stretchr/testify/require"
)

var (
	testClient xsql.DB
	dbname     = "admin"
)

func TestMain(m *testing.M) {
	var err error
	testClient, err = xsql.NewLogDB(&xsql.SQLConfig{
		Host:     "127.0.0.1",
		Username: "root",
		Password: "test",
		Port:     3306,
		Dbname:   dbname,
	})
	if err != nil {
		xlog.Error("mysql connect fail: %v.", err)
		xlog.Info("[Skip Test]")
		return
	}
	m.Run()
	os.Exit(0)
}

func TestDBParser_GetTables(t *testing.T) {
	parser := initTestParser()
	tables, err := parser.GetTables()
	require.Nil(t, err)
	fmt.Println(tables)
}

func TestDBParser_GetTableSchema(t *testing.T) {
	parser := initTestParser()
	columns, err := parser.GetTableSchema("t_user")
	require.Nil(t, err)
	fmt.Println(columns)
}

func initTestParser() *DBParser {
	return NewDBParser(testClient, dbname)
}
