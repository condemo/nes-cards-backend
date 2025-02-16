package config

type Config struct {
	GameRecordLimit int
	PlayerLimit     int
}

func newServerConfig() *Config {
	return &Config{
		GameRecordLimit: 15,
		PlayerLimit:     15,
	}
}

var ServerConfig = newServerConfig()
