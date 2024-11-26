package dbgenerator

import (
	"context"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gogozs/zlib/storage/xsql"
)

type (
	DBParser struct {
		db     xsql.DB
		dbname string
	}

	TableData struct {
		TableName string
		ModalName KeyString
		Columns   []Column
	}

	Column struct {
		ColumnName KeyString
		Type       string
	}

	DBColumn struct {
		ColumnName             string  `db:"COLUMN_NAME"`
		DataType               string  `db:"DATA_TYPE"`
		CharacterMaximumLength int     `db:"CHARACTER_MAXIMUM_LENGTH"`
		ColumnType             string  `db:"COLUMN_TYPE"`
		NumericPrecision       int     `db:"NUMERIC_PRECISION"`
		NumericScale           int     `db:"NUMERIC_SCALE"`
		ColumnComment          string  `db:"COLUMN_COMMENT"`
		ColumnDefault          *string `db:"COLUMN_DEFAULT"`
		IsNullable             string  `db:"IS_NULLABLE"`
		ColumnKey              string  `db:"COLUMN_KEY"`
		Extra                  string  `db:"EXTRA"`
		OrdinalPosition        string  `db:"ORDINAL_POSITION"`
	}
)

func NewDBParser(db xsql.DB, dbname string) *DBParser {
	return &DBParser{db: db, dbname: dbname}
}

const (
	tableSQL = `
select table_name
from information_schema.tables
where table_schema = ?
	`

	schemaSQL = `
select COLUMN_NAME,
       DATA_TYPE,
       IFNULL(CHARACTER_MAXIMUM_LENGTH, 0) as CHARACTER_MAXIMUM_LENGTH,
       COLUMN_TYPE,
       IFNULL(NUMERIC_PRECISION, 0) as NUMERIC_PRECISION,
       IFNULL(NUMERIC_SCALE, 0) as NUMERIC_SCALE,
       COLUMN_COMMENT,
       COLUMN_DEFAULT,
       IS_NULLABLE,
       COLUMN_KEY,
       EXTRA,
       ORDINAL_POSITION
from information_schema.COLUMNS
where table_schema = ?
  and table_name = ?
  `
)

func (p DBParser) GetTables() (tables []string, err error) {
	if err = p.db.SelectContext(context.Background(), &tables, tableSQL, p.dbname); err != nil {
		return nil, err
	}
	return
}

func (p DBParser) GetTableSchema(tableName string) (dbColumns []DBColumn, err error) {
	if err = p.db.SelectContext(context.Background(), &dbColumns, schemaSQL, p.dbname, tableName); err != nil {
		return nil, err
	}
	return
}

func (p DBParser) GetColumns(dbColumns []DBColumn, m Mapping) (columns []Column, err error) {
	columns = make([]Column, 0, len(dbColumns))
	for _, col := range dbColumns {
		columns = append(columns, Column{
			ColumnName: KeyString(col.ColumnName),
			Type:       m.GetType(SQLType(col.DataType)),
		})
	}

	return columns, nil
}

func (p DBParser) BuildTableData(tableName string, columns []Column, prefix string) (*TableData, error) {
	return &TableData{
		TableName: tableName,
		ModalName: KeyString(tableName[len(prefix):]),
		Columns:   columns,
	}, nil
}
