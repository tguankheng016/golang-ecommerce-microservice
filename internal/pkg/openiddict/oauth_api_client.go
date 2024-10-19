package openiddict

import (
	"context"
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type ConnectTokenResponseDto struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ConnectTokenErrorDto struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type ConnectUserInfoResponseDto struct {
	Sub               string `json:"sub"`
	PreferredUsername string `json:"preferred_username"`
	FamilyName        string `json:"family_name"`
	GivenName         string `json:"given_name"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	Picture           string `json:"picture"`
}

type IOAuthApiClient interface {
	ConnectToken(ctx context.Context, code string, redirectUri string) (*ConnectTokenResponseDto, error)
	ConnectUserInfo(ctx context.Context, accessToken string) (*ConnectUserInfoResponseDto, error)
}

type OAuthApiClient struct {
	client       *resty.Client
	oauthOptions *OAuthOptions
}

func NewOAuthApiClient(oauthOptions *OAuthOptions) IOAuthApiClient {
	client := resty.New()
	client.SetBaseURL(oauthOptions.BaseUrl)
	return &OAuthApiClient{
		client:       client,
		oauthOptions: oauthOptions,
	}
}

func (c *OAuthApiClient) ConnectToken(ctx context.Context, code string, redirectUri string) (*ConnectTokenResponseDto, error) {
	data := map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"client_id":     c.oauthOptions.ClientId,
		"client_secret": c.oauthOptions.ClientSecret,
		"redirect_uri":  redirectUri,
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetFormData(data).
		Post("connect/token")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		var errorDto ConnectTokenErrorDto
		if err := json.Unmarshal(resp.Body(), &errorDto); err != nil {
			return nil, err
		}
		return nil, errors.New(errorDto.ErrorDescription)
	}

	var result ConnectTokenResponseDto
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *OAuthApiClient) ConnectUserInfo(ctx context.Context, accessToken string) (*ConnectUserInfoResponseDto, error) {
	resp, err := c.client.R().
		SetContext(ctx).
		SetAuthToken(accessToken).
		Get("connect/userinfo")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		var errorDto ConnectTokenErrorDto
		if err := json.Unmarshal(resp.Body(), &errorDto); err != nil {
			return nil, err
		}
		return nil, errors.New(errorDto.ErrorDescription)
	}

	var result ConnectUserInfoResponseDto
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}

	return &result, nil
}
