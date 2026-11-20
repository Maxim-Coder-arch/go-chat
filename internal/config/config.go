package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerAddress string
	HistoryLimit  int
}

func Load() Config {
	return Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		HistoryLimit:  getEnvInt("HISTORY_LIMIT", 100),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}