package main

import (
	"log/slog"
	"os"
)

func errorLog(message string, err error) {
	if err != nil {
		slog.Error(message, "error", err)
		os.Exit(1)
	}
}

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		slog.Error("Environment missing", "key", key)
		os.Exit(1)
	}
	return value
}
