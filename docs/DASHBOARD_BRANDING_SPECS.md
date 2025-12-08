# Dashboard Branding Customization - Detailed Specifications

This document provides detailed specifications for implementing the dashboard branding customization feature across multiple work streams.

## Overview

The feature allows users to customize the dashboard display name (replacing "Talos" with a custom name) through machine configuration. Changes can be applied without rebooting the system.

### Example Configuration

```yaml
machine:
  dashboard:
    branding:
      name: "MyCustomOS"
```

---

## Stream A: Config Schema Definition (COMPLETED)

**Status**: ✅ Implemented

### Files Modified

| File | Changes |
|------|---------|
| `pkg/machinery/config/types/v1alpha1/v1alpha1_types.go` | Added `DashboardConfig` and `DashboardBrandingConfig` structs, added `MachineDashboard` field to `MachineConfig` |
| `pkg/machinery/config/types/v1alpha1/v1alpha1_examples.go` | Added `machineDashboardExample()` and `dashboardBrandingExample()` functions |
| `pkg/machinery/config/config/machine.go` | Added `Dashboard()` method to `MachineConfig` interface, added `Dashboard` and `DashboardBranding` interfaces |
| `pkg/machinery/config/types/v1alpha1/v1alpha1_provider.go` | Added `Dashboard()` method implementation for `MachineConfig` |
| `pkg/machinery/config/types/v1alpha1/v1alpha1_dashboard.go` | NEW FILE - Interface implementations for `DashboardConfig` and `DashboardBrandingConfig` |

### Type Definitions

```go
// DashboardConfig describes the dashboard configuration.
type DashboardConfig struct {
    DashboardBranding *DashboardBrandingConfig `yaml:"branding,omitempty"`
}

// DashboardBrandingConfig describes the dashboard branding configuration.
type DashboardBrandingConfig struct {
    BrandingName string `yaml:"name,omitempty"`
}
```

---

## Stream B: Config Wiring

**Status**: ✅ Implemented in Stream A

The `Dashboard()` method has been added to the `MachineConfig` interface and implemented in the provider.

---

## Stream C: Make Config Reloadable (COMPLETED)

**Status**: ✅ Implemented

### Files Modified

| File | Changes |
|------|---------|
| `internal/app/machined/pkg/runtime/v1alpha1/v1alpha1_runtime.go` | Added `.machine.dashboard` to reloadable config list in `CanApplyImmediate()` |

---

## Stream D: Resource Controller

**Status**: 🔲 Pending (Optional)

If a resource-based approach is needed for watching dashboard config changes:

### Files to Create/Modify

| File | Purpose |
|------|---------|
| `pkg/machinery/resources/config/dashboard.go` | Define dashboard config resource type |
| `internal/app/machined/pkg/controllers/config/dashboard.go` | Controller to watch dashboard config changes and restart dashboard service |

### Implementation Notes

- This stream is optional if the dashboard can read branding directly from machine config at startup
- If the dashboard service needs to react to config changes without restart, a controller approach is needed

---

## Stream E: Dashboard Service Integration

**Status**: 🔲 Pending

### Files to Modify

| File | Changes |
|------|---------|
| `internal/pkg/dashboard/options.go` | Add `branding` field to options struct, add `WithBranding(name string)` option |
| `internal/pkg/dashboard/dashboard.go` | Pass branding to components, store in Dashboard struct |
| `internal/app/dashboard/main.go` | Read branding from machine config at startup, pass to `dashboard.Run()` |

### Detailed Changes

#### `internal/pkg/dashboard/options.go`

```go
type options struct {
    interval      time.Duration
    allowExitKeys bool
    screens       []Screen
    branding      string  // ADD THIS
}

func defaultOptions() *options {
    return &options{
        interval:      5 * time.Second,
        allowExitKeys: true,
        screens: []Screen{
            ScreenSummary,
            ScreenMonitor,
            ScreenNetworkConfig,
        },
        branding: "Talos",  // ADD THIS - default branding
    }
}

// WithBranding sets the custom branding name.
func WithBranding(name string) Option {
    return func(o *options) {
        if name != "" {
            o.branding = name
        }
    }
}
```

#### `internal/app/dashboard/main.go`

```go
func dashboardMain() error {
    // ... existing code ...

    // Read branding from machine config
    cfg, err := c.COSI.Get(ctx, resource.NewMetadata(
        machineconfig.NamespaceName,
        machineconfig.MachineConfigType,
        machineconfig.V1Alpha1ID,
        resource.VersionUndefined,
    ))

    var branding string
    if err == nil {
        if mc := cfg.(*machineconfig.MachineConfig); mc != nil {
            if dashboard := mc.Config().Machine().Dashboard(); dashboard != nil {
                if b := dashboard.Branding(); b != nil {
                    branding = b.Name()
                }
            }
        }
    }

    return dashboard.Run(ctx, c,
        dashboard.WithAllowExitKeys(false),
        dashboard.WithScreens(screens...),
        dashboard.WithBranding(branding),  // ADD THIS
    )
}
```

---

## Stream F: UI Component Updates

**Status**: 🔲 Pending

### Files to Modify

| File | Changes |
|------|---------|
| `internal/pkg/dashboard/components/header.go` | Accept branding parameter, display custom name |
| `internal/pkg/dashboard/components/talosinfo.go` | Accept branding parameter (future use for widget titles) |
| `internal/pkg/dashboard/components/footer.go` | Accept branding parameter (future use) |

### Detailed Changes

#### `internal/pkg/dashboard/components/header.go`

Add branding support to the header component:

```go
type Header struct {
    tview.TextView
    selectedNode string
    nodeMap      map[string]*headerData
    spinnerPos   int
    branding     string  // ADD THIS
}

func NewHeader(branding string) *Header {
    header := &Header{
        TextView: *tview.NewTextView(),
        nodeMap:  make(map[string]*headerData),
        branding: branding,  // ADD THIS
    }
    // ... rest of initialization
}

func (widget *Header) redraw() {
    data := widget.getOrCreateNodeData(widget.selectedNode)
    spinnerPos := widget.spinnerPos % len(spinner)

    // Use branding in display (if applicable)
    text := fmt.Sprintf(
        "[green]%s [yellow::b]%s[-:-:-] (%s): uptime %s, %s, %s RAM, PROCS %s, CPU %s, RAM %s",
        spinner[spinnerPos],
        data.hostname,
        data.version,  // Could optionally replace version display with branding
        data.uptime,
        data.cpuFreq,
        data.totalMem,
        data.numProcesses,
        data.cpuUsagePercent,
        data.memUsagePercent,
    )
    widget.SetText(text)
}
```

---

## Stream G: Talosctl Dashboard

**Status**: 🔲 Pending

### Files to Modify

| File | Changes |
|------|---------|
| `cmd/talosctl/cmd/talos/dashboard.go` | Add `--branding` flag for remote viewing |

### Implementation

```go
var dashboardCmdFlags struct {
    interval time.Duration
    branding string  // ADD THIS
}

func init() {
    dashboardCmd.Flags().DurationVar(&dashboardCmdFlags.interval, "update-interval", 3*time.Second, "interval between updates")
    dashboardCmd.Flags().StringVar(&dashboardCmdFlags.branding, "branding", "", "custom branding name to display")  // ADD THIS
}

// In the run function, pass branding to dashboard.Run()
```

---

## Stream H: Testing

**Status**: 🔲 Pending

### Test Files to Create/Modify

| File | Purpose |
|------|---------|
| `pkg/machinery/config/types/v1alpha1/v1alpha1_dashboard_test.go` | Unit tests for dashboard config types |
| `internal/pkg/dashboard/dashboard_test.go` | Unit tests for branding option |

### Test Cases

1. **Config Parsing**: Verify YAML config with branding is parsed correctly
2. **Default Value**: Verify empty branding defaults to "Talos"
3. **Interface Implementation**: Verify `DashboardConfig` implements `config.Dashboard`
4. **Reloadable Config**: Verify `.machine.dashboard` can be applied without reboot

---

## Stream I: Documentation

**Status**: 🔲 Pending

### Files to Create/Modify

| File | Purpose |
|------|---------|
| `website/content/*/reference/configuration/v1alpha1/config.md` | Document new config option |
| `website/content/*/talos-guides/configuration/dashboard.md` | Guide for dashboard customization |

---

## Dependencies Graph

```
Phase 1:        [A: Config Schema] ────────────────────────────────┐
                         │                                          │
                         ▼                                          │
Phase 2:    ┌──[B: Config Wiring]──┐                               │
            │                      │                                │
            │   [C: Reloadable]    │                               │
            │                      │                                │
            └──[D: Controller]─────┘                               │
                         │                                          │
                         ▼                                          │
Phase 3:    ┌──[E: Dashboard Service]──┐     [G: Talosctl]         │
            │                          │          │                 │
            └──[F: UI Components]──────┘          │                 │
                         │                        │                 │
                         ▼                        ▼                 │
Phase 4:         [H: Testing] ◄───────────────────┘                │
                         │                                          │
                         ▼                                          │
                  [I: Documentation] ◄─────────────────────────────┘
```

---

## Code Generation

After modifying `v1alpha1_types.go`, run:

```bash
cd pkg/machinery/config/types/v1alpha1
go generate ./...
```

This regenerates:
- `v1alpha1_types_doc.go` - Documentation
- `zz_generated.deepcopy.go` - DeepCopy methods
