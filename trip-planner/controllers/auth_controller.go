package controllers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"trip-planner/services"
	"trip-planner/utils"
)

type AuthController struct {
	db        *sql.DB
	auth      *services.AuthService
	jwtSecret string
	domain    string // cookie domain (e.g. travel.geekmacsolutions.com)
}

func NewAuthController(db *sql.DB, auth *services.AuthService, jwtSecret string, cookieDomain string) *AuthController {
	return &AuthController{
		db:        db,
		auth:      auth,
		jwtSecret: jwtSecret,
		domain:    cookieDomain,
	}
}

// Request body: { "id_token": "<google id token jwt>" }
type googleLoginReq struct {
	IDToken string `json:"id_token"`
}

func (a *AuthController) GoogleLogin(c *gin.Context) {
	var req googleLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	if req.IDToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing_token"})
		return
	}

	u, err := a.auth.VerifyGoogleIDToken(c.Request.Context(), req.IDToken)
	if err != nil {
		// include details to debug; you can remove "details" later for production
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token", "details": err.Error()})
		return
	}

	now := time.Now().Unix()

	// upsert user (SQLite style shown; adapt if using other DB)
	_, _ = a.db.Exec(`INSERT OR IGNORE INTO users (id,email,name,picture,created_at) VALUES (?,?,?,?,?)`,
		u.Sub, u.Email, u.Name, u.Picture, now)
	_, _ = a.db.Exec(`UPDATE users SET email=?, name=?, picture=? WHERE id=?`,
		u.Email, u.Name, u.Picture, u.Sub)

	// ensure usage row
	_, _ = a.db.Exec(`INSERT OR IGNORE INTO usage (user_id,generations,updated_at) VALUES (?,?,?)`,
		u.Sub, 0, now)

	jwtToken, err := utils.SignUserJWT(a.jwtSecret, u.Sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "jwt_failed"})
		return
	}

	// Cookie settings
	// SameSite Lax works well for same-site requests
	c.SetSameSite(http.SameSiteLaxMode)

	// Set cookie domain explicitly (recommended).
	// If a.domain == "" -> host-only cookie (still ok).
	c.SetCookie(
		"session",
		jwtToken,
		60*60*24*30, // 30 days
		"/",
		a.domain,
		true, // Secure (HTTPS)
		true, // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"user": gin.H{
			"id":      u.Sub,
			"email":   u.Email,
			"name":    u.Name,
			"picture": u.Picture,
		},
	})
}

func (a *AuthController) Me(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var email, name, picture string
	err := a.db.QueryRow(`SELECT email,name,picture FROM users WHERE id=?`, uid.(string)).
		Scan(&email, &name, &picture)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      uid,
		"email":   email,
		"name":    name,
		"picture": picture,
	})
}

func (a *AuthController) Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("session", "", -1, "/", a.domain, true, true)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
