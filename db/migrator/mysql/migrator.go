package mysql

import (
	"github.com/elgris/sqrl"
	"github.com/xinzf/kit/container/klist"
	"github.com/xinzf/kit/container/kvar"
	"github.com/xinzf/kit/db/migrator"
	"gorm.io/gorm"
)

func New(tx *gorm.DB) migrator.Migrator {
	return &_migrator{tx: tx}
}

type _migrator struct {
	tx           *gorm.DB
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

		sq := sqrl.Select("*").
			From("information_schema.tables").
			Where("TABLE_SCHEMA = ?", this.tx.Migrator().CurrentDatabase())

		query, args, err = sq.ToSql()
		if err != nil {
			return nil, err
		}

		mp := make([]map[string]any, 0)
		err = this.tx.Raw(query, args...).Find(&mp).Error
		if err != nil {
			return nil, err
		}

		for _, t := range mp {
			name := kvar.New(t["TABLE_NAME"]).String()
			_table := &Table{
				mig:          this,
				TableName:    name,
				TableComment: kvar.New(t["TABLE_COMMENT"]).String(),
				TableColumns: []*Column{},
				TableIndexes: []*Index{},
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
