// Package sdk provides the hadron public sdk for plans
package sdk

import (
	"context"

	quark "github.com/farcloser/quark/sdk"
	"github.com/rs/zerolog"
)

// LoadEnv loads environment variables from a .env file.
// Wraps Hadron's LoadEnv for convenience.
//
//nolint:wrapcheck
func LoadEnv(path string) error {
	return quark.LoadEnv(path)
}

// GetEnv retrieves a required environment variable.
// Returns an error if the variable does not exist.
// Empty values (FOO="") are allowed and will not cause an error.
//
//nolint:wrapcheck
func GetEnv(key string) (string, error) {
	return quark.GetEnv(key)
}

// GetEnvWithFallback retrieves an environment variable or returns a default value.
// Wraps Hadron's GetEnvWithFallback for convenience.
func GetEnvWithFallback(key, defaultValue string) string {
	return quark.GetEnvWithFallback(key, defaultValue)
}

//nolint:wrapcheck
func GetSecret(ctx context.Context, itemRef string, fields []string) (map[string]string, error) {
	return quark.GetSecret(ctx, itemRef, fields)
}

// ConfigureDefaultLogger configures the global zerolog logger with sensible defaults.
// It uses a console writer with RFC3339 timestamps for human-readable output.
// If a log level is provided, it sets that level. Otherwise, it reads from the LOG_LEVEL
// environment variable (defaults to "info" if not set or invalid).
// Wraps Hadron's ConfigureDefaultLogger for convenience.
func ConfigureDefaultLogger(ctx context.Context, level ...zerolog.Level) {
	quark.ConfigureDefaultLogger(ctx, level...)
}
