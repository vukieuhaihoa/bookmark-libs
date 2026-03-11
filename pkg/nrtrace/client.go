package nrtrace

import "github.com/newrelic/go-agent/v3/newrelic"

func NewClient(envPrefix string) (*newrelic.Application, error) {
	cfg, err := newConfig(envPrefix)
	if err != nil {
		return nil, err
	}

	println(cfg.AppName)
	println(cfg.LicenseKey)
	println(cfg.LogForwardEnabled)

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
