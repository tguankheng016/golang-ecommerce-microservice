package seed

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	userGrpc "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/grpc_client/protos"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/models"
	userService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/services"
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

	userManager := userService.NewUserManager(u.db)

	for _, user := range res.Users {
		var newUser models.User
		if err := copier.Copy(&newUser, &user); err != nil {
			return err
		}

		if err := userManager.CreateUser(ctx, &newUser); err != nil {
			return err
		}
	}

	return nil
}
