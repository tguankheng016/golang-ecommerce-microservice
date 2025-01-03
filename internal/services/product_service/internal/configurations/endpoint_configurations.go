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
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
	creating_category "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/features/creating_category/v1"
	deleting_category "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/features/deleting_category/v1"
	getting_categories "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/features/getting_categories/v1"
	getting_category_by_id "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/features/getting_category_by_id/v1"
	updating_category "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/features/updating_category/v1"
	creating_product "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/features/creating_product/v1"
	deleting_product "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/features/deleting_product/v1"
	getting_product_by_id "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/features/getting_product_by_id/v1"
	getting_products "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/features/getting_products/v1"
	updating_product "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/features/updating_product/v1"
)

func ConfigureEndpoints(
	httpOptions *httpServer.ServerOptions,
	router *chi.Mux,
	pool *pgxpool.Pool,
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
		ConfigureAPIMiddlewares(api, pool, tokenHandler, permissionManager)
		MapV1Routes(api, pool, tokenHandler, permissionManager, cacheManager, publisher)
	})

	router.Route("/api/v2", func(r chi.Router) {
		config := DefineHumaConfig("API V2", "2.0.0")
		config.Servers = []*huma.Server{
			{URL: basePath + "/api/v2"},
		}
		api := humachi.New(r, config)
		ConfigureAPIMiddlewares(api, pool, tokenHandler, permissionManager)
		MapV2Routes(api, pool, tokenHandler, permissionManager, cacheManager, publisher)
	})
}

func MapLatestRoutes(
	api huma.API,
	pool *pgxpool.Pool,
	tokenHandler jwt.IJwtTokenHandler,
	permissionManager permissions.IPermissionManager,
	cacheManager *cache.Cache[string],
	publisher message.Publisher,
) {
	// Categories
	getting_category_by_id.MapRoute(api, pool)
	getting_categories.MapRoute(api, pool)
	creating_category.MapRoute(api, pool)
	updating_category.MapRoute(api, pool)
	deleting_category.MapRoute(api, pool)

	// Products
	getting_product_by_id.MapRoute(api, pool)
	getting_products.MapRoute(api, pool)
	creating_product.MapRoute(api, pool)
	updating_product.MapRoute(api, pool)
	deleting_product.MapRoute(api, pool)
}

func MapV1Routes(
	api huma.API,
	pool *pgxpool.Pool,
	tokenHandler jwt.IJwtTokenHandler,
	permissionManager permissions.IPermissionManager,
	cacheManager *cache.Cache[string],
	publisher message.Publisher,
) {
	MapLatestRoutes(api, pool, tokenHandler, permissionManager, cacheManager, publisher)
}

func MapV2Routes(
	api huma.API,
	pool *pgxpool.Pool,
	tokenHandler jwt.IJwtTokenHandler,
	permissionManager permissions.IPermissionManager,
	cacheManager *cache.Cache[string],
	publisher message.Publisher,
) {
	MapLatestRoutes(api, pool, tokenHandler, permissionManager, cacheManager, publisher)
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
	pool *pgxpool.Pool,
	tokenHandler jwt.IJwtTokenHandler,
	permissionManager permissions.IPermissionManager,
) {
	api.UseMiddleware(jwt.SetupJwtAuthentication(api, tokenHandler))
	api.UseMiddleware(permissions.SetupAuthorization(api, permissionManager))
}
