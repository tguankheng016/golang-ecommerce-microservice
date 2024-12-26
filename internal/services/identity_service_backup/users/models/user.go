package models

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/core/domain"
	roleModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/models"

	"gorm.io/gorm"
)

// User model
type User struct {
	Id                 int64            `json:"id" gorm:"primaryKey"`
	FirstName          string           `json:"firstName" gorm:"type:varchar(64);"`
	LastName           string           `json:"lastName" gorm:"type:varchar(64);"`
	UserName           string           `json:"userName" gorm:"type:varchar(256);not null"`
	NormalizedUserName string           `json:"normalizedUserName" gorm:"type:varchar(256);not null"`
	Email              string           `json:"email" gorm:"type:varchar(256);not null"`
	NormalizedEmail    string           `json:"normalizedEmail" gorm:"type:varchar(256);not null"`
	Password           string           `json:"password" gorm:"not null" copier:"-"`
	SecurityStamp      uuid.UUID        `json:"securityStamp" gorm:"not null"`
	ExternalUserId     uuid.UUID        `json:"externalUserId" gorm:"type:uuid;not null"`
	CreatedAt          time.Time        `json:"createdAt" gorm:"default:current_timestamp"`
	UpdatedAt          sql.NullTime     `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt   `json:"deletedAt"`
	Roles              []roleModel.Role `gorm:"many2many:user_roles;"`
	*domain.FullAuditedEntity
}

type UserRole struct {
	UserId    int64     `json:"userId" gorm:"primaryKey"`
	RoleId    int64     `json:"roleId" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt" gorm:"default:current_timestamp"`
	*domain.CreationAuditedEntity
}
