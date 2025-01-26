package azure

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

type AzureClient struct {
	client *azsecrets.Client
}

func NewAzureClient(options *AzureOptions) (*AzureClient, error) {
	tenantId := os.Getenv("TenantId")
	if tenantId != "" {
		options.TenantId = tenantId
	}

	clientId := os.Getenv("ClientId")
	if clientId != "" {
		options.ClientId = clientId
	}

	clientSecret := os.Getenv("ClientSecret")
	if clientSecret != "" {
		options.ClientSecret = clientSecret
	}

	vaultURI := fmt.Sprintf("https://%s.vault.azure.net/", options.KeyVaultName)

	// Create a credential using the NewDefaultAzureCredential type.
	cred, err := azidentity.NewClientSecretCredential(
		options.TenantId,
		options.ClientId,
		options.ClientSecret,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Establish a connection to the Key Vault client
	client, err := azsecrets.NewClient(vaultURI, cred, nil)
	if err != nil {
		return nil, err
	}

	return &AzureClient{client: client}, nil
}

func (client *AzureClient) GetSecret(secretName string) (string, error) {
	version := ""
	resp, err := client.client.GetSecret(context.TODO(), secretName, version, nil)
	if err != nil {
		return "", err
	}

	return *resp.Value, nil
}
