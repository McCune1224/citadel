package database

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, databaseURL string) (*DB, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{pool: pool}, nil
}

func (d *DB) Close() {
	d.pool.Close()
}

func (d *DB) Migrate(ctx context.Context) error {
	// Read migration file
	migration, err := os.ReadFile("migrations/001_initial_schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Split into individual statements
	statements := strings.Split(string(migration), ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		if _, err := d.pool.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}

	return nil
}

// User queries
func (d *DB) GetUserByUsername(ctx context.Context, username string) (map[string]interface{}, error) {
	row := d.pool.QueryRow(ctx, "SELECT id, username, password_hash, tailscale_ip FROM users WHERE username = $1", username)

	var id int
	var uname, phash, tip string
	err := row.Scan(&id, &uname, &phash, &tip)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":            id,
		"username":      uname,
		"password_hash": phash,
		"tailscale_ip":  tip,
	}, nil
}

func (d *DB) CreateUser(ctx context.Context, username, passwordHash string) (int, error) {
	var id int
	err := d.pool.QueryRow(
		ctx,
		"INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id",
		username, passwordHash,
	).Scan(&id)
	return id, err
}

func (d *DB) CreateSession(ctx context.Context, userID int, token string, expiresAt string) error {
	_, err := d.pool.Exec(
		ctx,
		"INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3)",
		userID, token, expiresAt,
	)
	return err
}

func (d *DB) GetSessionByToken(ctx context.Context, token string) (int, error) {
	var userID int
	err := d.pool.QueryRow(
		ctx,
		"SELECT user_id FROM sessions WHERE token = $1 AND expires_at > NOW()",
		token,
	).Scan(&userID)
	return userID, err
}

func (d *DB) DeleteSession(ctx context.Context, token string) error {
	_, err := d.pool.Exec(ctx, "DELETE FROM sessions WHERE token = $1", token)
	return err
}

// Widget data queries
func (d *DB) SaveWidgetData(ctx context.Context, userID int, widgetName, widgetKey, valueJSON string) error {
	_, err := d.pool.Exec(
		ctx,
		`INSERT INTO widget_data (user_id, widget_name, widget_key, value_json)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (user_id, widget_name, widget_key) DO UPDATE SET value_json = $4, updated_at = NOW()`,
		userID, widgetName, widgetKey, valueJSON,
	)
	return err
}

func (d *DB) GetWidgetData(ctx context.Context, userID int, widgetName, widgetKey string) (string, error) {
	var value string
	err := d.pool.QueryRow(
		ctx,
		"SELECT value_json FROM widget_data WHERE user_id = $1 AND widget_name = $2 AND widget_key = $3",
		userID, widgetName, widgetKey,
	).Scan(&value)
	return value, err
}

// Settings queries
func (d *DB) SaveSetting(ctx context.Context, userID int, key, value string) error {
	_, err := d.pool.Exec(
		ctx,
		`INSERT INTO settings (user_id, setting_key, setting_value)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (user_id, setting_key) DO UPDATE SET setting_value = $3, updated_at = NOW()`,
		userID, key, value,
	)
	return err
}

func (d *DB) GetSetting(ctx context.Context, userID int, key string) (string, error) {
	var value string
	err := d.pool.QueryRow(
		ctx,
		"SELECT setting_value FROM settings WHERE user_id = $1 AND setting_key = $2",
		userID, key,
	).Scan(&value)
	return value, err
}
