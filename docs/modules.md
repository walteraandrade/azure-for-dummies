# Modules

Each module is a self-contained package under `internal/` that implements the `Module` interface. All follow the same structure: module entry point, provider, list view, detail view with tabs.

## Container Apps

**Package:** `internal/containerapps/`

**Provider calls:**
- `ListContainerApps` — paginates `armappcontainers.ContainerAppsClient`
- `GetContainerApp` — fetches full detail including ingress, scale, containers
- `ListRevisions` — paginates `armappcontainers.ContainerAppsRevisionsClient`
- `StreamLogs` — shells out to `az containerapp logs show --follow`, streams lines over a channel

**Detail tabs:**

| Tab | Content |
|-----|---------|
| Overview | Name, state, FQDN, ingress type/port, scale range, container images + resources |
| Revisions | List of revisions with active/health/running state, traffic weight, replica count |
| Logs | Live-streamed console logs in a scrollable viewport. Activated lazily on tab switch. Canceled on back navigation. |
| Settings | Environment variables per container. Secret refs are masked with `••••••`. |

**Note:** Log streaming requires the Azure CLI — it's the only provider call that doesn't use the SDK directly.

## PostgreSQL

**Package:** `internal/postgres/`

**Provider calls:**
- `ListServers` — paginates `armpostgresqlflexibleservers.ServersClient`
- `GetServer` — fetches server detail (SKU, storage, FQDN, backup retention)
- `ListFirewallRules` — paginates `armpostgresqlflexibleservers.FirewallRulesClient`
- `GetMetrics` — queries `azquery.MetricsClient` for `cpu_percent`, `memory_percent`, `active_connections` (last 1h, 5min intervals)

**Detail tabs:**

| Tab | Content |
|-----|---------|
| Overview | Server properties, version, SKU/tier, storage, backup retention, connection string |
| Firewall | Table of firewall rules with name, start IP, end IP |
| Metrics | Sparkline charts (Unicode block characters) showing CPU%, memory%, active connections with current values |

## Storage

**Package:** `internal/storage/`

**Provider calls:**
- `ListStorageAccounts` — paginates `armstorage.AccountsClient`
- `GetStorageAccount` — fetches full properties (kind, SKU, access tier, TLS, HNS, network rules, endpoints)
- `ListBlobContainers` — paginates `armstorage.BlobContainersClient`

**Detail tabs:**

| Tab | Content |
|-----|---------|
| Overview | Account properties, endpoints (blob, file, table, queue), network/security settings |
| Containers | Table of blob containers with name, public access level, lease status, last modified |

## Adding a new module

1. Create a package under `internal/yourservice/`
2. Define a provider interface in `internal/provider/provider.go`
3. Implement `module.Module`:
   - `Name()` and `Icon()` for the home screen card
   - `ListView()` returns a `tea.Model` with a `bubbles/list`
   - `DetailView(id)` returns a `tea.Model` with tabs
4. Register in `app.go` inside the `AuthReadyMsg` handler

The list view should follow the existing pattern: spinner on load, fetch via provider, populate `bubbles/list`, push detail view on enter. The detail view should use the shared `tabs.Model` component.
