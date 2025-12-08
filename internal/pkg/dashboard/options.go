// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package dashboard

import (
	"time"
)

// DefaultBranding is the default branding name shown in the dashboard.
const DefaultBranding = "Talos"

type options struct {
	interval      time.Duration
	allowExitKeys bool
	screens       []Screen
	branding      string
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
		branding: DefaultBranding,
	}
}

// Option is a functional option for Dashboard.
type Option func(*options)

// WithInterval sets the interval for the dashboard.
func WithInterval(interval time.Duration) Option {
	return func(o *options) {
		o.interval = interval
	}
}

// WithAllowExitKeys sets whether the dashboard should allow exit keys (Ctrl + C).
func WithAllowExitKeys(allowExitKeys bool) Option {
	return func(o *options) {
		o.allowExitKeys = allowExitKeys
	}
}

// WithScreens sets the screens to display.
// The order is preserved.
func WithScreens(screens ...Screen) Option {
	return func(o *options) {
		o.screens = screens
	}
}

// WithBranding sets the custom branding name to display in the dashboard.
// If empty, the default "Talos" branding is used.
func WithBranding(name string) Option {
	return func(o *options) {
		if name != "" {
			o.branding = name
		}
	}
}
