package domain

import (
	"database/sql"
	"time"
)

type CreationAuditedEntity struct {
	CreatedAt time.Time     `json:"createdAt"`
	CreatedBy sql.NullInt64 `json:"createdBy"`
}

type UpdateAuditedEntity struct {
	UpdatedAt sql.NullTime  `json:"updatedAt"`
	UpdatedBy sql.NullInt64 `json:"updatedBy"`
}

type DeleteAuditedEntity struct {
	IsDeleted bool          `json:"isDeleted"`
	DeletedAt sql.NullTime  `json:"deletedAt"`
	DeletedBy sql.NullInt64 `json:"deletedBy"`
}

type FullAuditedEntity struct {
	CreationAuditedEntity
	UpdateAuditedEntity
	DeleteAuditedEntity
}
