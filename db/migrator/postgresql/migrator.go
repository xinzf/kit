package postgresql

import (
	"github.com/elgris/sqrl"
	"github.com/xinzf/kit/container/klist"
	"github.com/xinzf/kit/db/migrator"
	"gorm.io/gorm"
)

func New(tx *gorm.DB, schema string) migrator.Migrator {
	return &_migrator{tx: tx, schema: schema}
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

func (this *_migrator) Tables() (*klist.List[migrator.Table], error) {
	if !this.tableFetched {
		defer func() {
			this.tableFetched = true
		}()

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
			}
			_table.origin = _table.clone()
			this.tables = append(this.tables, _table)
		}
	}

	list := klist.New[migrator.Table]()
	for _, table := range this.tables {
		list.Add(table)
	}
	return list, nil
}

func (this *_migrator) NewTable(tableName string) migrator.Table {
	//TODO implement me
	panic("implement me")
}

func (this *_migrator) DropTable(tableName string) error {
	//TODO implement me
	panic("implement me")
}

func (this *_migrator) HasTable(tableName string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (this *_migrator) RenameTable(oldName, newName string) error {
	//TODO implement me
	panic("implement me")
}
