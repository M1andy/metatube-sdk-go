package scraper

import (
	"net/http"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"

	"github.com/metatube-community/metatube-sdk-go/common/fingerprint"
	"github.com/metatube-community/metatube-sdk-go/common/random"
)

type Option func(*Scraper) error

func WithAllowURLRevisit() Option {
	return func(s *Scraper) error {
		colly.AllowURLRevisit()(s.c)
		return nil
	}
}

func WithLogDebugger() Option {
	return func(s *Scraper) error {
		colly.Debugger(&debug.LogDebugger{})
		return nil
	}
}

func WithDetectCharset() Option {
	return func(s *Scraper) error {
		colly.DetectCharset()(s.c)
		return nil
	}
}

func WithIgnoreRobotsTxt() Option {
	return func(s *Scraper) error {
		colly.IgnoreRobotsTxt()(s.c)
		return nil
	}
}

func WithHeaders(headers map[string]string) Option {
	return func(s *Scraper) error {
		colly.Headers(headers)(s.c)
		return nil
	}
}

func WithUserAgent(ua string) Option {
	return func(s *Scraper) error {
		colly.UserAgent(ua)(s.c)
		return nil
	}
}

func WithRandomUserAgent() Option {
	return WithUserAgent(random.UserAgent())
}

func WithCookies(url string, cookies []*http.Cookie) Option {
	return func(s *Scraper) error {
		return s.c.SetCookies(url, cookies)
	}
}

func WithRequestTimeout(timeout time.Duration) Option {
	return func(s *Scraper) error {
		s.c.SetRequestTimeout(timeout)
		return nil
	}
}

func WithDisableCookies() Option {
	return func(s *Scraper) error {
		s.c.DisableCookies()
		return nil
	}
}

func WithDisableRedirects() Option {
	return func(s *Scraper) error {
		s.c.ParseHTTPErrorResponse = true
		s.c.SetRedirectHandler(func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		})
		return nil
	}
}

func WithTransport(transport http.RoundTripper) Option {
	return func(s *Scraper) error {
		s.c.WithTransport(transport)
		return nil
	}
}

func WithLimit(rule *colly.LimitRule) Option {
	return func(s *Scraper) error {
		return s.c.Limit(rule)
	}
}

// WithFingerprint configures browser fingerprint simulation for the scraper.
// It registers an OnRequest callback to inject browser headers and,
// if the fingerprinter provides a custom transport, replaces the default transport.
// Using fingerprint.DefaultFingerprinter applies the default mode (uTLS).
// Pass nil to disable fingerprint.
func WithFingerprint(f fingerprint.Fingerprinter) Option {
	return func(s *Scraper) error {
		if f == nil {
			return nil
		}
		// Register OnRequest callback for header injection.
		// This survives c.Clone() so ClonedCollector() retains fingerprint headers.
		s.c.OnRequest(func(r *colly.Request) {
			// Build a minimal *http.Request to apply fingerprint headers.
			req := &http.Request{
				Method: r.Method,
				URL:    r.URL,
				Header: r.Headers.Clone(),
			}
			_ = f.ApplyRequest(req)
			*r.Headers = req.Header
		})
		// Replace transport if fingerprinter provides one (e.g., uTLS).
		if t := f.Transport(); t != nil {
			s.c.WithTransport(t)
		}
		return nil
	}
}
