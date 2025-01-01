package dtos

import "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/core/domain"

type CategoryDto struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	domain.AuditedEntityDto
}
