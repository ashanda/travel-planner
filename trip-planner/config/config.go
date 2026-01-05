package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port string

	OpenAIKey   string
	OpenAIModel string

	GoogleMapsKey  string
	OpenWeatherKey string

	PlansFile string

	PlacesCacheHours  int
	WeatherCacheHours int
}

func Load() Config {
	return Config{
		Port: getenv("APP_PORT", "8080"),

		OpenAIKey:   mustGet("OPENAI_API_KEY"),
		OpenAIModel: getenv("OPENAI_MODEL", "gpt-5.2"),

		GoogleMapsKey:  mustGet("GOOGLE_MAPS_API_KEY"),
		OpenWeatherKey: getenv("OPENWEATHER_API_KEY", ""),

		PlansFile: getenv("PLANS_FILE", "storage/plans.json"),

		PlacesCacheHours:  atoi(getenv("PLACES_CACHE_HOURS", "168")),
		WeatherCacheHours: atoi(getenv("WEATHER_CACHE_HOURS", "2")),
	}
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
func mustGet(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic("Missing env: " + k)
	}
	return v
}
func atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}
