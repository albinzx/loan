package sql

import (
	"database/sql"
	"log"
	"time"
)

const (
	// defaultMaxOpen is default value for max open connection
	defaultMaxOpen = 10
	// defaultMaxIdle is default value for max idle connection
	defaultMaxIdle = 10
	// defaultMaxLifetime is default value for max connection lifetime
	defaultMaxLifetime = time.Hour
	// defaultMaxIdleTime is default value for max connection idle time
	defaultMaxIdleTime = 10 * time.Minute
)

// DataSource is database source
type DataSource interface {
	// Name returns driver name and data source name
	Name() (string, string, error)
}

// config is a group of options for sql db
type config struct {
	maxOpen     int
	maxIdle     int
	maxLifetime time.Duration
	maxIdleTime time.Duration
}

// DB return new sql db
func DB(ds DataSource, options ...Option) (*sql.DB, error) {
	cfg := &config{}
	defaults(cfg)

	for _, option := range options {
		option(cfg)
	}

	driverName, dsn, err := ds.Name()

	if err != nil {
		return nil, err
	}

	var db *sql.DB

	if db, err = sql.Open(driverName, dsn); err != nil {
		log.Printf("error while opening DB: %v", err)
		return nil, err
	}

	db.SetMaxOpenConns(cfg.maxOpen)
	db.SetMaxIdleConns(cfg.maxIdle)
	db.SetConnMaxLifetime(cfg.maxLifetime)
	db.SetConnMaxIdleTime(cfg.maxIdleTime)

	return db, nil
}

// Option applies option values for a config.
type Option func(*config)

// defaults applies default configuration for sql db
func defaults(config *config) {
	config.maxOpen = defaultMaxOpen
	config.maxIdle = defaultMaxIdle
	config.maxLifetime = defaultMaxLifetime
	config.maxIdleTime = defaultMaxIdleTime
}

// WithConnection returns an option to set db connection parameter
func WithConnection(maxOpen, maxIdle int, maxLifetime, maxIdleTime time.Duration) Option {
	return func(config *config) {
		if maxOpen > 0 {
			config.maxOpen = maxOpen
		}
		if maxIdle > 0 {
			config.maxIdle = min(maxIdle, config.maxOpen)
		}
		if maxLifetime >= time.Minute {
			config.maxLifetime = maxLifetime
		}
		if maxIdleTime >= time.Minute {
			config.maxIdleTime = maxIdleTime
		}
	}
}
