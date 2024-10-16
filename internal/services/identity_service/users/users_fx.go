package users

import (
	userService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/services"
	"go.uber.org/fx"
)

var (
	// Module provided to fx
	Module = fx.Module(
		"userfx",
		userProviders,
	)

	userProviders = fx.Options(
		fx.Provide(
			userService.NewUserRolePermissionManager,
			userService.NewUserPermissionDbManager,
		),
	)
)
