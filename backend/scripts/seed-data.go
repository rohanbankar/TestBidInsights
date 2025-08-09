package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"openrtb-insights/internal/config"
	"openrtb-insights/internal/database"

	_ "github.com/marcboeker/go-duckdb"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	// Run migrations first
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Seed data
	log.Println("Starting data seeding...")
	
	if err := seedPlatformStats(db); err != nil {
		log.Fatalf("Failed to seed platform stats: %v", err)
	}

	if err := seedContentHealth(db); err != nil {
		log.Fatalf("Failed to seed content health: %v", err)
	}

	if err := seedVideoHealth(db); err != nil {
		log.Fatalf("Failed to seed video health: %v", err)
	}

	log.Println("Data seeding completed successfully!")
}

func seedPlatformStats(db *sql.DB) error {
	log.Println("Seeding platform stats...")

	// Generate data for the last 30 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	stmt, err := db.Prepare(`
		INSERT OR REPLACE INTO platform_stats 
		(date, total_requests, multi_impression, big_guidance, addressable, 
		 compliance_strings, deals, tmax, invalid_requests, timeout_rate, bid_rate)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare platform stats statement: %w", err)
	}
	defer stmt.Close()

	for date := startDate; date.Before(endDate) || date.Equal(endDate); date = date.AddDate(0, 0, 1) {
		dateStr := date.Format("2006-01-02")
		
		// Generate realistic sample data with some randomness
		baseRequests := int64(8000 + rand.Intn(4000)) // 8K-12K requests per day
		
		_, err := stmt.Exec(
			dateStr,
			baseRequests,                              // total_requests
			int64(float64(baseRequests) * (0.15 + rand.Float64()*0.1)), // multi_impression (15-25%)
			int64(float64(baseRequests) * (0.30 + rand.Float64()*0.2)), // big_guidance (30-50%)
			int64(float64(baseRequests) * (0.70 + rand.Float64()*0.2)), // addressable (70-90%)
			int64(float64(baseRequests) * (0.80 + rand.Float64()*0.15)), // compliance_strings (80-95%)
			int64(200 + rand.Intn(300)),               // deals (200-500)
			int64(5000 + rand.Intn(15000)),           // tmax (5s-20s)
			int64(float64(baseRequests) * (0.01 + rand.Float64()*0.04)), // invalid_requests (1-5%)
			1.0 + rand.Float64()*4.0,                 // timeout_rate (1-5%)
			40.0 + rand.Float64()*30.0,               // bid_rate (40-70%)
		)
		if err != nil {
			return fmt.Errorf("failed to insert platform stats for %s: %w", dateStr, err)
		}
	}

	return nil
}

func seedContentHealth(db *sql.DB) error {
	log.Println("Seeding content health...")

	platforms := []string{"CTV", "Audio"}
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	stmt, err := db.Prepare(`
		INSERT OR REPLACE INTO content_health 
		(date, platform, total_requests, album, artist, cat, context, data, 
		 embeddable, episode, genre, id, kwarray, keywords, length, language, 
		 livestream, season, series, title, url, videoquality)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare content health statement: %w", err)
	}
	defer stmt.Close()

	for _, platform := range platforms {
		for date := startDate; date.Before(endDate) || date.Equal(endDate); date = date.AddDate(0, 0, 1) {
			dateStr := date.Format("2006-01-02")
			
			// Platform-specific base requests
			var baseRequests int64
			if platform == "CTV" {
				baseRequests = int64(3000 + rand.Intn(2000)) // 3K-5K for CTV
			} else {
				baseRequests = int64(2000 + rand.Intn(1500)) // 2K-3.5K for Audio
			}
			
			// Generate content field data based on platform
			albumCount := int64(0)
			artistCount := int64(0)
			if platform == "Audio" {
				albumCount = int64(float64(baseRequests) * (0.60 + rand.Float64()*0.3))
				artistCount = int64(float64(baseRequests) * (0.80 + rand.Float64()*0.15))
			}

			episodeCount := int64(0)
			seriesCount := int64(0)
			if platform == "CTV" {
				episodeCount = int64(float64(baseRequests) * (0.45 + rand.Float64()*0.25))
				seriesCount = int64(float64(baseRequests) * (0.40 + rand.Float64()*0.3))
			}
			
			_, err := stmt.Exec(
				dateStr,
				platform,
				baseRequests,                              // total_requests
				albumCount,                                // album
				artistCount,                               // artist
				int64(float64(baseRequests) * (0.70 + rand.Float64()*0.25)), // cat
				int64(float64(baseRequests) * (0.85 + rand.Float64()*0.10)), // context
				int64(float64(baseRequests) * (0.60 + rand.Float64()*0.30)), // data
				int64(float64(baseRequests) * (0.50 + rand.Float64()*0.30)), // embeddable
				episodeCount,                              // episode
				int64(float64(baseRequests) * (0.75 + rand.Float64()*0.20)), // genre
				int64(float64(baseRequests) * (0.90 + rand.Float64()*0.08)), // id
				int64(float64(baseRequests) * (0.30 + rand.Float64()*0.40)), // kwarray
				int64(float64(baseRequests) * (0.65 + rand.Float64()*0.25)), // keywords
				int64(float64(baseRequests) * (0.80 + rand.Float64()*0.15)), // length
				int64(float64(baseRequests) * (0.85 + rand.Float64()*0.10)), // language
				int64(float64(baseRequests) * (0.05 + rand.Float64()*0.15)), // livestream
				int64(float64(baseRequests) * (0.20 + rand.Float64()*0.30)), // season
				seriesCount,                               // series
				int64(float64(baseRequests) * (0.95 + rand.Float64()*0.04)), // title
				int64(float64(baseRequests) * (0.90 + rand.Float64()*0.08)), // url
				int64(float64(baseRequests) * (0.70 + rand.Float64()*0.25)), // videoquality
			)
			if err != nil {
				return fmt.Errorf("failed to insert content health for %s/%s: %w", platform, dateStr, err)
			}
		}
	}

	return nil
}

func seedVideoHealth(db *sql.DB) error {
	log.Println("Seeding video health...")

	platforms := []string{"CTV", "Display", "App"}
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	stmt, err := db.Prepare(`
		INSERT OR REPLACE INTO video_health 
		(date, platform, percent_ctv, api, boxing_allowed, delivery, h, linearity, 
		 max_bitrate, max_duration, mimes, min_bitrate, min_cpm_per_sec, min_duration, 
		 placement, play_backend, pod_dur, pod_id, pos, protocols, rqd_durs, skip, 
		 skip_after, skip_min, slot_in_pod, start_delay, w, max_seq, companion_ad, 
		 companion_type, protocol, placement_type)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare video health statement: %w", err)
	}
	defer stmt.Close()

	for _, platform := range platforms {
		for date := startDate; date.Before(endDate) || date.Equal(endDate); date = date.AddDate(0, 0, 1) {
			dateStr := date.Format("2006-01-02")
			
			// Platform-specific CTV percentage
			var percentCTV float64
			switch platform {
			case "CTV":
				percentCTV = 85.0 + rand.Float64()*10.0 // 85-95%
			case "Display":
				percentCTV = 5.0 + rand.Float64()*15.0  // 5-20%
			case "App":
				percentCTV = 25.0 + rand.Float64()*35.0 // 25-60%
			}
			
			// Generate base counts for video properties
			baseCount := int64(1000 + rand.Intn(2000)) // 1K-3K base
			
			_, err := stmt.Exec(
				dateStr,
				platform,
				percentCTV,                                // percent_ctv
				int64(500 + rand.Intn(1000)),             // api
				int64(float64(baseCount) * (0.60 + rand.Float64()*0.30)), // boxing_allowed
				int64(float64(baseCount) * (0.70 + rand.Float64()*0.25)), // delivery
				int64(300 + rand.Intn(900)),              // h (height)
				int64(float64(baseCount) * (0.80 + rand.Float64()*0.15)), // linearity
				int64(2000 + rand.Intn(6000)),            // max_bitrate (2-8Mbps)
				int64(15 + rand.Intn(45)),                // max_duration (15-60s)
				int64(float64(baseCount) * (0.85 + rand.Float64()*0.12)), // mimes
				int64(200 + rand.Intn(800)),              // min_bitrate
				int64(1 + rand.Intn(20)),                 // min_cpm_per_sec
				int64(5 + rand.Intn(10)),                 // min_duration (5-15s)
				int64(float64(baseCount) * (0.75 + rand.Float64()*0.20)), // placement
				int64(float64(baseCount) * (0.65 + rand.Float64()*0.25)), // play_backend
				int64(30 + rand.Intn(600)),               // pod_dur (30s-10min)
				int64(1 + rand.Intn(10)),                 // pod_id
				int64(1 + rand.Intn(5)),                  // pos
				int64(float64(baseCount) * (0.90 + rand.Float64()*0.08)), // protocols
				int64(float64(baseCount) * (0.45 + rand.Float64()*0.35)), // rqd_durs
				int64(float64(baseCount) * (0.40 + rand.Float64()*0.45)), // skip
				int64(3 + rand.Intn(7)),                  // skip_after (3-10s)
				int64(1 + rand.Intn(4)),                  // skip_min (1-5s)
				int64(1 + rand.Intn(8)),                  // slot_in_pod
				int64(-1 + rand.Intn(20)),                // start_delay (-1 to 18)
				int64(400 + rand.Intn(1200)),             // w (width)
				int64(1 + rand.Intn(5)),                  // max_seq
				int64(float64(baseCount) * (0.20 + rand.Float64()*0.30)), // companion_ad
				int64(1 + rand.Intn(4)),                  // companion_type
				int64(1 + rand.Intn(8)),                  // protocol
				int64(1 + rand.Intn(4)),                  // placement_type
			)
			if err != nil {
				return fmt.Errorf("failed to insert video health for %s/%s: %w", platform, dateStr, err)
			}
		}
	}

	return nil
}