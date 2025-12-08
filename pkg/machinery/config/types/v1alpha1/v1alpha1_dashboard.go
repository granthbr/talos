// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha1

import (
	"github.com/siderolabs/talos/pkg/machinery/config/config"
)

var (
	_ config.Dashboard        = (*DashboardConfig)(nil)
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
