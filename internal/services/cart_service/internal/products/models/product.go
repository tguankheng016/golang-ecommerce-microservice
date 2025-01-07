package models

const (
	ProductCollectionName = "products"
)

type Product struct {
	Id          int    `bson:"id"`
	Name        string `bson:"name"`
	Description string `bson:"description"`
	Price       string `bson:"price"`
}
