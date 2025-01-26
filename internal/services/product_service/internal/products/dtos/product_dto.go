package dtos

import "github.com/shopspring/decimal"

type ProductDto struct {
	Id            int
	Name          string
	Description   string
	Price         decimal.Decimal
	StockQuantity int
	CategoryId    int
	CategoryName  string
}
