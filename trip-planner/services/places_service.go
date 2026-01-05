package services

import (
	"context"
	"net/url"
	"strings"
	"sync"
	"time"

	"trip-planner/utils"
)

type PlacesService struct {
	apiKey string
	ttl    time.Duration

	mu    sync.RWMutex
	cache map[string]cacheItem
}

type cacheItem struct {
	expires time.Time
	value   any // cached RAW Google response
}

func NewPlacesService(apiKey string, cacheHours int) *PlacesService {
	return &PlacesService{
		apiKey: apiKey,
		ttl:    time.Duration(cacheHours) * time.Hour,
		cache:  map[string]cacheItem{},
	}
}

// GetPlacesByCity returns RAW Google response (cached)
func (s *PlacesService) GetPlacesByCity(ctx context.Context, city string) (any, error) {
	key := "places:" + strings.ToLower(strings.TrimSpace(city))

	s.mu.RLock()
	if item, ok := s.cache[key]; ok && time.Now().Before(item.expires) {
		s.mu.RUnlock()
		return item.value, nil
	}
	s.mu.RUnlock()

	q := url.QueryEscape("top attractions in " + city + " Sri Lanka")
	u := "https://maps.googleapis.com/maps/api/place/textsearch/json?query=" + q + "&key=" + s.apiKey

	var resp any
	if err := utils.GetJSON(ctx, u, &resp, nil); err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.cache[key] = cacheItem{expires: time.Now().Add(s.ttl), value: resp}
	s.mu.Unlock()

	return resp, nil
}

// SlimPlaces reduces RAW Google response to small JSON for AI prompt
// This massively reduces token usage + speeds up OpenAI.
func SlimPlaces(raw any, limit int) map[string]any {
	if limit <= 0 {
		limit = 12
	}

	m, ok := raw.(map[string]any)
	if !ok {
		return map[string]any{"top_places": []any{}}
	}

	status, _ := m["status"].(string)

	resultsAny, ok := m["results"].([]any)
	if !ok || len(resultsAny) == 0 {
		return map[string]any{
			"status":     status,
			"top_places": []any{},
		}
	}

	if len(resultsAny) < limit {
		limit = len(resultsAny)
	}

	top := make([]map[string]any, 0, limit)
	for i := 0; i < limit; i++ {
		r, ok := resultsAny[i].(map[string]any)
		if !ok {
			continue
		}

		place := map[string]any{
			"name":     r["name"],
			"rating":   r["rating"],
			"types":    r["types"],
			"place_id": r["place_id"],
		}

		// keep one address field if available
		if v, ok := r["formatted_address"]; ok {
			place["address"] = v
		} else if v, ok := r["vicinity"]; ok {
			place["address"] = v
		}

		// keep lat/lng if available
		if g, ok := r["geometry"].(map[string]any); ok {
			if loc, ok := g["location"].(map[string]any); ok {
				place["lat"] = loc["lat"]
				place["lng"] = loc["lng"]
			}
		}

		top = append(top, place)
	}

	return map[string]any{
		"status":     status,
		"top_places": top,
	}
}
