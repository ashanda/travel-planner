package routes

import (
	"github.com/gin-gonic/gin"
	"trip-planner/config"
	"trip-planner/controllers"
	"trip-planner/services"
)

func RegisterRoutes(r *gin.Engine, cfg config.Config) {
	placesSvc := services.NewPlacesService(cfg.GoogleMapsKey, cfg.PlacesCacheHours)
	weatherSvc := services.NewWeatherService(cfg.OpenWeatherKey, cfg.WeatherCacheHours)
	aiSvc := services.NewAIService(cfg.OpenAIKey, cfg.OpenAIModel)

	tripCtrl := controllers.NewTripController(cfg, aiSvc, placesSvc, weatherSvc)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", tripCtrl.Health)

		// MVP (cost-controlled)
		v1.POST("/trip/plan", tripCtrl.CreatePlan)              // generate if not exists
		v1.POST("/trip/plan/regenerate", tripCtrl.Regenerate)   // only if user edits / forces
		v1.GET("/trip/plan/:id", tripCtrl.GetPlan)              // fetch saved plan
		v1.GET("/trip/plans", tripCtrl.ListPlans)               // list all saved plans
	}
}
