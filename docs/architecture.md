# Architecture

## Overview

azure-for-dummies is a [BubbleTea](https://github.com/charmbracelet/bubbletea) application following the Elm architecture (Model-View-Update). The app is structured around a few core abstractions that keep modules decoupled from navigation and shared UI.

## Boot sequence

```
main.go
  └─ app.New()
       ├─ router.New(loading.New())   ← initial screen
       ├─ statusbar.New()
       └─ Init() → auth.ListSubscriptions()
                      │
                      ├─ 1 subscription → auth.ResolveWithSubscription()
                      └─ N subscriptions → subscriptionpicker screen
                                              │
                                              └─ AuthReadyMsg
                                                   ├─ register modules
                                                   └─ router.ReplaceRoot(home)
```

## Core patterns

### Router (Mediator)

Stack-based navigation. Screens never reference the router directly — they emit `PushMsg` or `PopMsg` and the router handles the rest.

```
router.go
  Push(screen)  → appends to stack
  Pop()         → removes top
  ReplaceRoot() → clears stack, sets new base
```

This keeps every screen testable in isolation.

### Module interface (Strategy)

```go
type Module interface {
    Name() string
    Icon() string
    ListView() tea.Model
    DetailView(id string) tea.Model
}
```

Each Azure service implements this. The home screen iterates over registered modules and calls `ListView()` on selection. The list view calls `DetailView(id)` when a resource is selected.

### Registry (Plugin)

```go
registry.Register(containerapps.New(ctx))
registry.Register(postgres.New(ctx))
registry.Register(storage.New(ctx))
```

Modules register after auth completes. `registry.All()` returns them in insertion order for the home screen.

### Provider (Repository)

Each module defines a provider interface for its Azure SDK calls:

```go
type ContainerAppsProvider interface {
    ListContainerApps(ctx context.Context) ([]ContainerApp, error)
    GetContainerApp(ctx context.Context, rg, name string) (ContainerApp, error)
    ListRevisions(ctx context.Context, rg, appName string) ([]RevisionInfo, error)
    StreamLogs(ctx context.Context, rg, appName string) (<-chan LogEntry, error)
}
```

The concrete implementation lives in `provider.go` inside each module package. Domain structs live in `internal/provider/` — shared across the app, decoupled from Azure SDK types.

## Package layout

```
main.go                         # entry point
internal/
  app/                          # root model — owns router + statusbar
  auth/                         # az CLI auth, subscription listing
  router/                       # stack-based screen navigation
  module/                       # Module interface + Registry
  provider/                     # shared domain structs + provider interfaces
  home/                         # home screen — module card grid
  statusbar/                    # bottom bar (subscription, user, errors)
  styles/                       # Catppuccin Mocha palette + lipgloss styles
  tabs/                         # reusable tab bar component
  loading/                      # loading spinner screen
  subscriptionpicker/           # multi-subscription selector
  azutil/                       # helpers (deref pointers, extract RG from ID)
  containerapps/                # Container Apps module
  postgres/                     # PostgreSQL module
  storage/                      # Storage module
```

## Data flow

```
Azure SDK → provider (concrete) → domain structs → BubbleTea model → lipgloss view
```

All Azure SDK types are mapped to plain Go structs in the provider layer. Views never import Azure SDK packages.

## Styling

Catppuccin Mocha palette defined in `styles/styles.go`. All colors and component styles are lipgloss variables — change the palette in one place.

## Auth

Uses `az account list` and `az account show` via CLI subprocess, then creates `azidentity.DefaultAzureCredential` for SDK calls. This means any auth method supported by `az login` works — device code, browser, service principal, managed identity.
