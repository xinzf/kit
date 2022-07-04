package db

import (
	"fmt"
	"github.com/xinzf/kit/container/kcfg"
	"github.com/xinzf/kit/klog"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var db *gorm.DB
var inited bool

func connect() error {
	if inited {
		return nil
	}

	inited = true
	var (
		host   = kcfg.Get[string]("db.host")
		user   = kcfg.Get[string]("db.user")
		pswd   = kcfg.Get[string]("db.pswd")
		name   = kcfg.Get[string]("db.name")
		driver = kcfg.Get[string]("db.driver")
		err    error
		dsn    string
		port   int
	)

	if driver == "" {
		driver = "mysql"
	}

	if driver == "postgresql" {
		strs := strings.Split(host, ":")
		if len(strs) == 2 {
			host = strs[0]
			port, _ = strconv.Atoi(strs[1])
		}
	}

	if driver == "mysql" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True", user, pswd, host, name)
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", host, user, pswd, name, port)
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,   // 慢 SQL 阈值
			LogLevel:                  logger.Silent, // 日志级别
			IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,         // 禁用彩色打印
		},
	)

	if driver == "mysql" {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
	} else {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
	}

	if err != nil {
		return fmt.Errorf("Connect DB failed, err: %s", err.Error())
	}

	db = db.Debug()
	klog.Args("dsn", fmt.Sprintf("[%s] %s", driver, dsn)).Info("Db connect success!")
	return nil
}

func DB() *gorm.DB {
	if !inited {
		if err := connect(); err != nil {
			panic(any(err))
		}
	}
	if db == nil {
		klog.Panic("DB has not connect, please connect db first")
	}
	return db
}
