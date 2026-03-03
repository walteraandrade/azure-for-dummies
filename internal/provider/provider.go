package provider

import (
	"context"
	"time"
)

type ContainerApp struct {
	ID                string
	Name              string
	ResourceGroup     string
	Region            string
	State             string
	Revision          string
	FQDN              string
	IngressExternal   bool
	IngressPort       int32
	Containers        []ContainerInfo
	ScaleMin          int32
	ScaleMax          int32
	ProvisioningState string
}

type ContainerInfo struct {
	Name   string
	Image  string
	CPU    float64
	Memory string
	Env    []EnvVar
}

type EnvVar struct {
	Name      string
	Value     string
	SecretRef string
}

type RevisionInfo struct {
	Name          string
	Active        bool
	CreatedTime   time.Time
	TrafficWeight int32
	Replicas      int32
	HealthState   string
	RunningState  string
}

type LogEntry struct {
	Timestamp time.Time
	Message   string
}

type ContainerAppsProvider interface {
	ListContainerApps(ctx context.Context, sub string) ([]ContainerApp, error)
	GetContainerApp(ctx context.Context, rg, name string) (ContainerApp, error)
	ListRevisions(ctx context.Context, rg, appName string) ([]RevisionInfo, error)
	StreamLogs(ctx context.Context, rg, appName string) (<-chan LogEntry, error)
}

type PostgresServer struct {
	ID              string
	Name            string
	ResourceGroup   string
	Region          string
	State           string
	Tier            string
	SKU             string
	Version         string
	FQDN            string
	StorageGB       int32
	BackupRetention int32
}

type FirewallRuleInfo struct {
	Name    string
	StartIP string
	EndIP   string
}

type MetricPoint struct {
	Timestamp time.Time
	Average   float64
}

type MetricSeries struct {
	Name   string
	Unit   string
	Points []MetricPoint
}

type PostgresProvider interface {
	ListServers(ctx context.Context, sub string) ([]PostgresServer, error)
	GetServer(ctx context.Context, rg, name string) (PostgresServer, error)
	ListFirewallRules(ctx context.Context, rg, serverName string) ([]FirewallRuleInfo, error)
	GetMetrics(ctx context.Context, resourceID string, names []string) ([]MetricSeries, error)
}

type StorageAccount struct {
	ID, Name, ResourceGroup, Region     string
	Kind, SKU, AccessTier               string
	ProvisioningState                   string
	CreationTime                        time.Time
	PrimaryBlobEndpoint                 string
	PrimaryFileEndpoint                 string
	PrimaryTableEndpoint                string
	PrimaryQueueEndpoint                string
	IsHnsEnabled, AllowBlobPublicAccess bool
	MinTLSVersion, NetworkDefaultAction string
}

type BlobContainer struct {
	Name, PublicAccess, LeaseStatus, LeaseState string
	LastModified                                time.Time
	HasLegalHold, HasImmutabilityPolicy         bool
}

type StorageProvider interface {
	ListStorageAccounts(ctx context.Context, sub string) ([]StorageAccount, error)
	GetStorageAccount(ctx context.Context, rg, name string) (StorageAccount, error)
	ListBlobContainers(ctx context.Context, rg, accountName string) ([]BlobContainer, error)
}

type Provider interface {
	ContainerAppsProvider
	PostgresProvider
	StorageProvider
}
