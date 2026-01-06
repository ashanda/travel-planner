package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"

	"trip-planner/config"
	"trip-planner/routes"
	"trip-planner/services"
	"trip-planner/storage"
)

func main() {
	cfg := config.Load()

	// ---- DB (SQLite) ----
	db, err := storage.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("db open failed: %v", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	// ---- Services ----
	aiSvc := services.NewAIService(cfg.OpenAIKey, cfg.OpenAIModel)
	placesSvc := services.NewPlacesService(cfg.GoogleMapsKey, cfg.PlacesCacheHours)
	weatherSvc := services.NewWeatherService(cfg.OpenWeatherKey, cfg.WeatherCacheHours)
	authSvc := services.NewAuthService(cfg.GoogleClientID)

	// ---- Gin ----
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// ---- Routes ----
	routes.RegisterRoutes(
		r,
		cfg,
		db,
		aiSvc,
		placesSvc,
		weatherSvc,
		authSvc,
	)

	log.Printf("Trip Planner API running on :%s", cfg.AppPort)
	_ = r.Run(":" + cfg.AppPort)
}
