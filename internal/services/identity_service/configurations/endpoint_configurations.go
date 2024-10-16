package configurations

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
	authenticate "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/identities/features/authenticating/v1/endpoints"
	identityService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/identities/services"
	creating_user "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/features/creating_user/v1/endpoints"
	deleting_user "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/features/deleting_user/v1/endpoints"
	getting_user_by_id "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/features/getting_user_by_id/v1/endpoints"
	getting_user_permissions "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/features/getting_user_permissions/v1/endpoints"
	getting_users "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/features/getting_users/v1/endpoints"
	resetting_user_permissions "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/features/reseting_user_permissions/v1/endpoints"
	updating_user "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/features/updating_user/v1/endpoints"
	updating_user_permissions "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/features/updating_user_permissions/v1/endpoints"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/services"
)

func ConfigEndpoints(
	permissionManager services.IUserRolePermissionManager,
	jwtTokenGenerator identityService.IJwtTokenGenerator,
	jwtTokenValidator jwt.IJwtTokenHandler,
	validator *validator.Validate,
	echo *echo.Echo,
) {
	// Identites
	authenticate.MapRoute(echo, validator, jwtTokenGenerator)

	// Users
	getting_users.MapRoute(echo, validator)
	getting_user_by_id.MapRoute(echo)
	creating_user.MapRoute(echo, validator)
	deleting_user.MapRoute(echo)
	getting_user_permissions.MapRoute(echo, permissionManager)
	resetting_user_permissions.MapRoute(echo, permissionManager)
	updating_user.MapRoute(echo, validator, permissionManager)
	updating_user_permissions.MapRoute(echo, permissionManager)
}
