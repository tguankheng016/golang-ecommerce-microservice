package dtos

import "github.com/shopspring/decimal"

type CreateOrEditProductDto struct {
	Id            *int            `json:"id"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	Price         decimal.Decimal `json:"price"`
	StockQuantity int             `json:"stockQuantity"`
	CategoryId    int             `json:"categoryId"`
}

type CreateProductDto struct {
	CreateOrEditProductDto
}

type EditProductDto struct {
	CreateOrEditProductDto
}
