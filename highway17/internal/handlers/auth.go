package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"citadel/highway17/internal/config"
	"citadel/highway17/internal/database"
	"citadel/highway17/internal/models"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	cfg *config.Config
	db  *database.DB
	log *zap.Logger
}

func NewAuthHandler(cfg *config.Config, db *database.DB, log *zap.Logger) *AuthHandler {
	return &AuthHandler{
		cfg: cfg,
		db:  db,
		log: log,
	}
}

type LoginRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

// Login handles user authentication
func (ah *AuthHandler) Login(c echo.Context) error {
	ctx := context.Background()
	req := new(LoginRequest)

	if err := c.Bind(req); err != nil {
		ah.log.Sugar().Errorw("failed to bind login request", "error", err)
		return c.JSON(400, map[string]string{"error": "invalid request"})
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		return c.JSON(400, map[string]string{"error": "username and password required"})
	}

	// Get user from database
	user, err := ah.db.GetUserByUsername(ctx, req.Username)
	if err != nil {
		// User doesn't exist - check against default password for first-time setup
		if req.Password == ah.cfg.LoginPassword {
			// Create default user on first login
			userID, err := ah.db.CreateUser(ctx, req.Username, hashPassword(ah.cfg.LoginPassword))
			if err != nil {
				ah.log.Sugar().Errorw("failed to create user", "error", err)
				return c.JSON(500, map[string]string{"error": "failed to create user"})
			}

			token, err := ah.createSession(ctx, userID)
			if err != nil {
				ah.log.Sugar().Errorw("failed to create session", "error", err)
				return c.JSON(500, map[string]string{"error": "failed to create session"})
			}

			ah.setSessionCookie(c, token)
			return c.JSON(200, LoginResponse{
				Token:    token,
				Username: req.Username,
				Message:  "User created and logged in",
			})
		}

		ah.log.Sugar().Warnw("login failed - user not found", "username", req.Username)
		return c.JSON(401, map[string]string{"error": "invalid credentials"})
	}

	// Verify password
	passwordHash := user["password_hash"].(string)
	if !verifyPassword(passwordHash, req.Password) {
		ah.log.Sugar().Warnw("login failed - invalid password", "username", req.Username)
		return c.JSON(401, map[string]string{"error": "invalid credentials"})
	}

	// Create session
	userID := user["id"].(int)
	token, err := ah.createSession(ctx, userID)
	if err != nil {
		ah.log.Sugar().Errorw("failed to create session", "error", err)
		return c.JSON(500, map[string]string{"error": "failed to create session"})
	}

	ah.setSessionCookie(c, token)
	ah.log.Sugar().Infow("user logged in", "username", req.Username)

	return c.JSON(200, LoginResponse{
		Token:    token,
		Username: req.Username,
		Message:  "logged in successfully",
	})
}

// Logout clears user session
func (ah *AuthHandler) Logout(c echo.Context) error {
	ctx := context.Background()

	// Get token from cookie
	cookie, err := c.Cookie("session_token")
	if err == nil && cookie.Value != "" {
		ah.db.DeleteSession(ctx, cookie.Value)
	}

	// Clear cookie
	http.SetCookie(c.Response(), &http.Cookie{
		Name:     "session_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	ah.log.Sugar().Info("user logged out")
	return c.JSON(200, map[string]string{"message": "logged out successfully"})
}

// createSession generates a new session token and stores it
func (ah *AuthHandler) createSession(ctx context.Context, userID int) (string, error) {
	token := generateToken()
	expiresAt := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	if err := ah.db.CreateSession(ctx, userID, token, expiresAt); err != nil {
		return "", fmt.Errorf("failed to save session: %w", err)
	}

	return token, nil
}

// setSessionCookie sets the session token cookie
func (ah *AuthHandler) setSessionCookie(c echo.Context, token string) {
	http.SetCookie(c.Response(), &http.Cookie{
		Name:     "session_token",
		Value:    token,
		MaxAge:   86400, // 24 hours
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// Helper functions

// generateToken creates a random 32-byte hex token
func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

// verifyPassword checks if a password matches its hash
func verifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetCurrentUser retrieves the current authenticated user from context
func GetCurrentUser(c echo.Context) (*models.User, error) {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return nil, fmt.Errorf("user not found in context")
	}
	return user, nil
}
