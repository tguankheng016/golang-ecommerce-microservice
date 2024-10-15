package services

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logger"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type customTokenKeyDBValidator struct {
	db     *gorm.DB
	client *redis.Client
}

func NewCustomTokenKeyDBValidator(db *gorm.DB, client *redis.Client) jwt.IJwtTokenKeyDbValidator {
	return &customTokenKeyDBValidator{
		db:     db,
		client: client,
	}
}

func (c *customTokenKeyDBValidator) ValidateTokenWithTokenKeyFromDb(ctx context.Context, cacheKey string, userId int64, tokenKey string) bool {
	var count int64
	if err := c.db.Model(&models.UserToken{}).Where("user_id = ? AND token_key = ? AND expiration_time > ?", userId, tokenKey, time.Now()).Count(&count).Error; err != nil || count == 0 {
		return false
	}

	if err := c.client.Set(ctx, cacheKey, tokenKey, jwt.DefaultCacheExpiration).Err(); err != nil {
		// Dont return just log
		logger.Logger.Error("error in setting cached token key", zap.Error(err))
	}

	return true
}
