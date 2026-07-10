package fingerprint

import "math/rand"

// BrowserProfile holds a set of HTTP headers that mimic a specific browser.
type BrowserProfile struct {
	// UserAgent is the full User-Agent string (can be empty to use external provider).
	UserAgent string

	// Sec-CH-UA is the User-Agent client hints header.
	SecCHUA string

	// Sec-CH-UA-Platform is the platform client hint.
	SecCHUAPlatform string

	// Accept is the HTTP Accept header.
	Accept string

	// AcceptLanguage is the Accept-Language header.
	AcceptLanguage string

	// AcceptEncoding is the Accept-Encoding header.
	AcceptEncoding string

	// SecFetchDest is the Sec-Fetch-Dest header.
	SecFetchDest string

	// SecFetchMode is the Sec-Fetch-Mode header.
	SecFetchMode string

	// SecFetchSite is the Sec-Fetch-Site header.
	SecFetchSite string

	// CacheControl is the Cache-Control header.
	CacheControl string
}

var profiles = []BrowserProfile{
	// Chrome 124 on Windows 10
	{
		SecCHUA:          `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`,
		SecCHUAPlatform:  `"Windows"`,
		Accept:           "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		AcceptLanguage:   "en-US,en;q=0.9,ja;q=0.8",
		AcceptEncoding:   "gzip, deflate, br",
		SecFetchDest:     "document",
		SecFetchMode:     "navigate",
		SecFetchSite:     "none",
		CacheControl:     "max-age=0",
	},
	// Chrome 124 on macOS
	{
		SecCHUA:          `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`,
		SecCHUAPlatform:  `"macOS"`,
		Accept:           "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		AcceptLanguage:   "en-US,en;q=0.9",
		AcceptEncoding:   "gzip, deflate, br",
		SecFetchDest:     "document",
		SecFetchMode:     "navigate",
		SecFetchSite:     "none",
		CacheControl:     "max-age=0",
	},
	// Chrome 124 on Linux
	{
		SecCHUA:          `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`,
		SecCHUAPlatform:  `"Linux"`,
		Accept:           "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		AcceptLanguage:   "en-US,en;q=0.9",
		AcceptEncoding:   "gzip, deflate, br",
		SecFetchDest:     "document",
		SecFetchMode:     "navigate",
		SecFetchSite:     "none",
		CacheControl:     "max-age=0",
	},
	// Firefox 126 on Windows
	{
		SecCHUA:          "", // Firefox does not send Sec-CH-UA
		SecCHUAPlatform:  "",
		Accept:           "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		AcceptLanguage:   "en-US,en;q=0.5",
		AcceptEncoding:   "gzip, deflate, br",
		SecFetchDest:     "document",
		SecFetchMode:     "navigate",
		SecFetchSite:     "none",
		CacheControl:     "max-age=0",
	},
	// Edge 124 on Windows
	{
		SecCHUA:          `"Chromium";v="124", "Microsoft Edge";v="124", "Not-A.Brand";v="99"`,
		SecCHUAPlatform:  `"Windows"`,
		Accept:           "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		AcceptLanguage:   "en-US,en;q=0.9",
		AcceptEncoding:   "gzip, deflate, br",
		SecFetchDest:     "document",
		SecFetchMode:     "navigate",
		SecFetchSite:     "none",
		CacheControl:     "max-age=0",
	},
	// Safari 17 on macOS
	{
		SecCHUA:          "", // Safari does not send Sec-CH-UA
		SecCHUAPlatform:  "",
		Accept:           "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		AcceptLanguage:   "en-US,en;q=0.9",
		AcceptEncoding:   "gzip, deflate, br",
		SecFetchDest:     "document",
		SecFetchMode:     "navigate",
		SecFetchSite:     "none",
		CacheControl:     "max-age=0",
	},
}

// RandomProfile returns a random BrowserProfile from the predefined set.
func RandomProfile() BrowserProfile {
	return profiles[rand.Intn(len(profiles))]
}

// Profiles returns all predefined browser profiles (read-only).
func Profiles() []BrowserProfile {
	return profiles
}
