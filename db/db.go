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

const (
	MYSQL      string = "mysql"
	POSTGRESQL string = "postgresql"
)

var inited map[string]*gorm.DB = map[string]*gorm.DB{}

func connect(config DbConfig) (client *gorm.DB, err error) {
	var found bool
	if client, found = inited[config.String()]; found {
		sqlDB, err := client.DB()
		if err == nil {
			err = sqlDB.Ping()
			if err == nil {
				return client, nil
			}
		}
		delete(inited, config.String())
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

	if config.Driver == "mysql" {
		client, err = gorm.Open(mysql.Open(config.String()), &gorm.Config{Logger: newLogger})
	} else if config.Driver == "postgresql" {
		client, err = gorm.Open(postgres.Open(config.String()), &gorm.Config{Logger: newLogger})
	}

	if err != nil {
		err = fmt.Errorf("Connect DB failed, err: %s\n", err.Error())
		return
	}

	//client = client.Debug()

	sqlDB, _ := client.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	klog.Args("dsn", fmt.Sprintf("[%s] %s", config.Driver, config.String())).Info("Db connect success!")
	inited[config.String()] = client
	return
}

func New(config DbConfig) (*gorm.DB, error) {
	return connect(config)
}

func DB(name ...string) *gorm.DB {
	connectName := "default"
	if len(name) > 0 {
		connectName = name[0]
	}

	var config DbConfig
	if connectName == "default" {
		config.Host = kcfg.Get[string]("db.host")
		config.User = kcfg.Get[string]("db.user")
		config.Pswd = kcfg.Get[string]("db.pswd")
		config.Name = kcfg.Get[string]("db.name")
		config.Driver = kcfg.Get[string]("db.driver")
	} else {
		config.Host = kcfg.Get[string](fmt.Sprintf("db.%s.host", connectName))
		config.User = kcfg.Get[string](fmt.Sprintf("db.%s.user", connectName))
		config.Pswd = kcfg.Get[string](fmt.Sprintf("db.%s.pswd", connectName))
		config.Name = kcfg.Get[string](fmt.Sprintf("db.%s.name", connectName))
		config.Driver = kcfg.Get[string](fmt.Sprintf("db.%s.driver", connectName))
	}

	db, err := connect(config)
	if err != nil {
		panic(err)
	}
	if db == nil {
		klog.Panic("DB connection has gone, please connect db first")
	}
	return db
}

type DbConfig struct {
	Host   string `json:"host"`
	User   string `json:"user"`
	Pswd   string `json:"pswd"`
	Name   string `json:"name"`
	Driver string `json:"driver"`
	port   int    `json:"-"`
}

func (d DbConfig) String() string {

	if d.Driver == "" {
		d.Driver = MYSQL
	}

	if d.Driver == POSTGRESQL {
		strs := strings.Split(d.Host, ":")
		if len(strs) == 2 {
			d.Host = strs[0]
			d.port, _ = strconv.Atoi(strs[1])
		}
	}

	var dsn string
	if d.Driver == MYSQL {
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&timeout=5s", d.User, d.Pswd, d.Host, d.Name)
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", d.Host, d.User, d.Pswd, d.Name, d.port)
	}

	return dsn
}
