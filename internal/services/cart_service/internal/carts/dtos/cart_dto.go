package dtos

import (
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type CartDto struct {
	Id           uuid.UUID
	UserId       int64
	ProductId    int
	ProductName  string
	ProductDesc  string
	ProductPrice decimal.Decimal
	Quantity     int
	IsOutOfStock bool
}
