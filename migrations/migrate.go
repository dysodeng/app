package main

import (
	"flag"
	"log"

	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/dysodeng/app/migrations/migration"
)

func main() {

	dbSql := db.DB()

	var migrate = flag.Bool("m", false, "执行迁移 -m")
	var rollback = flag.Bool("r", false, "执行迁移回滚 -r versionID")
	var rollbackLast = flag.Bool("rl", false, "执行最后一次迁移回滚 -rl")

	flag.Parse()

	if *migrate {
		if err := migration.Migrate(dbSql); err != nil {
			log.Fatalf("Could not migrate: %v", err)
		}
		log.Printf("Migration did run successfully")
	}

	if *rollback {
		arg := flag.Args()
		if len(arg) > 0 {
			if err := migration.Rollback(dbSql, arg[0]); err != nil {
				log.Fatalf("Could not rollback: %v", err)
			}
			log.Printf("Rollback to %s migrate successfully", arg)
		} else {
			log.Fatalf("请指定回滚版本号")
		}
	}

	if *rollbackLast {
		if err := migration.Rollback(dbSql); err != nil {
			log.Fatalf("Could not rollback: %v", err)
		}
		log.Printf("Rollback to last migrate successfully")
	}
}
