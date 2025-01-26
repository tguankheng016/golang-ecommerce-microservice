package models

import "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/core/domain"

type Category struct {
	Id             int
	Name           string
	NormalizedName string
	domain.FullAuditedEntity
}
