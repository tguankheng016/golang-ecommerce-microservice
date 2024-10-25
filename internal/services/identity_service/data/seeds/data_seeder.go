package seeds

import (
	"time"

	"github.com/pkg/errors"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/constants"
	roleModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/models"
	userModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/models"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/services"
	"gorm.io/gorm"
)

func DataSeeder(gorm *gorm.DB) error {
	if err := seedRole(gorm); err != nil {
		return err
	}

	if err := seedUser(gorm); err != nil {
		return err
	}

	return nil
}

func seedRole(gorm *gorm.DB) error {
	if (gorm.Find(&roleModel.Role{}).RowsAffected <= 0) {
		adminRole := &roleModel.Role{
			Name:      constants.DefaultAdminRoleName,
			CreatedAt: time.Now(),
		}

		if err := gorm.Create(adminRole).Error; err != nil {
			return errors.Wrap(err, "error in the inserting role into the database.")
		}
	}

	return nil
}

func seedUser(gorm *gorm.DB) error {
	if (gorm.Find(&userModel.User{}).RowsAffected <= 0) {
		userManager := services.NewUserManager(gorm)

		pass := "123qwe"

		adminUser := &userModel.User{
			FirstName: "admin",
			LastName:  "Tan",
			UserName:  constants.DefaultAdminUsername,
			Email:     "admin@testgk.com",
			CreatedAt: time.Now(),
		}

		if err := userManager.CreateUser(adminUser, pass); err != nil {
			return errors.Wrap(err, "error in the inserting admin user into the database.")
		}

		var adminRole roleModel.Role

		if err := gorm.Where("name = ?", constants.DefaultAdminRoleName).First(&adminRole).Error; err != nil {
			return errors.Wrap(err, "error in the selecting default admin role")
		}

		if err := gorm.Create(&userModel.UserRole{
			UserId: adminUser.Id,
			RoleId: adminRole.Id,
		}).Error; err != nil {
			return errors.Wrap(err, "error in the assigning admin role")
		}

		normalUser := &userModel.User{
			FirstName: "User",
			LastName:  "Tan",
			UserName:  "gkuser123",
			Email:     "user@testgk.com",
			CreatedAt: time.Now(),
		}

		if err := userManager.CreateUser(normalUser, pass); err != nil {
			return errors.Wrap(err, "error in the inserting normal user into the database.")
		}
	}

	return nil
}
