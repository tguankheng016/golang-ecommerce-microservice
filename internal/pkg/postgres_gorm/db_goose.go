package postgresgorm

import (
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

// RunGooseMigration runs the goose migration for the given *gorm.DB.
//
// It sets the goose dialect to "postgres" and then runs the Up migration.
// If there is an error, it is returned.
func RunGooseMigration(db *gorm.DB) error {
	dir := "../../data/migrations"

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	goose.SetBaseFS(nil)

	sqlDb, err := db.DB()
	if err != nil {
		return err
	}

	if err := goose.Up(sqlDb, dir); err != nil {
		return err
	}

	return nil
}
