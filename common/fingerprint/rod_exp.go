//go:build experimental

package fingerprint

import (
	"net/http"
)

// Importers: none by default (experimental build tag).
// Requires go-rod/rod + go-rod/stealth (not in default go.mod).
// No data files.

var _ Fingerprinter = (*RodFingerprinter)(nil)

// RodFingerprinter uses a headless Chromium browser for full
// browser fingerprint emulation. Requires Chrome/Chromium installed.
type RodFingerprinter struct {
	browserPath string
}

// NewRodFingerprinter creates a RodFingerprinter.
func NewRodFingerprinter(browserPath string) *RodFingerprinter {
	return &RodFingerprinter{browserPath: browserPath}
}

func (f *RodFingerprinter) ApplyRequest(req *http.Request) error { return nil }

func (f *RodFingerprinter) Transport() http.RoundTripper { return nil }

func (f *RodFingerprinter) Mode() Mode { return ModeRod }
