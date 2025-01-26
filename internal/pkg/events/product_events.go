package events

const (
	ProductCreatedTopicV1        = "product-created-v1"
	ProductUpdatedTopicV1        = "product-updated-v1"
	ProductDeletedTopicV1        = "product-deleted-v1"
	ChangeProductQuantityTopicV1 = "change-product-quantity-v1"
	ProductOutOfStockTopicV1     = "product-out-of-stock-v1"
)

type ProductCreatedEvent struct {
	Id            int
	Name          string
	Description   string
	Price         string
	StockQuantity int
}

type ProductUpdatedEvent struct {
	Id            int
	Name          string
	Description   string
	Price         string
	StockQuantity int
}

type ProductDeletedEvent struct {
	Id int
}

type ChangeProductQuantityEvent struct {
	Id              int
	QuantityChanged int
}

type ProductOutOfStockEvent struct {
	Id            int
	StockQuantity int
}
