package postgresql

import (
	"fmt"
	"github.com/xinzf/kit/db/migrator"
	"strings"
)

type Column struct {
	ColumnName         string `json:"name"`
	FullColumnDataType string `json:"fullColumnType"`
	ColumnDataType     string `json:"columnType"`
	IsPrimaryKey       bool   `json:"primaryKey"`
	IsAutoIncrement    bool   `json:"autoIncrement"`
	ColumnLength       int    `json:"length"`
	DecimalLength      int    `json:"decimalLength"`
	DecimalScale       int    `json:"decimalScale"`
	IsNullAble         bool   `json:"nullAble"`
	ColumnDefaultValue any    `json:"defaultValue"`
	ColumnComment      string `json:"comment"`
	table              *Table `json:"-"`
	origin             *Column
}

func (this *Column) Name() string {
	return this.ColumnName
}

func (this *Column) Table() migrator.Table {
	return this.table
}

func (this *Column) FullColumnType() string {
	return this.FullColumnDataType
}

func (this *Column) ColumnType() string {
	return this.ColumnDataType
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

func (this *Column) DecimalSize() (length, decimal int) {
	return this.DecimalLength, this.DecimalScale
}

func (this *Column) NullAble() bool {
	return this.IsNullAble
}

func (this *Column) DefaultValue() any {
	return this.ColumnDefaultValue
}

func (this *Column) Comment() string {
	return this.ColumnComment
}

func (this *Column) SetPrimaryKey() migrator.Column {
	this.IsPrimaryKey = true
	this.generateFullColumnType()
	return this
}

func (this *Column) SetAutoIncrement(auto bool) migrator.Column {
	this.IsAutoIncrement = auto
	this.generateFullColumnType()
	return this
}

func (this *Column) SetName(name string) migrator.Column {
	this.ColumnName = name
	this.generateFullColumnType()
	return this
}

func (this *Column) SetLength(length int) migrator.Column {
	this.ColumnLength = length
	this.generateFullColumnType()
	return this
}

func (this *Column) SetDecimal(length, scale int) migrator.Column {
	this.DecimalLength = length
	this.DecimalScale = scale
	this.generateFullColumnType()
	return this
}

func (this *Column) SetNull(nullAble bool) migrator.Column {
	this.IsNullAble = nullAble
	this.generateFullColumnType()
	return this
}

func (this *Column) SetDefaultValue(value any) migrator.Column {
	this.ColumnDefaultValue = value
	this.generateFullColumnType()
	return this
}

func (this *Column) SetComment(comment string) migrator.Column {
	this.ColumnComment = comment
	this.generateFullColumnType()
	return this
}

func (this *Column) clone() *Column {
	return &Column{
		ColumnName:         this.ColumnName,
		FullColumnDataType: this.FullColumnDataType,
		ColumnDataType:     this.ColumnDataType,
		IsPrimaryKey:       this.IsPrimaryKey,
		IsAutoIncrement:    this.IsAutoIncrement,
		ColumnLength:       this.ColumnLength,
		DecimalLength:      this.DecimalLength,
		DecimalScale:       this.DecimalScale,
		IsNullAble:         this.IsNullAble,
		ColumnDefaultValue: this.ColumnDefaultValue,
		ColumnComment:      this.ColumnComment,
		table:              this.table,
	}
}

func (this *Column) toSQL() string {
	sql := this.ColumnName + " " + this.FullColumnType()

	if this.ColumnDefaultValue != nil {
		sql += fmt.Sprintf(" default %v", this.ColumnDefaultValue)
	}

	if this.IsNullAble {
		sql += " null"
	} else {
		sql += " not null"
	}
	return sql
}

func (this *Column) generateFullColumnType() {
	if this.AutoIncrement() {
		this.ColumnDataType = "serial"
		this.FullColumnDataType = this.ColumnDataType
		return
	}

	switch strings.ToLower(this.ColumnDataType) {
	case "varchar", "char":
		this.FullColumnDataType = fmt.Sprintf("%s(%d)", this.ColumnDataType, this.ColumnLength)
	case "numeric", "decimal":
		this.FullColumnDataType = fmt.Sprintf("%s(%d,%d)", this.ColumnDataType, this.DecimalLength, this.DecimalScale)
	default:
		this.FullColumnDataType = this.ColumnDataType
	}
}
