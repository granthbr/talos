# Dashboard Branding - Agent Prompts for Parallel Execution

This document contains ready-to-use prompts for running multiple agents in parallel to implement the dashboard branding feature.

## Prerequisites

Before running the agent prompts, ensure Phase 1 (Config Schema) is complete. The following files should already be modified:
- `pkg/machinery/config/types/v1alpha1/v1alpha1_types.go`
- `pkg/machinery/config/types/v1alpha1/v1alpha1_examples.go`
- `pkg/machinery/config/config/machine.go`
- `pkg/machinery/config/types/v1alpha1/v1alpha1_provider.go`
- `pkg/machinery/config/types/v1alpha1/v1alpha1_dashboard.go` (new file)
- `internal/app/machined/pkg/runtime/v1alpha1/v1alpha1_runtime.go`

---

## Agent 1: Dashboard Options & Service Integration

**Focus**: Stream E - Dashboard service integration

```
You are implementing dashboard branding support in the Talos dashboard.

## Task
Modify the dashboard package to support custom branding through options.

## Files to Modify

### 1. `internal/pkg/dashboard/options.go`
Add branding support:
- Add `branding string` field to the `options` struct
- Set default value to "Talos" in `defaultOptions()`
- Create `WithBranding(name string) Option` function

### 2. `internal/pkg/dashboard/dashboard.go`
- Add `branding string` field to the `Dashboard` struct
- Store branding from options in the Dashboard struct during initialization
- Pass branding to components that need it (Header, TalosInfo)

### 3. `internal/app/dashboard/main.go`
Read branding from machine config:
- Import necessary packages for reading machine config
- After creating the client, read the machine config
- Extract branding name from `config.Machine().Dashboard().Branding().Name()`
- Pass branding to `dashboard.Run()` using `dashboard.WithBranding(branding)`

## Implementation Pattern
Follow existing patterns in the codebase:
- Look at `WithInterval()` and `WithScreens()` for option patterns
- The branding should default to "Talos" if not configured

## Do NOT
- Modify any config/machinery files (already done)
- Create new resource types
- Add tests (separate agent handles this)

Return a summary of changes made and any issues encountered.
```

---

## Agent 2: UI Components Update

**Focus**: Stream F - Dashboard UI components

```
You are updating Talos dashboard UI components to display custom branding.

## Task
Modify dashboard components to accept and display custom branding.

## Files to Modify

### 1. `internal/pkg/dashboard/components/header.go`
- Add `branding string` field to the `Header` struct
- Modify `NewHeader()` to accept branding parameter: `NewHeader(branding string) *Header`
- Store branding in the struct
- Optionally use branding in the `redraw()` method display

### 2. `internal/pkg/dashboard/components/talosinfo.go`
- Add `branding string` field to the `TalosInfo` struct
- Modify `NewTalosInfo()` to accept branding parameter: `NewTalosInfo(branding string) *TalosInfo`
- Store branding for potential future use in widget display

### 3. Update callers in `internal/pkg/dashboard/`
After modifying the component constructors, update all places that create these components:
- `summary.go` - Update `NewHeader()` and `NewTalosInfo()` calls
- `dashboard.go` - Update component creation if done there
- Any other files that instantiate these components

## Pattern to Follow
Look at how other parameters are passed to components. The branding should flow from:
1. `dashboard.Run()` options
2. Dashboard struct
3. Component constructors

## Do NOT
- Modify config/machinery files
- Add complex branding logic yet (just wire through the name)
- Add tests

Return a summary of changes made and the component call chain.
```

---

## Agent 3: Talosctl Dashboard Command

**Focus**: Stream G - CLI integration

```
You are adding branding support to the talosctl dashboard command.

## Task
Add a `--branding` flag to the talosctl dashboard command for custom branding when viewing dashboards remotely.

## Files to Modify

### 1. `cmd/talosctl/cmd/talos/dashboard.go`
- Add `branding string` field to `dashboardCmdFlags` struct
- Register the flag in `init()`: `--branding` with empty default and description "custom branding name to display"
- In the dashboard run function, pass branding to `dashboard.Run()` using `dashboard.WithBranding()`
- If branding flag is empty, try to read from the target node's config (optional enhancement)

## Implementation Example

```go
var dashboardCmdFlags struct {
    interval time.Duration
    branding string
}

func init() {
    dashboardCmd.Flags().DurationVar(&dashboardCmdFlags.interval, "update-interval", 3*time.Second, "interval between updates")
    dashboardCmd.Flags().StringVar(&dashboardCmdFlags.branding, "branding", "", "custom branding name to display instead of 'Talos'")
    // ... rest of init
}
```

## Pattern to Follow
Look at existing flags like `--update-interval` for the pattern.

## Do NOT
- Modify the dashboard package itself
- Add tests
- Modify machine config types

Return a summary of the flag implementation.
```

---

## Agent 4: Code Generation & Validation

**Focus**: Regenerate code and validate

```
You are running code generation and validating the dashboard branding implementation.

## Task
Run code generation for the modified config types and verify everything compiles.

## Steps

### 1. Run Code Generation
```bash
cd /home/user/talos/pkg/machinery/config/types/v1alpha1
go generate ./...
```

### 2. Verify Compilation
```bash
cd /home/user/talos
go build ./...
```

### 3. Check for Missing Implementations
Look for any interface methods that need to be implemented:
- Check that `DashboardConfig` implements `config.Dashboard`
- Check that `DashboardBrandingConfig` implements `config.DashboardBranding`

### 4. Verify Deep Copy Generation
Check that `zz_generated.deepcopy.go` includes:
- `DeepCopyInto` for `DashboardConfig`
- `DeepCopyInto` for `DashboardBrandingConfig`

## If Compilation Fails
- Report the specific errors
- Identify which files need fixes
- Do NOT attempt to fix - just report

Return:
1. Code generation output
2. Build output (success or errors)
3. List of any missing implementations
```

---

## Agent 5: Testing Implementation

**Focus**: Stream H - Unit tests

```
You are implementing tests for the dashboard branding feature.

## Task
Create unit tests for the new dashboard branding configuration.

## Files to Create

### 1. `pkg/machinery/config/types/v1alpha1/v1alpha1_dashboard_test.go`

```go
package v1alpha1_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"
)

func TestDashboardConfig_Branding(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name     string
        config   *v1alpha1.DashboardConfig
        expected string
    }{
        {
            name:     "nil config returns empty branding",
            config:   nil,
            expected: "",
        },
        {
            name:     "nil branding returns empty name",
            config:   &v1alpha1.DashboardConfig{},
            expected: "",
        },
        {
            name: "custom branding name",
            config: &v1alpha1.DashboardConfig{
                DashboardBranding: &v1alpha1.DashboardBrandingConfig{
                    BrandingName: "MyCustomOS",
                },
            },
            expected: "MyCustomOS",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()

            var branding string
            if tt.config != nil {
                branding = tt.config.Branding().Name()
            }
            assert.Equal(t, tt.expected, branding)
        })
    }
}

func TestMachineConfig_Dashboard(t *testing.T) {
    t.Parallel()

    mc := &v1alpha1.MachineConfig{
        MachineDashboard: &v1alpha1.DashboardConfig{
            DashboardBranding: &v1alpha1.DashboardBrandingConfig{
                BrandingName: "TestOS",
            },
        },
    }

    dashboard := mc.Dashboard()
    require.NotNil(t, dashboard)

    branding := dashboard.Branding()
    require.NotNil(t, branding)

    assert.Equal(t, "TestOS", branding.Name())
}
```

### 2. Test YAML Parsing
Add a test to verify the config can be parsed from YAML:

```go
func TestDashboardConfig_YAML(t *testing.T) {
    t.Parallel()

    yamlConfig := `
machine:
  dashboard:
    branding:
      name: "MyBrand"
`
    // Parse and verify
}
```

## Run Tests
```bash
go test -v ./pkg/machinery/config/types/v1alpha1/... -run Dashboard
```

Return test file contents and test results.
```

---

## Parallel Execution Guide

### Phase 2 (Run in Parallel)
After Phase 1 is complete, run these agents simultaneously:
- **Agent 1**: Dashboard Options & Service Integration
- **Agent 2**: UI Components Update
- **Agent 3**: Talosctl Dashboard Command

### Phase 3 (Sequential)
After Phase 2:
- **Agent 4**: Code Generation & Validation

### Phase 4 (After Validation)
- **Agent 5**: Testing Implementation

---

## Quick Start Commands

To spawn all Phase 2 agents in parallel using Claude Code:

```
Run these three tasks in parallel:
1. Implement dashboard branding options in internal/pkg/dashboard/options.go and wire through main.go
2. Update header.go and talosinfo.go components to accept branding parameter
3. Add --branding flag to cmd/talosctl/cmd/talos/dashboard.go
```

---

## Expected Final Config Example

After all agents complete, this config should work:

```yaml
version: v1alpha1
machine:
  type: controlplane
  dashboard:
    branding:
      name: "MyCustomOS"
cluster:
  # ... cluster config
```

And the dashboard should display "MyCustomOS" instead of "Talos" in the header.
