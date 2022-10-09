package postgresql

import (
	"fmt"
	"github.com/xinzf/kit/container/klist"
	"github.com/xinzf/kit/db/migrator"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Table struct {
	mig             *_migrator
	TableName       string    `json:"name"`
	TableComment    string    `json:"comment"`
	TableCreated    time.Time `json:"created"`
	TableUpdated    time.Time `json:"updated"`
	TableUnicode    string    `json:"unicode"`
	TableColumns    []*Column `json:"columns"`
	TableIndexes    []*Index  `json:"indexes"`
	Tx              *gorm.DB  `json:"-"`
	alters          *klist.List[tableAlter]
	newTableName    string
	isNew           bool
	hasFetchColumns bool
	hasFetchIndex   bool
}

func (this *Table) Name() string {
	return this.TableName
}

func (this *Table) FullName() string {
	return fmt.Sprintf("%s.%s", this.mig.schema, this.TableName)
}

func (this *Table) Comment() string {
	return this.TableComment
}

func (this *Table) Engine() string {
	return ""
}

func (this *Table) CreateTime() time.Time {
	return this.TableCreated
}

func (this *Table) UpdateTime() time.Time {
	return this.TableUpdated
}

func (this *Table) Collation() string {
	return ""
}

func (this *Table) Columns() (*klist.List[migrator.Column], error) {
	if !this.hasFetchColumns {
		tableColumns, err := this.Tx.Migrator().ColumnTypes(this.FullName())
		if err != nil {
			return nil, err
		}
		for _, column := range tableColumns {
			comment, _ := column.Comment()

			columnType, _ := column.ColumnType()

			isPk, _ := column.PrimaryKey()

			isAutoIncre, _ := column.AutoIncrement()

			length, _ := column.Length()

			_len, decimal, ok := column.DecimalSize()
			if ok {
				length = _len
			}

			isNull, _ := column.Nullable()

			defaultVal, _ := column.DefaultValue()

			this.TableColumns = append(this.TableColumns, &Column{
				ColumnName:          column.Name(),
				ColumnDatabaseType:  column.DatabaseTypeName(),
				ColumnFullType:      columnType,
				IsPrimaryKey:        isPk,
				IsAutoIncrement:     isAutoIncre,
				ColumnLength:        int(length),
				ColumnDecimal:       int(decimal),
				IsNullAble:          isNull,
				ColumnDefaultValue:  defaultVal,
				GeneratedExpression: "",
				ColumnComment:       comment,
				table:               this,
			})
		}
	}

	columns := klist.New[migrator.Column]()
	for _, column := range this.TableColumns {
		columns.Add(column)
	}
	return columns, nil
}

func (this *Table) Indexes() (*klist.List[migrator.Index], error) {
	//TODO implement me
	panic("implement me")
}

func (this *Table) SetName(name string) migrator.Table {
	//TODO implement me
	panic("implement me")
}

func (this *Table) SetComment(comment string) migrator.Table {
	//TODO implement me
	panic("implement me")
}

func (this *Table) AddColumn(columnName string, databaseType migrator.DatabaseType) migrator.Column {
	//TODO implement me
	panic("implement me")
}

func (this *Table) DropColumn(columnName ...string) migrator.Table {
	//TODO implement me
	panic("implement me")
}

func (this *Table) HasColumn(columnName string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (this *Table) RenameColumn(oldName, newName string) migrator.Table {
	//TODO implement me
	panic("implement me")
}

func (this *Table) AddIndex(columnNames ...string) migrator.Index {
	//TODO implement me
	panic("implement me")
}

func (this *Table) DropIndex(indexName ...string) migrator.Table {
	//TODO implement me
	panic("implement me")
}

func (this *Table) HasIndex(name string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (this *Table) Save() error {
	//TODO implement me
	panic("implement me")
}

type tableAlter interface {
	toSQL() string
	name() string
}

type renameTableName struct {
	newName string
}

func (this *renameTableName) name() string {
	return "tableName"
}

func (this *renameTableName) toSQL() string {
	return fmt.Sprintf("RENAME TO %s", this.newName)
}

type commentTable struct {
	comment string
}

func (this *commentTable) name() string {
	return "tableComment"
}

func (this *commentTable) toSQL() string {
	return fmt.Sprintf("COMMENT '%s'", this.comment)
}

type newColumn struct {
	column *Column
}

func (this *newColumn) name() string {
	return "newColumn"
}

func (this *newColumn) toSQL() string {
	sql := this.column.ColumnName + " " + this.column.ColumnType()
	if this.column.IsAutoIncrement {
		sql += " auto_increment"
	}
	if this.column.IsPrimaryKey {
		sql += " primary key"
	}

	if this.column.ColumnDefaultValue != nil {
		sql += fmt.Sprintf(" default %v", this.column.ColumnDefaultValue)
	}

	if this.column.IsNullAble {
		sql += " null"
	} else {
		sql += " not null"
	}

	if this.column.ColumnComment != "" {
		sql += fmt.Sprintf(" comment '%s'", this.column.ColumnComment)
	}
	return sql
}

type modifyColumn struct {
	column *Column
}

func (this *modifyColumn) toSQL() string {
	originName := this.column.ColumnName
	if this.column.originName != "" {
		originName = this.column.originName
	}
	sql := fmt.Sprintf("change %s %s %s", originName, this.column.ColumnName, this.column.ColumnFullType)

	if this.column.ColumnDefaultValue != nil {
		sql += fmt.Sprintf(" default %v", this.column.ColumnDefaultValue)
	}

	if this.column.IsNullAble {
		sql += " null"
	} else {
		sql += " not null"
	}

	sql += fmt.Sprintf(" comment '%s'", this.column.ColumnComment)
	return sql
}

func (this *modifyColumn) name() string {
	return "modifyColumn"
}

type dropColumn struct {
	columnName string
}

func (this *dropColumn) name() string {
	return "dropColumn"
}

func (this *dropColumn) toSQL() string {
	return fmt.Sprintf("DROP %s", this.columnName)
}

type newIndex struct {
	index *Index
}

func (this *newIndex) name() string {
	return "newIndex"
}

func (this *newIndex) toSQL() string {
	columns := strings.Join(this.index.ContainColumns, ", ")
	sql := ""
	if this.index.IsUnique {
		sql = fmt.Sprintf("unique index %s (%s)", this.index.IndexName, columns)
	} else {
		sql = fmt.Sprintf("index %s (%s)", this.index.IndexName, columns)
	}
	return sql
}

type modifyIndex struct {
	index *Index
}

func (this *modifyIndex) toSQL() string {
	//TODO implement me
	panic("implement me")
}

func (this *modifyIndex) name() string {
	return "modifyIndex"
}

type dropIndex struct {
	indexName string
}

func (this *dropIndex) name() string {
	return "dropIndex"
}

func (this *dropIndex) toSQL() string {
	return fmt.Sprintf("drop index %s", this.indexName)
}
