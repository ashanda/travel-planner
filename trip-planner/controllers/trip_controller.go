package controllers

import (
	"crypto/sha256"
	"database/sql"
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
	db  *sql.DB

	ai      *services.AIService
	places  *services.PlacesService
	weather *services.WeatherService

	store *utils.JSONStore[models.TripPlan]
}

func NewTripController(
	cfg config.Config,
	db *sql.DB,
	ai *services.AIService,
	places *services.PlacesService,
	weather *services.WeatherService,
) *TripController {
	return &TripController{
		cfg:     cfg,
		db:      db,
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
// Returns ONLY the logged user's plans
func (t *TripController) ListPlans(c *gin.Context) {
	uid := c.GetString("uid")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	all, err := t.store.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "read_failed", "details": err.Error()})
		return
	}

	// filter by user
	out := make([]models.TripPlan, 0, len(all))
	for _, p := range all {
		if p.UserID == uid {
			out = append(out, p)
		}
	}
	c.JSON(http.StatusOK, out)
}

// GET /api/v1/trip/plan/:id
func (t *TripController) GetPlan(c *gin.Context) {
	uid := c.GetString("uid")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")

	all, err := t.store.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "read_failed", "details": err.Error()})
		return
	}

	for _, p := range all {
		if p.ID == id && p.UserID == uid {
			c.JSON(http.StatusOK, p)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
}

// POST /api/v1/trip/plan
// Cost control:
// - Same request hash => return saved plan (NO AI call)
// - Otherwise: check quota then generate AI once and save
func (t *TripController) CreatePlan(c *gin.Context) {
	uid := c.GetString("uid")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

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

	// ✅ If same hash for same user => return saved plan (NO AI)
	for _, p := range all {
		if p.UserID == uid && p.InputHash == hash {
			c.JSON(http.StatusOK, p)
			return
		}
	}

	// ✅ Enforce quota only when we REALLY need AI
	if err := t.ensureFreeQuota(c, uid); err != nil {
		// ensureFreeQuota already wrote response
		return
	}

	// Places (cached by city)
	places, err := t.places.GetPlacesByCity(c.Request.Context(), req.Destination)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "places_failed", "details": err.Error()})
		return
	}

	// Weather (cached)
	weather, err := t.weather.GetCityWeather(c.Request.Context(), req.Destination)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "weather_failed", "details": err.Error()})
		return
	}

	// AI (only once)
	itinerary, err := t.ai.GenerateTrip(c.Request.Context(), req, places, weather)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "ai_failed", "details": err.Error()})
		return
	}

	// ✅ Attach contexts so frontend can show WeatherCard etc.
	// itinerary is map[string]any (recommended). If your AI returns map, great.
	itinerary["weather"] = weather
	itinerary["places"] = places

	now := time.Now().Unix()
	plan := models.TripPlan{
		ID:        uuid.NewString(),
		UserID:    uid,
		InputHash: hash,
		Request:   req,
		Itinerary: itinerary,
		CreatedAt: now,
		UpdatedAt: now,
	}

	all = append(all, plan)
	if err := t.store.WriteAll(all); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save_failed", "details": err.Error()})
		return
	}

	// ✅ count usage only after successful save
	_ = t.incrementUsage(uid)

	c.JSON(http.StatusOK, plan)
}

// POST /api/v1/trip/plan/regenerate
// Only called when user edits / forces regenerate
func (t *TripController) Regenerate(c *gin.Context) {
	uid := c.GetString("uid")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

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

	// ✅ Enforce quota (regen also consumes a generation)
	if err := t.ensureFreeQuota(c, uid); err != nil {
		return
	}

	places, err := t.places.GetPlacesByCity(c.Request.Context(), req.Destination)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "places_failed", "details": err.Error()})
		return
	}

	weather, err := t.weather.GetCityWeather(c.Request.Context(), req.Destination)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "weather_failed", "details": err.Error()})
		return
	}

	itinerary, err := t.ai.GenerateTrip(c.Request.Context(), req, places, weather)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "ai_failed", "details": err.Error()})
		return
	}

	itinerary["weather"] = weather
	itinerary["places"] = places

	now := time.Now().Unix()

	// If hash exists for same user, update it
	for i := range all {
		if all[i].UserID == uid && all[i].InputHash == hash {
			all[i].Request = req
			all[i].Itinerary = itinerary
			all[i].UpdatedAt = now
			if err := t.store.WriteAll(all); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "save_failed", "details": err.Error()})
				return
			}
			_ = t.incrementUsage(uid)
			c.JSON(http.StatusOK, all[i])
			return
		}
	}

	// Otherwise create new
	plan := models.TripPlan{
		ID:        uuid.NewString(),
		UserID:    uid,
		InputHash: hash,
		Request:   req,
		Itinerary: itinerary,
		CreatedAt: now,
		UpdatedAt: now,
	}
	all = append(all, plan)
	if err := t.store.WriteAll(all); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save_failed", "details": err.Error()})
		return
	}
	_ = t.incrementUsage(uid)

	c.JSON(http.StatusOK, plan)
}

func (t *TripController) ensureFreeQuota(c *gin.Context, uid string) error {
	limit := t.cfg.FreeLimit
	if limit <= 0 {
		limit = 2 // default
	}

	// If DB not configured, allow (but you should configure)
	if t.db == nil {
		return nil
	}

	// Ensure usage row exists
	now := time.Now().Unix()
	_, _ = t.db.Exec(`INSERT OR IGNORE INTO usage (user_id, generations, updated_at) VALUES (?, ?, ?)`, uid, 0, now)

	var gens int
	err := t.db.QueryRow(`SELECT generations FROM usage WHERE user_id = ?`, uid).Scan(&gens)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "usage_read_failed", "details": err.Error()})
		return err
	}

	if gens >= limit {
		c.JSON(http.StatusPaymentRequired, gin.H{
			"error":   "limit_reached",
			"details": "Free limit reached. Please upgrade to generate more itineraries.",
			"limit":   limit,
			"used":    gens,
		})
		return sql.ErrNoRows // any non-nil to stop flow
	}

	return nil
}

func (t *TripController) incrementUsage(uid string) error {
	if t.db == nil {
		return nil
	}
	now := time.Now().Unix()
	_, err := t.db.Exec(`UPDATE usage SET generations = generations + 1, updated_at = ? WHERE user_id = ?`, now, uid)
	return err
}

func hashTripRequest(req models.TripRequest) string {
	// stable hash: serialize important fields
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
