package configurations

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/core"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
	"gorm.io/gorm"
)

func ConfigMiddlewares(
	appOptions *core.AppOptions,
	e *echo.Echo,
	db *gorm.DB,
	validator *validator.Validate,
	jwtTokenHandler jwt.IJwtTokenHandler,
	permissionManager permissions.IPermissionManager,
) {
	skipper := func(c echo.Context) bool {
		fmt.Println(c.Request().URL.Path)
		return c.Request().URL.Path == "/" ||
			strings.Contains(c.Request().URL.Path, "swagger") ||
			strings.Contains(c.Request().URL.Path, "metrics") ||
			strings.Contains(c.Request().URL.Path, "health") ||
			strings.Contains(c.Request().URL.Path, "favicon.ico")
	}

	e.HideBanner = false

	e.Use(middlewares.SetupLogger())
	e.HTTPErrorHandler = middlewares.ProblemDetailsHandler
	e.Use(middleware.RequestID())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level:   5,
		Skipper: skipper,
	}))

	e.Use(middleware.BodyLimit("2M"))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(appOptions.CorsOrigins, ","), // Adjust the origin as needed
	}))

	e.Use(middlewares.SetupAuthenticate(skipper, jwtTokenHandler))
	e.Use(middlewares.SetupAuthorize(skipper, permissionManager))
	e.Use(middlewares.SetupTransaction(skipper, db))
}
