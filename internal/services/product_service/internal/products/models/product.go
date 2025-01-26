package models

import (
	"github.com/shopspring/decimal"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/core/domain"
	categoryModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/models"
)

type Product struct {
	Id                    int
	Name                  string
	NormalizedName        string
	Description           string
	NormalizedDescription string
	Price                 decimal.Decimal
	StockQuantity         int
	CategoryId            int
	domain.FullAuditedEntity
}

type ProductWithCategory struct {
	Product
	CategoryFK categoryModel.Category
}
