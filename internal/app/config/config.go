package config

import (
	"flag"
	"os"
)

var FlagOptions struct {
	RunAddr          string
	ShortLinkBaseURL string
}

func ParseFlags() {
	flag.StringVar(&FlagOptions.RunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&FlagOptions.ShortLinkBaseURL, "b", "http://localhost:8080", "short link base URL")

	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		FlagOptions.RunAddr = envRunAddr
	}

	if shortLinkBaseUrl := os.Getenv("BASE_URL"); shortLinkBaseUrl != "" {
		FlagOptions.ShortLinkBaseURL = shortLinkBaseUrl
	}
}
