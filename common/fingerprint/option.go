package fingerprint

// Option configures a Fingerprinter.
type Option func(any)

// WithProfileRotation enables or disables random profile rotation
// for each request. When enabled, a new BrowserProfile is selected
// for each ApplyRequest call.
func WithProfileRotation(rotate bool) Option {
	return func(v any) {
		switch f := v.(type) {
		case *HeaderFingerprinter:
			f.rotate = rotate
		case *UTLSFingerprinter:
			f.rotate = rotate
		}
	}
}

// WithBrowser selects a specific browser for the fingerprint profile.
// Supported values: "chrome", "firefox", "edge", "safari".
// An empty string or "random" selects a random profile.
func WithBrowser(browser string) Option {
	return func(v any) {
		switch f := v.(type) {
		case *HeaderFingerprinter:
			f.browser = browser
		case *UTLSFingerprinter:
			f.browser = browser
		}
	}
}
