package fingerprint

import (
	"net/http"
	"strings"

	"github.com/metatube-community/metatube-sdk-go/common/random"
)

// Importers: utls.go (composition), scraper/option.go, fetch/fetch.go.
// Pure HTTP header injector. No data files.

var _ Fingerprinter = (*HeaderFingerprinter)(nil)

// HeaderFingerprinter implements Fingerprinter by injecting browser-like
// HTTP headers into each request. It selects a BrowserProfile at creation
// time and applies its headers to every request.
type HeaderFingerprinter struct {
	profile BrowserProfile
	rotate  bool
	browser string
}

// NewHeaderFingerprinter creates a HeaderFingerprinter.
func NewHeaderFingerprinter(opts ...Option) *HeaderFingerprinter {
	f := &HeaderFingerprinter{
		profile: resolveProfile(""),
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

// ApplyRequest injects browser headers into the request.
func (f *HeaderFingerprinter) ApplyRequest(req *http.Request) error {
	p := f.profile
	if f.rotate {
		p = resolveProfile(f.browser)
	}

	// User-Agent.
	if p.UserAgent != "" {
		req.Header.Set("User-Agent", p.UserAgent)
	}
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", random.UserAgent())
	}

	// Standard browser headers.
	setIfNotEmpty(req.Header, "Accept", p.Accept)
	setIfNotEmpty(req.Header, "Accept-Language", p.AcceptLanguage)
	setIfNotEmpty(req.Header, "Accept-Encoding", p.AcceptEncoding)
	setIfNotEmpty(req.Header, "Cache-Control", p.CacheControl)

	// Client hints (Chromium-based browsers).
	setIfNotEmpty(req.Header, "Sec-CH-UA", p.SecCHUA)
	setIfNotEmpty(req.Header, "Sec-CH-UA-Platform", p.SecCHUAPlatform)

	// Fetch metadata headers.
	if p.SecFetchDest != "" {
		req.Header.Set("Sec-Fetch-Dest", p.SecFetchDest)
	}
	if p.SecFetchMode != "" {
		req.Header.Set("Sec-Fetch-Mode", p.SecFetchMode)
	}
	if p.SecFetchSite != "" {
		req.Header.Set("Sec-Fetch-Site", p.SecFetchSite)
	}

	// Common browser behaviors.
	if req.Header.Get("Upgrade-Insecure-Requests") == "" {
		req.Header.Set("Upgrade-Insecure-Requests", "1")
	}
	if req.Header.Get("DNT") == "" {
		req.Header.Set("DNT", "1")
	}

	return nil
}

func (f *HeaderFingerprinter) Transport() http.RoundTripper { return nil }

func (f *HeaderFingerprinter) Mode() Mode { return ModeHeader }

func setIfNotEmpty(h http.Header, key, value string) {
	if value != "" {
		h.Set(key, value)
	}
}

func resolveProfile(browser string) BrowserProfile {
	if browser == "" || strings.EqualFold(browser, "random") {
		return RandomProfile()
	}
	lower := strings.ToLower(browser)
	for _, p := range Profiles() {
		pLower := strings.ToLower(p.SecCHUA)
		switch {
		case lower == "chrome" && strings.Contains(pLower, "chrome") && !strings.Contains(pLower, "edge"):
			return p
		case lower == "firefox" && p.SecCHUA == "" && p.AcceptLanguage != "":
			return p
		case lower == "edge" && strings.Contains(pLower, "edge"):
			return p
		case lower == "safari" && p.SecCHUA == "" && strings.Contains(p.Accept, "*/*;q=0.8"):
			return p
		}
	}
	return RandomProfile()
}
