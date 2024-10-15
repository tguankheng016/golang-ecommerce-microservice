package jwt

import (
	"context"
	"fmt"

	jwtGo "github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/constants"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logger"
	"go.uber.org/zap"
)

type IJwtSecurityStampValidator interface {
	ValidateTokenWithStamp(ctx context.Context, userId int64, claims jwtGo.MapClaims) error
}

type IJwtSecurityStampDbValidator interface {
	ValidateTokenWithStampFromDb(ctx context.Context, cacheKey string, userId int64, securityStamp string) bool
}

type jwtSecurityStampValidator struct {
	dbValidator IJwtSecurityStampDbValidator
	client      *redis.Client
}

func NewSecurityStampValidator(dbValidator IJwtSecurityStampDbValidator, client *redis.Client) IJwtSecurityStampValidator {
	return &jwtSecurityStampValidator{
		dbValidator: dbValidator,
		client:      client,
	}
}

func (j *jwtSecurityStampValidator) ValidateTokenWithStamp(ctx context.Context, userId int64, claims jwtGo.MapClaims) error {
	securityStamp := claims[constants.SecurityStampKey]
	invalidSecurityStampErr := errors.New("Invalid stamp")

	if securityStamp == nil {
		return invalidSecurityStampErr
	}

	securityStampStr, ok := securityStamp.(string)
	if !ok {
		return invalidSecurityStampErr
	}

	isValid := j.validateTokenWithStampFromCache(ctx, userId, securityStampStr)

	if !isValid {
		isValid = j.dbValidator.ValidateTokenWithStampFromDb(ctx, GenerateStampCacheKey(userId), userId, securityStampStr)
	}

	if !isValid {
		return invalidSecurityStampErr
	}

	return nil
}

func (j *jwtSecurityStampValidator) validateTokenWithStampFromCache(ctx context.Context, userId int64, securityStamp string) bool {
	cacheKey := GenerateStampCacheKey(userId)

	cachedStamp, err := j.client.Get(ctx, cacheKey).Result()
	if err != nil {
		logger.Logger.Error("error in getting cached stamp", zap.Error(err))
		return false
	}

	return cachedStamp != "" && cachedStamp == securityStamp
}

func GenerateStampCacheKey(userId int64) string {
	return fmt.Sprintf("%s.%d", constants.SecurityStampKey, userId)
}
