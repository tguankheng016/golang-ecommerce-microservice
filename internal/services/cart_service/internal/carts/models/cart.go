package models

const (
	CartCollectionName = "carts"
)

type Cart struct {
	Id           string `bson:"id"`
	UserId       int64  `bson:"user_id"`
	ProductId    int    `bson:"product_id"`
	ProductName  string `bson:"product_name"`
	ProductDesc  string `bson:"product_desc"`
	ProductPrice string `bson:"product_price"`
	Quantity     int    `bson:"quantity"`
	IsOutOfStock bool   `bson:"is_out_of_stock"`
}
