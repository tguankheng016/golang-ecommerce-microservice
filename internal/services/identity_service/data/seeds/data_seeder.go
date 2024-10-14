package seeds

import (
	"time"

	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/constants"
	rolemodel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/models"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/managers"
	usermodel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/models"

	"github.com/pkg/errors"

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
	if (gorm.Find(&rolemodel.Role{}).RowsAffected <= 0) {
		adminRole := &rolemodel.Role{
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
	if (gorm.Find(&usermodel.User{}).RowsAffected <= 0) {
		userManager := managers.NewUserManager(gorm)

		pass := "123qwe"

		adminUser := &usermodel.User{
			FirstName: "admin",
			LastName:  "Tan",
			UserName:  constants.DefaultAdminUsername,
			Email:     "admin@testgk.com",
			CreatedAt: time.Now(),
		}

		if err := userManager.CreateUser(adminUser, pass); err != nil {
			return errors.Wrap(err, "error in the inserting admin user into the database.")
		}

		var adminRole rolemodel.Role

		if err := gorm.Where("name = ?", constants.DefaultAdminRoleName).First(&adminRole).Error; err != nil {
			return errors.Wrap(err, "error in the selecting default admin role")
		}

		if err := gorm.Model(&adminUser).Association("Roles").Append(&adminRole); err != nil {
			return errors.Wrap(err, "error in the assigning admin role")
		}

		normalUser := &usermodel.User{
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
