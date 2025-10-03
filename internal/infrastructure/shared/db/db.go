package db

import (
	"log"

	"gorm.io/gorm"

	"github.com/dysodeng/app/internal/infrastructure/config"
)

func Initialize(cfg *config.Config) (*gorm.DB, error) {
	return initMainDB(cfg), nil
}

func DB() *gorm.DB {
	return db
}

func Driver() string {
	return dbDriver
}

func Close() {
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Printf("failed to close database connection: %+v", err)
		return
	}
	log.Println("database connection closed")
}
