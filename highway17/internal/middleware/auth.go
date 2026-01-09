package middleware

import (
	"citadel/highway17/internal/config"
	"citadel/highway17/internal/database"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	cfg *config.Config
	db  *database.DB
	log *zap.Logger
}

func NewAuthMiddleware(cfg *config.Config, db *database.DB, log *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		cfg: cfg,
		db:  db,
		log: log,
	}
}

// CheckAuth middleware validates session tokens and Tailscale IPs
func (am *AuthMiddleware) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Skip auth for health check and login routes
		path := c.Request().URL.Path
		if path == "/health" || path == "/api/login" {
			return next(c)
		}

		// TODO: Check session token from cookie
		// TODO: Verify against database
		// TODO: Check Tailscale IP if enabled

		return next(c)
	}
}
