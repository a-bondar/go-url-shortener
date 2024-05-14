package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddr          string
	ShortLinkBaseURL string
}

func NewConfig() *Config {
	config := &Config{}

	flag.StringVar(&config.RunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&config.ShortLinkBaseURL, "b", "http://localhost:8080", "short link base URL")
	flag.Parse()

	if envRunAddr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		config.RunAddr = envRunAddr
	}

	if shortLinkBaseURL, ok := os.LookupEnv("BASE_URL"); ok {
		config.ShortLinkBaseURL = shortLinkBaseURL
	}

	return config
}
