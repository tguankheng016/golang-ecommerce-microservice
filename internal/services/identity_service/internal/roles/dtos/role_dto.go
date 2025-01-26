package dtos

import "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/core/domain"

type RoleDto struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"isDefault"`
	IsStatic  bool   `json:"isStatic"`
	domain.AuditedEntityDto
}
