package seed

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/models"
	categoryService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/services"
)

type CategorySeeder struct {
	db *pgxpool.Pool
}

func NewCategorySeeder(db *pgxpool.Pool) CategorySeeder {
	return CategorySeeder{
		db: db,
	}
}

func (u CategorySeeder) SeedCategories(ctx context.Context) error {
	tx, err := u.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			// Rollback the transaction in case of error
			tx.Rollback(ctx)
		} else {
			// Commit the transaction if no error occurs
			err = tx.Commit(ctx)
			if err != nil {
				err = fmt.Errorf("unable to commit transaction: %w", err)
			}
		}
	}()

	categoryManager := categoryService.NewCategoryManager(tx)

	count, err := categoryManager.GetCategoriesCount(ctx)
	if err != nil {
		return err
	}

	if count == 0 {
		categoriesToSeed := []struct {
			categoryName string
		}{
			{
				categoryName: "Furniture",
			},
			{
				categoryName: "Clothing",
			},
			{
				categoryName: "Electronics",
			},
			{
				categoryName: "Travel",
			},
			{
				categoryName: "Books",
			},
			{
				categoryName: "Kitchen",
			},
		}

		for _, categoryToSeed := range categoriesToSeed {
			newCategory := models.Category{
				Name: categoryToSeed.categoryName,
			}

			if err := categoryManager.CreateCategory(ctx, &newCategory); err != nil {
				return err
			}
		}
	}

	return nil
}
