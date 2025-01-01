package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	httpServer "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/models"
)

func (u categoryManager) CreateCategory(ctx context.Context, category *models.Category) error {
	query := `
		INSERT INTO categories (
			name, 
			normalized_name, 
			created_by,
			is_deleted
		) 
		VALUES (
			@name, 
			@normalized_name, 
			@created_by,
			false
		)
		RETURNING id;
	`
	currentUserId, ok := httpServer.GetCurrentUser(ctx)
	if ok {
		category.CreatedBy.Int64 = currentUserId
		category.CreatedBy.Valid = true
	}

	args := pgx.NamedArgs{
		"name":            category.Name,
		"normalized_name": strings.ToUpper(category.Name),
		"created_by":      category.CreatedBy,
	}

	// Variable to store the returned ID
	var insertedID int

	// Execute the insert query and retrieve the inserted ID
	if err := u.db.QueryRow(ctx, query, args).Scan(&insertedID); err != nil {
		return fmt.Errorf("unable to insert category: %w", err)
	}

	category.Id = insertedID

	return nil
}

func (u categoryManager) UpdateCategory(ctx context.Context, category *models.Category) error {
	query := `
		UPDATE categories
		SET 
			name = @name, 
			normalized_name = @normalized_name, 
			updated_at = @updated_at,
			updated_by = @updated_by
		WHERE 
			id = @id
	`
	currentUserId, ok := httpServer.GetCurrentUser(ctx)
	if ok {
		category.UpdatedBy.Int64 = currentUserId
		category.UpdatedBy.Valid = true
	}

	args := pgx.NamedArgs{
		"name":            category.Name,
		"normalized_name": strings.ToUpper(category.Name),
		"id":              category.Id,
		"updated_at":      time.Now(),
		"updated_by":      category.UpdatedBy,
	}

	if _, err := u.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("unable to update category: %w", err)
	}

	return nil
}

func (u categoryManager) DeleteCategory(ctx context.Context, categoryId int) error {
	query := `
		UPDATE categories
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
		"id":         categoryId,
		"deleted_at": time.Now(),
		"deleted_by": deletedUserId,
	}

	if _, err := u.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("unable to delete category: %w", err)
	}

	return nil
}
