package mysql

import (
	"fmt"
	"github.com/xinzf/kit/container/kvar"
	"github.com/xinzf/kit/db/migrator"
)

type Column struct {
	ColumnName          string `json:"name"`           // 字段名称
	ColumnDatabaseType  string `json:"database_type"`  // 数据库字段类型 varchar
	ColumnFullType      string `json:"column_type"`    // varchar(255)
	IsPrimaryKey        bool   `json:"primary_key"`    // 是否为主键
	IsAutoIncrement     bool   `json:"auto_increment"` // 是否自增
	ColumnLength        int    `json:"length"`         // 数据列长度
	ColumnDecimal       int    `json:"decimal"`        // 小数点后位数
	IsNullAble          bool   `json:"null_able"`      // 默认 null
	ColumnDefaultValue  any    `json:"default_value"`  // 默认值
	GeneratedExpression string `json:"expression"`     // 虚拟列表达式
	ColumnComment       string `json:"comment"`        // 字段描述
	table               *Table
	isNew               bool
	originName          string
}

func (this *Column) Name() string {
	return this.ColumnName
}

func (this *Column) Table() migrator.Table {
	return this.table
}

func (this *Column) ColumnType() string {
	return this.ColumnFullType
}

func (this *Column) DatabaseType() migrator.DatabaseType {
	return migrator.DatabaseType(this.ColumnDatabaseType)
}

func (this *Column) PrimaryKey() bool {
	return this.IsPrimaryKey
}

func (this *Column) AutoIncrement() bool {
	return this.IsAutoIncrement
}

func (this *Column) Length() int {
	return this.ColumnLength
}

func (this *Column) Decimal() int {
	return this.ColumnDecimal
}

func (this *Column) NullAble() bool {
	return this.IsNullAble
}

func (this *Column) DefaultValue() any {
	return this.ColumnDefaultValue
}

func (this *Column) Generated() string {
	return this.GeneratedExpression
}

func (this *Column) Comment() string {
	return this.ColumnComment
}

func (this *Column) SetPrimaryKey() migrator.Column {
	if !this.IsPrimaryKey {
		this.IsPrimaryKey = true
		this.addAlter()
	}
	if this.IsPrimaryKey {
		this.IsNullAble = false
	}
	return this
}

func (this *Column) SetAutoIncrement(auto bool) migrator.Column {
	if this.IsAutoIncrement != auto {
		this.IsAutoIncrement = auto
		this.addAlter()
	}
	return this
}

func (this *Column) SetName(name string) migrator.Column {
	if this.ColumnName != name {
		this.originName = this.ColumnName
		this.ColumnName = name
		this.addAlter()
	}
	return this
}

func (this *Column) SetLength(length int) migrator.Column {
	if this.ColumnLength != length {
		this.ColumnLength = length
		this.generateColumnType()
		this.addAlter()
	}
	return this
}

func (this *Column) SetDecimal(decimal int) migrator.Column {
	if this.ColumnDecimal != decimal {
		this.ColumnDecimal = decimal
		this.generateColumnType()
		this.addAlter()
	}
	return this
}

func (this *Column) SetNull(nullAble bool) migrator.Column {
	if this.IsNullAble != nullAble {
		this.IsNullAble = nullAble
		this.addAlter()
	}
	return this
}

func (this *Column) SetDefaultValue(value any) migrator.Column {
	if !kvar.New(this.ColumnDefaultValue).Equal(kvar.New(value)) {
		this.ColumnDefaultValue = value
		this.addAlter()
	}
	return this
}

func (this *Column) SetComment(comment string) migrator.Column {
	if this.ColumnComment != comment {
		this.ColumnComment = comment
		this.addAlter()
	}
	return this
}

func (this *Column) generateColumnType() {
	switch this.ColumnDatabaseType {
	case string(migrator.VARCHAR), string(migrator.INT), string(migrator.BIGINT), string(migrator.CHAR):
		this.ColumnFullType = fmt.Sprintf("%s(%d)", this.ColumnDatabaseType, this.Length())
	case string(migrator.DECIMAL):
		this.ColumnFullType = fmt.Sprintf("%s(%d,%d)", this.ColumnDatabaseType, this.Length(), this.Decimal())
	default:
		this.ColumnFullType = this.ColumnDatabaseType
	}
}

func (this *Column) addAlter() {
	if this.isNew {
		return
	}

	idx, _ := this.table.alters.Find(func(alt tableAlter) bool {
		return alt.name() == "modifyColumn" && alt.(*modifyColumn).column.ColumnName == this.ColumnName
	})
	if idx == -1 {
		this.table.alters.Add(&modifyColumn{column: this})
	}
}
