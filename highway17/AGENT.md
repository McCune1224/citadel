# Highway 17 Dashboard - Agent Guide

This document provides context for AI agents working on the Highway 17 Dashboard project.

## Project Overview

**Highway 17 Dashboard** is a personal dashboard + service hub for city17 home lab server, built with:
- **Backend:** Go with Echo web framework
- **Frontend:** templ components with HTMX + TailwindCSS
- **Database:** PostgreSQL with pgx + sqlc
- **Styling:** Valve Half-Life 2 aesthetic (Orange #FF8C00, Dark #0D0D0D, Cyan #00FFFF, Green #00FF00)
- **Monitoring:** System stats with gopsutil, weather via Open-Meteo API
- **Authentication:** Tailscale IP whitelist + password fallback

## Quick Start

### Prerequisites
- Go 1.24+
- Node.js 18+ (for TailwindCSS)
- PostgreSQL 16+ (local or via Railway)
- Docker (optional, for containerization)

### Local Development

1. **Setup environment:**
   ```bash
   cd citadel/highway17
   cp .env.example .env
   # Edit .env with your DATABASE_URL and preferences
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   npm install
   ```

3. **Generate templ components:**
   ```bash
   templ generate
   ```

4. **Build TailwindCSS:**
   ```bash
   npm run css
   # For watch mode: npm run css:watch
   ```

5. **Run database migrations:**
   - Migrations run automatically on app startup via `db.Migrate(ctx)`
   - Located in `migrations/001_initial_schema.sql`

6. **Start the server:**
   ```bash
   go run ./cmd/dashboard
   ```

7. **Access the dashboard:**
   - Login: http://localhost:8080 (default password: checkpoint)

### Docker

```bash
docker-compose -f deployments/docker-compose.yml up --build
```

## Project Structure

```
highway17/
├── cmd/dashboard/main.go              # Application entry point
├── internal/
│   ├── app/app.go                     # Echo app setup
│   ├── config/config.go               # Configuration loader
│   ├── logger/logger.go               # Zap logging setup
│   ├── database/
│   │   ├── db.go                      # Database connection & queries
│   │   └── schema.sql                 # Database schema
│   ├── models/models.go               # Domain models (User, Session, etc)
│   ├── middleware/auth.go             # Authentication middleware
│   ├── handlers/
│   │   ├── auth.go                    # Login/logout handlers
│   │   └── dashboard.go               # Dashboard & widget handlers
│   └── services/
│       ├── weather.go                 # Open-Meteo weather API
│       └── system.go                  # System stats via gopsutil
├── web/
│   ├── components/                    # templ components
│   │   ├── layout.templ               # Base HTML layout
│   │   ├── login.templ                # Login page
│   │   ├── dashboard.templ            # Dashboard container
│   │   └── widgets.templ              # Weather, System, Uptime widgets
│   └── static/
│       ├── css/
│       │   ├── input.css              # TailwindCSS input
│       │   └── style.css              # Generated CSS (auto-generated)
│       ├── js/                        # JavaScript (if needed)
│       └── fonts/                     # Font assets
├── migrations/001_initial_schema.sql  # Database schema migration
├── deployments/
│   ├── Dockerfile                     # Multi-stage Docker build
│   └── docker-compose.yml             # Local dev environment
├── go.mod & go.sum                    # Go dependencies
├── package.json & package-lock.json   # Node dependencies
├── tailwind.config.js                 # TailwindCSS configuration
├── postcss.config.js                  # PostCSS configuration
└── .env.example                       # Environment variables template
```

## Key Concepts

### Authentication Flow
1. User submits username/password at `/api/login`
2. If user doesn't exist AND password matches default, create user with bcrypt hash
3. If user exists, verify password hash with bcrypt
4. Generate 32-byte hex session token
5. Store token in database with 24-hour expiry
6. Return token in session cookie
7. Middleware checks cookie on subsequent requests

### Widgets
- **Weather:** Polls Open-Meteo API every 10 minutes, cached for performance
- **System Stats:** Queries gopsutil every 5 seconds for CPU/Memory/Disk usage
- **Uptime:** Shows system uptime in readable format

### HTMX Integration
- Dashboard uses HTMX triggers for automatic widget polling
- `hx-trigger="load, every 10m"` for weather
- `hx-trigger="load, every 5s"` for system stats
- Widget endpoints return JSON (can be enhanced with HTMX fragments)

### Database
- PostgreSQL with pgx connection pool
- Schema: users, sessions, widget_data, settings tables
- Migrations auto-run on app startup
- UGC: Upsert on conflict for widget/settings data

## Configuration

Environment variables in `.env`:

```env
# Database
DATABASE_URL=postgresql://user:pass@host:5432/highway17

# Application
APP_PORT=8080
APP_ENV=development
LOG_LEVEL=info

# Tailscale
TAILSCALE_ENABLED=true
TAILSCALE_API_KEY=tskey-...

# Authentication
LOGIN_PASSWORD=checkpoint

# Weather
WEATHER_LATITUDE=43.1629      # Rochester, NY
WEATHER_LONGITUDE=-77.6099
WEATHER_CACHE_TTL=600         # seconds

# Polling
STATS_POLL_INTERVAL=5         # seconds
WEATHER_POLL_INTERVAL=600     # seconds

# Features
SYSTEM_STATS_ENABLED=true
UPTIME_ENABLED=true
```

## Development Tips

### Add a New Widget
1. Create service in `internal/services/` (e.g., `uptime.go`)
2. Add handler method in `internal/handlers/dashboard.go`
3. Create templ component in `web/components/widgets.templ`
4. Register route in `internal/app/app.go`
5. Add HTMX trigger to dashboard in `web/components/dashboard.templ`

### Rebuild After Changes
```bash
# After Go code changes
go build ./cmd/dashboard

# After templ changes
templ generate
go build ./cmd/dashboard

# After CSS changes
npm run css

# After migrations
# (migrations auto-run on startup, or manually run:)
# psql $DATABASE_URL < migrations/001_initial_schema.sql
```

### Debugging
- Logs go to stdout in JSON format (Zap)
- Set `LOG_LEVEL=debug` for verbose output
- Database errors logged with context
- All HTTP requests logged with status codes

## Common Tasks

### Connect to Database
```bash
psql $DATABASE_URL
# View users table
SELECT id, username, tailscale_ip FROM users;
```

### Reset Database
```bash
# Drop and recreate (careful!)
psql $DATABASE_URL << EOF
DROP TABLE IF EXISTS settings, widget_data, sessions, users;
\q
EOF

# Restart app to re-run migrations
```

### Test Weather API
```bash
curl "https://api.open-meteo.com/v1/forecast?latitude=43.1629&longitude=-77.6099&current=temperature_2m"
```

### Build Docker Image
```bash
docker build -f deployments/Dockerfile -t highway17:latest .
```

## Deployment

### To city17 Server
1. Build and test locally with docker-compose
2. SSH to city17 server
3. Clone repo or pull latest changes
4. Set production DATABASE_URL in `.env`
5. Run: `docker-compose -f deployments/docker-compose.yml up -d`
6. Verify: `curl http://localhost:8080/health`

### Health Check Endpoint
- `GET /health` returns `{"status":"ok"}`
- Use for monitoring and readiness checks

## Git Workflow

All changes committed directly to main (no feature branches):
```bash
git add .
git commit -m "feat: add widget X" # or "fix:", "refactor:", "docs:"
git push origin main
```

## Performance Notes

- Weather data cached for 10 minutes (configurable via `WEATHER_CACHE_TTL`)
- System stats cached for 5 seconds (configurable via `STATS_POLL_INTERVAL`)
- Database connection pool: pgx default settings
- TailwindCSS built once at build time, not runtime
- templ generates Go code at build time

## Known Limitations

- Tailscale IP check not yet implemented (middleware placeholder)
- Dashboard page returns JSON (should return templ component)
- Widget data storage not yet wired to persistent storage
- Settings not yet exposed in UI
- No support for multiple users on same dashboard (could add)

## Future Enhancements

- [ ] Complete Tailscale IP whitelist check
- [ ] Return proper HTML from dashboard (use templ components)
- [ ] Add widget settings UI
- [ ] Support multiple dashboard layouts
- [ ] Add more widgets (e.g., Docker container status, Git repos)
- [ ] Add webhook integrations for alerts
- [ ] Implement dark/light theme toggle (beyond Valve color scheme)
- [ ] Add CI/CD pipeline with GitHub Actions
- [ ] Add automated tests and coverage reporting

## Support & Debugging

For questions or issues:
1. Check `.env` configuration first
2. Look at server logs: `docker logs highway17-dashboard`
3. Check database connectivity: `psql $DATABASE_URL -c "SELECT 1"`
4. Verify PostgreSQL is running: `docker ps | grep postgres`
5. Test API endpoints directly: `curl http://localhost:8080/api/widgets/weather`

## References

- Echo: https://echo.labstack.com/
- templ: https://templ.guide/
- HTMX: https://htmx.org/
- TailwindCSS: https://tailwindcss.com/
- Open-Meteo: https://open-meteo.com/
- gopsutil: https://github.com/shirou/gopsutil
- pgx: https://github.com/jackc/pgx
