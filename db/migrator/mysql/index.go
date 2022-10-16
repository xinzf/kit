package mysql

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
	//TODO implement me
	panic("implement me")
}

func (this *Index) SetColumns(columns ...string) migrator.Index {
	//TODO implement me
	panic("implement me")
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
	str := strings.Join(this.IndexColumns, "_")
	if !this.IsUnique {
		this.IndexName = fmt.Sprintf("%s_%s_index", this.table.TableName, str)
	} else {
		this.IndexName = fmt.Sprintf("%s_%s_uindex", this.table.TableName, str)
	}
}
