package db

import (
	"github.com/xinzf/kit/db/migrator"
	"github.com/xinzf/kit/klog"
	"gorm.io/gorm"
	"testing"
)

var (
	mysqlTX   *gorm.DB
	postgreTx *gorm.DB
)

func init() {
	var err error
	mysqlTX, err = New(DbConfig{
		Host:        "127.0.0.1:3306",
		User:        "root",
		Pswd:        "111111",
		Name:        "jurun_v4.2.1",
		Driver:      MYSQL,
		MaxIdleCons: 10,
		MaxOpenCons: 10,
		Debug:       true,
	})
	if err != nil {
		panic(err)
	}

	postgreTx, err = New(DbConfig{
		Host:        "127.0.0.1:5432",
		User:        "postgres",
		Pswd:        "111111",
		Name:        "bi",
		Driver:      POSTGRESQL,
		MaxIdleCons: 10,
		MaxOpenCons: 10,
		Debug:       true,
	})
	if err != nil {
		panic(err)
	}
}

func TestMigrator_Tables(t *testing.T) {
	type args struct {
		tx     *gorm.DB
		schema []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: MYSQL,
			args: args{
				tx:     mysqlTX,
				schema: []string{},
			},
		},
		{
			name: POSTGRESQL,
			args: args{
				tx:     postgreTx,
				schema: []string{"repository"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mig := Migrator(tt.args.tx, tt.args.schema...)
			tables, err := mig.Tables()
			if err != nil {
				t.Error(err)
				return
			}

			klog.Args("tables", tables.List()).Debug(tt.name)
		})
	}
}

func TestMigrator_Table_Columns(t *testing.T) {
	type args struct {
		tx     *gorm.DB
		schema []string
		table  string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: MYSQL,
			args: args{
				tx:     mysqlTX,
				schema: []string{},
				table:  "activity",
			},
		},
		{
			name: POSTGRESQL,
			args: args{
				tx:     postgreTx,
				schema: []string{"repository"},
				table:  "category",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mig := Migrator(tt.args.tx, tt.args.schema...)
			tables, err := mig.Tables()
			if err != nil {
				t.Error(err)
				return
			}

			_, table := tables.Find(func(tb migrator.Table) bool {
				return tb.Name() == tt.args.table
			})
			columns, err := table.Columns()
			if err != nil {
				t.Error(err)
				return
			}

			klog.Args("columns", columns.List()).Debug(tt.name)
		})
	}
}
