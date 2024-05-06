package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddr          string
	ShortLinkBaseURL string
}

func GetConfig() *Config {
	config := Config{}

	flag.StringVar(&config.RunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&config.ShortLinkBaseURL, "b", "http://localhost:8080", "short link base URL")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		config.RunAddr = envRunAddr
	}

	if shortLinkBaseURL := os.Getenv("BASE_URL"); shortLinkBaseURL != "" {
		config.ShortLinkBaseURL = shortLinkBaseURL
	}

	return &config
}
