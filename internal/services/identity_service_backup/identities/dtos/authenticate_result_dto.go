package dtos

type AuthenticateResult struct {
	AccessToken                 string `json:"access_token"`
	ExpireInSeconds             int    `json:"expire_in_seconds"`
	RefreshToken                string `json:"refresh_token"`
	RefreshTokenExpireInSeconds int    `json:"refresh_token_expire_in_seconds"`
} // @name AuthenticateResult
