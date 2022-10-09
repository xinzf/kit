package db

import (
	"fmt"
	"github.com/xinzf/kit/db/migrator"
	"github.com/xinzf/kit/db/migrator/mysql"
	"github.com/xinzf/kit/db/migrator/postgresql"
	"gorm.io/gorm"
)

func Migrator(tx *gorm.DB, schema string) migrator.Migrator {
	if tx.Name() == "mysql" {
		return mysql.Migrator(tx, schema)
	} else if tx.Name() == "postgres" {
		return postgresql.Migrator(tx, schema)
	} else {
		panic(fmt.Sprintf("unsupport driver: %s", tx.Name()))
	}
}
