package server

type Config struct {
	address string
}

func DefaultCfg() *Config {
	return &Config{
		address: "localhost:8000",
	}
}

func NewCfg(addr string) *Config {
	return &Config{
		address: addr,
	}
}
