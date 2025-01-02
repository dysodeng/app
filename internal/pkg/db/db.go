package db

import (
	"fmt"
	"log"
	"time"

	"github.com/dysodeng/app/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func init() {
	var err error
	var dsn string
	var driver string

	var maxIdleConn, maxOpenConn, connMaxLifetime int
	switch config.Database.Default {
	case "main":
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			config.Database.Main.Username,
			config.Database.Main.Password,
			config.Database.Main.Host,
			config.Database.Main.Port,
			config.Database.Main.Database,
		) + "?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"
		maxIdleConn = config.Database.Main.MaxIdleConns
		maxOpenConn = config.Database.Main.MaxOpenConns
		connMaxLifetime = config.Database.Main.MaxConnLifetime
		driver = config.Database.Main.Driver
		break
	default:
		log.Fatalln("database source not found")
	}

	var dbConnector gorm.Dialector
	switch driver {
	case "mysql":
		dbConnector = mysql.Open(dsn)
	}

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
		log.Fatalf("failed to connect database %+v", err)
	}

	sqlDB, _ := db.DB()

	// 连接池
	sqlDB.SetMaxIdleConns(maxIdleConn)                                     // 连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于该值，超过的连接会被连接池关闭。
	sqlDB.SetMaxOpenConns(maxOpenConn)                                     // 连接池最大连接数
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(connMaxLifetime)) // 连接空闲超时

	log.Println("database connection successful")
}

func Initialize() {}

func DB() *gorm.DB {
	return db
}

func Close() {
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Printf("failed to close database connection: %+v", err)
		return
	}
	log.Println("database connection closed")
}
