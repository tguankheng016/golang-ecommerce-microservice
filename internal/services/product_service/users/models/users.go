package models

type User struct {
	Id        int64  `json:"id" gorm:"primaryKey"`
	FirstName string `json:"firstName" gorm:"type:varchar(64);"`
	LastName  string `json:"lastName" gorm:"type:varchar(64);"`
	UserName  string `json:"userName" gorm:"type:varchar(256);not null"`
}
