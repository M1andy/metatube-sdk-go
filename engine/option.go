package engine

import (
	"time"

	"github.com/metatube-community/metatube-sdk-go/common/fingerprint"
	mt "github.com/metatube-community/metatube-sdk-go/provider"
)

type Option func(*Engine)

func WithEngineName(name string) Option {
	return func(e *Engine) {
		e.name = name
	}
}

func WithRequestTimeout(timeout time.Duration) Option {
	return func(e *Engine) {
		e.timeout = timeout
	}
}

func WithActorProviderConfig(name string, config mt.Config) Option {
	return func(e *Engine) {
		e.actorProviderConfigs.Set(name, config)
	}
}

func WithMovieProviderConfig(name string, config mt.Config) Option {
	return func(e *Engine) {
		e.movieProviderConfigs.Set(name, config)
	}
}

// WithFingerprintMode sets the browser fingerprint simulation mode.
// Supported modes: "utls" (default), "header", "off".
func WithFingerprintMode(mode fingerprint.Mode) Option {
	return func(e *Engine) {
		e.fingerprintMode = mode
	}
}
