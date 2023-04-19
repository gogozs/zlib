package dbgenerator

type Mapping interface {
	GetType(tp SQLType) string
}

type SQLType string

type TypeMap map[string]string

func (m TypeMap) GetType(tp SQLType) string {
	v, ok := m[string(tp)]
	if !ok {
		return "UNKNOWN"
	}
	return v
}

const (
	TinyInt   SQLType = "tinyint"
	SmallInt  SQLType = "smallint"
	MediumInt SQLType = "mediumint"
	Int       SQLType = "int"
	Integer   SQLType = "integer"
	BigInt    SQLType = "bigint"
	float     SQLType = "float"
	Double    SQLType = "double"
	Decimal   SQLType = "decimal"
	Datetime  SQLType = "datetime"
	Date      SQLType = "date"
	TimeStamp SQLType = "timestamp"
	Char      SQLType = "char"
	Varchar   SQLType = "varchar"
	Bit       SQLType = "bit"
	Numeric   SQLType = "numeric"
	Text      SQLType = "text"
	LongText  SQLType = "longtext"
)
