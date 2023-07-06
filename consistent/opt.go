package consistent

func newDefaultConfig() *Config {
	return &Config{
		NumberOfReplicas: 20,
	}
}

type Config struct {
	NumberOfReplicas int
}

type Option func(*Config)
