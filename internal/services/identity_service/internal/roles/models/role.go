package models

import "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/core/domain"

type Role struct {
	Id             int64
	Name           string
	NormalizedName string
	IsDefault      bool
	IsStatic       bool
	domain.FullAuditedEntity
}
