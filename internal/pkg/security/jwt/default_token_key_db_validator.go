package jwt

import (
	"context"

	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logger"
	identity_service "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt/grpc_client/protos"
	"go.uber.org/zap"
)

type defaultTokenKeyDBValidator struct {
	identityGrpcClient identity_service.IdentityGrpcServiceClient
}

func NewDefaultTokenKeyDBValidator(identityGrpcClient identity_service.IdentityGrpcServiceClient) IJwtTokenKeyDbValidator {
	return &defaultTokenKeyDBValidator{
		identityGrpcClient: identityGrpcClient,
	}
}

func (d *defaultTokenKeyDBValidator) ValidateTokenWithTokenKeyFromDb(ctx context.Context, cacheKey string, userId int64, tokenKey string) bool {
	response, err := d.identityGrpcClient.ValidateKey(ctx, &identity_service.GetValidateTokenKeyRequest{
		CacheKey:      cacheKey,
		UserId:        userId,
		TokenKey: tokenKey,
	})

	if err != nil {
		logger.Logger.Error("error in validating token key from grpc", zap.Error(err))
		return false
	}

	return response.IsValid
}
