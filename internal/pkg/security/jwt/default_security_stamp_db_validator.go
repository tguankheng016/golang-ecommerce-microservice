package jwt

import (
	"context"

	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logging"
	identity_service "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt/grpc_client/protos"
	"go.uber.org/zap"
)

type defaultStampDBValidator struct {
	identityGrpcClient identity_service.IdentityGrpcServiceClient
}

func NewDefaultStampDBValidator(identityGrpcClient identity_service.IdentityGrpcServiceClient) IJwtSecurityStampDbValidator {
	return &defaultStampDBValidator{
		identityGrpcClient: identityGrpcClient,
	}
}

func (d *defaultStampDBValidator) ValidateTokenWithStampFromDb(ctx context.Context, cacheKey string, userId int64, securityStamp string) bool {
	response, err := d.identityGrpcClient.ValidateSecurityStamp(ctx, &identity_service.GetValidateSecurityStampRequest{
		CacheKey:      cacheKey,
		UserId:        userId,
		SecurityStamp: securityStamp,
	})

	if err != nil {
		logging.Logger.Error("error in validating stamp from grpc", zap.Error(err))
		return false
	}

	return response.IsValid
}
