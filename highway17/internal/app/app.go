package app

import (
	"time"

	"citadel/highway17/internal/config"
	"citadel/highway17/internal/database"
	"citadel/highway17/internal/handlers"
	"citadel/highway17/internal/middleware"
	"citadel/highway17/internal/services"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func New(cfg *config.Config, db *database.DB, log *zap.Logger) (*echo.Echo, error) {
	e := echo.New()

	// Middleware
	e.Use(echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogError:  true,
	}))
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Custom middleware for authentication
	authMW := middleware.NewAuthMiddleware(cfg, db, log)
	e.Use(authMW.CheckAuth)

	// Initialize services
	weatherService := services.NewWeatherService(cfg, log)
	systemStatsService := services.NewSystemStatsService(log, time.Duration(cfg.StatsPollInterval)*time.Second)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(cfg, db, log)
	dashboardHandler := handlers.NewDashboardHandler(cfg, db, log, weatherService, systemStatsService)

	// Routes
	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Auth routes
	e.POST("/api/login", authHandler.Login)
	e.POST("/api/logout", authHandler.Logout)

	// Dashboard routes
	e.GET("/", dashboardHandler.Dashboard)
	e.GET("/dashboard", dashboardHandler.Dashboard)

	// Widget API routes
	e.GET("/api/widgets/weather", dashboardHandler.GetWeatherWidget)
	e.GET("/api/widgets/system", dashboardHandler.GetSystemStatsWidget)
	e.GET("/api/widgets/uptime", dashboardHandler.GetUptimeWidget)

	// Widget data routes
	e.POST("/api/widgets/save", dashboardHandler.SaveWidgetData)
	e.GET("/api/widgets/data", dashboardHandler.GetWidgetData)

	return e, nil
}
