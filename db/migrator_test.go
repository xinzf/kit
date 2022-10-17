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
		tables []string
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
				tables: []string{"activity"},
			},
		},
		{
			name: POSTGRESQL,
			args: args{
				tx:     postgreTx,
				schema: []string{"repository"},
				tables: []string{"category"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mig := Migrator(tt.args.tx, tt.args.schema...)
			tables, err := mig.Tables(tt.args.tables...)
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

			if err := table.Describe(); err != nil {
				t.Error(err)
				return
			}

			klog.Args("columns", table.Columns().List()).Debug(tt.name)
		})
	}
}

func TestMigrator_Table_Indexes(t *testing.T) {
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
				table:  "ttt",
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

			if err := table.Describe(); err != nil {
				t.Error(err)
				return
			}

			klog.Args("indexes", table.Indexes().List()).Debug(tt.name)
		})
	}
}

func TestMigrator_NewTable(t *testing.T) {
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
				table:  "ttt",
			},
		},
		{
			name: POSTGRESQL,
			args: args{
				tx:     postgreTx,
				schema: []string{"repository"},
				table:  "ttt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mig := Migrator(tt.args.tx, tt.args.schema...)
			table := mig.NewTable(tt.args.table)
			table.SetComment("这是描述")
			table.AddColumn("id", "int").
				SetPrimaryKey().
				SetAutoIncrement(true)

			table.AddColumn("name", "varchar").
				SetNull(false).
				SetLength(30).
				SetComment("姓名")

			table.AddColumn("password", "varchar").
				SetLength(32).
				SetComment("密码")

			table.AddIndex("name", "password").SetUnique(true)

			if err := table.Save(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestMigrator_UpdateTable(t *testing.T) {
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
				table:  "ttt",
			},
		},
		{
			name: POSTGRESQL,
			args: args{
				tx:     postgreTx,
				schema: []string{"repository"},
				table:  "ttt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mig := Migrator(tt.args.tx, tt.args.schema...)
			tables, err := mig.Tables(tt.args.table)
			if err != nil {
				t.Error(err)
				return
			}

			klog.Args("tables", tables.List()).Debug(tt.name)
			table := tables.First()
			if err = table.Describe(); err != nil {
				t.Error(err)
				return
			}

			_, column := table.Columns().Find(func(col migrator.Column) bool {
				return col.Name() == "password"
			})
			column.SetName("password").
				SetNull(false).
				SetLength(30).
				SetDefaultValue("sss")

			name, _ := table.GetColumn("name")
			name.SetName("username")
			table.SetComment("修改了注释1")

			//table.AddColumn("id", "int").
			//    SetPrimaryKey().
			//    SetAutoIncrement(true)

			//table.DropColumn("name")

			table.AddColumn("name", "varchar").
				SetNull(false).
				SetLength(30).
				SetComment("姓名")

			//table.AddColumn("password", "varchar").
			//    SetLength(32).
			//    SetComment("密码")

			//table.AddIndex("name", "password").SetUnique(true)

			if err := table.Save(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestMigrator_HasTable(t *testing.T) {
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
				table:  "ttt",
			},
		},
		{
			name: POSTGRESQL,
			args: args{
				tx:     postgreTx,
				schema: []string{"repository"},
				table:  "tt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mig := Migrator(tt.args.tx, tt.args.schema...)
			has, err := mig.HasTable(tt.args.table)
			if err != nil {
				t.Error(err)
				return
			}
			klog.Args("has", has).Debug(tt.name)
		})
	}
}
