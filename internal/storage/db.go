package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var DB *sql.DB

// Init Initialize the database connection and verifies the connectivity
func Init() error {
	if DB != nil {
		return nil
	}
	host := viper.GetString("DB_HOST")
	port := viper.GetString("DB_PORT")
	user := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASSWORD")
	dbname := viper.GetString("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("sql.Open failed: %w", err)
	}

	// Connection pool settings
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(30 * time.Minute)

	// 🔍 Verify connectivity on startup
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("db ping failed: %w", err)
	}

	fmt.Println("✅ Postgres DB connected successfully")
	return nil
}

// GetDB exposes the reusable DB handle
func GetDB() *sql.DB {
	if DB == nil {
		panic("DB not initialized. Call db.Init() first")
	}
	return DB
}

// Close shuts down the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
