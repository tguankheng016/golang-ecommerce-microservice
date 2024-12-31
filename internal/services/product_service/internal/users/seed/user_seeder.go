package seed

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	userGrpc "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/grpc_client/protos"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserSeeder struct {
	db                    *pgxpool.Pool
	userGrpcClientService userGrpc.UserGrpcServiceClient
}

func NewUserSeeder(db *pgxpool.Pool, userGrpcClientService userGrpc.UserGrpcServiceClient) UserSeeder {
	return UserSeeder{
		db:                    db,
		userGrpcClientService: userGrpcClientService,
	}
}

func (u UserSeeder) SeedUsers(ctx context.Context) error {
	t := time.Now().AddDate(0, 0, -4)
	timestamp := timestamppb.New(t)
	res, err := u.userGrpcClientService.GetAllUsers(context.Background(), &userGrpc.GetAllUsersRequest{CreationDate: timestamp})
	if err != nil {
		return err
	}

	if len(res.Users) == 0 {
		return nil
	}

	for _, user := range res.Users {
		if err := u.insertUser(ctx, user); err != nil {
			return err
		}
	}

	return nil
}

func (u UserSeeder) insertUser(ctx context.Context, user *userGrpc.UserModel) error {
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
