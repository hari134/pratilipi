package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// DB is the global database instance type alias for *bun.DB.
type DB = bun.DB

// Config holds the database configuration details.
type Config struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

// InitDB initializes the Bun database connection.
func InitDB(cfg Config) *DB {
	// Construct the DSN (Data Source Name)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	// Create an sql.DB instance using pgdriver
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// Create a Bun DB instance using the alias
	db := bun.NewDB(sqldb, pgdialect.New())

	// Set connection pool settings
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Println("Database connection established with Bun ORM!")
	return db
}

// CloseDB closes the Bun database connection.
func CloseDB(db *DB) {
	if err := db.Close(); err != nil {
		log.Printf("Error closing the database connection: %v", err)
	}
	log.Println("Database connection closed.")
}
