package storage

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/smarthow/azure-for-dummies/internal/auth"
	"github.com/smarthow/azure-for-dummies/internal/azutil"
	"github.com/smarthow/azure-for-dummies/internal/provider"
)

type azureProvider struct {
	auth *auth.Context
}

func newAzureProvider(ctx *auth.Context) *azureProvider {
	return &azureProvider{auth: ctx}
}

func (p *azureProvider) ListStorageAccounts(ctx context.Context) ([]provider.StorageAccount, error) {
	client, err := armstorage.NewAccountsClient(p.auth.SubscriptionID, p.auth.Credential, nil)
	if err != nil {
		return nil, fmt.Errorf("new storage client: %w", err)
	}

	var accounts []provider.StorageAccount
	pager := client.NewListPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("list storage accounts: %w", err)
		}
		for _, a := range page.Value {
			accounts = append(accounts, mapAccount(a))
		}
	}
	return accounts, nil
}

func (p *azureProvider) GetStorageAccount(ctx context.Context, rg, name string) (provider.StorageAccount, error) {
	client, err := armstorage.NewAccountsClient(p.auth.SubscriptionID, p.auth.Credential, nil)
	if err != nil {
		return provider.StorageAccount{}, fmt.Errorf("new storage client: %w", err)
	}

	resp, err := client.GetProperties(ctx, rg, name, nil)
	if err != nil {
		return provider.StorageAccount{}, fmt.Errorf("get storage account: %w", err)
	}
	return mapAccount(&resp.Account), nil
}

func (p *azureProvider) ListBlobContainers(ctx context.Context, rg, accountName string) ([]provider.BlobContainer, error) {
	client, err := armstorage.NewBlobContainersClient(p.auth.SubscriptionID, p.auth.Credential, nil)
	if err != nil {
		return nil, fmt.Errorf("new blob containers client: %w", err)
	}

	var containers []provider.BlobContainer
	pager := client.NewListPager(rg, accountName, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("list blob containers: %w", err)
		}
		for _, c := range page.Value {
			containers = append(containers, mapContainer(c))
		}
	}
	return containers, nil
}

func mapAccount(a *armstorage.Account) provider.StorageAccount {
	sa := provider.StorageAccount{
		ID:            azutil.DerefStr(a.ID),
		Name:          azutil.DerefStr(a.Name),
		ResourceGroup: azutil.ExtractRG(azutil.DerefStr(a.ID)),
		Region:        azutil.DerefStr(a.Location),
	}
	if a.Kind != nil {
		sa.Kind = string(*a.Kind)
	}
	if a.SKU != nil && a.SKU.Name != nil {
		sa.SKU = string(*a.SKU.Name)
	}
	if a.Properties != nil {
		if a.Properties.AccessTier != nil {
			sa.AccessTier = string(*a.Properties.AccessTier)
		}
		if a.Properties.ProvisioningState != nil {
			sa.ProvisioningState = string(*a.Properties.ProvisioningState)
		}
		if a.Properties.CreationTime != nil {
			sa.CreationTime = *a.Properties.CreationTime
		}
		sa.IsHnsEnabled = azutil.DerefBool(a.Properties.IsHnsEnabled)
		sa.AllowBlobPublicAccess = azutil.DerefBool(a.Properties.AllowBlobPublicAccess)
		if a.Properties.MinimumTLSVersion != nil {
			sa.MinTLSVersion = string(*a.Properties.MinimumTLSVersion)
		}
		if a.Properties.NetworkRuleSet != nil && a.Properties.NetworkRuleSet.DefaultAction != nil {
			sa.NetworkDefaultAction = string(*a.Properties.NetworkRuleSet.DefaultAction)
		}
		if a.Properties.PrimaryEndpoints != nil {
			sa.PrimaryBlobEndpoint = azutil.DerefStr(a.Properties.PrimaryEndpoints.Blob)
			sa.PrimaryFileEndpoint = azutil.DerefStr(a.Properties.PrimaryEndpoints.File)
			sa.PrimaryTableEndpoint = azutil.DerefStr(a.Properties.PrimaryEndpoints.Table)
			sa.PrimaryQueueEndpoint = azutil.DerefStr(a.Properties.PrimaryEndpoints.Queue)
		}
	}
	return sa
}

func mapContainer(c *armstorage.ListContainerItem) provider.BlobContainer {
	bc := provider.BlobContainer{
		Name: azutil.DerefStr(c.Name),
	}
	if c.Properties != nil {
		if c.Properties.PublicAccess != nil {
			sa := string(*c.Properties.PublicAccess)
			if sa == "None" {
				bc.PublicAccess = "Private"
			} else {
				bc.PublicAccess = sa
			}
		}
		if c.Properties.LeaseStatus != nil {
			bc.LeaseStatus = string(*c.Properties.LeaseStatus)
		}
		if c.Properties.LeaseState != nil {
			bc.LeaseState = string(*c.Properties.LeaseState)
		}
		if c.Properties.LastModifiedTime != nil {
			bc.LastModified = *c.Properties.LastModifiedTime
		}
		if c.Properties.LegalHold != nil {
			bc.HasLegalHold = azutil.DerefBool(c.Properties.LegalHold.HasLegalHold)
		}
		bc.HasImmutabilityPolicy = c.Properties.ImmutabilityPolicy != nil
	}
	return bc
}
