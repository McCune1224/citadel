package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"citadel/highway17/internal/config"
	"citadel/highway17/internal/database"
	"citadel/highway17/internal/services"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type DashboardHandler struct {
	cfg                *config.Config
	db                 *database.DB
	log                *zap.Logger
	weatherService     *services.WeatherService
	systemStatsService *services.SystemStatsService
}

func NewDashboardHandler(
	cfg *config.Config,
	db *database.DB,
	log *zap.Logger,
	ws *services.WeatherService,
	ss *services.SystemStatsService,
) *DashboardHandler {
	return &DashboardHandler{
		cfg:                cfg,
		db:                 db,
		log:                log,
		weatherService:     ws,
		systemStatsService: ss,
	}
}

// Dashboard serves the main dashboard page (HTML)
func (dh *DashboardHandler) Dashboard(c echo.Context) error {
	// TODO: Return templ component with proper authentication
	// For now, return JSON indicating dashboard is loaded
	return c.JSON(200, map[string]string{
		"message": "dashboard loaded",
		"status":  "ok",
	})
}

// GetWeatherWidget returns current weather as JSON
func (dh *DashboardHandler) GetWeatherWidget(c echo.Context) error {
	ctx := context.Background()

	weather, err := dh.weatherService.GetWeather(ctx)
	if err != nil {
		dh.log.Sugar().Errorw("failed to get weather", "error", err)
		return c.JSON(500, map[string]string{"error": "failed to fetch weather"})
	}

	return c.JSON(200, weather)
}

// GetSystemStatsWidget returns system statistics as JSON
func (dh *DashboardHandler) GetSystemStatsWidget(c echo.Context) error {
	ctx := context.Background()

	stats, err := dh.systemStatsService.GetStats(ctx)
	if err != nil {
		dh.log.Sugar().Errorw("failed to get system stats", "error", err)
		return c.JSON(500, map[string]string{"error": "failed to fetch system stats"})
	}

	return c.JSON(200, stats)
}

// GetUptimeWidget returns uptime information
func (dh *DashboardHandler) GetUptimeWidget(c echo.Context) error {
	ctx := context.Background()

	stats, err := dh.systemStatsService.GetStats(ctx)
	if err != nil {
		dh.log.Sugar().Errorw("failed to get uptime", "error", err)
		return c.JSON(500, map[string]string{"error": "failed to fetch uptime"})
	}

	// Format uptime
	days := stats.UptimeSeconds / 86400
	hours := (stats.UptimeSeconds % 86400) / 3600
	minutes := (stats.UptimeSeconds % 3600) / 60

	uptimeData := map[string]interface{}{
		"uptime_seconds": stats.UptimeSeconds,
		"uptime_text":    formatUptime(days, hours, minutes),
		"days":           days,
		"hours":          hours,
		"minutes":        minutes,
		"last_updated":   stats.LastUpdated,
	}

	return c.JSON(200, uptimeData)
}

// SaveWidgetData saves user widget data
func (dh *DashboardHandler) SaveWidgetData(c echo.Context) error {
	ctx := context.Background()

	type SaveRequest struct {
		UserID     int    `json:"user_id"`
		WidgetName string `json:"widget_name"`
		WidgetKey  string `json:"widget_key"`
		ValueJSON  string `json:"value_json"`
	}

	var req SaveRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "invalid request"})
	}

	// Validate required fields
	if req.UserID == 0 || req.WidgetName == "" || req.WidgetKey == "" {
		return c.JSON(400, map[string]string{"error": "missing required fields"})
	}

	// Ensure ValueJSON is valid JSON
	var jsonData interface{}
	if err := json.Unmarshal([]byte(req.ValueJSON), &jsonData); err != nil {
		return c.JSON(400, map[string]string{"error": "invalid JSON in value_json"})
	}

	if err := dh.db.SaveWidgetData(ctx, req.UserID, req.WidgetName, req.WidgetKey, req.ValueJSON); err != nil {
		dh.log.Sugar().Errorw("failed to save widget data", "error", err)
		return c.JSON(500, map[string]string{"error": "failed to save widget data"})
	}

	return c.JSON(200, map[string]string{"message": "widget data saved"})
}

// GetWidgetData retrieves saved widget data
func (dh *DashboardHandler) GetWidgetData(c echo.Context) error {
	ctx := context.Background()

	userID := c.QueryParam("user_id")
	widgetName := c.QueryParam("widget_name")
	widgetKey := c.QueryParam("widget_key")

	if userID == "" || widgetName == "" || widgetKey == "" {
		return c.JSON(400, map[string]string{"error": "missing query parameters"})
	}

	// Convert userID to int (simple parsing)
	var uid int
	if err := json.Unmarshal([]byte(userID), &uid); err != nil {
		return c.JSON(400, map[string]string{"error": "invalid user_id"})
	}

	data, err := dh.db.GetWidgetData(ctx, uid, widgetName, widgetKey)
	if err != nil {
		dh.log.Sugar().Warnw("widget data not found", "user_id", uid, "widget_name", widgetName)
		return c.JSON(404, map[string]string{"error": "widget data not found"})
	}

	// Parse and return JSON
	var jsonData interface{}
	json.Unmarshal([]byte(data), &jsonData)

	return c.JSON(200, jsonData)
}

// Helper functions

// formatUptime formats uptime into a readable string
func formatUptime(days, hours, minutes uint64) string {
	if days > 0 {
		return formatDuration(days, hours, minutes)
	}
	if hours > 0 {
		return formatDuration(0, hours, minutes)
	}
	return formatDuration(0, 0, minutes)
}

func formatDuration(days, hours, minutes uint64) string {
	parts := []string{}

	if days > 0 {
		if days == 1 {
			parts = append(parts, "1 day")
		} else {
			parts = append(parts, formatNum(days)+" days")
		}
	}

	if hours > 0 {
		if hours == 1 {
			parts = append(parts, "1 hour")
		} else {
			parts = append(parts, formatNum(hours)+" hours")
		}
	}

	if minutes > 0 {
		if minutes == 1 {
			parts = append(parts, "1 minute")
		} else {
			parts = append(parts, formatNum(minutes)+" minutes")
		}
	}

	// Join with commas and "and"
	if len(parts) == 0 {
		return "less than a minute"
	}
	if len(parts) == 1 {
		return parts[0]
	}

	result := ""
	for i, p := range parts {
		if i == len(parts)-1 {
			result += " and " + p
		} else if i == 0 {
			result = p
		} else {
			result += ", " + p
		}
	}
	return result
}

func formatNum(n uint64) string {
	// Simple formatting using sprintf
	return fmt.Sprintf("%d", n)
}
