package configurations

import (
	"strings"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	httpServer "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/openiddict"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
	authenticating "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/identities/features/authenticating/v1"
	authenticating_v2 "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/identities/features/authenticating/v2"
	getting_all_permissions "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/identities/features/getting_all_permissions/v1"
	getting_current_session "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/identities/features/getting_current_session/v1"
	oauth_authenticating "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/identities/features/oauth_authenticating/v1"
	refreshing_token "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/identities/features/refreshing_token/v1"
	sign_out "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/identities/features/signing_out/v1"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/identities/services"
	creating_role "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/roles/features/creating_role/v1"
	deleting_role "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/roles/features/deleting_role/v1"
	getting_role_by_id "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/roles/features/getting_role_by_id/v1"
	getting_roles "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/roles/features/getting_roles/v1"
	updating_role "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/roles/features/updating_roles/v1"
	creating_user "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/features/creating_user/v1"
	deleting_user "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/features/deleting_user/v1"
	getting_user_by_id "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/features/getting_user_by_id/v1"
	getting_user_permissions "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/features/getting_user_permissions/v1"
	getting_users "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/features/getting_users/v1"
	resetting_user_permissions "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/features/resetting_user_permissions/v1"
	updating_user "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/features/updating_user/v1"
	updating_user_permissions "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/features/updating_user_permissions/v1"
	userService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/services"
)

func ConfigureEndpoints(
	httpOptions *httpServer.ServerOptions,
	router *chi.Mux,
	pool *pgxpool.Pool,
	jwtTokenGenerator services.IJwtTokenGenerator,
	tokenHandler jwt.IJwtTokenHandler,
	permissionManager permissions.IPermissionManager,
	userRolePermissionManager userService.IUserRolePermissionManager,
	cacheManager *cache.Cache[string],
	publisher message.Publisher,
	oAuthApiClient openiddict.IOAuthApiClient,
) {
	router.Use(middleware.RequestID)
	router.Use(httpServer.SetupLogger())
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(httpOptions.CorsOrigins, ","), // explicitly allow the frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))
	router.Use(httpServer.SetupTracing(httpOptions.Name, router)...)

	basePath := httpOptions.GetBasePath()

	router.Route("/api/v1", func(r chi.Router) {
		config := DefineHumaConfig("API V1", "1.0.0")
		config.Servers = []*huma.Server{
			{URL: basePath + "/api/v1"},
		}
		api := humachi.New(r, config)
		ConfigureAPIMiddlewares(api, pool, tokenHandler, permissionManager)
		MapV1Routes(api, pool, jwtTokenGenerator, tokenHandler, permissionManager, userRolePermissionManager, cacheManager, publisher, oAuthApiClient)
	})

	router.Route("/api/v2", func(r chi.Router) {
		config := DefineHumaConfig("API V2", "2.0.0")
		config.Servers = []*huma.Server{
			{URL: basePath + "/api/v2"},
		}
		api := humachi.New(r, config)
		ConfigureAPIMiddlewares(api, pool, tokenHandler, permissionManager)
		MapV2Routes(api, pool, jwtTokenGenerator, tokenHandler, permissionManager, userRolePermissionManager, cacheManager, publisher, oAuthApiClient)
	})
}

func DefineHumaConfig(title string, version string) huma.Config {
	config := huma.DefaultConfig(title, version)
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearer": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}

	return config
}

func MapLatestRoutes(
	api huma.API,
	pool *pgxpool.Pool,
	jwtTokenGenerator services.IJwtTokenGenerator,
	tokenHandler jwt.IJwtTokenHandler,
	permissionManager permissions.IPermissionManager,
	userRolePermissionManager userService.IUserRolePermissionManager,
	cacheManager *cache.Cache[string],
	publisher message.Publisher,
	oAuthApiClient openiddict.IOAuthApiClient,
) {
	getting_users.MapRoute(api, pool)
	getting_user_by_id.MapRoute(api, pool)
	creating_user.MapRoute(api, pool, publisher)
	updating_user.MapRoute(api, pool, userRolePermissionManager, publisher)
	deleting_user.MapRoute(api, pool, cacheManager)
	getting_user_permissions.MapRoute(api, pool, userRolePermissionManager)
	resetting_user_permissions.MapRoute(api, pool, userRolePermissionManager)
	updating_user_permissions.MapRoute(api, pool, userRolePermissionManager)

	getting_roles.MapRoute(api, pool)
	getting_role_by_id.MapRoute(api, pool, userRolePermissionManager)
	creating_role.MapRoute(api, pool)
	updating_role.MapRoute(api, pool, userRolePermissionManager)
	deleting_role.MapRoute(api, pool, userRolePermissionManager)

	sign_out.MapRoute(api, jwtTokenGenerator)
	refreshing_token.MapRoute(api, pool, jwtTokenGenerator, tokenHandler)
	getting_all_permissions.MapRoute(api)
	getting_current_session.MapRoute(api, pool, userRolePermissionManager)
	oauth_authenticating.MapRoute(api, pool, jwtTokenGenerator, oAuthApiClient, publisher)
}

func MapV1Routes(
	api huma.API,
	pool *pgxpool.Pool,
	jwtTokenGenerator services.IJwtTokenGenerator,
	tokenHandler jwt.IJwtTokenHandler,
	permissionManager permissions.IPermissionManager,
	userRolePermissionManager userService.IUserRolePermissionManager,
	cacheManager *cache.Cache[string],
	publisher message.Publisher,
	oAuthApiClient openiddict.IOAuthApiClient,
) {
	MapLatestRoutes(api, pool, jwtTokenGenerator, tokenHandler, permissionManager, userRolePermissionManager, cacheManager, publisher, oAuthApiClient)
	authenticating.MapRoute(api, jwtTokenGenerator)
}

func MapV2Routes(
	api huma.API,
	pool *pgxpool.Pool,
	jwtTokenGenerator services.IJwtTokenGenerator,
	tokenHandler jwt.IJwtTokenHandler,
	permissionManager permissions.IPermissionManager,
	userRolePermissionManager userService.IUserRolePermissionManager,
	cacheManager *cache.Cache[string],
	publisher message.Publisher,
	oAuthApiClient openiddict.IOAuthApiClient,
) {
	MapLatestRoutes(api, pool, jwtTokenGenerator, tokenHandler, permissionManager, userRolePermissionManager, cacheManager, publisher, oAuthApiClient)
	authenticating_v2.MapRoute(api, pool, jwtTokenGenerator)
}

func ConfigureAPIMiddlewares(
	api huma.API,
	pool *pgxpool.Pool,
	tokenHandler jwt.IJwtTokenHandler,
	permissionManager permissions.IPermissionManager,
) {
	api.UseMiddleware(jwt.SetupJwtAuthentication(api, tokenHandler))
	api.UseMiddleware(permissions.SetupAuthorization(api, permissionManager))
}
