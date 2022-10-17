package postgresql

import (
	"fmt"
	"github.com/elgris/sqrl"
	"github.com/xinzf/kit/container/klist"
	"github.com/xinzf/kit/db/migrator"
	"gorm.io/gorm"
)

func New(tx *gorm.DB, schema string) migrator.Migrator {
	return &_migrator{tx: tx, schema: schema, tables: []*Table{}}
}

type _migrator struct {
	tx           *gorm.DB
	schema       string
	tables       []*Table
	tableFetched bool
}

func (this *_migrator) Schema() string {
	return this.tx.Migrator().CurrentDatabase()
}

func (this *_migrator) Tables(name ...string) (*klist.List[migrator.Table], error) {
	//if !this.tableFetched {
	//    defer func() {
	//        this.tableFetched = true
	//    }()

	this.tables = []*Table{}

	var (
		query string
		args  []any
		err   error
	)

	type tableStruct struct {
		Schemaname string `json:"schemaname"`
		Tablename  string `json:"tablename"`
		Comment    string `json:"comment"`
	}
	tableRows := make([]*tableStruct, 0)

	sq := sqrl.Select("*", "obj_description(relfilenode,'pg_class') as comment").
		From("pg_tables", "pg_class").
		Where("pg_tables.tablename = pg_class.relname AND pg_tables.schemaname = ?", this.schema)
	if len(name) > 0 {
		sq = sq.Where("pg_tables.tablename IN ?", name)
	}

	query, args, err = sq.ToSql()
	if err != nil {
		return nil, err
	}

	err = this.tx.Raw(query, args...).Find(&tableRows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range tableRows {
		_table := &Table{
			TableName:    row.Tablename,
			TableComment: row.Comment,
			TableColumns: []*Column{},
			TableIndexes: []*Index{},
			mig:          this,
			dropColumns:  []*Column{},
		}
		_table.origin = _table.clone()
		this.tables = append(this.tables, _table)
	}
	//}

	list := klist.New[migrator.Table]()
	for _, table := range this.tables {
		list.Add(table)
	}
	return list, nil
}

func (this *_migrator) NewTable(tableName string) migrator.Table {
	table := &Table{
		TableName:    tableName,
		TableColumns: []*Column{},
		TableIndexes: []*Index{},
		mig:          this,
	}
	this.tables = append(this.tables, table)
	return table
}

func (this *_migrator) DropTable(tableName string) error {
	return this.tx.Exec(fmt.Sprintf("drop table if exists %s.%s cascade", this.schema, tableName)).Error
}

func (this *_migrator) HasTable(tableName string) (bool, error) {
	sql := fmt.Sprintf("SELECT count(1) AS TOTAL FROM pg_tables WHERE schemaname = '%s' AND tablename = '%s'", this.schema, tableName)

	type resultStruct struct {
		TOTAL int `json:"TOTAL"`
	}

	var result resultStruct
	err := this.tx.Raw(sql).First(&result).Error
	if err != nil {
		return false, err
	}

	return result.TOTAL > 0, nil
}

func (this *_migrator) RenameTable(oldName, newName string) error {
	sql := fmt.Sprintf("alter table if exists %s.%s rename to %s", this.schema, oldName, newName)
	if err := this.tx.Exec(sql).Error; err != nil {
		return err
	}

	if this.tables != nil {
		for _, table := range this.tables {
			if table.TableName == oldName {
				table.origin.TableName = newName
				table.TableName = newName
			}
		}
	}

	return nil
}
