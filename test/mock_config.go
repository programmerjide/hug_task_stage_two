package test

import "hng_stage_two_task/config"

// MockConfig returns a mock configuration for testing
func MockConfig() config.AppConfig {
	return config.AppConfig{
		ServerPort: "8098",
		Dsn:        "host=127.0.0.1 user=postgres password=postgres dbname=organization port=5432 sslmode=disable",
		AppSecret:  "hng_organization_stage2_task",
	}
}
