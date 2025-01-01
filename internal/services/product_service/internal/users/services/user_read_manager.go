package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/models"
)

type IUserManager interface {
	GetUserById(ctx context.Context, userId int64) (*models.User, error)

	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error

	DeleteUser(ctx context.Context, userId int64) error
}

type userManager struct {
	db postgres.IPgxDbConn
}

func NewUserManager(db postgres.IPgxDbConn) IUserManager {
	return userManager{
		db: db,
	}
}

func (u userManager) GetUserById(ctx context.Context, userId int64) (*models.User, error) {
	query := `SELECT * FROM users where id = @userId LIMIT 1`

	args := pgx.NamedArgs{
		"userId": userId,
	}
	rows, err := u.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("unable to query user by id: %w", err)
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}
