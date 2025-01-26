package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	httpServer "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/models"
)

func (u productManager) CreateProduct(ctx context.Context, product *models.Product) error {
	query := `
		INSERT INTO products (
			name, 
			normalized_name,
			description,
			normalized_description,
			price,
			stock_quantity,
			category_id,
			created_by,
			is_deleted
		) 
		VALUES (
			@name, 
			@normalized_name,
			@description,
			@normalized_description,
			@price,
			@stock_quantity,
			@category_id,
			@created_by,
			false
		)
		RETURNING id;
	`
	currentUserId, ok := httpServer.GetCurrentUser(ctx)
	if ok {
		product.CreatedBy.Int64 = currentUserId
		product.CreatedBy.Valid = true
	}

	args := pgx.NamedArgs{
		"name":                   product.Name,
		"normalized_name":        strings.ToUpper(product.Name),
		"description":            product.Description,
		"normalized_description": strings.ToUpper(product.Description),
		"price":                  product.Price,
		"stock_quantity":         product.StockQuantity,
		"category_id":            product.CategoryId,
		"created_by":             product.CreatedBy,
	}

	// Variable to store the returned ID
	var insertedID int

	// Execute the insert query and retrieve the inserted ID
	if err := u.db.QueryRow(ctx, query, args).Scan(&insertedID); err != nil {
		return fmt.Errorf("unable to insert product: %w", err)
	}

	product.Id = insertedID

	return nil
}

func (u productManager) UpdateProduct(ctx context.Context, product *models.Product) error {
	query := `
		UPDATE products
		SET 
			name = @name, 
			normalized_name = @normalized_name, 
			description = @description, 
			normalized_description = @normalized_description,
			price = @price,
			stock_quantity = @stock_quantity,
			category_id = @category_id,
			updated_at = @updated_at,
			updated_by = @updated_by
		WHERE 
			id = @id
	`
	currentUserId, ok := httpServer.GetCurrentUser(ctx)
	if ok {
		product.UpdatedBy.Int64 = currentUserId
		product.UpdatedBy.Valid = true
	}

	args := pgx.NamedArgs{
		"name":                   product.Name,
		"normalized_name":        strings.ToUpper(product.Name),
		"description":            product.Description,
		"normalized_description": strings.ToUpper(product.Description),
		"price":                  product.Price,
		"stock_quantity":         product.StockQuantity,
		"category_id":            product.CategoryId,
		"id":                     product.Id,
		"updated_at":             time.Now(),
		"updated_by":             product.UpdatedBy,
	}

	if _, err := u.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("unable to update product: %w", err)
	}

	return nil
}

func (u productManager) DeleteProduct(ctx context.Context, productId int) error {
	query := `
		UPDATE products
		SET 
			is_deleted = true,
			deleted_at = @deleted_at,
			deleted_by = @deleted_by
		WHERE 
			id = @id
	`

	deletedUserId := &sql.NullInt64{}

	currentUserId, ok := httpServer.GetCurrentUser(ctx)
	if ok {
		deletedUserId.Int64 = currentUserId
		deletedUserId.Valid = true
	}

	args := pgx.NamedArgs{
		"id":         productId,
		"deleted_at": time.Now(),
		"deleted_by": deletedUserId,
	}

	if _, err := u.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("unable to delete product: %w", err)
	}

	return nil
}
