package models

import (
	"database/sql"
	"time"
)

type User struct {
	Id        int64
	FirstName string
	LastName  string
	UserName  string
	CreatedAt time.Time
	UpdatedAt sql.NullTime `json:"updatedAt"`
}
