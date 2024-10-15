package jwt

import (
	"context"
	"fmt"
	"time"

	jwtGo "github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/constants"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logger"
	"go.uber.org/zap"
)

const (
	DefaultCacheExpiration = 1 * time.Hour
)

type IJwtTokenKeyValidator interface {
	ValidateTokenWithTokenKey(ctx context.Context, userId int64, claims jwtGo.MapClaims) error
}

type IJwtTokenKeyDbValidator interface {
	ValidateTokenWithTokenKeyFromDb(ctx context.Context, cacheKey string, userId int64, tokenKey string) bool
}

type jwtTokenKeyValidator struct {
	dbValidator IJwtTokenKeyDbValidator
	client      redis.UniversalClient
}

func NewTokenKeyValidator(dbValidator IJwtTokenKeyDbValidator, client redis.UniversalClient) IJwtTokenKeyValidator {
	return &jwtTokenKeyValidator{
		dbValidator: dbValidator,
		client:      client,
	}
}

func (j *jwtTokenKeyValidator) ValidateTokenWithTokenKey(ctx context.Context, userId int64, claims jwtGo.MapClaims) error {
	tokenKey := claims[constants.TokenValidityKey]
	invalidTokenKeyErr := errors.New("Invalid token key")

	if tokenKey == nil {
		return invalidTokenKeyErr
	}

	tokenKeyStr, ok := tokenKey.(string)
	if !ok {
		return invalidTokenKeyErr
	}

	isValid := j.validateTokenWithTokenKeyFromCache(ctx, userId, tokenKeyStr)

	if !isValid {
		isValid = j.dbValidator.ValidateTokenWithTokenKeyFromDb(ctx, GenerateTokenValidityCacheKey(userId, tokenKeyStr), userId, tokenKeyStr)
	}

	if !isValid {
		return invalidTokenKeyErr
	}

	return nil
}

func (j *jwtTokenKeyValidator) validateTokenWithTokenKeyFromCache(ctx context.Context, userId int64, tokenKey string) bool {
	tokenCacheKey := GenerateTokenValidityCacheKey(userId, tokenKey)

	cachedTokenKey, err := j.client.Get(ctx, tokenCacheKey).Result()
	if err != nil {
		logger.Logger.Error("error in getting cached token key", zap.Error(err))
		return false
	}

	return cachedTokenKey != ""
}

func GenerateTokenValidityCacheKey(userId int64, tokenKey string) string {
	return fmt.Sprintf("%s.%d.%s", constants.TokenValidityKey, userId, tokenKey)
}
