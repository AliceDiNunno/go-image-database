package postgres

import (
	"fmt"
 	"github.com/AliceDiNunno/go-image-database/src/config"
	migrate "github.com/golang-migrate/migrate/v4"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"

	migrator "github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func StartGormDatabase(config config.GormConfig) *gorm.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DbName)
	db, err := gorm.Open(pg.Open(psqlInfo), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(err)
	}
	db.Logger.LogMode(logger.Info)
	return db
}

func Migrate(db *gorm.DB, migrationsPath string, migrationsTable string) error {
	log.Println("Preparing for migration...")

	sqldb, err := db.DB()

	if err != nil {
		println("Unable to open DB for migration")
		return err
	}

	driver, err := migrator.WithInstance(sqldb, &migrator.Config{
		MigrationsTable: migrationsTable,
	})
	if err != nil {
		println("Failed to create postgres instance")
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver)

	if err != nil {
		println("Failed to prepare migration")
		return err
	}

	log.Println("Migration...")

	err = m.Up()

	if err != nil && err.Error() != "no change" {
		println("An error has happened while migrating")
		return err
	}

	println("Migration has succeeded")
	return nil
}
