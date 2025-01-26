package openiddict

type OAuthOptions struct {
	ClientId     string `mapstructure:"clientId"`
	ClientSecret string `mapstructure:"clientSecret"`
	BaseUrl      string `mapstructure:"baseUrl"`
}
