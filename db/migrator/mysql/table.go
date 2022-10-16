package mysql

import (
	"fmt"
	"github.com/xinzf/kit/container/klist"
	"github.com/xinzf/kit/db/migrator"
)

type Table struct {
	TableName    string     `json:"name"`
	TableComment string     `json:"comment"`
	TableColumns []*Column  `json:"columns"`
	TableIndexes []*Index   `json:"indexes"`
	mig          *_migrator `json:"-"`
	origin       *Table     `json:"-"`
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
	return fmt.Sprintf("%s.%s", this.mig.tx.Migrator().CurrentDatabase(), this.TableName)
}

func (this *Table) Comment() string {
	return this.TableComment
}

func (this *Table) loadColumns() error {
	this.TableColumns = []*Column{}
	cols, err := this.mig.tx.Migrator().ColumnTypes(this.TableName)
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
			column.Decimal.Decimal = int(decimalDec)
			column.Decimal.Length = int(decimalLen)
		}
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

func (this *Table) loadIndexes() error {
	this.TableIndexes = []*Index{}

	idxes, err := this.mig.tx.Migrator().GetIndexes(this.TableName)
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
	col := &Column{
		ColumnName:         columnName,
		FullColumnDataType: ColumnType,
		ColumnDataType:     ColumnType,
		IsPrimaryKey:       false,
		IsAutoIncrement:    false,
		ColumnLength:       0,
		Decimal: struct {
			Length  int `json:"length"`
			Decimal int `json:"decimal"`
		}{},
		IsNullAble:         false,
		ColumnDefaultValue: nil,
		ColumnComment:      "",
		table:              this,
		origin:             nil,
	}

	this.TableColumns = append(this.TableColumns, col)
	return col
}

func (this *Table) DropColumn(columnName ...string) migrator.Table {
	this.TableColumns = []*Column{}
	this.Columns().Select(func(_ int, col migrator.Column) bool {
		var found bool
		for _, name := range columnName {
			if name == col.Name() {
				found = true
				break
			}
		}

		return !found
	}).Each(func(_ int, col migrator.Column) {
		this.TableColumns = append(this.TableColumns, col.(*Column))
	})
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
	this.TableIndexes = []*Index{}
	this.Indexes().Select(func(_ int, idx migrator.Index) bool {
		var found bool
		for _, name := range indexName {
			if idx.Name() == name {
				found = true
				break
			}
		}
		return !found
	}).Each(func(_ int, idx migrator.Index) {
		this.TableIndexes = append(this.TableIndexes, idx.(*Index))
	})
	return this
}

func (this *Table) HasIndex(name string) bool {
	idx, _ := this.Indexes().Find(func(idx migrator.Index) bool {
		return idx.Name() == name
	})
	return idx != -1
}

func (this *Table) Save() error {
	if this.origin == nil {
		// create
	} else {
		// update
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
	table := &Table{
		TableName:    this.TableName,
		TableComment: this.TableComment,
		TableColumns: columns,
		TableIndexes: indexes,
		mig:          this.mig,
		origin:       nil,
		hasDescribe:  this.hasDescribe,
	}

	return table
}
