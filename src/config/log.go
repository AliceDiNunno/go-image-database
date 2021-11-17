package config

import GoLoggerClient "github.com/AliceDiNunno/go-logger-client"

func LoadLogConfiguration() GoLoggerClient.ClientConfiguration {
	version, err := GetEnvString("APP_VERSION")
	if err != nil {
		version = "unknown"
	}

	env, err := GetEnvString("APP_ENVIRONMENT")
	if err != nil {
		env = "unknown"
	}

	return GoLoggerClient.ClientConfiguration{
		Url:         RequireEnvString("LOGGER_URL"),
		Port:        RequireEnvInt("LOGGER_PORT"),
		ProjectId:   RequireEnvString("LOGGER_PROJECT_ID"),
		Key:         RequireEnvString("LOGGER_PROJECT_KEY"),
		Environment: env,
		Version:     version,
	}
}
