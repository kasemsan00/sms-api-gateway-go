package config

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var db *sqlx.DB

// Database holds the database connection pool
type Database struct {
	*sqlx.DB
}

// InitDatabase initializes the MySQL database connection pool
func InitDatabase(cfg *Config) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&multiStatements=true",
		cfg.MySQL.User,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.Database,
	)

	var err error
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().Msgf("Connected to MySQL database: %s:%s/%s", cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database)

	// Start keep-alive goroutine
	go keepAlive(db)

	return &Database{db}, nil
}

// GetDB returns the database connection
func GetDB() *sqlx.DB {
	return db
}

// keepAlive periodically pings the database to maintain connections
func keepAlive(db *sqlx.DB) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := db.PingContext(ctx); err != nil {
			log.Error().Err(err).Msg("Database keep-alive ping failed")
		}
		cancel()
	}
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.DB != nil {
		log.Info().Msg("Closing database connection pool")
		return d.DB.Close()
	}
	return nil
}

// Health checks the database health
func (d *Database) Health(ctx context.Context) error {
	return d.PingContext(ctx)
}

// BeginTx starts a new transaction
func (d *Database) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return d.BeginTxx(ctx, nil)
}
