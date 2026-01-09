package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	// Database
	DatabaseURL string

	// Application
	AppPort  int
	AppEnv   string
	LogLevel string

	// Tailscale
	TailscaleEnabled bool
	TailscaleAPIKey  string

	// Auth
	LoginPassword string

	// Weather
	WeatherLatitude  float64
	WeatherLongitude float64
	WeatherCacheTTL  int

	// Polling
	StatsPollInterval   int
	WeatherPollInterval int

	// Features
	SystemStatsEnabled bool
	UptimeEnabled      bool
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		AppPort:             getEnvInt("APP_PORT", 8080),
		AppEnv:              getEnv("APP_ENV", "development"),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		TailscaleEnabled:    getEnvBool("TAILSCALE_ENABLED", true),
		TailscaleAPIKey:     getEnv("TAILSCALE_API_KEY", ""),
		LoginPassword:       getEnv("LOGIN_PASSWORD", "checkpoint"),
		WeatherLatitude:     getEnvFloat64("WEATHER_LATITUDE", 43.1629),   // Rochester, NY default
		WeatherLongitude:    getEnvFloat64("WEATHER_LONGITUDE", -77.6099), // Rochester, NY default
		WeatherCacheTTL:     getEnvInt("WEATHER_CACHE_TTL", 600),
		StatsPollInterval:   getEnvInt("STATS_POLL_INTERVAL", 5),
		WeatherPollInterval: getEnvInt("WEATHER_POLL_INTERVAL", 600),
		SystemStatsEnabled:  getEnvBool("SYSTEM_STATS_ENABLED", true),
		UptimeEnabled:       getEnvBool("UPTIME_ENABLED", true),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

func getEnvFloat64(key string, defaultVal float64) float64 {
	valStr := getEnv(key, "")
	if val, err := strconv.ParseFloat(valStr, 64); err == nil {
		return val
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}
	return valStr == "true" || valStr == "1" || valStr == "yes"
}
