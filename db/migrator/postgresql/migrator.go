package postgresql

import (
	"fmt"
	sq "github.com/elgris/sqrl"
	"github.com/xinzf/kit/container/klist"
	"github.com/xinzf/kit/db/migrator"
	"gorm.io/gorm"
	"time"
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
		TableCreated: time.Time{},
		TableUpdated: time.Time{},
	}
}

func (this *_migrator) DropTable(tableName string) error {
	sql := fmt.Sprintf("drop table %s.%s", this.schema, tableName)
	return this.tx.Exec(sql).Error
}

func (this *_migrator) HasTable(tableName string) (bool, error) {
	_, err := this.Tables()
	if err != nil {
		return false, err
	}

	for _, table := range this.tables {
		if table.Name() == tableName {
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
	type tableStruct struct {
		Schemaname string `json:"schemaname"`
		Tablename  string `json:"tablename"`
		Comment    string `json:"comment"`
	}
	tableRows := make([]*tableStruct, 0)

	_sq := sq.Select("*", "obj_description(relfilenode,'pg_class') as comment").
		From("pg_tables", "pg_class")

	where := "pg_tables.tablename = pg_class.relname and pg_tables.schemaname = ?"
	if tableNames != nil && len(tableNames) > 0 {
		where += " and pg_tables.tablename in ?"
		_sq.Where(where, this.schema, tableNames)
	} else {
		_sq.Where(where, this.schema)
	}

	sql, args, err := _sq.ToSql()
	if err != nil {
		return nil, err
	}

	err = this.tx.Raw(sql, args...).Find(&tableRows).Error
	if err != nil {
		return nil, err
	}

	tables := make([]*Table, 0)
	for _, row := range tableRows {
		tables = append(tables, &Table{
			mig:          this,
			TableName:    row.Tablename,
			TableComment: row.Comment,
			TableCreated: time.Time{},
			TableUpdated: time.Time{},
			TableUnicode: "",
			Tx:           this.tx,
		})
	}
	return tables, nil
}
