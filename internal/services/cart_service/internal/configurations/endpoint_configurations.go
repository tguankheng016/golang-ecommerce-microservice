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
	httpServer "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
	getting_carts "github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/carts/features/getting_carts/v1"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func ConfigureEndpoints(
	httpOptions *httpServer.ServerOptions,
	router *chi.Mux,
	client *mongo.Client,
	database *mongo.Database,
	tokenHandler jwt.IJwtTokenHandler,
	permissionManager permissions.IPermissionManager,
	cacheManager *cache.Cache[string],
	publisher message.Publisher,
) {
	router.Use(middleware.RequestID)
	router.Use(httpServer.SetupLogger())
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(httpOptions.CorsOrigins, ","), // explicitly allow the frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

	basePath := httpOptions.GetBasePath()

	router.Route("/api/v1", func(r chi.Router) {
		config := DefineHumaConfig("API V1", "1.0.0")
		config.Servers = []*huma.Server{
			{URL: basePath + "/api/v1"},
		}
		api := humachi.New(r, config)
		ConfigureAPIMiddlewares(api, client, tokenHandler, permissionManager)
		MapV1Routes(api, client, database, tokenHandler, permissionManager, cacheManager, publisher)
	})

	router.Route("/api/v2", func(r chi.Router) {
		config := DefineHumaConfig("API V2", "2.0.0")
		config.Servers = []*huma.Server{
			{URL: basePath + "/api/v2"},
		}
		api := humachi.New(r, config)
		ConfigureAPIMiddlewares(api, client, tokenHandler, permissionManager)
		MapV2Routes(api, client, database, tokenHandler, permissionManager, cacheManager, publisher)
	})
}

func MapLatestRoutes(
	api huma.API,
	client *mongo.Client,
	database *mongo.Database,
	tokenHandler jwt.IJwtTokenHandler,
	permissionManager permissions.IPermissionManager,
	cacheManager *cache.Cache[string],
	publisher message.Publisher,
) {
	// Carts
	cartCollection := database.Collection("carts")
	getting_carts.MapRoute(api, cartCollection)
}

func MapV1Routes(
	api huma.API,
	client *mongo.Client,
	database *mongo.Database,
	tokenHandler jwt.IJwtTokenHandler,
	permissionManager permissions.IPermissionManager,
	cacheManager *cache.Cache[string],
	publisher message.Publisher,
) {
	MapLatestRoutes(api, client, database, tokenHandler, permissionManager, cacheManager, publisher)
}

func MapV2Routes(
	api huma.API,
	client *mongo.Client,
	database *mongo.Database,
	tokenHandler jwt.IJwtTokenHandler,
	permissionManager permissions.IPermissionManager,
	cacheManager *cache.Cache[string],
	publisher message.Publisher,
) {
	MapLatestRoutes(api, client, database, tokenHandler, permissionManager, cacheManager, publisher)
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

func ConfigureAPIMiddlewares(
	api huma.API,
	client *mongo.Client,
	tokenHandler jwt.IJwtTokenHandler,
	permissionManager permissions.IPermissionManager,
) {
	api.UseMiddleware(jwt.SetupJwtAuthentication(api, tokenHandler))
	api.UseMiddleware(permissions.SetupAuthorization(api, permissionManager))
}
