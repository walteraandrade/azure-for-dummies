# Contributing

## Dev setup

```bash
git clone https://github.com/walteraandrade/azure-for-dummies.git
cd azure-for-dummies
go mod download
go build -o azfd .
```

You need a real Azure subscription — there's no mock layer. Make sure `az login` works before running.

## Running

```bash
az login
./azfd
```

## Project structure

Everything lives under `internal/`. Each Azure service is a self-contained module package. Shared types are in `provider/`, shared styles in `styles/`, shared components in `tabs/`.

See [architecture.md](architecture.md) for the full breakdown.

## Adding a new Azure service

1. Add domain structs + provider interface to `internal/provider/provider.go`
2. Create `internal/yourservice/` with:
   - `module.go` — implements `module.Module`
   - `provider.go` — Azure SDK calls, maps to domain structs
   - `listview.go` — `bubbles/list` with spinner, fetch, and enter-to-detail
   - `detailview.go` — `tabs.Model` with overview + service-specific tabs
   - `overview.go` — key-value display of resource properties
3. Register the module in `internal/app/app.go` inside `AuthReadyMsg`
4. Add the Azure SDK dependency: `go get github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/...`

Follow the existing modules as templates — they all share the same patterns.

## Code style

- No mocks — all builds against real Azure
- Domain structs in `provider/`, never import Azure SDK types outside the provider layer
- Pointer derefs go through `azutil.DerefStr`, `azutil.DerefBool`, etc.
- Resource group and name extracted from Azure resource IDs via `azutil.ExtractRG`, `azutil.ExtractName`
- Catppuccin Mocha palette only — all colors defined in `styles/styles.go`

## Tests

```bash
go test ./...
```

Currently covers `azutil` and `router` packages.
