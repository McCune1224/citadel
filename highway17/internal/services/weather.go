package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"citadel/highway17/internal/config"
	"citadel/highway17/internal/models"

	"go.uber.org/zap"
)

type WeatherService struct {
	cfg    *config.Config
	log    *zap.Logger
	client *http.Client

	// Cache
	mu       sync.RWMutex
	cached   *models.WeatherData
	cacheTTL time.Duration
}

// Open-Meteo API response structure
type openMeteoResponse struct {
	Current struct {
		Temperature      float64 `json:"temperature"`
		RelativeHumidity int     `json:"relative_humidity_2m"`
		WeatherCode      int     `json:"weather_code"`
		WindSpeed        float64 `json:"wind_speed_10m"`
		CloudCover       int     `json:"cloud_cover"`
		Pressure         float64 `json:"pressure_msl"`
		Precipitation    float64 `json:"precipitation"`
		Visibility       float64 `json:"visibility"`
	} `json:"current"`
}

func NewWeatherService(cfg *config.Config, log *zap.Logger) *WeatherService {
	return &WeatherService{
		cfg: cfg,
		log: log,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		cacheTTL: time.Duration(cfg.WeatherCacheTTL) * time.Second,
	}
}

// GetWeather fetches current weather data with caching
func (ws *WeatherService) GetWeather(ctx context.Context) (*models.WeatherData, error) {
	ws.mu.RLock()
	if ws.cached != nil && time.Now().Before(ws.cached.NextUpdate) {
		defer ws.mu.RUnlock()
		return ws.cached, nil
	}
	ws.mu.RUnlock()

	// Fetch from Open-Meteo
	data, err := ws.fetchOpenMeteo(ctx)
	if err != nil {
		ws.log.Sugar().Errorw("failed to fetch weather", "error", err)
		// Return cached data if available, even if expired
		ws.mu.RLock()
		if ws.cached != nil {
			defer ws.mu.RUnlock()
			return ws.cached, nil
		}
		ws.mu.RUnlock()
		return nil, err
	}

	// Update cache
	ws.mu.Lock()
	ws.cached = data
	ws.mu.Unlock()

	return data, nil
}

// fetchOpenMeteo fetches weather from Open-Meteo API
func (ws *WeatherService) fetchOpenMeteo(ctx context.Context) (*models.WeatherData, error) {
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m,relative_humidity_2m,weather_code,wind_speed_10m,cloud_cover,pressure_msl,precipitation,visibility",
		ws.cfg.WeatherLatitude,
		ws.cfg.WeatherLongitude,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := ws.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var omResp openMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&omResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	now := time.Now()
	return &models.WeatherData{
		Temperature:   omResp.Current.Temperature,
		Condition:     weatherCodeToCondition(omResp.Current.WeatherCode),
		Humidity:      omResp.Current.RelativeHumidity,
		WindSpeed:     omResp.Current.WindSpeed,
		Precipitation: omResp.Current.Precipitation,
		CloudCover:    omResp.Current.CloudCover,
		Pressure:      omResp.Current.Pressure,
		Visibility:    omResp.Current.Visibility / 1000, // Convert to km
		LastUpdated:   now,
		NextUpdate:    now.Add(ws.cacheTTL),
	}, nil
}

// weatherCodeToCondition converts WMO weather code to readable condition
func weatherCodeToCondition(code int) string {
	// WMO Weather interpretation codes
	switch {
	case code == 0:
		return "Clear sky"
	case code == 1 || code == 2:
		return "Mainly clear"
	case code == 3:
		return "Overcast"
	case code == 45 || code == 48:
		return "Foggy"
	case code == 51 || code == 53 || code == 55:
		return "Drizzle"
	case code == 61 || code == 63 || code == 65:
		return "Rain"
	case code == 71 || code == 73 || code == 75 || code == 77 || code == 80 || code == 81 || code == 82:
		return "Snow"
	case code == 80 || code == 81 || code == 82:
		return "Rain showers"
	case code == 85 || code == 86:
		return "Snow showers"
	case code == 95 || code == 96 || code == 99:
		return "Thunderstorm"
	default:
		return "Unknown"
	}
}

// ClearCache clears the weather cache
func (ws *WeatherService) ClearCache() {
	ws.mu.Lock()
	ws.cached = nil
	ws.mu.Unlock()
}
