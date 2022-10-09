package mysql

import (
	"fmt"
	sq "github.com/elgris/sqrl"
	"github.com/xinzf/kit/container/klist"
	"github.com/xinzf/kit/container/kvar"
	"github.com/xinzf/kit/db/migrator"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Table struct {
	mig             *_migrator
	TableName       string    `json:"name"`
	TableComment    string    `json:"comment"`
	TableEngine     string    `json:"engine"`
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

func (this *Table) HasColumn(columnName string) (bool, error) {
	_, err := this.Columns()
	if err != nil {
		return false, err
	}

	for _, column := range this.TableColumns {
		if column.ColumnName == columnName {
			return true, nil
		}
	}
	return false, nil
}

func (this *Table) RenameColumn(oldName, newName string) migrator.Table {
	if oldName == newName {
		return this
	}

	columns, _ := this.Columns()
	if columns == nil || columns.Size() == 0 {
		return this
	}

	idx, col := columns.Find(func(col migrator.Column) bool {
		return col.Name() == oldName
	})
	if idx == -1 {
		return nil
	}

	col.SetName(newName)
	return this
}

func (this *Table) HasIndex(name string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (this *Table) Name() string {
	return this.TableName
}

func (this *Table) FullName() string {
	return this.Name()
}

func (this *Table) Comment() string {
	return this.TableComment
}

func (this *Table) Engine() string {
	return this.TableEngine
}

func (this *Table) CreateTime() time.Time {
	return this.TableCreated
}

func (this *Table) UpdateTime() time.Time {
	return this.TableUpdated
}

func (this *Table) Collation() string {
	return this.TableUnicode
}

func (this *Table) Columns() (*klist.List[migrator.Column], error) {
	if !this.hasFetchColumns {
		this.TableColumns = []*Column{}
		query, args, err := sq.Select(
			"TABLE_SCHEMA",
			"TABLE_NAME",
			"COLUMN_NAME",
			"ORDINAL_POSITION",
			"IS_NULLABLE",
			"DATA_TYPE",
			"CHARACTER_MAXIMUM_LENGTH",
			"NUMERIC_PRECISION",
			"NUMERIC_SCALE",
			"CHARACTER_SET_NAME",
			"COLLATION_NAME",
			"COLUMN_DEFAULT",
			"COLUMN_TYPE",
			"COLUMN_KEY",
			"EXTRA",
			"COLUMN_COMMENT",
		).From("information_schema.columns").
			Where("TABLE_SCHEMA = ? AND TABLE_NAME = ?", this.mig.schema, this.TableName).
			OrderBy("ORDINAL_POSITION").
			ToSql()
		if err != nil {
			return nil, err
		}

		mp := make([]map[string]any, 0)
		err = this.Tx.Raw(query, args...).Find(&mp).Error
		if err != nil {
			return nil, err
		}

		for _, m := range mp {
			columnName := kvar.New(m["COLUMN_NAME"]).String()
			databaseType := kvar.New(m["DATA_TYPE"]).String()
			var length = 0
			if databaseType == "varchar" {
				length = kvar.New(m["CHARACTER_MAXIMUM_LENGTH"]).Int()
			} else {
				length = kvar.New(m["NUMERIC_PRECISION"]).Int()
			}

			expr := kvar.New(m["GENERATION_EXPRESSION"]).String()
			this.TableColumns = append(this.TableColumns, &Column{
				table:               this,
				ColumnName:          columnName,
				ColumnDatabaseType:  databaseType,
				ColumnFullType:      kvar.New(m["COLUMN_TYPE"]).String(),
				IsPrimaryKey:        kvar.New(m["COLUMN_KEY"]).String() == "PRI",
				IsAutoIncrement:     kvar.New(m["EXTRA"]).String() == "auto_increment",
				ColumnLength:        length,
				ColumnDecimal:       kvar.New(m["NUMERIC_SCALE"]).Int(),
				IsNullAble:          kvar.New(m["IS_NULLABLE"]).String() == "YES",
				ColumnDefaultValue:  kvar.New(m["COLUMN_DEFAULT"]).Interface(),
				GeneratedExpression: expr,
				ColumnComment:       kvar.New(m["COLUMN_COMMENT"]).String(),
			})
		}
		this.hasFetchColumns = true
	}

	columns := klist.New[migrator.Column]()
	for _, column := range this.TableColumns {
		columns.Add(column)
	}

	return columns, nil
}

func (this *Table) Indexes() (*klist.List[migrator.Index], error) {
	if !this.hasFetchIndex {
		this.TableIndexes = []*Index{}
		query, args, err := sq.Select(
			"TABLE_SCHEMA",
			"TABLE_NAME",
			"NON_UNIQUE",
			"INDEX_NAME",
			"COLUMN_NAME",
		).From("INFORMATION_SCHEMA.STATISTICS").
			Where("TABLE_SCHEMA = ? AND TABLE_NAME = ?", this.mig.schema, this.TableName).
			OrderBy("SEQ_IN_INDEX").
			ToSql()
		if err != nil {
			return nil, err
		}

		mp := make([]map[string]any, 0)
		err = this.Tx.Raw(query, args...).Find(&mp).Error
		if err != nil {
			return nil, err
		}

		finder := func(idxName string) (idx *Index, found bool) {
			for _, idx := range this.TableIndexes {
				if idx.Name() == idxName {
					return idx, true
				}
			}
			return nil, false
		}

		for _, m := range mp {
			idxName := kvar.New(m["INDEX_NAME"]).String()
			index, found := finder(idxName)
			if !found {
				index = &Index{
					table:          this,
					IndexName:      idxName,
					ContainColumns: []string{},
					IsPrimaryKey:   idxName == "PRIMARY",
					IsUnique:       !kvar.New(m["NON_UNIQUE"]).Bool(),
				}
			}
			index.ContainColumns = append(index.ContainColumns, kvar.New(m["COLUMN_NAME"]).String())

			if !found {
				this.TableIndexes = append(this.TableIndexes, index)
			}
		}
		this.hasFetchIndex = true
	}

	indexs := klist.New[migrator.Index]()
	for _, idx := range this.TableIndexes {
		indexs.Add(idx)
	}
	return indexs, nil
}

func (this *Table) SetName(name string) migrator.Table {
	if this.TableName != name {
		this.newTableName = name
		this.alters.Add(&renameTableName{newName: name})
	}
	return this
}

func (this *Table) SetComment(comment string) migrator.Table {
	if this.TableComment != comment {
		this.alters.Add(&commentTable{comment: comment})
	}
	return this
}

func (this *Table) AddColumn(columnName string, databaseType migrator.DatabaseType) migrator.Column {
	this.alters.RemoveByFilter(func(alt tableAlter) bool {
		return alt.name() == "dropColumn" && alt.(*dropColumn).columnName == columnName
	})

	col := &Column{
		table:              this,
		ColumnName:         columnName,
		ColumnDatabaseType: string(databaseType),
		IsNullAble:         true,
		isNew:              true,
	}
	switch databaseType {
	case migrator.BIGINT:
		col.SetLength(20)
	case migrator.INT:
		col.SetLength(11)
	case migrator.TINYINT:
		col.SetLength(4)
	case migrator.MEDIUMINT:
		col.SetLength(8)
	case migrator.VARCHAR, migrator.CHAR:
		col.SetLength(255)
	}

	idx, _ := this.alters.Find(func(alt tableAlter) bool {
		return alt.name() == "newColumn" && alt.(*newColumn).column.ColumnName == col.ColumnName
	})
	if idx == -1 {
		this.alters.Add(&newColumn{column: col})
	} else {
		this.alters.Set(idx, &newColumn{column: col})
	}

	return col
}

func (this *Table) DropColumn(columnName ...string) migrator.Table {
	for _, name := range columnName {
		this.alters.RemoveByFilter(func(alt tableAlter) bool {
			return (alt.name() == "newColumn" && alt.(*newColumn).column.ColumnName == name) || (alt.name() == "modifyColumn" && alt.(*modifyColumn).column.ColumnName == name)
		})
		this.alters.Add(&dropColumn{name})
	}
	return this
}

func (this *Table) AddIndex(columnNames ...string) migrator.Index {
	index := &Index{
		table:          this,
		ContainColumns: columnNames,
	}
	index.SetUnique(false)
	this.alters.Add(&newIndex{index: index})
	return index
}

func (this *Table) DropIndex(indexName ...string) migrator.Table {
	for _, s := range indexName {
		this.alters.Add(&dropIndex{s})
	}
	return this
}

func (this *Table) Save() (err error) {
	defer func() {
		if err != nil {
			return
		}

		name := this.TableName
		if this.newTableName != "" {
			name = this.newTableName
		}

		var tables []*Table
		tables, err = this.mig.loadTables(name)
		if err != nil {
			return
		}
		this.mig.tables = tables

		var (
			hasFetchColumns = this.hasFetchColumns
			hasFetchIndex   = this.hasFetchIndex
		)
		*this = *tables[0]

		if hasFetchColumns {
			if cols, er := this.Columns(); er != nil {
				err = er
			} else {
				this.TableColumns = klist.Transform(cols, func(t1 migrator.Column) *Column {
					return t1.(*Column)
				}).List()
			}
		}

		if hasFetchIndex {
			if indexes, er := this.Indexes(); er != nil {
				err = er
			} else {
				this.TableIndexes = klist.Transform(indexes, func(t1 migrator.Index) *Index {
					return t1.(*Index)
				}).List()
			}
		}
	}()
	if this.CreateTime().IsZero() {
		err = this.create()
	} else {
		err = this.modify()
	}
	return
}

func (this *Table) create() error {
	if err := this.beforeSave(); err != nil {
		return err
	}
	queries := make([]string, 0)
	this.alters.Each(func(_ int, alt tableAlter) {
		if alt.name() == "newColumn" {
			queries = append(queries, alt.toSQL())
		}
	})
	this.alters.Each(func(_ int, alt tableAlter) {
		if alt.name() == "newIndex" {
			queries = append(queries, alt.toSQL())
		}
	})

	sql := fmt.Sprintf("CREATE table %s (%s)", this.FullName(), strings.Join(queries, ","))
	idx, alt := this.alters.Find(func(alt tableAlter) bool {
		return alt.name() == "tableComment"
	})
	if idx != -1 {
		if comment := alt.toSQL(); comment != "" {
			sql += " " + comment
		}
	}

	if err := this.Tx.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}

func (this *Table) modify() error {
	if this.alters.Size() == 0 {
		return nil
	}

	if err := this.beforeSave(); err != nil {
		return err
	}

	queries := make([]string, 0)
	this.alters.Select(func(_ int, alt tableAlter) bool {
		return alt.name() == "newColumn"
	}).Each(func(_ int, alt tableAlter) {
		queries = append(queries, fmt.Sprintf("add %s", alt.toSQL()))
	})
	this.alters.Select(func(_ int, alt tableAlter) bool {
		return alt.name() != "newColumn" && alt.name() != "newIndex"
	}).Each(func(_ int, alt tableAlter) {
		queries = append(queries, alt.toSQL())
	})
	this.alters.Select(func(_ int, alt tableAlter) bool {
		return alt.name() == "newIndex"
	}).Each(func(_ int, alt tableAlter) {
		queries = append(queries, fmt.Sprintf("add %s", alt.toSQL()))
	})

	sql := fmt.Sprintf("alter table %s ", this.TableName) + strings.Join(queries, ",")
	if err := this.Tx.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}

func (this *Table) beforeSave() error {
	return nil
	//columns:=this.alters.Select(func(_ int,alt tableAlter) bool {
	//    return alt.name()=="newColumn" || alt.name()=="modifyColumn"
	//})
	//if columns.Size() == 0 {
	//    return errors.New("数据表必须包含一个字段")
	//}
	//
	//pkNum := 0
	//mp := make(map[string]bool)
	//err := columns.Iterator(func(_ int, alt tableAlter) error {
	//    var col *Column
	//    switch alt.(type) {
	//    case *newColumn:
	//        col = alt.(*newColumn).column
	//    case *modifyColumn:
	//        col = alt.(*newColumn).column
	//    }
	//    _, found := mp[col.ColumnName]
	//    if found {
	//        return fmt.Errorf("字段：%s 冲突", col.ColumnName)
	//    }
	//    mp[col.ColumnName] = true
	//
	//    if col.IsPrimaryKey {
	//        pkNum++
	//    }
	//    if pkNum > 1 {
	//        return errors.New("最多只能有一个主键字段")
	//    }
	//    return nil
	//})
	//if err != nil {
	//    return err
	//}
	//
	//diffList := make([]*Index, 0)
	//err = this.alters.Select(func(_ int,alt tableAlter) bool {
	//    return alt.name()=="newIndex" || alt.name()=="modifyIndex"
	//}).Iterator(func(_ int, alt tableAlter) error {
	//    var idx *Index
	//    switch alt.(type) {
	//    case *newIndex:
	//        idx = alt.(*newIndex).index
	//    case *modifyIndex:
	//        idx = alt.(*modifyIndex).index
	//    }
	//    for _, index := range diffList {
	//        if index.equal(idx) {
	//            return fmt.Errorf("索引：%s 冲突", idx.IndexName)
	//        }
	//    }
	//    diffList = append(diffList, idx)
	//    return nil
	//})
	//return err
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
