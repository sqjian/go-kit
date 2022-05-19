package consistent

type Config struct {
	NumberOfReplicas int
}

type Option func(*Config)

func WithNumberOfReplicas(NumberOfReplicas int) Option {
	return func(c *Config) {
		c.NumberOfReplicas = NumberOfReplicas
	}
}
