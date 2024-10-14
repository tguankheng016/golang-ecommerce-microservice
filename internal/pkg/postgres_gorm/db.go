package postgresgorm

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	appConsts "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/constants"
	"github.com/uptrace/bun/driver/pgdriver"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewGormDB creates a new instance of Gorm DB.
//
// It takes the following steps:
//
// 1. Connect to the default database using the given options.
// 2. Check if the target database exists.
// 3. If it does not exist, create it.
// 4. Connect to the target database using the given options.
// 5. Register callbacks for the db.
// 6. Return the Gorm DB instance.
//
// It will return an error if anything goes wrong.
func NewGormDB(options *GormOptions) (*gorm.DB, error) {
	datasource := options.GetDatasource()

	if options.DBName == "" {
		return nil, errors.New("Database name is required in config.json")
	}

	err := createDb(options)

	if err != nil {
		return nil, err
	}

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 10 * time.Second // Maximum time to retry
	maxRetries := 3                      // Number of retries (including the initial attempt)

	var gormDb *gorm.DB

	err = backoff.Retry(func() error {

		gormDb, err = gorm.Open(gorm_postgres.Open(datasource), &gorm.Config{})

		if err != nil {
			return errors.Errorf("failed to connect postgres: %v and connection information: %s", err, datasource)
		}

		return nil

	}, backoff.WithMaxRetries(bo, uint64(maxRetries-1)))

	registerCallBacks(gormDb)

	return gormDb, err
}

// createDb creates a new database if it does not exist.
//
// It takes the following steps:
//
// 1. Connect to the default database using the given options.
// 2. Check if the target database exists.
// 3. If it does not exist, create it.
// 4. Close the database connection.
//
// It will return an error if anything goes wrong.
func createDb(options *GormOptions) error {
	datasource := options.GetPostgresDatasource()

	// Create Db If Not Exist
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(datasource)))

	var exists bool
	err := sqldb.QueryRow(fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_catalog.pg_database WHERE datname='%s')", options.DBName)).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	_, err = sqldb.Exec(fmt.Sprintf("CREATE DATABASE %s", options.DBName))
	if err != nil {
		return err
	}

	defer sqldb.Close()

	return nil
}

// Migrate runs the auto migration for the given types on the given
// *gorm.DB instance. It returns an error if the migration fails.
//
// The types that are passed in are the types that are to be migrated.
//
// For example, if you have a User type and a Product type, you would
// call this function like this:
//
// Migrate(db, &User{}, &Product{})
//
// The function will then run the auto migration for the User and Product
// types on the given database.
func Migrate(gorm *gorm.DB, types ...interface{}) error {
	for _, t := range types {
		err := gorm.AutoMigrate(t)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetTxFromContext retrieves the current transaction from the echo context.
// The transaction is added to the context by the middleware in this package.
// It returns an error if the transaction is not found.
func GetTxFromContext(c echo.Context) (*gorm.DB, error) {
	tx, ok := c.Get(appConsts.DbContextKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("Transaction not found in context")
	}

	return tx, nil
}
