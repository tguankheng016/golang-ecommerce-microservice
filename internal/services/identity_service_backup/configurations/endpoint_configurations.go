package configurations

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/openiddict"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/rabbitmq"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
	authenticate "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/identities/features/authenticating/v1/endpoints"
	getting_all_permissions "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/identities/features/getting_all_permissions/v1/endpoints"
	getting_current_session "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/identities/features/getting_current_session/v1/endpoints"
	oauth_authenticating "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/identities/features/oauth_authenticating/v1/endpoints"
	refreshToken "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/identities/features/refreshing_token/v1/endpoints"
	sign_out "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/identities/features/signing_out/v1/endpoints"
	identityService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/identities/services"
	creating_role "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/features/creating_role/v1/endpoints"
	deleting_role "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/features/deleting_role/v1/endpoints"
	getting_role_by_id "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/features/getting_role_by_id/v1/endpoints"
	getting_roles "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/features/getting_roles/v1/endpoints"
	updating_role "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/features/updating_role/v1/endpoints"
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
	rabbitMQPublisher rabbitmq.IPublisher,
	oAuthApiClient openiddict.IOAuthApiClient,
) {
	// Identities
	authenticate.MapRoute(echo, validator, jwtTokenGenerator)
	getting_current_session.MapRoute(echo, permissionManager)
	getting_all_permissions.MapRoute(echo)
	refreshToken.MapRoute(echo, validator, jwtTokenGenerator, jwtTokenValidator)
	sign_out.MapRoute(echo, jwtTokenGenerator)
	oauth_authenticating.MapRoute(echo, validator, jwtTokenGenerator, oAuthApiClient, rabbitMQPublisher)

	// Users
	getting_users.MapRoute(echo, validator)
	getting_user_by_id.MapRoute(echo)
	creating_user.MapRoute(echo, validator, rabbitMQPublisher)
	deleting_user.MapRoute(echo, rabbitMQPublisher)
	getting_user_permissions.MapRoute(echo, permissionManager)
	resetting_user_permissions.MapRoute(echo, permissionManager)
	updating_user.MapRoute(echo, validator, permissionManager, rabbitMQPublisher)
	updating_user_permissions.MapRoute(echo, permissionManager)

	// Roles
	getting_roles.MapRoute(echo, validator)
	getting_role_by_id.MapRoute(echo, permissionManager)
	creating_role.MapRoute(echo, validator)
	updating_role.MapRoute(echo, validator, permissionManager)
	deleting_role.MapRoute(echo, permissionManager)
}
