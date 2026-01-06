package services

import (
	"context"
	"time"

	"google.golang.org/api/idtoken"
)

type GoogleUser struct {
	Sub     string
	Email   string
	Name    string
	Picture string
}

type AuthService struct {
	GoogleClientID string
}

func NewAuthService(cid string) *AuthService {
	return &AuthService{GoogleClientID: cid}
}

func (a *AuthService) VerifyGoogleIDToken(ctx context.Context, token string) (*GoogleUser, error) {
	payload, err := idtoken.Validate(ctx, token, a.GoogleClientID)
	if err != nil { return nil, err }

	u := &GoogleUser{
		Sub:     payload.Subject,
		Email:   payload.Claims["email"].(string),
	}
	if v, ok := payload.Claims["name"].(string); ok { u.Name = v }
	if v, ok := payload.Claims["picture"].(string); ok { u.Picture = v }
	_ = time.Now()

	return u, nil
}
