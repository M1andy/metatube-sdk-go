package fingerprint

import "net/http"

// Mode represents the fingerprint simulation mode.
type Mode string

const (
	// ModeOff disables fingerprint simulation.
	ModeOff Mode = "off"

	// ModeHeader enables header-level browser fingerprint only.
	ModeHeader Mode = "header"

	// ModeUTLS enables uTLS (TLS ClientHello spoofing) + header fingerprint.
	// This is the default mode.
	ModeUTLS Mode = "utls"

	// ModeRod enables headless Chromium browser (requires experimental build tag).
	ModeRod Mode = "rod"
)

// Fingerprinter applies browser fingerprint simulation to HTTP requests.
type Fingerprinter interface {
	// ApplyRequest modifies an http.Request to include browser fingerprint headers.
	ApplyRequest(req *http.Request) error

	// Transport returns a wrapped http.RoundTripper that applies
	// TLS-level fingerprint spoofing. Returns nil if no wrapping is needed.
	Transport() http.RoundTripper

	// Mode returns the fingerprint mode.
	Mode() Mode
}

// DefaultFingerprinter is the package-level default fingerprint instance.
// It is initialized by the uTLS implementation (utls.go) to ModeUTLS.
// Falls back to ModeHeader if uTLS initialization fails.
var DefaultFingerprinter Fingerprinter
