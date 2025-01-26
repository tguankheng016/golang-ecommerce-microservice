package events

const (
	UserCreatedTopicV1 = "user-created-v1"
	UserUpdatedTopicV1 = "user-updated-v1"
	UserDeletedTopicV1 = "user-deleted-v1"
)

type UserCreatedEvent struct {
	Id        int64
	UserName  string
	FirstName string
	LastName  string
}

type UserUpdatedEvent struct {
	Id        int64
	UserName  string
	FirstName string
	LastName  string
}

type UserDeletedEvent struct {
	Id int64
}
