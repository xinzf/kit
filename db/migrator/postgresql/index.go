package postgresql

import (
	"fmt"
	"github.com/xinzf/kit/db/migrator"
	"strings"
)

type Index struct {
	IndexName      string   `json:"name"`
	ContainColumns []string `json:"columns"`
	IsPrimaryKey   bool     `json:"primary_key"`
	IsUnique       bool     `json:"unique"`
	table          *Table
	isNew          bool
}

func (this *Index) Table() migrator.Table {
	return this.table
}

func (this *Index) Name() string {
	return this.IndexName
}

func (this *Index) Columns() []string {
	if this.ContainColumns == nil {
		this.ContainColumns = []string{}
	}
	return this.ContainColumns
}

func (this *Index) Unique() bool {
	return this.IsUnique
}

func (this *Index) SetUnique(unique bool) migrator.Index {
	str := strings.Join(this.ContainColumns, "_")
	this.IsUnique = unique
	if !this.IsUnique {
		this.IndexName = fmt.Sprintf("%s_%s_index", this.table.TableName, str)
	} else {
		this.IndexName = fmt.Sprintf("%s_%s_uindex", this.table.TableName, str)
	}
	return this
}

func (this *Index) SetColumns(columns ...string) migrator.Index {
	if columns == nil || len(columns) == 0 {
		return this
	}

	if this.ContainColumns == nil {
		this.ContainColumns = []string{}
	}
	this.ContainColumns = append(this.ContainColumns, columns...)
	return this
}

func (this *Index) equal(idx *Index) bool {
	if this.IsUnique != idx.IsUnique {
		return false
	}

	if len(this.ContainColumns) != len(idx.ContainColumns) {
		return false
	}
	same := len(this.Columns())

	for _, column := range this.ContainColumns {
		for _, containColumn := range idx.ContainColumns {
			if column == containColumn {
				same--
				break
			}
		}
	}

	if same == 0 {
		return true
	}
	return false
}
