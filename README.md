# azure-for-dummies

A full-screen terminal UI for browsing your Azure resources. Think k9s, but for Azure.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go) ![Azure](https://img.shields.io/badge/Azure-SDK-0078D4?logo=microsoftazure) ![License](https://img.shields.io/badge/License-MIT-green)

## What it does

Launch a single command and get a navigable dashboard of your Azure subscription. Browse Container Apps, PostgreSQL Flexible Servers, and Storage Accounts — complete with detail views, live log streaming, metrics sparklines, and firewall rules.

No YAML. No JSON piping. Just arrow keys.

## Modules

| Module | List | Detail tabs |
|--------|------|-------------|
| **Container Apps** | All apps in subscription | Overview · Revisions · Logs (live stream) · Settings (env vars) |
| **PostgreSQL** | Flexible Servers | Overview · Firewall Rules · Metrics (CPU, memory, connections) |
| **Storage** | Storage Accounts | Overview · Blob Containers |

## Requirements

- Go 1.25+
- [Azure CLI](https://learn.microsoft.com/en-us/cli/azure/install-azure-cli) (`az`) installed and logged in
- An active Azure subscription

## Installation

**From source:**

```bash
git clone https://github.com/walteraandrade/azure-for-dummies.git
cd azure-for-dummies
go build -o azfd .
./azfd
```

**Go install:**

```bash
go install github.com/walteraandrade/azure-for-dummies@latest
```

## Quick start

```bash
# Make sure you're logged in
az login

# Run it
./azfd
```

If you have multiple subscriptions, a picker appears. Otherwise it drops you straight into the home screen.

## Navigation

| Key | Action |
|-----|--------|
| `←` `→` / `h` `l` | Navigate modules on home screen |
| `Enter` | Open module / select resource |
| `Esc` / `Backspace` | Go back |
| `Tab` / `Shift+Tab` | Switch detail tabs |
| `1`-`4` | Jump to tab by number |
| `/` | Filter lists |
| `q` | Quit |

## Architecture

Stack-based router with a module registry. Each Azure service is a self-contained module implementing a common interface — add a new service by dropping in a new package.

See [docs/architecture.md](docs/architecture.md) for the full breakdown.

## Docs

- [Architecture](docs/architecture.md) — patterns, router, module system
- [Modules](docs/modules.md) — per-module details and Azure SDK usage
- [Contributing](docs/contributing.md) — dev setup, adding new modules
- [Roadmap](plans/azure-tui-roadmap.md) — sprint plan and open questions

## License

MIT
