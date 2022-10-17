package postgresql

import (
	"fmt"
	"github.com/xinzf/kit/db/migrator"
	"strings"
)

type Index struct {
	IndexName    string   `json:"name"`
	IndexColumns []string `json:"columns"`
	IsPK         bool     `json:"primaryKey"`
	table        *Table   `json:"-"`
	IsUnique     bool     `json:"unique"`
	origin       *Index
	deleted      bool
}

func (this *Index) Table() migrator.Table {
	return this.table
}

func (this *Index) Name() string {
	return this.IndexName
}

func (this *Index) Columns() []string {
	if this.IndexColumns == nil {
		this.IndexColumns = []string{}
	}
	return this.IndexColumns
}

func (this *Index) Unique() bool {
	return this.IsUnique
}

func (this *Index) PrimaryKey() bool {
	return this.IsPK
}

func (this *Index) SetUnique(unique bool) migrator.Index {
	this.IsUnique = unique
	this.generateName()
	return this
}

func (this *Index) SetColumns(columns ...string) migrator.Index {
	this.IndexColumns = columns
	this.generateName()
	return this
}

func (this *Index) clone() *Index {
	return &Index{
		IndexName:    this.IndexName,
		IndexColumns: this.IndexColumns,
		table:        this.table,
		IsUnique:     this.IsUnique,
	}
}

func (this *Index) generateName() {
	idxTpe := "index"
	if this.IsUnique {
		idxTpe = "uindex"
	}
	this.IndexName = fmt.Sprintf("%s_%s_%s", this.table.Name(), strings.Join(this.IndexColumns, "_"), idxTpe)
}

func (this *Index) toSQL() string {
	idxType := "index"
	if this.IsUnique {
		idxType = "unique index"
	}
	sql := fmt.Sprintf("create %s %s on %s (%s)", idxType, this.Name(), this.table.FullName(), strings.Join(this.IndexColumns, ","))
	return sql
}
