package fingerprint

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderFingerprinter_ApplyRequest(t *testing.T) {
	fp := NewHeaderFingerprinter()
	req, _ := http.NewRequest("GET", "https://example.com/", nil)
	err := fp.ApplyRequest(req)
	require.NoError(t, err)

	assert.NotEmpty(t, req.Header.Get("User-Agent"))
	assert.NotEmpty(t, req.Header.Get("Accept"))
	assert.NotEmpty(t, req.Header.Get("Accept-Language"))
	assert.NotEmpty(t, req.Header.Get("Accept-Encoding"))
	assert.Equal(t, "1", req.Header.Get("DNT"))
}

func TestHeaderFingerprinter_Mode(t *testing.T) {
	fp := NewHeaderFingerprinter()
	assert.Equal(t, ModeHeader, fp.Mode())
}

func TestHeaderFingerprinter_Transport_Nil(t *testing.T) {
	fp := NewHeaderFingerprinter()
	assert.Nil(t, fp.Transport())
}

func TestHeaderFingerprinter_Rotate(t *testing.T) {
	fp := NewHeaderFingerprinter(WithProfileRotation(true))
	req, _ := http.NewRequest("GET", "https://example.com/", nil)
	_ = fp.ApplyRequest(req)
	ua1 := req.Header.Get("User-Agent")

	req2, _ := http.NewRequest("GET", "https://example.com/", nil)
	_ = fp.ApplyRequest(req2)
	ua2 := req2.Header.Get("User-Agent")
	t.Logf("UA1=%s UA2=%s", ua1, ua2)
}

func TestUTLSFingerprinter_Transport(t *testing.T) {
	fp := NewUTLSFingerprinter()
	tr := fp.Transport()
	assert.NotNil(t, tr)
}

func TestUTLSFingerprinter_Mode(t *testing.T) {
	fp := NewUTLSFingerprinter()
	assert.Equal(t, ModeUTLS, fp.Mode())
}

func TestUTLSFingerprinter_ApplyRequest(t *testing.T) {
	fp := NewUTLSFingerprinter()
	req, _ := http.NewRequest("GET", "https://example.com/", nil)
	err := fp.ApplyRequest(req)
	require.NoError(t, err)
	assert.NotEmpty(t, req.Header.Get("User-Agent"))
	assert.NotEmpty(t, req.Header.Get("Accept"))
}

func TestUTLSFingerprinter_SkipVerify(t *testing.T) {
	fp := NewUTLSFingerprinter(WithSkipVerify(true))
	assert.True(t, fp.skipVerify)
}

func TestDefaultFingerprinter_NotNil(t *testing.T) {
	assert.NotNil(t, DefaultFingerprinter)
	assert.Equal(t, ModeUTLS, DefaultFingerprinter.Mode())
}

func TestParseMode(t *testing.T) {
	tests := []struct {
		input string
		want  Mode
	}{
		{"utls", ModeUTLS},
		{"UTLS", ModeUTLS},
		{"tls", ModeUTLS},
		{"header", ModeHeader},
		{"headers", ModeHeader},
		{"off", ModeOff},
		{"none", ModeOff},
		{"disable", ModeOff},
		{"disabled", ModeOff},
		{"rod", ModeRod},
		{"chrome", ModeRod},
		{"", ModeUTLS},
		{"unknown", ModeUTLS},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, ParseMode(tt.input))
		})
	}
}

func TestRandomProfile(t *testing.T) {
	p := RandomProfile()
	assert.NotEmpty(t, p.Accept)
}

func TestProfiles(t *testing.T) {
	ps := Profiles()
	assert.GreaterOrEqual(t, len(ps), 4)
	for i, p := range ps {
		assert.NotEmpty(t, p.Accept, "profile %d accept empty", i)
	}
}
