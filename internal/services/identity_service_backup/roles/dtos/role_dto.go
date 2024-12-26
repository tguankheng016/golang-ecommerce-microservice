package dtos

import "time"

type RoleDto struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	IsStatic  bool      `json:"isStatic"`
	IsDefault bool      `json:"isDefault"`
	CreatedAt time.Time `json:"createdAt"`
} // @name RoleDto
