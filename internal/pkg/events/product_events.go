package events

const (
	ProductCreatedTopicV1 = "product-created-v1"
	ProductUpdatedTopicV1 = "product-updated-v1"
	ProductDeletedTopicV1 = "product-deleted-v1"
)

type ProductCreatedEvent struct {
	Id          int
	Name        string
	Description string
	Price       string
}

type ProductUpdatedEvent struct {
	Id          int
	Name        string
	Description string
	Price       string
}

type ProductDeletedEvent struct {
	Id int
}
