package migration

import (
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// 定义数据库迁移
var migration []*gormigrate.Migration

func margeMigration() {
	migration = append(migration, commonMigration...)
	migration = append(migration, userMigration...)
	migration = append(migration, fileMigration...)
}

func Migrate(db *gorm.DB, version ...string) (err error) {
	log.Println("Migration in progress...")
	margeMigration()
	if len(migration) == 0 {
		return
	}
	if len(version) > 0 {
		err = gormigrate.New(db, gormigrate.DefaultOptions, migration).MigrateTo(version[0])
	} else {
		err = gormigrate.New(db, gormigrate.DefaultOptions, migration).Migrate()
	}
	return
}

func Rollback(db *gorm.DB, version ...string) (err error) {
	margeMigration()
	if len(migration) == 0 {
		return
	}
	if len(version) > 0 {
		err = gormigrate.New(db, gormigrate.DefaultOptions, migration).RollbackTo(version[0])
	} else {
		err = gormigrate.New(db, gormigrate.DefaultOptions, migration).RollbackLast()
	}
	return
}
