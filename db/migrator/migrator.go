package migrator

import (
	"github.com/xinzf/kit/container/klist"
)

//
//type ColumnType string
//
//const (
//    BIT               ColumnType = "bit"
//    BLOB              ColumnType = "blob"
//    DATE              ColumnType = "date"
//    DATETIME          ColumnType = "datetime"
//    DECIMAL           ColumnType = "decimal"
//    DOUBLE            ColumnType = "double"
//    ENUM              ColumnType = "enum"
//    FLOAT             ColumnType = "float"
//    GEOMETRY          ColumnType = "geometry"
//    MEDIUMINT         ColumnType = "mediumint"
//    JSON              ColumnType = "json"
//    INT               ColumnType = "int"
//    LONGTEXT          ColumnType = "longtext"
//    LONGBLOB          ColumnType = "longblob"
//    BIGINT            ColumnType = "bigint"
//    MEDIUMBLOB        ColumnType = "mediumblob"
//    SET               ColumnType = "set"
//    SMALLINT          ColumnType = "smallint"
//    CHAR              ColumnType = "char"
//    TIME              ColumnType = "time"
//    TIMESTAMP         ColumnType = "timestamp"
//    TINYINT           ColumnType = "tinyint"
//    TINYBLOB          ColumnType = "tinyblob"
//    VARCHAR           ColumnType = "varchar"
//    YEAR              ColumnType = "year"    // mysql end
//    INTEGER           ColumnType = "integer" // postgresql start
//    NUMERIC           ColumnType = "numeric"
//    REAL              ColumnType = "real"
//    DOUBLE_PRECISION  ColumnType = "double precision"
//    SMALLSERIAL       ColumnType = "smallserial"
//    SERIAL            ColumnType = "serial"
//    BIGSERIAL         ColumnType = "bigserial"
//    MONEY             ColumnType = "money"
//    CHARACTER_VARYING ColumnType = "character varying"
//    CHARACTER         ColumnType = "character"
//    TEXT              ColumnType = "text"
//    TIMESTAMPTZ       ColumnType = "timestamptz"
//)

type Table interface {
	Describe() error
	Name() string
	FullName() string
	Comment() string
	Columns() *klist.List[Column]
	Indexes() *klist.List[Index]
	SetName(name string) Table
	SetComment(comment string) Table
	AddColumn(columnName string, ColumnType string) Column
	DropColumn(columnName ...string) Table
	HasColumn(columnName string) bool
	AddIndex(columnNames ...string) Index
	DropIndex(indexName ...string) Table
	HasIndex(name string) bool
	Save() error
}

type Column interface {
	Name() string
	Table() Table
	FullColumnType() string // varchar(255)
	ColumnType() string     // varchar
	PrimaryKey() bool
	AutoIncrement() bool
	Length() int
	DecimalSize() (length, decimal int)
	NullAble() bool
	DefaultValue() any
	Comment() string
	SetPrimaryKey() Column
	SetAutoIncrement(auto bool) Column
	SetName(name string) Column
	SetLength(length int) Column
	SetDecimal(length, scale int) Column
	SetNull(nullAble bool) Column
	SetDefaultValue(value any) Column
	SetComment(comment string) Column
}

type Index interface {
	Table() Table
	Name() string
	Columns() []string
	Unique() bool
	PrimaryKey() bool
	SetUnique(unique bool) Index
	SetColumns(columns ...string) Index
}

type Migrator interface {
	Schema() string
	Tables(name ...string) (*klist.List[Table], error)
	NewTable(tableName string) Table
	DropTable(tableName string) error
	HasTable(tableName string) (bool, error)
	RenameTable(oldName, newName string) error
}
