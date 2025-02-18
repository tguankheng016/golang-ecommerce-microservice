package identities

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/caching"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
	authenticate "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/identities/features/authenticating/v2"
	userConsts "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/constants"
)

const (
	signOutEndpoint = "/api/v1/identities/sign-out"
)

func (suite *IdentityTestSuite) TestShouldSignOutSuccess() {
	request := authenticate.AuthenticateRequest{}
	request.UsernameOrEmailAddress = userConsts.DefaultAdminUserName
	request.Password = "123qwe"

	resp, err := suite.Client.R().
		SetContext(suite.Ctx).
		SetBody(request).
		SetResult(&authenticate.AuthenticateResult{}).
		Post(authenticateEndpoint)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

	result := resp.Result().(*authenticate.AuthenticateResult)
	suite.NotNil(result)

	userId, claims, err := suite.JwtTokenHandler.ValidateToken(suite.Ctx, result.AccessToken, jwt.AccessToken)
	suite.NoError(err)

	resp, err = suite.Client.R().
		SetContext(suite.Ctx).
		SetAuthToken(result.AccessToken).
		Post(signOutEndpoint)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode())

	tokenValidityKey, ok := claims[jwt.TokenValidityKey].(string)
	suite.True(ok)

	_, err = suite.CacheManager.Get(suite.Ctx, jwt.GenerateTokenValidityCacheKey(userId, tokenValidityKey))
	suite.True(caching.CheckIsCacheValueNotFound(err))

	refreshTokenValidityKey, ok := claims[jwt.TokenValidityKey].(string)
	suite.True(ok)

	_, err = suite.CacheManager.Get(suite.Ctx, jwt.GenerateTokenValidityCacheKey(userId, refreshTokenValidityKey))
	suite.True(caching.CheckIsCacheValueNotFound(err))
}
