package nrtrace

import "github.com/kelseyhightower/envconfig"

// config defines the configuration structure for the New Relic client.
type config struct {
	AppName           string `envconfig:"NR_APP_NAME" default:"bookmark_service"`
	LicenseKey        string `envconfig:"NR_LICENSE_KEY" default:""`
	LogForwardEnabled bool   `envconfig:"NR_LOG_FORWARD_ENABLED" default:"true"`
}

// newConfig reads the configuration from environment variables and returns a config struct.
// It takes an envPrefix parameter to allow for different prefixes when reading environment variables.
// Returns:
//   - *config: A pointer to the config struct containing the New Relic configuration
//   - error: An error object if there was an issue processing the environment variables, otherwise nil
func newConfig(envPrefix string) (*config, error) {
	cfg := &config{}
	err := envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
