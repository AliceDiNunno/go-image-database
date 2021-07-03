package config

type GinConfig struct {
	Host string
	Port int
	Mode string
	Tls  bool
}

func LoadGinConfiguration() GinConfig {
	return GinConfig{
		Host: RequireEnvString("GIN_LISTEN_URL"),
		Port: RequireEnvInt("GIN_PORT"),
		Mode: RequireEnvString("GIN_MODE"),
		Tls:  RequireEnvBool("GIN_TLS"),
	}
}
