package db

import (
	"github.com/xinzf/kit/db/migrator"
	"github.com/xinzf/kit/klog"
	"gorm.io/gorm"
	"testing"
)

var (
	pgTx    *gorm.DB = nil
	mysqlTx *gorm.DB = nil
)

func init() {
	var err error
	mysqlTx, err = New(DbConfig{
		Host:        "127.0.0.1",
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

	pgTx, err = New(DbConfig{
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
	type fields struct {
		tx     *gorm.DB
		schema string
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		//{
		//    name: mysqlTx.Name(),
		//    fields: fields{
		//        tx:     mysqlTx,
		//        schema: "jurun_v4.2.1",
		//    },
		//},
		{
			name: pgTx.Name(),
			fields: fields{
				tx:     pgTx,
				schema: "repository",
			},
		},
	}

	for _, tt := range tests {
		mig := Migrator(tt.fields.tx, tt.fields.schema)

		tables, err := mig.Tables()
		if err != nil {
			t.Error(err)
			return
		}

		klog.Args("Tables", tables.List()).Debug(tt.name)
	}
}

func TestMigrator_DropTable(t *testing.T) {
	type fields struct {
		tx     *gorm.DB
		schema string
		table  string
	}

	test := []struct {
		name   string
		fields fields
	}{
		{
			name: mysqlTx.Name(),
			fields: fields{
				tx:     mysqlTx,
				schema: "jurun_v4.2.1",
				table:  "ttt",
			},
		},
	}

	for _, tt := range test {
		mig := Migrator(tt.fields.tx, tt.fields.schema)
		err := mig.DropTable(tt.fields.table)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestMigrator_NewTable(t *testing.T) {
	type fields struct {
		tx     *gorm.DB
		schema string
		table  string
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: mysqlTx.Name(),
			fields: fields{
				tx:     mysqlTx,
				schema: "jurun_v4.2.1",
				table:  "ttt",
			},
		},
	}

	for _, tt := range tests {
		mig := Migrator(tt.fields.tx, tt.fields.schema)
		table := mig.NewTable(tt.fields.table).SetComment("测试表")
		table.AddColumn("id", migrator.BIGINT).SetAutoIncrement(true).SetPrimaryKey().SetComment("自增ID")
		table.AddColumn("username", migrator.VARCHAR).SetLength(255).SetComment("登录账号")
		table.AddColumn("password", migrator.CHAR).SetLength(32).SetComment("密码")
		table.AddIndex("username", "password")
		table.AddIndex("username").SetUnique(true)
		if err := table.Save(); err != nil {
			t.Error(err)
		} else {
			klog.Args("table", table).Debug(tt.name)
		}
	}
}

func TestMigrator_RenameTable(t *testing.T) {
	type fields struct {
		tx      *gorm.DB
		schema  string
		oldName string
		newName string
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: mysqlTx.Name(),
			fields: fields{
				tx:      mysqlTx,
				schema:  "jurun_v4.2.1",
				oldName: "tt",
				newName: "ttt",
			},
		},
	}

	for _, tt := range tests {
		mig := Migrator(tt.fields.tx, tt.fields.schema)
		err := mig.RenameTable(tt.fields.oldName, tt.fields.newName)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestMigrator_ModifyTable(t *testing.T) {
	type fields struct {
		tx     *gorm.DB
		schema string
		table  string
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: mysqlTx.Name(),
			fields: fields{
				tx:     mysqlTx,
				schema: "jurun_v4.2.1",
				table:  "ttt",
			},
		},
	}

	for _, tt := range tests {
		mig := Migrator(tt.fields.tx, tt.fields.schema)
		tables, err := mig.Tables()
		if err != nil {
			t.Error(err)
			return
		}

		_, table := tables.Find(func(tb migrator.Table) bool {
			return tb.Name() == tt.fields.table
		})

		table.DropIndex("ttt_name_passwd_index", "ttt_username_uindex")
		if err = table.Save(); err != nil {
			t.Error(err)
		}
	}
}

func TestMigrator_Columns(t *testing.T) {
	type fields struct {
		tx     *gorm.DB
		schema string
		table  string
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: pgTx.Name(),
			fields: fields{
				tx:     pgTx,
				schema: "repository",
				table:  "category",
			},
		},
	}

	for _, tt := range tests {
		mig := Migrator(tt.fields.tx, tt.fields.schema)
		tables, err := mig.Tables()
		if err != nil {
			t.Error(err)
			return
		}

		_, table := tables.Find(func(tb migrator.Table) bool {
			return tb.Name() == tt.fields.table
		})

		columns, err := table.Columns()
		if err != nil {
			t.Error(err)
			return
		}
		klog.Args("columns", columns.List()).Debug(tt.name)
	}
}
