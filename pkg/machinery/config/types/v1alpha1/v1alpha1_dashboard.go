// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha1

import (
	"github.com/siderolabs/talos/pkg/machinery/config/config"
)

var (
	_ config.Dashboard         = (*DashboardConfig)(nil)
	_ config.DashboardBranding = (*DashboardBrandingConfig)(nil)
)

// Branding implements config.Dashboard interface.
func (d *DashboardConfig) Branding() config.DashboardBranding {
	if d.DashboardBranding == nil {
		return &DashboardBrandingConfig{}
	}

	return d.DashboardBranding
}

// Name implements config.DashboardBranding interface.
func (b *DashboardBrandingConfig) Name() string {
	return b.BrandingName
}

// DeepCopyInto is a deepcopy function, copying the receiver, writing into out.
func (in *DashboardConfig) DeepCopyInto(out *DashboardConfig) {
	*out = *in
	if in.DashboardBranding != nil {
		in, out := &in.DashboardBranding, &out.DashboardBranding
		*out = new(DashboardBrandingConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is a deepcopy function, copying the receiver, creating a new DashboardConfig.
func (in *DashboardConfig) DeepCopy() *DashboardConfig {
	if in == nil {
		return nil
	}
	out := new(DashboardConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is a deepcopy function, copying the receiver, writing into out.
func (in *DashboardBrandingConfig) DeepCopyInto(out *DashboardBrandingConfig) {
	*out = *in
}

// DeepCopy is a deepcopy function, copying the receiver, creating a new DashboardBrandingConfig.
func (in *DashboardBrandingConfig) DeepCopy() *DashboardBrandingConfig {
	if in == nil {
		return nil
	}
	out := new(DashboardBrandingConfig)
	in.DeepCopyInto(out)
	return out
}
