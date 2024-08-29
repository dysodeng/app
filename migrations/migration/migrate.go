package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// 定义数据库迁移
var migration []*gormigrate.Migration

func margeMigration() {
	migration = append(migration, userMigration...)
}

func Migrate(db *gorm.DB, version ...string) (err error) {
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