package controllers

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"trip-planner/services"
	"trip-planner/utils"
)

type AuthController struct {
	db *sql.DB
	auth *services.AuthService
	jwtSecret string
}

func NewAuthController(db *sql.DB, auth *services.AuthService, jwtSecret string) *AuthController {
	return &AuthController{db: db, auth: auth, jwtSecret: jwtSecret}
}

type googleLoginReq struct {
	id_token string `json:"id_token"`
}

func (a *AuthController) GoogleLogin(c *gin.Context) {
	var body map[string]string
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error":"invalid_request"}); return
	}
	idToken := body["id_token"]
	if idToken == "" { c.JSON(400, gin.H{"error":"missing_token"}); return }

	u, err := a.auth.VerifyGoogleIDToken(c.Request.Context(), idToken)
	if err != nil { c.JSON(401, gin.H{"error":"invalid_token", "details": err.Error()}); return }

	// upsert user
	now := time.Now().Unix()
	_, _ = a.db.Exec(`INSERT OR IGNORE INTO users (id,email,name,picture,created_at) VALUES (?,?,?,?,?)`,
		u.Sub, u.Email, u.Name, u.Picture, now)
	_, _ = a.db.Exec(`UPDATE users SET email=?, name=?, picture=? WHERE id=?`,
		u.Email, u.Name, u.Picture, u.Sub)

	// ensure usage row
	_, _ = a.db.Exec(`INSERT OR IGNORE INTO usage (user_id,generations,updated_at) VALUES (?,?,?)`,
		u.Sub, 0, now)

	jwtToken, err := utils.SignUserJWT(a.jwtSecret, u.Sub)
	if err != nil { c.JSON(500, gin.H{"error":"jwt_failed"}); return }

	// secure cookie (set SameSite Lax)
	c.SetCookie("session", jwtToken, 60*60*24*30, "/", "", true, true)

	c.JSON(200, gin.H{"ok": true, "user": gin.H{
		"id": u.Sub, "email": u.Email, "name": u.Name, "picture": u.Picture,
	}})
}

func (a *AuthController) Me(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok { c.JSON(401, gin.H{"error":"unauthorized"}); return }

	var email, name, picture string
	err := a.db.QueryRow(`SELECT email,name,picture FROM users WHERE id=?`, uid.(string)).Scan(&email,&name,&picture)
	if err != nil { c.JSON(401, gin.H{"error":"unauthorized"}); return }

	c.JSON(200, gin.H{"id": uid, "email": email, "name": name, "picture": picture})
}

func (a *AuthController) Logout(c *gin.Context) {
	c.SetCookie("session", "", -1, "/", "", true, true)
	c.JSON(200, gin.H{"ok": true})
}
