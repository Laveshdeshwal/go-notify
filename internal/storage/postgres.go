package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Postgres struct {
	DB *sql.DB
}

func (p *Postgres) Init() error {
	if p.DB != nil {
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
	p.DB.SetMaxOpenConns(10)
	p.DB.SetMaxIdleConns(5)
	p.DB.SetConnMaxLifetime(30 * time.Minute)

	// 🔍 Verify connectivity on startup
	if err = p.DB.Ping(); err != nil {
		return fmt.Errorf("db ping failed: %w", err)
	}

	fmt.Println("✅ Postgres DB connected successfully")
	return nil
}

func (p *Postgres) Close() error {
	if p.DB != nil {
		return p.DB.Close()
	}
	return nil
}

// GetDB exposes the reusable DB handle
func (p *Postgres) GetDB() *sql.DB {
	if p.DB == nil {
		panic("DB not initialized. Call db.Init() first")
	}
	return p.DB
}

func (p *Postgres) CreateSchema() error {
	schema := `
	CREATE EXTENSION IF NOT EXISTS "pgcrypto";

	CREATE TYPE IF NOT EXISTS notification_channel
	AS ENUM ('email','sms','push');

	CREATE TYPE IF NOT EXISTS notification_status
	AS ENUM ('pending','processing','sent','failed');

	CREATE TABLE IF NOT EXISTS notifications (
		notification_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id VARCHAR NOT NULL,
		channel notification_channel NOT NULL,
		payload JSONB NOT NULL,
		status notification_status NOT NULL DEFAULT 'pending',
		created_at TIMESTAMPTZ DEFAULT now()
	);

	CREATE TABLE IF NOT EXISTS delivery_attempts (
		attempt_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		notification_id UUID NOT NULL,
		attempt_number INT NOT NULL,
		provider VARCHAR NOT NULL,
		result BOOLEAN NOT NULL,
		error_message TEXT,
		created_at TIMESTAMPTZ DEFAULT now(),
		CONSTRAINT fk_notification
			FOREIGN KEY(notification_id)
			REFERENCES notifications(notification_id)
			ON DELETE CASCADE
	);
	`
	_, err := p.DB.Exec(schema)
	return err
}
