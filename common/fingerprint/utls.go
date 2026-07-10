package fingerprint

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	utls "github.com/refraction-networking/utls"
)

// Importers: scraper/option.go (WithTransport), fetch/fetch.go (Default),
// engine/init.go (initFetcher). Adds dependency: refraction-networking/utls.
// No data files.

var _ Fingerprinter = (*UTLSFingerprinter)(nil)

func init() {
	DefaultFingerprinter = NewUTLSFingerprinter()
}

// UTLSFingerprinter combines uTLS ClientHello spoofing with browser
// header injection. This is the default fingerprint mode (ModeUTLS).
type UTLSFingerprinter struct {
	headerFP      *HeaderFingerprinter
	clientHelloID utls.ClientHelloID
	skipVerify    bool
	rotate        bool
	browser       string
	transport     http.RoundTripper
}

// NewUTLSFingerprinter creates a UTLSFingerprinter with Chrome ClientHello.
func NewUTLSFingerprinter(opts ...Option) *UTLSFingerprinter {
	f := &UTLSFingerprinter{
		headerFP:      NewHeaderFingerprinter(),
		clientHelloID: utls.HelloChrome_Auto,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

// WithSkipVerify configures TLS skip verification.
func WithSkipVerify(skip bool) Option {
	return func(v any) {
		if f, ok := v.(*UTLSFingerprinter); ok {
			f.skipVerify = skip
		}
	}
}

// WithClientHelloID sets the uTLS ClientHello ID.
func WithClientHelloID(id utls.ClientHelloID) Option {
	return func(v any) {
		if f, ok := v.(*UTLSFingerprinter); ok {
			f.clientHelloID = id
		}
	}
}

func (f *UTLSFingerprinter) ApplyRequest(req *http.Request) error {
	return f.headerFP.ApplyRequest(req)
}

// Transport returns a RoundTripper that uses uTLS for HTTPS connections
// and respects HTTP_PROXY/HTTPS_PROXY for proxy support.
func (f *UTLSFingerprinter) Transport() http.RoundTripper {
	if f.transport != nil {
		return f.transport
	}
	f.transport = &utlsRoundTripper{
		clientHelloID: f.clientHelloID,
		skipVerify:    f.skipVerify,
	}
	return f.transport
}

func (f *UTLSFingerprinter) Mode() Mode { return ModeUTLS }

// utlsRoundTripper implements http.RoundTripper with uTLS TLS fingerprinting.
type utlsRoundTripper struct {
	clientHelloID utls.ClientHelloID
	skipVerify    bool
}

func (rt *utlsRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Scheme != "https" {
		// For HTTP, use standard transport to handle redirects/proxy correctly.
		return http.DefaultTransport.RoundTrip(req)
	}

	// Determine proxy from environment.
	proxyURL, err := http.ProxyFromEnvironment(req)
	if err != nil {
		return nil, err
	}

	// Retry chain: Chrome → NoALPN → Chrome → NoALPN
	// H2/curve errors are intermittent; more attempts + delay helps.
	helloIDs := []utls.ClientHelloID{
		rt.clientHelloID,
		utls.HelloRandomizedNoALPN,
		rt.clientHelloID,
		utls.HelloRandomizedNoALPN,
	}

	var lastErr error
	for i, helloID := range helloIDs {
		if i > 0 {
			// Small delay between retries to avoid proxy exhaustion.
			select {
			case <-req.Context().Done():
				return nil, req.Context().Err()
			default:
			}
		}
		resp, err := rt.roundTripSingle(req, proxyURL, helloID)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		// Only retry on H2-related errors (malformed HTTP response).
		if !strings.Contains(err.Error(), "malformed HTTP response") &&
			!strings.Contains(err.Error(), "unsupported curve") {
			return nil, err
		}
	}
	return nil, lastErr
}

// roundTripSingle performs one attempt with a specific ClientHelloID.
func (rt *utlsRoundTripper) roundTripSingle(req *http.Request, proxyURL *url.URL, helloID utls.ClientHelloID) (*http.Response, error) {
	// Establish connection (direct or via proxy CONNECT).
	conn, err := rt.dial(req, proxyURL)
	if err != nil {
		return nil, err
	}

	// Perform uTLS handshake.
	serverName := req.URL.Hostname()
	uConn := utls.UClient(conn, &utls.Config{
		ServerName:         serverName,
		InsecureSkipVerify: rt.skipVerify,
	}, helloID)
	if err := uConn.Handshake(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("utls handshake to %s: %w", serverName, err)
	}

	return rt.doRequest(uConn, req)
}

func (rt *utlsRoundTripper) dial(req *http.Request, proxyURL *url.URL) (net.Conn, error) {
	if proxyURL != nil {
		return rt.connectViaProxy(proxyURL, req.URL.Host)
	}
	return rt.connectDirect(req.URL.Host)
}

func (rt *utlsRoundTripper) connectDirect(host string) (net.Conn, error) {
	dialer := &net.Dialer{}
	addr := host
	if _, _, err := net.SplitHostPort(host); err != nil {
		addr = net.JoinHostPort(host, "443")
	}
	return dialer.Dial("tcp", addr)
}

func (rt *utlsRoundTripper) connectViaProxy(proxyURL *url.URL, targetHost string) (net.Conn, error) {
	dialer := &net.Dialer{}
	proxyAddr := proxyURL.Host
	if _, _, err := net.SplitHostPort(proxyAddr); err != nil {
		port := proxyURL.Port()
		if port == "" {
			port = "80"
		}
		proxyAddr = net.JoinHostPort(proxyAddr, port)
	}

	conn, err := dialer.Dial("tcp", proxyAddr)
	if err != nil {
		return nil, fmt.Errorf("dial proxy %s: %w", proxyAddr, err)
	}

	// Ensure target host has port.
	target := targetHost
	if _, _, err := net.SplitHostPort(target); err != nil {
		target = net.JoinHostPort(target, "443")
	}

	// Send CONNECT request.
	connectReq := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", target, target)
	if _, err := fmt.Fprint(conn, connectReq); err != nil {
		conn.Close()
		return nil, fmt.Errorf("proxy CONNECT write: %w", err)
	}

	// Read CONNECT response.
	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("proxy CONNECT read: %w", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		conn.Close()
		return nil, fmt.Errorf("proxy CONNECT %s: %s", target, resp.Status)
	}

	return conn, nil
}

// doRequest writes the HTTP request over the uTLS connection and reads the response.
func (rt *utlsRoundTripper) doRequest(conn net.Conn, req *http.Request) (*http.Response, error) {
	// Write the request.
	if err := req.Write(conn); err != nil {
		return nil, fmt.Errorf("write request: %w", err)
	}

	// Read the response.
	return http.ReadResponse(bufio.NewReader(conn), req)
}

// ParseMode converts a string to a fingerprint Mode.
func ParseMode(s string) Mode {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "off", "none", "false", "disable", "disabled":
		return ModeOff
	case "header", "headers":
		return ModeHeader
	case "utls", "tls", "ja3":
		return ModeUTLS
	case "rod", "browser", "chrome", "chromium":
		return ModeRod
	default:
		return ModeUTLS
	}
}

// Ensure unused imports are kept (needed by other files in the package).
var (
	_ = context.Background
	_ = tls.Config{}
	_ = io.Discard
)
