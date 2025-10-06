package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/dysodeng/app/internal/infrastructure/config"
)

var db *gorm.DB
var dbDriver string

func initMainDB(cfg *config.Config) *gorm.DB {
	var err error
	var dsn string
	var driver string

	var maxIdleConn, maxOpenConn, connMaxLifetime int

	var dbConnector gorm.Dialector
	switch cfg.Database.Driver {
	case "mysql":
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s",
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Database,
		) + "?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"
		dbConnector = mysql.Open(dsn)
	case "postgres":
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			cfg.Database.Host,
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Database,
			cfg.Database.Port,
		)
		dbConnector = postgres.Open(dsn)
	}

	dbDriver = driver

	db, err = gorm.Open(dbConnector, &gorm.Config{
		SkipDefaultTransaction:                   true, // 禁用默认事务
		PrepareStmt:                              true, // 预编译sql
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用创建外键约束
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 禁止表名复数
		},
		Logger: NewGormLogger(), // db日志
	})
	if err != nil {
		log.Fatalf("failed to connect main database %+v", err)
	}

	sqlDB, _ := db.DB()

	// 连接池
	sqlDB.SetMaxIdleConns(maxIdleConn)                                     // 连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于该值，超过的连接会被连接池关闭。
	sqlDB.SetMaxOpenConns(maxOpenConn)                                     // 连接池最大连接数
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(connMaxLifetime)) // 连接空闲超时

	log.Println("main database connection successful")

	return db
}
