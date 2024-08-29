package db

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/dysodeng/app/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
		panic("database connection not found")
	}

	// db日志
	logFilename := config.LogPath + "/db.log"
	dbLogFile, _ := os.OpenFile(logFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	dbLogger := logger.New(
		log.New(io.MultiWriter(os.Stdout, dbLogFile), "", log.LstdFlags),
		logger.Config{
			SlowThreshold: 200 * time.Millisecond, // 慢查询时间
			LogLevel:      logger.Warn,
			Colorful:      false,
		},
	)

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
		Logger: dbLogger, // db日志
	})
	if err != nil {
		panic("failed to connect database " + err.Error())
	}

	sqlDB, _ := db.DB()

	// 连接池
	sqlDB.SetMaxIdleConns(maxIdleConn)                                     // 连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于该值，超过的连接会被连接池关闭。
	sqlDB.SetMaxOpenConns(maxOpenConn)                                     // 连接池最大连接数
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(connMaxLifetime)) // 连接空闲超时
}

func Initialize() {}

func DB() *gorm.DB {
	return db
}
