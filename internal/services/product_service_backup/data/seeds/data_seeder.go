package seeds

import (
	"context"
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	categoryModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/categories/models"
	user_service "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/users/grpc_client/protos"
	userModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/users/models"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func DataSeeder(gorm *gorm.DB, userGrpcClientService user_service.UserGrpcServiceClient) error {
	if err := seedCategory(gorm); err != nil {
		return err
	}

	if err := seedUser(gorm, userGrpcClientService); err != nil {
		return err
	}

	return nil
}

func seedCategory(gorm *gorm.DB) error {
	if (gorm.Find(&categoryModel.Category{}).RowsAffected <= 0) {
		categoriesToSeed := []struct {
			name string
		}{
			{
				name: "Clothing",
			},
			{
				name: "Shoes",
			},
			{
				name: "Computers",
			},
			{
				name: "Furniture",
			},
		}

		for _, category := range categoriesToSeed {
			newCategory := &categoryModel.Category{
				Name:      category.name,
				CreatedAt: time.Now(),
			}

			if err := gorm.Create(newCategory).Error; err != nil {
				return errors.Wrap(err, "error in the inserting category into the database.")
			}
		}
	}

	return nil
}

func seedUser(gorm *gorm.DB, userGrpcClientService user_service.UserGrpcServiceClient) error {
	t := time.Now().AddDate(0, 0, -4)
	timestamp := timestamppb.New(t)
	res, err := userGrpcClientService.GetAllUsers(context.Background(), &user_service.GetAllUsersRequest{CreationDate: timestamp})
	if err != nil {
		return err
	}

	if len(res.Users) == 0 {
		return nil
	}

	for _, user := range res.Users {
		var count int64
		if err := gorm.Model(&userModel.User{}).Where("id = ?", user.Id).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			// New User
			var newUser userModel.User
			if err := copier.Copy(&newUser, &user); err != nil {
				return err
			}
			if err := gorm.Create(&newUser).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
