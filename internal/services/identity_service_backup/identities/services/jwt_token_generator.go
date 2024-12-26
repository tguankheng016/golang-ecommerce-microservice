package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	jwtGo "github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/constants"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logger"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	AccessTokenExpirationTime  = 24 * time.Hour
	RefreshTokenExpirationTime = 30 * 24 * time.Hour
)

type IJwtTokenGenerator interface {
	GenerateAccessToken(ctx context.Context, user *models.User, refreshTokenKey string) (string, int, error)
	GenerateRefreshToken(ctx context.Context, user *models.User) (string, string, int, error)
	RemoveUserTokens(ctx context.Context, userId int64, claims jwtGo.MapClaims) error
}

type jwtTokenGenerator struct {
	secretKey string
	issuer    string
	audience  string
	db        *gorm.DB
	client    redis.UniversalClient
}

func NewJwtTokenGenerator(db *gorm.DB, client redis.UniversalClient, authOptions *jwt.AuthOptions) IJwtTokenGenerator {
	return &jwtTokenGenerator{
		secretKey: authOptions.SecretKey,
		issuer:    authOptions.Issuer,
		audience:  authOptions.Audience,
		db:        db,
		client:    client,
	}
}

// GenerateAccessToken generates an access token for the given user.
//
// The function takes in a user model and the refresh token key, and returns a JWT token, the expiration time in seconds, and an error.
// The expiration time is set to the AccessTokenExpirationTime duration.
// The function will also insert a new row into the user_tokens table with the user's ID, the token key, and the expiration time.
// Finally, the function will cache the token validity key in Redis with the key "user_token:<user_id>:<token_key>" and the expiration time set to the DefaultCacheExpiration.
// If there is an error caching the token validity key, the function will log the error but not return it.
// If there is an error, the function will return an empty string, 0, and the error.
func (j *jwtTokenGenerator) GenerateAccessToken(ctx context.Context, user *models.User, refreshTokenKey string) (string, int, error) {
	claims, err := j.createJwtClaims(ctx, user, jwt.AccessToken, refreshTokenKey)

	if err != nil {
		return "", 0, err
	}

	accessToken, err := j.createToken(claims)

	return accessToken, int(AccessTokenExpirationTime.Seconds()), err
}

// GenerateRefreshToken generates a refresh token for the given user.
// The function takes in a user model and returns a JWT token, the token key, and the expiration time in seconds.
// If there is an error, the function will return an empty string, an empty string, 0, and the error.
func (j *jwtTokenGenerator) GenerateRefreshToken(ctx context.Context, user *models.User) (string, string, int, error) {
	claims, err := j.createJwtClaims(ctx, user, jwt.RefreshToken, "")

	if err != nil {
		return "", "", 0, err
	}

	refreshToken, err := j.createToken(claims)

	refreshTokenKey := claims[constants.TokenValidityKey]
	refreshTokenStr := fmt.Sprintf("%s", refreshTokenKey)

	return refreshToken, refreshTokenStr, int(RefreshTokenExpirationTime.Seconds()), err
}

// RemoveUserTokens removes a user's access token and refresh token from the database and Redis cache.
//
// It takes in the user's ID and the claims map from the JWT token. It checks if the claims map contains
// the token validity key and removes the token from the database and Redis cache. If the claims map also
// contains the refresh token validity key, it removes the refresh token from the database and Redis cache
// as well.
//
// If there is an error removing the tokens, it returns the error.
func (j *jwtTokenGenerator) RemoveUserTokens(ctx context.Context, userId int64, claims jwtGo.MapClaims) error {
	tokenKey, ok := claims[constants.TokenValidityKey]
	if !ok {
		return errors.New("Invalid token key")
	}

	if err := j.removeToken(ctx, userId, tokenKey.(string)); err != nil {
		return err
	}

	refreshTokenKey, ok := claims[constants.RefreshTokenValidityKey]
	if ok {
		if err := j.removeToken(ctx, userId, refreshTokenKey.(string)); err != nil {
			return err
		}
	}

	return nil
}

// createToken creates a new JWT token with the given claims.
//
// The function takes in a MapClaims object and returns a signed string
// representation of the token. The signing method used is HS256.
//
// The function also sets the "iss" and "aud" headers of the token to the
// issuer and audience of the JWT, respectively.
//
// If there is an error signing the token, the function will return the error.
func (j *jwtTokenGenerator) createToken(claims jwtGo.MapClaims) (string, error) {
	token := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, claims)
	token.Header["iss"] = j.issuer
	token.Header["aud"] = j.audience

	return token.SignedString([]byte(j.secretKey))
}

// createJwtClaims creates a new JWT claims map from the given user and token type.
//
// The claims map will contain the user's ID, the JWT ID, the issued at time, the not before time,
// the expiration time, the token validity key, the security stamp, and the token type.
//
// If the token type is RefreshToken, the expiration time will be set to the RefreshTokenExpirationTime.
// Otherwise, the expiration time will be set to the AccessTokenExpirationTime.
//
// The function will also insert a new row into the user_tokens table with the user's ID, the token key,
// and the expiration time.
//
// Finally, the function will cache the token validity key in Redis with the key
// "user_token:<user_id>:<token_key>" and the expiration time set to the DefaultCacheExpiration.
// If there is an error caching the token validity key, the function will log the error but not return
// it.
//
// Returns the claims map and an error if there is one.
func (j *jwtTokenGenerator) createJwtClaims(ctx context.Context, user *models.User, tokenType jwt.TokenType, refreshTokenKey string) (jwtGo.MapClaims, error) {
	tokenValidityKey, err := uuid.NewV4()

	if err != nil {
		return nil, err
	}

	claims := jwtGo.MapClaims{}

	claims["jti"], err = uuid.NewV4()

	if err != nil {
		return nil, err
	}

	now := time.Now()
	var expiration time.Duration
	if tokenType == jwt.RefreshToken {
		expiration = RefreshTokenExpirationTime
	} else {
		expiration = AccessTokenExpirationTime
	}

	claims["sub"] = strconv.FormatInt(user.Id, 10)
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	claims["exp"] = now.Add(expiration).Unix()
	claims[constants.TokenValidityKey] = tokenValidityKey
	claims[constants.SecurityStampKey] = user.SecurityStamp
	claims["token_type"] = strconv.Itoa(int(tokenType))

	if refreshTokenKey != "" {
		claims[constants.RefreshTokenValidityKey] = refreshTokenKey
	}

	// Add User Token
	userToken := models.UserToken{
		UserId:         user.Id,
		TokenKey:       tokenValidityKey.String(),
		ExpirationTime: now.Add(expiration),
	}

	if err := j.db.Create(&userToken).Error; err != nil {
		return nil, errors.Wrap(err, "error when inserting user token into the database.")
	}

	if err := j.client.Set(ctx, jwt.GenerateTokenValidityCacheKey(user.Id, tokenValidityKey.String()), tokenValidityKey.String(), jwt.DefaultCacheExpiration).Err(); err != nil {
		// Dont return just log
		logger.Logger.Error("error in setting cached token key", zap.Error(err))
	}

	return claims, nil
}

// removeToken deletes a user's token with the given tokenKey from the database and the redis cache.
// It is used to invalidate a user's token, which is necessary for logging out a user.
func (j *jwtTokenGenerator) removeToken(ctx context.Context, userId int64, tokenKey string) error {
	if err := j.db.Where("token_key = ? AND user_id = ?", tokenKey, userId).Delete(&models.UserToken{}).Error; err != nil {
		return err
	}

	if err := j.client.Del(ctx, jwt.GenerateTokenValidityCacheKey(userId, tokenKey)).Err(); err != nil {
		// Dont return just log
		logger.Logger.Error("error in deleting cached token key", zap.Error(err))
	}

	return nil
}
