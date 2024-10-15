package services

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logger"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type customStampDBValidator struct {
	db     *gorm.DB
	client *redis.Client
}

func NewCustomStampDBValidator(db *gorm.DB, client *redis.Client) jwt.IJwtSecurityStampDbValidator {
	return &customStampDBValidator{
		db:     db,
		client: client,
	}
}

func (c *customStampDBValidator) ValidateTokenWithStampFromDb(ctx context.Context, cacheKey string, userId int64, securityStamp string) bool {
	var user models.User
	if err := c.db.First(&user, userId).Error; err != nil {
		return false
	}

	if err := c.client.Set(ctx, cacheKey, user.SecurityStamp.String(), jwt.DefaultCacheExpiration).Err(); err != nil {
		// Dont return just log
		logger.Logger.Error("error in setting cached security stamp", zap.Error(err))
	}

	if user.SecurityStamp.String() != securityStamp {
		return false
	}

	return true
}
