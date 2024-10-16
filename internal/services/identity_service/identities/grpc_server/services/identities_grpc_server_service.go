package services

import (
	"context"

	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
	identity_service "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/identities/grpc_server/protos"
)

type IdentityGrpcServerService struct {
	securityStampvalidator jwt.IJwtSecurityStampDbValidator
	tokenKeyValidator      jwt.IJwtTokenKeyDbValidator
}

func NewIdentityGrpcServerService(securityStampvalidator jwt.IJwtSecurityStampDbValidator, tokenKeyValidator jwt.IJwtTokenKeyDbValidator) *IdentityGrpcServerService {
	return &IdentityGrpcServerService{
		securityStampvalidator: securityStampvalidator,
		tokenKeyValidator:      tokenKeyValidator,
	}
}

func (i *IdentityGrpcServerService) ValidateSecurityStamp(ctx context.Context, req *identity_service.GetValidateSecurityStampRequest) (*identity_service.GetValidateSecurityStampResponse, error) {
	isValid := i.securityStampvalidator.ValidateTokenWithStampFromDb(ctx, req.CacheKey, req.UserId, req.SecurityStamp)

	var result = &identity_service.GetValidateSecurityStampResponse{
		IsValid: isValid,
	}

	return result, nil
}

func (i *IdentityGrpcServerService) ValidateKey(ctx context.Context, req *identity_service.GetValidateTokenKeyRequest) (*identity_service.GetValidateTokenKeyResponse, error) {
	isValid := i.tokenKeyValidator.ValidateTokenWithTokenKeyFromDb(ctx, req.CacheKey, req.UserId, req.TokenKey)

	var result = &identity_service.GetValidateTokenKeyResponse{
		IsValid: isValid,
	}

	return result, nil
}
