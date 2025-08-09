package database

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			role VARCHAR(20) NOT NULL CHECK (role IN ('Viewer', 'Analyst', 'Admin')),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		`CREATE TABLE IF NOT EXISTS platform_stats (
			date DATE NOT NULL,
			total_requests BIGINT,
			multi_impression BIGINT,
			big_guidance BIGINT,
			addressable BIGINT,
			compliance_strings BIGINT,
			deals BIGINT,
			tmax BIGINT,
			invalid_requests BIGINT,
			timeout_rate DECIMAL(5,2),
			bid_rate DECIMAL(5,2),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (date)
		)`,
		
		`CREATE TABLE IF NOT EXISTS content_health (
			date DATE NOT NULL,
			platform VARCHAR(20) NOT NULL,
			total_requests BIGINT,
			album BIGINT,
			artist BIGINT,
			cat BIGINT,
			context BIGINT,
			data BIGINT,
			embeddable BIGINT,
			episode BIGINT,
			genre BIGINT,
			id BIGINT,
			kwarray BIGINT,
			keywords BIGINT,
			length BIGINT,
			language BIGINT,
			livestream BIGINT,
			season BIGINT,
			series BIGINT,
			title BIGINT,
			url BIGINT,
			videoquality BIGINT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (date, platform)
		)`,
		
		`CREATE TABLE IF NOT EXISTS video_health (
			date DATE NOT NULL,
			platform VARCHAR(20) NOT NULL,
			percent_ctv DECIMAL(5,2),
			api BIGINT,
			boxing_allowed BIGINT,
			delivery BIGINT,
			h BIGINT,
			linearity BIGINT,
			max_bitrate BIGINT,
			max_duration BIGINT,
			mimes BIGINT,
			min_bitrate BIGINT,
			min_cpm_per_sec BIGINT,
			min_duration BIGINT,
			placement BIGINT,
			play_backend BIGINT,
			pod_dur BIGINT,
			pod_id BIGINT,
			pos BIGINT,
			protocols BIGINT,
			rqd_durs BIGINT,
			skip BIGINT,
			skip_after BIGINT,
			skip_min BIGINT,
			slot_in_pod BIGINT,
			start_delay BIGINT,
			w BIGINT,
			max_seq BIGINT,
			companion_ad BIGINT,
			companion_type BIGINT,
			protocol BIGINT,
			placement_type BIGINT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (date, platform)
		)`,
		
		// Create indexes for better query performance
		`CREATE INDEX IF NOT EXISTS idx_platform_stats_date ON platform_stats(date)`,
		`CREATE INDEX IF NOT EXISTS idx_content_health_date_platform ON content_health(date, platform)`,
		`CREATE INDEX IF NOT EXISTS idx_video_health_date_platform ON video_health(date, platform)`,
		
		// Insert default admin user (password: admin123)
		`INSERT OR IGNORE INTO users (id, username, password_hash, role) VALUES 
		(1, 'admin', '$2a$10$ek0nw8RvUHOhqP9y48t6uusr3NUq0Zt8rLHKCn.UMVRmzyGEqZ..m', 'Admin')`,
		
		// Insert demo users
		`INSERT OR IGNORE INTO users (id, username, password_hash, role) VALUES 
		(2, 'analyst', '$2a$10$ek0nw8RvUHOhqP9y48t6uusr3NUq0Zt8rLHKCn.UMVRmzyGEqZ..m', 'Analyst')`,
		
		`INSERT OR IGNORE INTO users (id, username, password_hash, role) VALUES 
		(3, 'viewer', '$2a$10$ek0nw8RvUHOhqP9y48t6uusr3NUq0Zt8rLHKCn.UMVRmzyGEqZ..m', 'Viewer')`,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			log.Printf("Migration %d failed: %v", i+1, err)
			return err
		}
	}

	log.Println("All database migrations completed successfully")
	return nil
}