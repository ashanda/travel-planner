package services

import (
	"context"
	"net/url"
	"strings"
	"sync"
	"time"

	"trip-planner/utils"
)

type WeatherService struct {
	apiKey string
	ttl    time.Duration

	mu    sync.RWMutex
	cache map[string]cacheItem
}

func NewWeatherService(apiKey string, cacheHours int) *WeatherService {
	return &WeatherService{
		apiKey: apiKey,
		ttl:    time.Duration(cacheHours) * time.Hour,
		cache:  map[string]cacheItem{},
	}
}

// GetCityWeather returns RAW OpenWeather response (cached)
func (s *WeatherService) GetCityWeather(ctx context.Context, city string) (any, error) {
	if strings.TrimSpace(s.apiKey) == "" {
		// No key => disabled
		return map[string]any{"enabled": false}, nil
	}

	key := "weather:" + strings.ToLower(strings.TrimSpace(city))

	s.mu.RLock()
	if item, ok := s.cache[key]; ok && time.Now().Before(item.expires) {
		s.mu.RUnlock()
		return item.value, nil
	}
	s.mu.RUnlock()

	q := url.QueryEscape(strings.TrimSpace(city) + ",LK")
	u := "https://api.openweathermap.org/data/2.5/weather?q=" + q + "&appid=" + s.apiKey + "&units=metric"

	var resp any
	if err := utils.GetJSON(ctx, u, &resp, nil); err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.cache[key] = cacheItem{expires: time.Now().Add(s.ttl), value: resp}
	s.mu.Unlock()

	return resp, nil
}

// ✅ SlimWeather reduces raw OpenWeather response to a small JSON object for AI + frontend
// Output format example:
// {"enabled":true,"temp_c":23.6,"feels_like_c":24.1,"humidity":80,"condition":"Clouds","wind_mps":2.5}
func SlimWeather(raw any) map[string]any {
	// if service returned {"enabled": false}
	if m0, ok := raw.(map[string]any); ok {
		if v, ok := m0["enabled"].(bool); ok && v == false {
			return map[string]any{"enabled": false}
		}
	}

	m, ok := raw.(map[string]any)
	if !ok {
		return map[string]any{"enabled": false}
	}

	main, _ := m["main"].(map[string]any)
	wind, _ := m["wind"].(map[string]any)

	condition := ""
	if arr, ok := m["weather"].([]any); ok && len(arr) > 0 {
		if w0, ok := arr[0].(map[string]any); ok {
			if s, ok := w0["main"].(string); ok {
				condition = s
			}
		}
	}

	return map[string]any{
		"enabled":       true,
		"temp_c":        main["temp"],
		"feels_like_c":  main["feels_like"],
		"humidity":      main["humidity"],
		"condition":     condition,
		"wind_mps":      wind["speed"],
	}
}
