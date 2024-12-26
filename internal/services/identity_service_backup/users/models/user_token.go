package models

import "time"

type UserToken struct {
	Id             int64     `gorm:"primarykey"`
	UserId         int64     `gorm:"column:user_id;index"`
	TokenKey       string    `gorm:"column:token_key;type:varchar(64)"`
	ExpirationTime time.Time `gorm:"column:expiration_time"`
}
