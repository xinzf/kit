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
	deleted            bool `json:"deleted"`
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
		deleted:            this.deleted,
	}
}

func (this *Column) alterSqls() []string {
	queries := make([]string, 0)
	if this.deleted {
		if this.origin != nil {
			queries = append(queries, fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s cascade", this.table.FullName(), this.ColumnName))
		}
	} else {
		if this.origin == nil {
			queries = append(queries, fmt.Sprintf("ALTER TABLE %s ADD %s", this.table.FullName(), this.createSQL()))
			if this.ColumnComment != "" {
				queries = append(queries, fmt.Sprintf("comment on column %s.%s is '%s'", this.table.FullName(), this.ColumnName, this.ColumnComment))
			}
		} else {
			if this.ColumnName != this.origin.ColumnName {
				queries = append(queries, fmt.Sprintf("alter table %s rename column %s to %s", this.table.FullName(), this.origin.ColumnName, this.ColumnName))
			}
			if this.FullColumnDataType != this.origin.FullColumnDataType {
				queries = append(queries, fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s using %s::%s", this.table.FullName(), this.ColumnName, this.FullColumnDataType, this.ColumnName, this.origin.FullColumnType()))
			}
			if this.NullAble() != this.origin.NullAble() {
				if this.NullAble() {
					queries = append(queries, fmt.Sprintf("alter table %s alter column %s drop not null", this.table.FullName(), this.ColumnName))
				} else {
					queries = append(queries, fmt.Sprintf("alter table %s alter column %s set not null", this.table.FullName(), this.ColumnName))
				}
			}
			if this.ColumnDefaultValue != this.origin.ColumnDefaultValue {
				if this.ColumnDefaultValue != nil {
					queries = append(queries, fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s SET DEFAULT '%v'", this.table.FullName(), this.ColumnName, this.ColumnDefaultValue))
				} else {
					queries = append(queries, fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s DROP DEFAULT", this.table.FullName(), this.ColumnName))
				}
			}
			if this.ColumnComment != this.origin.ColumnComment {
				queries = append(queries, fmt.Sprintf("comment on column %s.%s is '%s'", this.table.FullName(), this.ColumnName, this.ColumnComment))
			}
		}

	}
	return queries
}

func (this *Column) createSQL() string {
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
	}

	switch strings.ToLower(this.ColumnDataType) {
	case "varchar", "char", "varying", "character":
		this.FullColumnDataType = fmt.Sprintf("%s(%d)", this.ColumnDataType, this.ColumnLength)
	case "numeric", "decimal":
		this.FullColumnDataType = fmt.Sprintf("%s(%d,%d)", this.ColumnDataType, this.DecimalLength, this.DecimalScale)
	default:
		this.FullColumnDataType = this.ColumnDataType
	}
}
