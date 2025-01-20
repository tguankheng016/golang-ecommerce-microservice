package models

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/core/domain"
)

type User struct {
	Id                 int64
	FirstName          string
	LastName           string
	UserName           string
	NormalizedUserName string
	Email              string
	NormalizedEmail    string
	PasswordHash       string
	SecurityStamp      uuid.UUID
	ExternalUserId     sql.NullString
	domain.FullAuditedEntity
}

type UserRole struct {
	UserId int64
	RoleId int64
	domain.CreationAuditedEntity
}
