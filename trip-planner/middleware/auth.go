package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"trip-planner/utils"
)

func RequireAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session")
		if err != nil || token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error":"unauthorized"})
			c.Abort()
			return
		}
		uid, err := utils.ParseUserJWT(jwtSecret, token)
		if err != nil || uid == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error":"unauthorized"})
			c.Abort()
			return
		}
		c.Set("uid", uid)
		c.Next()
	}
}
