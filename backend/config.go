package main

import (
	"log"
	"os"
)

var (
	// general
	port string
	frontendURL string

	// Database
	dbHost string
	dbPort string
	dbUser string
	dbName string
	dbPassword string
	
	// otel
	serviceName string
	otelEndpoint string
	insecureMode string
)

func loadConfig() {
	envVars := map[string]*string{
		"PORT":                        &port,
		"FRONTEND_URL":               &frontendURL,
		"DB_NAME":                    &dbName,
		"DB_PASSWORD":                &dbPassword,
		"DB_USER":                    &dbUser,
		"DB_PORT":                    &dbPort,
		"DB_HOST":                    &dbHost,
		"SERVICE_NAME":               &serviceName,
		"OTEL_EXPORTER_OTLP_ENDPOINT": &otelEndpoint,
		"INSECURE_MODE":              &insecureMode,
	}

	var missing []string
	for key, ref := range envVars {
		*ref = os.Getenv(key)
		if *ref == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		log.Fatalf("Missing required environment variables: %v", missing)
	}
}
