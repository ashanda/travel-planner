package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	// Server
	AppPort string

	// AI
	OpenAIKey   string
	OpenAIModel string

	// Google Places
	GoogleMapsKey string

	// Weather
	OpenWeatherKey string

	// Storage
	PlansFile string
	DBPath    string

	// Cache
	PlacesCacheHours  int
	WeatherCacheHours int

	// Auth
	GoogleClientID string
	JWTSecret      string

	// Limits
	FreeLimit int
}

func Load() Config {
	return Config{
		AppPort: getEnv("APP_PORT", "8080"),

		OpenAIKey:   mustEnv("OPENAI_API_KEY"),
		OpenAIModel: getEnv("OPENAI_MODEL", "gpt-5.2"),

		GoogleMapsKey:  mustEnv("GOOGLE_MAPS_API_KEY"),
		OpenWeatherKey: getEnv("OPENWEATHER_API_KEY", ""),

		PlansFile: getEnv("PLANS_FILE", "/app/storage/plans.json"),
		DBPath:    getEnv("DB_PATH", "/app/storage/app.db"),

		PlacesCacheHours:  getEnvInt("PLACES_CACHE_HOURS", 168),
		WeatherCacheHours: getEnvInt("WEATHER_CACHE_HOURS", 2),

		GoogleClientID: mustEnv("GOOGLE_CLIENT_ID"),
		JWTSecret:      mustEnv("JWT_SECRET"),

		FreeLimit: getEnvInt("FREE_LIMIT", 2),
	}
}

// ---------- helpers ----------

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Missing required env var: %s", key)
	}
	return v
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getEnvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}
