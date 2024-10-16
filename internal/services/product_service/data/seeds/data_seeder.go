package seeds

import (
	"time"

	"github.com/pkg/errors"
	categoryModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/categories/models"
	"gorm.io/gorm"
)

func DataSeeder(gorm *gorm.DB) error {
	if err := seedCategory(gorm); err != nil {
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
