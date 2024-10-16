package models

import (
	"database/sql"
	"time"

	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/core/domain"
	"gorm.io/gorm"
)

// Category model
type Category struct {
	Id        int64          `json:"id" gorm:"primaryKey"`
	Name      string         `json:"Name" gorm:"type:varchar(64);"`
	CreatedAt time.Time      `json:"createdAt" gorm:"default:current_timestamp"`
	UpdatedAt sql.NullTime   `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
	*domain.FullAuditedEntity
}
