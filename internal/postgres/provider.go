package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresqlflexibleservers/v4"
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

func (p *azureProvider) ListServers(ctx context.Context, subID string) ([]provider.PostgresServer, error) {
	client, err := armpostgresqlflexibleservers.NewServersClient(subID, p.auth.Credential, nil)
	if err != nil {
		return nil, fmt.Errorf("new postgres client: %w", err)
	}

	var servers []provider.PostgresServer
	pager := client.NewListPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("list postgres servers: %w", err)
		}
		for _, s := range page.Value {
			servers = append(servers, mapServer(s))
		}
	}
	return servers, nil
}

func (p *azureProvider) GetServer(ctx context.Context, rg, name string) (provider.PostgresServer, error) {
	client, err := armpostgresqlflexibleservers.NewServersClient(p.auth.SubscriptionID, p.auth.Credential, nil)
	if err != nil {
		return provider.PostgresServer{}, fmt.Errorf("new postgres client: %w", err)
	}

	resp, err := client.Get(ctx, rg, name, nil)
	if err != nil {
		return provider.PostgresServer{}, fmt.Errorf("get postgres server: %w", err)
	}
	return mapServer(&resp.Server), nil
}

func (p *azureProvider) ListFirewallRules(ctx context.Context, rg, serverName string) ([]provider.FirewallRuleInfo, error) {
	client, err := armpostgresqlflexibleservers.NewFirewallRulesClient(p.auth.SubscriptionID, p.auth.Credential, nil)
	if err != nil {
		return nil, fmt.Errorf("new firewall rules client: %w", err)
	}

	var rules []provider.FirewallRuleInfo
	pager := client.NewListByServerPager(rg, serverName, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("list firewall rules: %w", err)
		}
		for _, r := range page.Value {
			rule := provider.FirewallRuleInfo{Name: azutil.DerefStr(r.Name)}
			if r.Properties != nil {
				rule.StartIP = azutil.DerefStr(r.Properties.StartIPAddress)
				rule.EndIP = azutil.DerefStr(r.Properties.EndIPAddress)
			}
			rules = append(rules, rule)
		}
	}
	return rules, nil
}

func (p *azureProvider) GetMetrics(ctx context.Context, resourceID string, names []string) ([]provider.MetricSeries, error) {
	client, err := azquery.NewMetricsClient(p.auth.Credential, nil)
	if err != nil {
		return nil, fmt.Errorf("new metrics client: %w", err)
	}

	metricNames := strings.Join(names, ",")
	interval := "PT5M"
	avg := azquery.AggregationTypeAverage
	resp, err := client.QueryResource(ctx, resourceID, &azquery.MetricsClientQueryResourceOptions{
		MetricNames: &metricNames,
		Timespan:    toTimeInterval("PT1H"),
		Interval:    &interval,
		Aggregation: []*azquery.AggregationType{&avg},
	})
	if err != nil {
		return nil, fmt.Errorf("query metrics: %w", err)
	}

	var series []provider.MetricSeries
	for _, m := range resp.Value {
		ms := provider.MetricSeries{
			Name: azutil.DerefStr(m.Name.Value),
		}
		if m.Unit != nil {
			ms.Unit = string(*m.Unit)
		}
		for _, ts := range m.TimeSeries {
			for _, dp := range ts.Data {
				if dp.TimeStamp == nil {
					continue
				}
				ms.Points = append(ms.Points, provider.MetricPoint{
					Timestamp: *dp.TimeStamp,
					Average:   azutil.DerefFloat64(dp.Average),
				})
			}
		}
		series = append(series, ms)
	}
	return series, nil
}

func toTimeInterval(s string) *azquery.TimeInterval {
	t := azquery.TimeInterval(s)
	return &t
}

func mapServer(s *armpostgresqlflexibleservers.Server) provider.PostgresServer {
	srv := provider.PostgresServer{
		ID:            azutil.DerefStr(s.ID),
		Name:          azutil.DerefStr(s.Name),
		ResourceGroup: azutil.ExtractRG(azutil.DerefStr(s.ID)),
		Region:        azutil.DerefStr(s.Location),
	}
	if s.SKU != nil {
		srv.SKU = azutil.DerefStr(s.SKU.Name)
		if s.SKU.Tier != nil {
			srv.Tier = string(*s.SKU.Tier)
		}
	}
	if s.Properties != nil {
		if s.Properties.State != nil {
			srv.State = string(*s.Properties.State)
		}
		if s.Properties.Version != nil {
			srv.Version = string(*s.Properties.Version)
		}
		srv.FQDN = azutil.DerefStr(s.Properties.FullyQualifiedDomainName)
		if s.Properties.Storage != nil {
			srv.StorageGB = azutil.DerefInt32(s.Properties.Storage.StorageSizeGB)
		}
		if s.Properties.Backup != nil {
			srv.BackupRetention = azutil.DerefInt32(s.Properties.Backup.BackupRetentionDays)
		}
	}
	return srv
}
