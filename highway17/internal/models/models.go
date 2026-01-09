package models

import "time"

// User represents a dashboard user
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	TailscaleIP  string    `json:"tailscale_ip,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Session represents a user session
type Session struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// WeatherData represents current weather conditions
type WeatherData struct {
	Temperature   float64   `json:"temperature"`
	Condition     string    `json:"condition"`
	Humidity      int       `json:"humidity"`
	WindSpeed     float64   `json:"wind_speed"`
	Precipitation float64   `json:"precipitation"`
	CloudCover    int       `json:"cloud_cover"`
	Pressure      float64   `json:"pressure"`
	Visibility    float64   `json:"visibility"`
	LastUpdated   time.Time `json:"last_updated"`
	NextUpdate    time.Time `json:"next_update"`
}

// SystemStats represents system resource usage
type SystemStats struct {
	CPUPercent    float64    `json:"cpu_percent"`
	MemoryPercent float64    `json:"memory_percent"`
	MemoryUsedGB  float64    `json:"memory_used_gb"`
	MemoryTotalGB float64    `json:"memory_total_gb"`
	DiskPercent   float64    `json:"disk_percent"`
	DiskUsedGB    float64    `json:"disk_used_gb"`
	DiskTotalGB   float64    `json:"disk_total_gb"`
	UptimeSeconds uint64     `json:"uptime_seconds"`
	ProcessCount  int        `json:"process_count"`
	LoadAverage   [3]float64 `json:"load_average"`
	LastUpdated   time.Time  `json:"last_updated"`
}

// Widget represents a dashboard widget
type Widget struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Key       string    `json:"key"`
	ValueJSON string    `json:"value_json"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DashboardSettings represents user dashboard preferences
type DashboardSettings struct {
	UserID              int    `json:"user_id"`
	RefreshIntervalSecs int    `json:"refresh_interval_secs"`
	Theme               string `json:"theme"`
	TemperatureUnit     string `json:"temperature_unit"` // C or F
	TimeFormat          string `json:"time_format"`      // 12h or 24h
}
