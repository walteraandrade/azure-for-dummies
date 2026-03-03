package containerapps

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers/v3"
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

func (p *azureProvider) ListContainerApps(ctx context.Context, subID string) ([]provider.ContainerApp, error) {
	client, err := armappcontainers.NewContainerAppsClient(subID, p.auth.Credential, nil)
	if err != nil {
		return nil, fmt.Errorf("new container apps client: %w", err)
	}

	var apps []provider.ContainerApp
	pager := client.NewListBySubscriptionPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("list container apps: %w", err)
		}
		for _, a := range page.Value {
			apps = append(apps, mapContainerAppSummary(a))
		}
	}
	return apps, nil
}

func (p *azureProvider) GetContainerApp(ctx context.Context, rg, name string) (provider.ContainerApp, error) {
	client, err := armappcontainers.NewContainerAppsClient(p.auth.SubscriptionID, p.auth.Credential, nil)
	if err != nil {
		return provider.ContainerApp{}, fmt.Errorf("new container apps client: %w", err)
	}

	resp, err := client.Get(ctx, rg, name, nil)
	if err != nil {
		return provider.ContainerApp{}, fmt.Errorf("get container app: %w", err)
	}

	app := mapContainerAppSummary(&resp.ContainerApp)
	props := resp.Properties
	if props == nil {
		return app, nil
	}

	if props.Configuration != nil && props.Configuration.Ingress != nil {
		ing := props.Configuration.Ingress
		app.FQDN = azutil.DerefStr(ing.Fqdn)
		app.IngressExternal = azutil.DerefBool(ing.External)
		app.IngressPort = azutil.DerefInt32(ing.TargetPort)
	}

	if props.ProvisioningState != nil {
		app.ProvisioningState = string(*props.ProvisioningState)
	}

	if props.Template != nil {
		if props.Template.Scale != nil {
			app.ScaleMin = azutil.DerefInt32(props.Template.Scale.MinReplicas)
			app.ScaleMax = azutil.DerefInt32(props.Template.Scale.MaxReplicas)
		}
		for _, c := range props.Template.Containers {
			ci := provider.ContainerInfo{
				Name:  azutil.DerefStr(c.Name),
				Image: azutil.DerefStr(c.Image),
			}
			if c.Resources != nil {
				ci.CPU = azutil.DerefFloat64(c.Resources.CPU)
				ci.Memory = azutil.DerefStr(c.Resources.Memory)
			}
			for _, e := range c.Env {
				ci.Env = append(ci.Env, provider.EnvVar{
					Name:      azutil.DerefStr(e.Name),
					Value:     azutil.DerefStr(e.Value),
					SecretRef: azutil.DerefStr(e.SecretRef),
				})
			}
			app.Containers = append(app.Containers, ci)
		}
	}

	return app, nil
}

func (p *azureProvider) ListRevisions(ctx context.Context, rg, appName string) ([]provider.RevisionInfo, error) {
	client, err := armappcontainers.NewContainerAppsRevisionsClient(p.auth.SubscriptionID, p.auth.Credential, nil)
	if err != nil {
		return nil, fmt.Errorf("new revisions client: %w", err)
	}

	var revisions []provider.RevisionInfo
	pager := client.NewListRevisionsPager(rg, appName, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("list revisions: %w", err)
		}
		for _, r := range page.Value {
			rev := provider.RevisionInfo{
				Name: azutil.DerefStr(r.Name),
			}
			if r.Properties != nil {
				rev.Active = azutil.DerefBool(r.Properties.Active)
				rev.Replicas = azutil.DerefInt32(r.Properties.Replicas)
				if r.Properties.HealthState != nil {
					rev.HealthState = string(*r.Properties.HealthState)
				}
				if r.Properties.RunningState != nil {
					rev.RunningState = string(*r.Properties.RunningState)
				}
				rev.TrafficWeight = azutil.DerefInt32(r.Properties.TrafficWeight)
				if r.Properties.CreatedTime != nil {
					rev.CreatedTime = *r.Properties.CreatedTime
				}
			}
			revisions = append(revisions, rev)
		}
	}
	return revisions, nil
}

func (p *azureProvider) StreamLogs(ctx context.Context, rg, appName string) (<-chan provider.LogEntry, error) {
	cmd := exec.CommandContext(ctx,
		"az", "containerapp", "logs", "show",
		"--follow",
		"--name", appName,
		"--resource-group", rg,
		"--type", "console",
		"--output", "tsv",
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start log stream: %w", err)
	}

	ch := make(chan provider.LogEntry, 64)
	go func() {
		defer close(ch)
		defer cmd.Wait() //nolint:errcheck
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case ch <- provider.LogEntry{Timestamp: time.Now(), Message: scanner.Text()}:
			}
		}
	}()
	return ch, nil
}

func mapContainerAppSummary(a *armappcontainers.ContainerApp) provider.ContainerApp {
	app := provider.ContainerApp{
		ID:   azutil.DerefStr(a.ID),
		Name: azutil.DerefStr(a.Name),
	}
	if a.Location != nil {
		app.Region = *a.Location
	}
	app.ResourceGroup = azutil.ExtractRG(app.ID)
	if a.Properties != nil {
		if a.Properties.LatestRevisionName != nil {
			app.Revision = *a.Properties.LatestRevisionName
		}
		if a.Properties.RunningStatus != nil {
			app.State = string(*a.Properties.RunningStatus)
		}
	}
	return app
}
