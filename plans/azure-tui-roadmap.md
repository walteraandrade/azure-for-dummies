# Azure TUI — Roadmap

## Concept

Full-screen terminal interface for Azure operations. No commands to memorize.
Launch binary → navigate services → inspect and modify resources.

---

## Landscape

| Project | Stack | Scope | Status |
|---|---|---|---|
| [az-tui](https://github.com/IAL32/az-tui) | Go + BubbleTea | Container Apps only | v0.4.0, ~16 stars |
| [k9s](https://k9scli.io/) | Go + BubbleTea | Kubernetes | Gold standard UX |
| [aws-tui](https://github.com/rfc2119/aws-tui) | Python | AWS (WIP) | Abandoned |

**Differentiators:** broader Azure scope, "dummies" UX (inline help, no CLI knowledge needed), module plugin system.

---

## Tech Stack

| Layer | Choice | Rationale |
|---|---|---|
| Language | Go | Goroutines fit log streaming; single binary dist |
| TUI | BubbleTea (Elm arch) | Battle-tested for this exact pattern; user familiar |
| Azure auth | `azidentity.DefaultAzureCredential` | Picks up `az login` automatically |
| Azure data | `azure-sdk-for-go` (GA, 2025) | Typed clients for all services |
| Styling | Lip Gloss | BubbleTea companion for layout/colors |

---

## Navigation Model

```
Home Screen
│  [Container Apps]   [PostgreSQL]   [Storage]   [+ more modules]
│
├─ Container Apps List
│   name · state · resource group · region · revision
│   │
│   └─ App Detail
│       ├─ Overview   (ingress, env vars, resource limits)
│       ├─ Revisions  (traffic split, active/inactive)
│       ├─ Logs       (streaming, filterable)
│       └─ Settings   (edit env vars, scaling rules)
│
├─ PostgreSQL List
│   name · tier · state · region
│   │
│   └─ DB Detail
│       ├─ Overview   (connection string, firewall)
│       ├─ Metrics    (CPU, mem, connections)
│       └─ Settings
│
└─ Storage (future)
```

**Stack-based router:** `enter` pushes, `esc`/`backspace` pops. Global `q` to quit.

---

## Architecture: Module Registry (Strategy pattern)

Each Azure service = a **Module** implementing one interface.
Home screen renders whatever modules are registered — adding a new service = zero changes to existing code.

```go
type Module interface {
    Name()                      string
    Icon()                      string
    Fetch(ctx context.Context)  tea.Cmd
    ListView()                  tea.Model
    DetailView(id string)       tea.Model
}
```

Router maintains a `[]tea.Model` stack. Modules register themselves at startup.

---

## Sprints

### Sprint 0 — Foundation `~1 week`
- [ ] Go module scaffold, BubbleTea entry point
- [ ] Stack-based screen router
- [ ] `DefaultAzureCredential` auth setup
- [ ] Module registry
- [ ] Mock data layer (dev without live Azure)
- [ ] Status bar: subscription / user / last refresh

### Sprint 1 — MVP: Home + Container Apps List `~2 weeks`
- [ ] Home screen: service grid from registered modules
- [ ] Container Apps module: list view (name, state, RG, region, revision)
- [ ] Keyboard nav: arrows, `enter`, `esc`, `q`, `/` filter
- [ ] Subscription picker on launch (if multiple)

**Deliverable: launch binary → see your container apps**

### Sprint 2 — Container Apps Detail `~1 week`
- [ ] Detail view with tabs: Overview / Revisions / Logs / Settings
- [ ] Log streaming (goroutine → `tea.Cmd` → scrollable viewport)
- [ ] Env var viewer

### Sprint 3 — PostgreSQL Module `~1 week`
- [ ] List view: name, tier, state, region
- [ ] Detail view: connection info, firewall rules, metrics

### Sprint 4 — Storage Module `~1 week`
- [ ] Blob containers list + detail view

### Backlog
- Edit settings in-place (env vars, scaling rules, traffic splits)
- Global fuzzy search across all resources
- Multi-subscription support
- Cost / billing overlay
- Azure Container Registry module
- AKS module

---

## Open Questions

1. Binary name? (`azf`, `azt`, `azm`...)
2. Auth: `az login` only, or also service principal via env vars?
3. Detail: tabs within a screen, or drill-down sub-screens?
4. Multi-subscription: MVP or Sprint 3+?
