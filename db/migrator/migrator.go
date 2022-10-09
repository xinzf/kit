package migrator

import (
	"github.com/xinzf/kit/container/klist"
	"time"
)

type DatabaseType string

const (
	BIT               DatabaseType = "bit"
	BLOB              DatabaseType = "blob"
	DATE              DatabaseType = "date"
	DATETIME          DatabaseType = "datetime"
	DECIMAL           DatabaseType = "decimal"
	DOUBLE            DatabaseType = "double"
	ENUM              DatabaseType = "enum"
	FLOAT             DatabaseType = "float"
	GEOMETRY          DatabaseType = "geometry"
	MEDIUMINT         DatabaseType = "mediumint"
	JSON              DatabaseType = "json"
	INT               DatabaseType = "int"
	LONGTEXT          DatabaseType = "longtext"
	LONGBLOB          DatabaseType = "longblob"
	BIGINT            DatabaseType = "bigint"
	MEDIUMBLOB        DatabaseType = "mediumblob"
	SET               DatabaseType = "set"
	SMALLINT          DatabaseType = "smallint"
	CHAR              DatabaseType = "char"
	TIME              DatabaseType = "time"
	TIMESTAMP         DatabaseType = "timestamp"
	TINYINT           DatabaseType = "tinyint"
	TINYBLOB          DatabaseType = "tinyblob"
	VARCHAR           DatabaseType = "varchar"
	YEAR              DatabaseType = "year"    // mysql end
	INTEGER           DatabaseType = "integer" // postgresql start
	NUMERIC           DatabaseType = "numeric"
	REAL              DatabaseType = "real"
	DOUBLE_PRECISION  DatabaseType = "double precision"
	SMALLSERIAL       DatabaseType = "smallserial"
	SERIAL            DatabaseType = "serial"
	BIGSERIAL         DatabaseType = "bigserial"
	MONEY             DatabaseType = "money"
	CHARACTER_VARYING DatabaseType = "character varying"
	CHARACTER         DatabaseType = "character"
	TEXT              DatabaseType = "text"
	TIMESTAMPTZ       DatabaseType = "timestamptz"
)

type Table interface {
	Name() string
	FullName() string
	Comment() string
	Engine() string
	CreateTime() time.Time
	UpdateTime() time.Time
	Collation() string
	Columns() (*klist.List[Column], error)
	Indexes() (*klist.List[Index], error)
	SetName(name string) Table
	SetComment(comment string) Table
	AddColumn(columnName string, databaseType DatabaseType) Column
	DropColumn(columnName ...string) Table
	HasColumn(columnName string) (bool, error)
	RenameColumn(oldName, newName string) Table
	AddIndex(columnNames ...string) Index
	DropIndex(indexName ...string) Table
	HasIndex(name string) (bool, error)
	Save() error
}

type Column interface {
	Name() string
	Table() Table
	ColumnType() string         // varchar(255)
	DatabaseType() DatabaseType // varchar
	PrimaryKey() bool
	AutoIncrement() bool
	Length() int
	Decimal() int
	NullAble() bool
	DefaultValue() any
	Generated() string
	Comment() string
	SetPrimaryKey() Column
	SetAutoIncrement(auto bool) Column
	SetName(name string) Column
	SetLength(length int) Column
	SetDecimal(decimal int) Column
	SetNull(nullAble bool) Column
	SetDefaultValue(value any) Column
	SetComment(comment string) Column
}

type Index interface {
	Table() Table
	Name() string
	Columns() []string
	Unique() bool
	SetUnique(unique bool) Index
	SetColumns(columns ...string) Index
}

type Migrator interface {
	Schema() string
	Tables() (*klist.List[Table], error)
	NewTable(tableName string) Table
	DropTable(tableName string) error
	HasTable(tableName string) (bool, error)
	RenameTable(oldName, newName string) error
}
