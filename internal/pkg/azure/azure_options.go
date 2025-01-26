package azure

type AzureOptions struct {
	Enabled      bool   `mapstructure:"enabled"`
	KeyVaultName string `mapstructure:"keyVaultName"`
	TenantId     string `mapstructure:"tenantId"`
	ClientId     string `mapstructure:"clientId"`
	ClientSecret string `mapstructure:"clientSecret"`
}
