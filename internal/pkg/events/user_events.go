package events

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
