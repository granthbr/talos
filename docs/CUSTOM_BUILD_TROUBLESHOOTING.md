# Custom Talos Build Troubleshooting Guide

This document describes the issues encountered while building a custom Talos fork with dashboard branding changes, and the solutions applied.

## Overview

When adding custom dashboard branding configuration to Talos, we encountered several build issues related to code generation, Docker BuildKit compatibility, and version validation.

## Issue 1: Go Version Mismatch in `make generate`

### Symptom
```
go: go.mod requires go >= 1.25.3 (running go 1.23.x)
```

The `make generate` command failed because the code generation tools require a specific Go toolchain version that wasn't available in the build container.

### Root Cause
Talos uses `go generate` directives that invoke tools like `docgen` and `deepcopy-gen`. These tools have strict Go version requirements specified in `go.mod`.

### Solution
We bypassed the code generation requirement by:

1. **Temporarily disabled docgen** in `pkg/machinery/config/types/v1alpha1/v1alpha1_types.go`:
   ```go
   // NOTE: docgen temporarily disabled for custom branding build
   ////go:generate go tool github.com/siderolabs/talos/tools/docgen ...
   ```

2. **Added manual DeepCopy methods** in `pkg/machinery/config/types/v1alpha1/v1alpha1_dashboard.go`:
   ```go
   func (in *DashboardConfig) DeepCopyInto(out *DashboardConfig) {
       *out = *in
       if in.DashboardBranding != nil {
           in, out := &in.DashboardBranding, &out.DashboardBranding
           *out = new(DashboardBrandingConfig)
           (*in).DeepCopyInto(*out)
       }
   }

   func (in *DashboardConfig) DeepCopy() *DashboardConfig {
       if in == nil {
           return nil
       }
       out := new(DashboardConfig)
       in.DeepCopyInto(out)
       return out
   }
   ```

   Similar methods were added for `DashboardBrandingConfig`.

This allows the build to proceed without running `make generate`.

---

## Issue 2: Docker BuildKit `rewrite-timestamp` Conflict

### Symptom
```
ERROR: failed to build: failed to solve: exporter option "rewrite-timestamp" conflicts with "unpack"
```

### Root Cause
The Makefile's `registry-%` target uses `rewrite-timestamp=true` for reproducible builds. This option conflicts with the `unpack` behavior in newer Docker Desktop/BuildKit versions.

### Solution
Removed the `rewrite-timestamp=true` option from the Makefile (line 332):

**Before:**
```makefile
registry-%:
    @$(MAKE) target-$* TARGET_ARGS="--output type=image,name=...,rewrite-timestamp=true $(TARGET_ARGS)"
```

**After:**
```makefile
registry-%:
    @$(MAKE) target-$* TARGET_ARGS="--output type=image,name=... $(TARGET_ARGS)"
```

### Trade-off
This removes reproducible build timestamps, which is acceptable for custom builds. For official releases requiring reproducibility, use `push=true` which disables unpack:
```makefile
--output type=image,name=...,rewrite-timestamp=true,push=true
```

---

## Issue 3: Invalid Semver Version for Extensions

### Symptom
```
error validating extension "nebula": Invalid character(s) found in major number "1a4484a0a"
```

### Root Cause
When building without specifying a `TAG`, the build system uses the git commit SHA (e.g., `1a4484a0a`) as the version. System extensions require a valid semver version for compatibility validation.

### Solution
Build with an explicit `TAG` variable:

```bash
make imager TAG=v1.11.5-custom
```

Then use the properly tagged imager:
```bash
docker run --rm -t -v $PWD/_out:/out \
  ghcr.io/siderolabs/imager:v1.11.5-custom \
  installer \
  --system-extension-image ghcr.io/siderolabs/nebula:1.9.6 \
  --system-extension-image ghcr.io/siderolabs/nonfree-kmod-nvidia-production:570.172.08-v1.11.5 \
  --system-extension-image ghcr.io/siderolabs/nvidia-container-toolkit-production:570.172.08-v1.17.8
```

---

## Complete Build Process

After applying all fixes, the build process is:

```bash
# 1. Skip make generate (manual DeepCopy methods are already in place)

# 2. Build the imager with a proper version tag
make imager TAG=v1.11.5-custom

# 3. Generate the installer with extensions
docker run --rm -t -v $PWD/_out:/out \
  ghcr.io/siderolabs/imager:v1.11.5-custom \
  installer \
  --system-extension-image ghcr.io/siderolabs/nebula:1.9.6 \
  --system-extension-image ghcr.io/siderolabs/nonfree-kmod-nvidia-production:570.172.08-v1.11.5 \
  --system-extension-image ghcr.io/siderolabs/nvidia-container-toolkit-production:570.172.08-v1.17.8
```

---

## Commits Applied

1. `6605fb5` - feat: add dashboard branding configuration (Phase 1)
2. `7f5f3d6` - feat: wire dashboard branding through UI components (Phase 2)
3. `8fcd959` - fix: correct safe.StateGetByID usage in dashboard main
4. `0dc2b68` - build: temporarily disable docgen for custom build
5. `1585dd6` - fix: add manual DeepCopy methods for dashboard config types
6. `1a4484a` - build: remove rewrite-timestamp option for Docker compatibility

---

## Applying Dashboard Branding

Once deployed, apply the branding configuration to your cluster:

```yaml
machine:
  dashboard:
    branding:
      name: "YourCustomOS"
```

This can be patched onto existing nodes without requiring a reboot.
