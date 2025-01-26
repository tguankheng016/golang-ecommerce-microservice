package seed

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	categoryService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/services"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/models"
	productService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/services"
)

type ProductSeeder struct {
	db *pgxpool.Pool
}

func NewProductSeeder(db *pgxpool.Pool) ProductSeeder {
	return ProductSeeder{
		db: db,
	}
}

func (u ProductSeeder) SeedProducts(ctx context.Context) error {
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

	productManager := productService.NewProductManager(tx)
	categoryManager := categoryService.NewCategoryManager(tx)

	count, err := productManager.GetProductsCount(ctx)
	if err != nil {
		return err
	}

	if count == 0 {
		productsToSeed := []struct {
			productName        string
			productDescription string
			categoryName       string
			price              decimal.Decimal
			stockQuantity      int
		}{
			{
				productName:        "Modern Sofa",
				productDescription: "A stylish and comfortable sofa for your living room",
				categoryName:       "Furniture",
				price:              decimal.NewFromFloat(1200),
				stockQuantity:      5,
			},
			{
				productName:        "Summer Dress",
				productDescription: "A lightweight and stylish dress for the summer season",
				categoryName:       "Clothing",
				price:              decimal.NewFromFloat(75),
				stockQuantity:      10,
			},
			{
				productName:        "Smartphone",
				productDescription: "A high-end smartphone with advanced features and a sleek design",
				categoryName:       "Electronics",
				price:              decimal.NewFromFloat(800),
				stockQuantity:      2,
			},
			{
				productName:        "Backpack",
				productDescription: "A durable and spacious backpack for your travels",
				categoryName:       "Travel",
				price:              decimal.NewFromFloat(50),
				stockQuantity:      15,
			},
			{
				productName:        "Novel",
				productDescription: "A thrilling novel by your favorite author",
				categoryName:       "Books",
				price:              decimal.NewFromFloat(20),
				stockQuantity:      20,
			},
			{
				productName:        "Kitchen Appliances",
				productDescription: "A set of kitchen appliances for your modern kitchen",
				categoryName:       "Kitchen",
				price:              decimal.NewFromFloat(500),
				stockQuantity:      3,
			},
		}

		for _, productToSeed := range productsToSeed {
			newProduct := models.Product{
				Name:          productToSeed.productName,
				Description:   productToSeed.productDescription,
				Price:         productToSeed.price,
				StockQuantity: productToSeed.stockQuantity,
			}

			category, err := categoryManager.GetCategoryByName(ctx, productToSeed.categoryName)
			if err != nil {
				return err
			}
			if category == nil {
				return errors.New("category not found")
			}

			newProduct.CategoryId = category.Id

			if err := productManager.CreateProduct(ctx, &newProduct); err != nil {
				return err
			}
		}
	}

	return nil
}
