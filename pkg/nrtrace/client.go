package nrtrace

import "github.com/newrelic/go-agent/v3/newrelic"

// NewClient initializes and returns a New Relic application instance based on the provided environment variable prefix.
//
// The function reads the New Relic configuration from environment variables using the specified prefix.
// It expects the following environment variables to be set:
//   - {envPrefix}_APP_NAME: The name of the New Relic application.
//   - {envPrefix}_LICENSE_KEY: The New Relic license key.
//   - {envPrefix}_LOG_FORWARD_ENABLED: A boolean value indicating whether log forwarding is enabled.
//
// Parameters:
//   - envPrefix: A string prefix used to identify the relevant environment variables for New Relic configuration.
//
// Returns:
//   - *newrelic.Application: A pointer to the initialized New Relic application instance.
//   - error: An error object if there was an issue during initialization, or nil if successful.
func NewClient(envPrefix string) (*newrelic.Application, error) {
	cfg, err := newConfig(envPrefix)
	if err != nil {
		return nil, err
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(cfg.AppName),
		newrelic.ConfigLicense(cfg.LicenseKey),
		newrelic.ConfigAppLogForwardingEnabled(cfg.LogForwardEnabled),
	)
	if err != nil {
		return nil, err
	}

	return app, nil
}
