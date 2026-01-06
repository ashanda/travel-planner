package routes

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"trip-planner/config"
	"trip-planner/controllers"
	"trip-planner/middleware"
	"trip-planner/services"
)

func RegisterRoutes(
	r *gin.Engine,
	cfg config.Config,
	db *sql.DB,
	ai *services.AIService,
	places *services.PlacesService,
	weather *services.WeatherService,
	authSvc *services.AuthService,
) {
	// -------- Base middleware (recommended) --------
	r.Use(gin.Recovery())

	// If you are behind reverse proxy (nginx), set trusted proxies properly
	// For now, disable trusting all proxies:
	_ = r.SetTrustedProxies(nil)

	// -------- Versioned API group --------
	api := r.Group("/api")
	v1 := api.Group("/v1")

	// -------- Health --------
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// -------- Controllers --------
	authCtrl := controllers.NewAuthController(db, authSvc, cfg.JWTSecret)
	tripCtrl := controllers.NewTripController(cfg, db, ai, places, weather)

	// -------- Auth routes --------
	// Frontend sends Google "id_token"
	v1.POST("/auth/google", authCtrl.GoogleLogin)

	// Logged-in session info
	v1.GET("/auth/me", middleware.RequireAuth(cfg.JWTSecret), authCtrl.Me)

	// Logout clears cookie
	v1.POST("/auth/logout", middleware.RequireAuth(cfg.JWTSecret), authCtrl.Logout)

	// -------- Protected Trip routes --------
	trip := v1.Group("/trip")
	trip.Use(middleware.RequireAuth(cfg.JWTSecret))

	trip.GET("/plans", tripCtrl.ListPlans)
	trip.GET("/plan/:id", tripCtrl.GetPlan)

	// Generate plan (2-free limit enforced inside controller)
	trip.POST("/plan", tripCtrl.CreatePlan)

	// Force regenerate (consumes generation)
	trip.POST("/plan/regenerate", tripCtrl.Regenerate)
}
