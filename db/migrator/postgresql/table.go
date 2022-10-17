package postgresql

import (
	"fmt"
	"github.com/xinzf/kit/container/klist"
	"github.com/xinzf/kit/db/migrator"
	"gorm.io/gorm"
	"strings"
)

type Table struct {
	TableName    string     `json:"name"`
	TableComment string     `json:"comment"`
	TableColumns []*Column  `json:"columns"`
	TableIndexes []*Index   `json:"indexes"`
	mig          *_migrator `json:"-"`
	origin       *Table     `json:"-"`
	dropColumns  []*Column  `json:"-"`
	hasDescribe  bool
}

func (this *Table) Describe() error {
	if !this.hasDescribe {
		defer func() {
			this.hasDescribe = true
		}()

		if err := this.loadColumns(); err != nil {
			return err
		}

		return this.loadIndexes()
	}
	return nil
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

func (this *Table) loadColumns() error {
	this.TableColumns = []*Column{}
	cols, err := this.mig.tx.Migrator().ColumnTypes(this.FullName())
	if err != nil {
		return err
	}

	for _, col := range cols {
		tpe, _ := col.ColumnType()

		pk, _ := col.PrimaryKey()
		auto, _ := col.AutoIncrement()
		length, _ := col.Length()
		decimalLen, decimalDec, hasDecimal := col.DecimalSize()
		nullAble, _ := col.Nullable()
		dft, _ := col.DefaultValue()
		com, _ := col.Comment()
		column := &Column{
			ColumnName:         col.Name(),
			FullColumnDataType: tpe,
			ColumnDataType:     col.DatabaseTypeName(),
			IsPrimaryKey:       pk,
			IsAutoIncrement:    auto,
			ColumnLength:       int(length),
			IsNullAble:         nullAble,
			ColumnDefaultValue: dft,
			ColumnComment:      com,
			table:              this,
		}
		if hasDecimal {
			column.DecimalScale = int(decimalDec)
			column.DecimalLength = int(decimalLen)
		}

		column.generateFullColumnType()
		column.origin = column.clone()

		this.TableColumns = append(this.TableColumns, column)
	}

	return nil
}

func (this *Table) Columns() *klist.List[migrator.Column] {
	if this.TableColumns == nil {
		this.TableColumns = []*Column{}
	}
	columns := klist.New[migrator.Column]()
	for _, column := range this.TableColumns {
		columns.Add(column)
	}
	return columns
}

func (this *Table) GetColumn(name string) (migrator.Column, bool) {
	idx, column := this.Columns().Find(func(col migrator.Column) bool {
		return col.Name() == name
	})
	if idx == -1 {
		return nil, false
	}

	return column, true
}

func (this *Table) loadIndexes() error {

	this.TableIndexes = []*Index{}

	idxes, err := this.mig.tx.Migrator().GetIndexes(this.Name())
	if err != nil {
		return err
	}

	for _, idx := range idxes {
		unique, _ := idx.Unique()
		pk, _ := idx.PrimaryKey()
		if pk {
			col := idx.Columns()[0]
			for _, _col := range this.TableColumns {
				if _col.Name() == col {
					_col.SetPrimaryKey()
					break
				}
			}
		}

		_index := &Index{
			IndexName:    idx.Name(),
			IndexColumns: idx.Columns(),
			table:        this,
			IsUnique:     unique,
			IsPK:         pk,
		}

		_index.generateName()
		_index.origin = _index.clone()

		this.TableIndexes = append(this.TableIndexes, _index)
	}
	return nil
}

func (this *Table) Indexes() *klist.List[migrator.Index] {
	if this.TableIndexes == nil {
		this.TableIndexes = []*Index{}
	}
	indexes := klist.New[migrator.Index]()
	for _, index := range this.TableIndexes {
		if index.deleted {
			continue
		}
		indexes.Add(index)
	}
	return indexes
}

func (this *Table) SetName(name string) migrator.Table {
	this.TableName = name
	return this
}

func (this *Table) SetComment(comment string) migrator.Table {
	this.TableComment = comment
	return this
}

func (this *Table) AddColumn(columnName string, ColumnType string) migrator.Column {

	column := &Column{
		ColumnName:         columnName,
		FullColumnDataType: "",
		ColumnDataType:     ColumnType,
		IsPrimaryKey:       false,
		IsAutoIncrement:    false,
		ColumnLength:       0,
		DecimalLength:      0,
		DecimalScale:       0,
		IsNullAble:         false,
		ColumnDefaultValue: nil,
		ColumnComment:      "",
		table:              this,
	}

	switch strings.ToLower(ColumnType) {
	case "varchar", "char", "character", "varying":
		column.SetLength(500)
	case "numeric", "decimal":
		column.SetDecimal(13, 2)
	}

	this.TableColumns = append(this.TableColumns, column)
	return column
}

func (this *Table) DropColumn(columnName ...string) migrator.Table {
	columns := []*Column{}
	for _, name := range columnName {
		for _, column := range this.TableColumns {
			if column.ColumnName == name {
				column.deleted = true
			} else {
				columns = append(columns, column)
			}
		}
	}
	this.TableColumns = columns
	return this
}

func (this *Table) HasColumn(columnName string) bool {
	idx, _ := this.Columns().Find(func(col migrator.Column) bool {
		return col.Name() == columnName
	})
	return idx != -1
}

func (this *Table) AddIndex(columnNames ...string) migrator.Index {
	idx := &Index{
		IndexName:    "",
		IndexColumns: columnNames,
		IsPK:         false,
		table:        this,
		IsUnique:     false,
		origin:       nil,
	}
	idx.generateName()
	this.TableIndexes = append(this.TableIndexes, idx)
	return idx
}

func (this *Table) DropIndex(indexName ...string) migrator.Table {
	//TODO implement me
	panic("implement me")
}

func (this *Table) HasIndex(name string) bool {
	//TODO implement me
	panic("implement me")
}

func (this *Table) Save() error {
	if this.origin == nil {
		return this.create()
	}
	return this.update()
}

func (this *Table) create() error {
	transactions := make([]func(tx *gorm.DB) error, 0)
	transactions = append(transactions, func(tx *gorm.DB) error {
		sql := fmt.Sprintf("CREATE TABLE %s", this.FullName())
		columnsQuery := make([]string, 0)
		for _, column := range this.TableColumns {
			columnsQuery = append(columnsQuery, column.createSQL())
		}

		pkColumns := make([]string, 0)
		for _, index := range this.TableIndexes {
			if index.IsPK {
				pkColumns = append(pkColumns, index.IndexColumns...)
			}
		}
		if len(pkColumns) > 0 {
			columnsQuery = append(columnsQuery, fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(pkColumns, ",")))
		}

		sql += fmt.Sprintf("(%s)", strings.Join(columnsQuery, ","))

		return tx.Exec(sql).Error
	})

	if this.TableComment != "" {
		transactions = append(transactions, func(tx *gorm.DB) error {
			sql := fmt.Sprintf("comment on table %s is '%s'", this.FullName(), this.TableComment)
			return tx.Exec(sql).Error
		})
	}

	for _, column := range this.TableColumns {
		if column.ColumnComment != "" {
			transactions = append(transactions, func(tx *gorm.DB) error {
				sql := fmt.Sprintf("comment on column %s.%s is '%s'", this.FullName(), column.ColumnName, column.ColumnComment)
				return tx.Exec(sql).Error
			})
		}
	}

	for _, index := range this.TableIndexes {
		transactions = append(transactions, func(tx *gorm.DB) error {
			return tx.Exec(index.toSQL()).Error
		})
	}

	return this.mig.tx.Transaction(func(tx *gorm.DB) error {
		for _, transaction := range transactions {
			if err := transaction(tx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (this *Table) update() error {
	transactions := make([]func(tx *gorm.DB) error, 0)
	generateTransaction := func(sql string) func(tx *gorm.DB) error {
		return func(tx *gorm.DB) error {
			return tx.Exec(sql).Error
		}
	}

	for _, column := range this.dropColumns {
		for _, query := range column.alterSqls() {
			transactions = append(transactions, generateTransaction(query))
		}
	}

	for _, column := range this.TableColumns {
		for _, query := range column.alterSqls() {
			transactions = append(transactions, generateTransaction(query))
		}
	}

	if this.TableComment != this.origin.TableComment {
		transactions = append(transactions, func(tx *gorm.DB) error {
			sql := fmt.Sprintf("comment on table %s is '%s'", this.FullName(), this.TableComment)
			return tx.Exec(sql).Error
		})
	}

	if this.TableName != this.origin.TableName {
		transactions = append(transactions, func(tx *gorm.DB) error {
			sql := fmt.Sprintf("alter table %s rename to %s", this.origin.FullName(), this.TableName)
			return tx.Exec(sql).Error
		})
	}

	if len(transactions) > 0 {
		defer func() {
			this.dropColumns = []*Column{}
		}()
		return this.mig.tx.Transaction(func(tx *gorm.DB) error {
			for _, transaction := range transactions {
				if err := transaction(tx); err != nil {
					return err
				}
			}
			return nil
		})
	}
	return nil
}

func (this *Table) clone() *Table {
	columns := make([]*Column, 0, len(this.TableColumns))
	{
		for _, column := range this.TableColumns {
			columns = append(columns, column.clone())
		}
	}

	indexes := make([]*Index, 0, len(this.TableIndexes))
	{
		for _, index := range this.TableIndexes {
			indexes = append(indexes, index.clone())
		}
	}

	dropColumns := make([]*Column, 0)
	{
		for _, column := range this.dropColumns {
			dropColumns = append(dropColumns, column.clone())
		}
	}

	table := &Table{
		TableName:    this.TableName,
		TableComment: this.TableComment,
		TableColumns: columns,
		TableIndexes: indexes,
		mig:          this.mig,
		origin:       nil,
		hasDescribe:  this.hasDescribe,
		dropColumns:  dropColumns,
	}

	return table
}
