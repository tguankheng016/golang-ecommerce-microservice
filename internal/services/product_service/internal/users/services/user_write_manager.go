package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/models"
)

func (u userManager) CreateUser(ctx context.Context, user *models.User) error {
	var count int

	query := "SELECT COUNT(*) FROM users where id = $1"

	if err := u.db.QueryRow(ctx, query, user.Id).Scan(&count); err != nil {
		return fmt.Errorf("unable to count row: %w", err)
	}

	if count > 0 {
		return nil
	}

	query = `
		INSERT INTO users (
			id,
			first_name, 
			last_name, 
			user_name
		) 
		VALUES (
			@id,
			@first_name, 
			@last_name, 
			@user_name
		)
	`

	args := pgx.NamedArgs{
		"id":         user.Id,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"user_name":  user.UserName,
	}

	if _, err := u.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("unable to insert user: %w", err)
	}

	return nil
}

func (u userManager) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET 
			first_name = @first_name, 
			last_name = @last_name, 
			user_name = @user_name, 
			updated_at = @updated_at
		WHERE 
			id = @id
	`

	args := pgx.NamedArgs{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"user_name":  user.UserName,
		"id":         user.Id,
		"updated_at": time.Now(),
	}

	if _, err := u.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("unable to update user: %w", err)
	}

	return nil
}

func (u userManager) DeleteUser(ctx context.Context, userId int64) error {
	query := "DELETE FROM users WHERE id = $1"

	if _, err := u.db.Exec(ctx, query, userId); err != nil {
		return fmt.Errorf("unable to delete user: %w", err)
	}

	return nil
}
