package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"trip-planner/config"
	"trip-planner/models"
	"trip-planner/services"
	"trip-planner/utils"
)

type TripController struct {
	cfg config.Config

	ai      *services.AIService
	places  *services.PlacesService
	weather *services.WeatherService

	store *utils.JSONStore[models.TripPlan]
}

func NewTripController(cfg config.Config, ai *services.AIService, places *services.PlacesService, weather *services.WeatherService) *TripController {
	return &TripController{
		cfg:     cfg,
		ai:      ai,
		places:  places,
		weather: weather,
		store:   utils.NewJSONStore[models.TripPlan](cfg.PlansFile),
	}
}

func (t *TripController) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// GET /api/v1/trip/plans
func (t *TripController) ListPlans(c *gin.Context) {
	all, err := t.store.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "read_failed", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, all)
}

// GET /api/v1/trip/plan/:id
func (t *TripController) GetPlan(c *gin.Context) {
	id := c.Param("id")

	all, err := t.store.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "read_failed", "details": err.Error()})
		return
	}

	for _, p := range all {
		if p.ID == id {
			c.JSON(http.StatusOK, p)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
}

// POST /api/v1/trip/plan
// Cost control: if same request hash exists => return saved plan (NO AI)
func (t *TripController) CreatePlan(c *gin.Context) {
	var req models.TripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "details": err.Error()})
		return
	}

	// defaults
	if req.Budget == "" {
		req.Budget = "mid"
	}
	if req.Pace == "" {
		req.Pace = "balanced"
	}

	hash := hashTripRequest(req)

	all, err := t.store.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "read_failed", "details": err.Error()})
		return
	}

	// ✅ Return existing plan for same hash (NO AI call)
	for _, p := range all {
		if p.InputHash == hash {
			c.JSON(http.StatusOK, p)
			return
		}
	}

	// -------------------------
	// Places (RAW -> SLIM)
	// -------------------------
	placesRaw, err := t.places.GetPlacesByCity(c.Request.Context(), req.Destination)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "places_failed", "details": err.Error()})
		return
	}
	placesSlim := services.SlimPlaces(placesRaw, 12)

	// -------------------------
	// Weather (RAW -> SLIM)
	// -------------------------
	weatherRaw, err := t.weather.GetCityWeather(c.Request.Context(), req.Destination)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "weather_failed", "details": err.Error()})
		return
	}
	weatherSlim := services.SlimWeather(weatherRaw)

	// -------------------------
	// AI (only once)
	// -------------------------
	itinerary, err := t.ai.GenerateTrip(c.Request.Context(), req, placesSlim, weatherSlim)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "ai_failed", "details": err.Error()})
		return
	}

	now := time.Now().Unix()
	plan := models.TripPlan{
		ID:        uuid.NewString(),
		InputHash: hash,
		Request:   req,
		Itinerary: itinerary,

		// ✅ Save for frontend view + debugging
		Weather: weatherSlim,
		Places:  placesSlim,

		CreatedAt: now,
		UpdatedAt: now,
	}

	all = append(all, plan)
	if err := t.store.WriteAll(all); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save_failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plan)
}

// POST /api/v1/trip/plan/regenerate
// Only called when user edits / forces regenerate
func (t *TripController) Regenerate(c *gin.Context) {
	var req models.TripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "details": err.Error()})
		return
	}

	if req.Budget == "" {
		req.Budget = "mid"
	}
	if req.Pace == "" {
		req.Pace = "balanced"
	}

	hash := hashTripRequest(req)

	all, err := t.store.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "read_failed", "details": err.Error()})
		return
	}

	// Places RAW -> SLIM
	placesRaw, err := t.places.GetPlacesByCity(c.Request.Context(), req.Destination)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "places_failed", "details": err.Error()})
		return
	}
	placesSlim := services.SlimPlaces(placesRaw, 12)

	// Weather RAW -> SLIM
	weatherRaw, err := t.weather.GetCityWeather(c.Request.Context(), req.Destination)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "weather_failed", "details": err.Error()})
		return
	}
	weatherSlim := services.SlimWeather(weatherRaw)

	// Force AI call
	itinerary, err := t.ai.GenerateTrip(c.Request.Context(), req, placesSlim, weatherSlim)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "ai_failed", "details": err.Error()})
		return
	}

	now := time.Now().Unix()

	// If hash exists, update it
	for i := range all {
		if all[i].InputHash == hash {
			all[i].Request = req
			all[i].Itinerary = itinerary
			all[i].Weather = weatherSlim
			all[i].Places = placesSlim
			all[i].UpdatedAt = now

			if err := t.store.WriteAll(all); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "save_failed", "details": err.Error()})
				return
			}
			c.JSON(http.StatusOK, all[i])
			return
		}
	}

	// Otherwise create new
	plan := models.TripPlan{
		ID:        uuid.NewString(),
		InputHash: hash,
		Request:   req,
		Itinerary: itinerary,
		Weather:   weatherSlim,
		Places:    placesSlim,
		CreatedAt: now,
		UpdatedAt: now,
	}
	all = append(all, plan)
	if err := t.store.WriteAll(all); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save_failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plan)
}

func hashTripRequest(req models.TripRequest) string {
	data := req.Destination + "|" + req.StartDate + "|" + itoa(req.Days) + "|" + req.Budget + "|" + req.Pace + "|" + join(req.Interests) + "|" + req.Notes
	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}
	buf := make([]byte, 0, 12)
	for n > 0 {
		d := byte(n%10) + '0'
		buf = append([]byte{d}, buf...)
		n /= 10
	}
	return sign + string(buf)
}

func join(a []string) string {
	if len(a) == 0 {
		return ""
	}
	out := a[0]
	for i := 1; i < len(a); i++ {
		out += "," + a[i]
	}
	return out
}
