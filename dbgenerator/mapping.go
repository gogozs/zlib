package dbgenerator

type Mapping interface {
	GetType(tp SQLType) string
}

var (
	goTypeMap Mapping
)

func init() {
	goTypeMap = TypeMap{
		TinyInt:   "int",
		SmallInt:  "int",
		MediumInt: "int",
		Int:       "int",
		Integer:   "int",
		BigInt:    "uint64",
		Float:     "float64",
		Double:    "float64",
		Decimal:   "float64",
		Datetime:  "time.Time",
		Date:      "time.Time",
		TimeStamp: "uint64",
		Char:      "string",
		Varchar:   "string",
		Bit:       "bool",
		Numeric:   "float64",
		Text:      "string",
		LongText:  "string",
	}
}

func GetGoTypeMap() Mapping {
	return goTypeMap
}

type SQLType string

type TypeMap map[SQLType]string

func (m TypeMap) GetType(tp SQLType) string {
	v, ok := m[tp]
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
	Float     SQLType = "float"
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
