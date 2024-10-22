package postgresgorm

import (
	"os"
	"path/filepath"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

// RunGooseMigration runs the goose migration for the given *gorm.DB.
//
// It sets the goose dialect to "postgres" and then runs the Up migration.
// If there is an error, it is returned.
func RunGooseMigration(db *gorm.DB) error {
	dir, err := getDataMigrationsPath()
	if err != nil {
		return err
	}

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

func getDataMigrationsPath() (string, error) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Traverse up to find the go.mod file
	rootPath := cwd
	for {
		if _, err := os.Stat(filepath.Join(rootPath, "go.mod")); err == nil {
			// Found the go.mod file
			break
		}
		parent := filepath.Dir(rootPath)
		if parent == rootPath {
			return "", err
		}
		rootPath = parent
	}

	// Get the path to the "data/migrations" folder within the project directory
	migrationsPath := filepath.Join(rootPath, "data/migrations")

	return migrationsPath, nil
}
