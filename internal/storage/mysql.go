package storage

import "database/sql"

type MySQL struct {
	DB *sql.DB
}

func (m *MySQL) Init() error  { return nil }
func (m *MySQL) Close() error { return m.DB.Close() }

func (m *MySQL) CreateSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS notifications (
		notification_id CHAR(36) PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL,
		channel ENUM('email','sms','push') NOT NULL,
		payload JSON NOT NULL,
		status ENUM('pending','processing','sent','failed') DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS delivery_attempts (
		attempt_id CHAR(36) PRIMARY KEY,
		notification_id CHAR(36) NOT NULL,
		attempt_number INT NOT NULL,
		provider VARCHAR(255) NOT NULL,
		result BOOLEAN NOT NULL,
		error_message TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (notification_id)
			REFERENCES notifications(notification_id)
			ON DELETE CASCADE
	);
	`
	_, err := m.DB.Exec(schema)
	return err
}
