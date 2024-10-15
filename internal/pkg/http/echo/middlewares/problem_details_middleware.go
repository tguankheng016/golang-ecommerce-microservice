package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/problem-details"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logger"
	"go.uber.org/zap"
)

func ProblemDetailsHandler(error error, c echo.Context) {
	if !c.Response().Committed {
		if _, err := problem.ResolveProblemDetails(c.Response(), c.Request(), error); err != nil {
			logger.Logger.Error("Resolve problem details error: ", zap.Error(err))
		}
	}
}
