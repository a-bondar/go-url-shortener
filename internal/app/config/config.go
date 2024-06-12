package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddr          string
	ShortLinkBaseURL string
	FileStoragePath  string
	DatabaseDSN      string
}

func NewConfig() *Config {
	config := &Config{}

	flag.StringVar(&config.RunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&config.ShortLinkBaseURL, "b", "http://localhost:8080", "short link base URL")
	flag.StringVar(&config.FileStoragePath, "f", "/tmp/short-url-db.json", "file storage path")
	flag.StringVar(&config.DatabaseDSN, "d", "", "database data source name")
	flag.Parse()

	if envRunAddr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		config.RunAddr = envRunAddr
	}

	if shortLinkBaseURL, ok := os.LookupEnv("BASE_URL"); ok {
		config.ShortLinkBaseURL = shortLinkBaseURL
	}

	if fileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		config.FileStoragePath = fileStoragePath
	}

	if databaseDSN, ok := os.LookupEnv("DATABASE_DSN"); ok {
		config.DatabaseDSN = databaseDSN
	}

	return config
}
