package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres"
	userGrpc "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/grpc_server/protos"
	userModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/models"
	userService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/services"
)

type UserGrpcServerService struct {
	db postgres.IPgxDbConn
}

func NewUserGrpcServerService(db *pgxpool.Pool) *UserGrpcServerService {
	return &UserGrpcServerService{db: db}
}

func (u *UserGrpcServerService) GetAllUsers(ctx context.Context, req *userGrpc.GetAllUsersRequest) (*userGrpc.GetAllUsersResponse, error) {
	creationDate := req.CreationDate.AsTime()

	query := `SELECT * FROM users WHERE is_deleted = false and created_at >= $1`

	rows, err := u.db.Query(ctx, query, creationDate)
	if err != nil {
		return nil, fmt.Errorf("unable to query users: %w", err)
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[userModel.User])
	if err != nil {
		return nil, err
	}

	var grpcUsers []*userGrpc.UserModel
	if err := copier.Copy(&grpcUsers, &users); err != nil {
		return nil, err
	}

	var result = &userGrpc.GetAllUsersResponse{
		Users: grpcUsers,
	}

	return result, nil
}

func (u *UserGrpcServerService) GetUserById(ctx context.Context, req *userGrpc.GetUserByIdRequest) (*userGrpc.GetUserByIdResponse, error) {
	userManager := userService.NewUserManager(u.db)

	user, err := userManager.GetUserById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	var grpcUser *userGrpc.UserModel
	if err := copier.Copy(&grpcUser, &user); err != nil {
		return nil, err
	}

	var result = &userGrpc.GetUserByIdResponse{
		User: grpcUser,
	}

	return result, nil
}
