package mysql

import (
	"fmt"
	"github.com/xinzf/kit/container/klist"
	"github.com/xinzf/kit/container/kvar"
	"github.com/xinzf/kit/db/migrator"
	"gorm.io/gorm"
	"strings"
)

func Migrator(tx *gorm.DB, schema string) migrator.Migrator {
	return &_migrator{tx: tx, schema: schema, tables: []*Table{}}
}

type _migrator struct {
	tx            *gorm.DB
	schema        string
	tablesFetched bool
	tables        []*Table
}

func (this *_migrator) Schema() string {
	return this.schema
}

func (this *_migrator) Tables() (*klist.List[migrator.Table], error) {
	var err error
	if !this.tablesFetched {
		defer func() {
			this.tablesFetched = true
		}()

		this.tables, err = this.loadTables()
		if err != nil {
			return nil, err
		}
	}

	list := klist.New[migrator.Table]()
	for _, table := range this.tables {
		list.Add(table)
	}
	return list, nil
}

func (this *_migrator) NewTable(tableName string) migrator.Table {
	return &Table{
		mig:          this,
		TableName:    tableName,
		TableEngine:  "InnoDB",
		TableUnicode: "utf8mb4_0900_ai_ci",
		TableColumns: []*Column{},
		TableIndexes: []*Index{},
		Tx:           this.tx,
		alters:       klist.New[tableAlter](),
		isNew:        true,
	}
}

func (this *_migrator) DropTable(tableName string) error {
	sql := fmt.Sprintf("drop table if exists `%s`.%s cascade", this.schema, tableName)
	return this.tx.Exec(sql).Error
}

func (this *_migrator) HasTable(tableName string) (bool, error) {
	_, err := this.Tables()
	if err != nil {
		return false, err
	}
	for _, table := range this.tables {
		if table.TableName == tableName {
			return true, nil
		}
	}
	return false, nil
}

func (this *_migrator) RenameTable(oldName, newName string) error {
	if oldName == newName {
		return nil
	}
	_, err := this.Tables()
	if err != nil {
		return err
	}

	var table *Table
	for _, tb := range this.tables {
		if tb.TableName == oldName {
			table = tb
			break
		}
	}
	if table == nil {
		return fmt.Errorf("Not found table: %s\n", oldName)
	}

	return table.SetName(newName).Save()
}

func (this *_migrator) loadTables(tableNames ...string) ([]*Table, error) {
	query := fmt.Sprintf("SELECT * FROM information_schema.tables WHERE TABLE_SCHEMA = '%s'", this.schema)
	if tableNames != nil && len(tableNames) > 0 {
		names := ""
		for _, name := range tableNames {
			names += fmt.Sprintf("'%s',", name)
		}
		names = strings.TrimRight(names, ",")
		query += fmt.Sprintf(" AND TABLE_NAME in (%s)", names)
	}

	mp := make([]map[string]any, 0)
	err := this.tx.Raw(query).Find(&mp).Error
	if err != nil {
		return nil, err
	}

	var list = make([]*Table, 0)
	for _, t := range mp {
		name := kvar.New(t["TABLE_NAME"]).String()
		_table := &Table{
			mig:          this,
			TableName:    name,
			TableComment: kvar.New(t["TABLE_COMMENT"]).String(),
			TableEngine:  kvar.New(t["ENGINE"]).String(),
			TableCreated: kvar.New(t["CREATE_TIME"]).Time(),
			TableUpdated: kvar.New(t["UPDATE_TIME"]).Time(),
			TableUnicode: kvar.New(t["TABLE_COLLATION"]).String(),
			TableColumns: []*Column{},
			TableIndexes: []*Index{},
			Tx:           this.tx,
			alters:       klist.New[tableAlter](),
		}
		list = append(list, _table)
	}

	return list, nil
}
