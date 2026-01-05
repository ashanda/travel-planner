package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"trip-planner/config"
	"trip-planner/routes"
)

func main() {
	cfg := config.Load()

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	routes.RegisterRoutes(r, cfg)

	log.Printf("Trip Planner API running on :%s", cfg.Port)
	_ = r.Run(":" + cfg.Port)
}
