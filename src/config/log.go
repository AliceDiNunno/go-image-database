package config

type LogConfig struct {
	RollbarToken string
}

func LoadLogConfiguration() LogConfig {
	return LogConfig{
		RollbarToken: RequireEnvString("ROLLBAR_TOKEN"),
	}
}
