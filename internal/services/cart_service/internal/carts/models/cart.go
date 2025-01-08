package models

const (
	CartCollectionName = "carts"
)

type Cart struct {
	Id           string
	UserId       int64
	ProductId    int
	ProductName  string
	ProductDesc  string
	ProductPrice string
	Quantity     int
	IsOutOfStock bool
}
