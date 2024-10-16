package services

import (
	"context"

	"github.com/jinzhu/copier"
	user_service "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/grpc_server/protos"
	userModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/models"
	"gorm.io/gorm"
)

type UserGrpcService struct {
	db *gorm.DB
}

func NewUserGrpcService(db *gorm.DB) *UserGrpcService {
	return &UserGrpcService{db: db}
}

func (u UserGrpcService) GetAllUsers(ctx context.Context, req *user_service.GetAllUsersRequest) (*user_service.GetAllUsersResponse, error) {
	creationDate := req.CreationDate.AsTime()

	var users []userModel.User
	if err := u.db.Where("created_at >= ?", creationDate).Find(&users).Error; err != nil {
		return nil, err
	}

	var grpcUsers []*user_service.UserModel
	if err := copier.Copy(&grpcUsers, &users); err != nil {
		return nil, err
	}

	var result = &user_service.GetAllUsersResponse{
		Users: grpcUsers,
	}

	return result, nil
}

func (u UserGrpcService) GetUserById(ctx context.Context, req *user_service.GetUserByIdRequest) (*user_service.GetUserByIdResponse, error) {
	var user userModel.User
	if err := u.db.Where("id = ?", req.Id).First(&user).Error; err != nil {
		return nil, err
	}

	var grpcUser *user_service.UserModel
	if err := copier.Copy(&grpcUser, &user); err != nil {
		return nil, err
	}

	var result = &user_service.GetUserByIdResponse{
		User: grpcUser,
	}

	return result, nil
}
