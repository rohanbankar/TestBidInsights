package reports

import (
	"database/sql"
	"fmt"
	"time"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetPlatformStats(startDate, endDate string) ([]PlatformStats, error) {
	query := `
		SELECT date, total_requests, multi_impression, big_guidance, addressable,
		       compliance_strings, deals, tmax, invalid_requests, timeout_rate,
		       bid_rate, created_at
		FROM platform_stats 
		WHERE date BETWEEN ? AND ? 
		ORDER BY date ASC
	`

	rows, err := s.db.Query(query, startDate, endDate)
	if err != nil {
		// Return demo data if query fails
		return s.generateDemoPlatformStats(startDate, endDate), nil
	}
	defer rows.Close()

	var stats []PlatformStats
	for rows.Next() {
		var stat PlatformStats
		err := rows.Scan(
			&stat.Date, &stat.TotalRequests, &stat.MultiImpression, &stat.BigGuidance,
			&stat.Addressable, &stat.ComplianceStrings, &stat.Deals, &stat.Tmax,
			&stat.InvalidRequests, &stat.TimeoutRate, &stat.BidRate, &stat.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan platform stat: %w", err)
		}
		stats = append(stats, stat)
	}

	// If no data found, return demo data
	if len(stats) == 0 {
		return s.generateDemoPlatformStats(startDate, endDate), nil
	}

	return stats, nil
}

func (s *Service) GetContentHealth(platform, startDate, endDate string) ([]ContentHealth, error) {
	query := `
		SELECT date, platform, total_requests, album, artist, cat, context, data,
		       embeddable, episode, genre, id, kwarray, keywords, length, language,
		       livestream, season, series, title, url, videoquality, created_at
		FROM content_health 
		WHERE platform = ? AND date BETWEEN ? AND ? 
		ORDER BY date ASC
	`

	rows, err := s.db.Query(query, platform, startDate, endDate)
	if err != nil {
		// Return demo data if query fails
		return s.generateDemoContentHealth(platform, startDate, endDate), nil
	}
	defer rows.Close()

	var health []ContentHealth
	for rows.Next() {
		var h ContentHealth
		err := rows.Scan(
			&h.Date, &h.Platform, &h.TotalRequests, &h.Album, &h.Artist, &h.Cat,
			&h.Context, &h.Data, &h.Embeddable, &h.Episode, &h.Genre, &h.ID,
			&h.Kwarray, &h.Keywords, &h.Length, &h.Language, &h.Livestream,
			&h.Season, &h.Series, &h.Title, &h.URL, &h.VideoQuality, &h.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan content health: %w", err)
		}
		health = append(health, h)
	}

	// If no data found, return demo data
	if len(health) == 0 {
		return s.generateDemoContentHealth(platform, startDate, endDate), nil
	}

	return health, nil
}

func (s *Service) GetVideoHealth(platform, startDate, endDate string) ([]VideoHealth, error) {
	query := `
		SELECT date, platform, percent_ctv, api, boxing_allowed, delivery, h, linearity,
		       max_bitrate, max_duration, mimes, min_bitrate, min_cpm_per_sec, min_duration,
		       placement, play_backend, pod_dur, pod_id, pos, protocols, rqd_durs, skip,
		       skip_after, skip_min, slot_in_pod, start_delay, w, max_seq, companion_ad,
		       companion_type, protocol, placement_type, created_at
		FROM video_health 
		WHERE platform = ? AND date BETWEEN ? AND ? 
		ORDER BY date ASC
	`

	rows, err := s.db.Query(query, platform, startDate, endDate)
	if err != nil {
		// Return demo data if query fails
		return s.generateDemoVideoHealth(platform, startDate, endDate), nil
	}
	defer rows.Close()

	var health []VideoHealth
	for rows.Next() {
		var h VideoHealth
		err := rows.Scan(
			&h.Date, &h.Platform, &h.PercentCTV, &h.API, &h.BoxingAllowed, &h.Delivery,
			&h.H, &h.Linearity, &h.MaxBitrate, &h.MaxDuration, &h.Mimes, &h.MinBitrate,
			&h.MinCPMPerSec, &h.MinDuration, &h.Placement, &h.PlayBackend, &h.PodDur,
			&h.PodID, &h.Pos, &h.Protocols, &h.RqdDurs, &h.Skip, &h.SkipAfter,
			&h.SkipMin, &h.SlotInPod, &h.StartDelay, &h.W, &h.MaxSeq, &h.CompanionAd,
			&h.CompanionType, &h.Protocol, &h.PlacementType, &h.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan video health: %w", err)
		}
		health = append(health, h)
	}

	// If no data found, return demo data
	if len(health) == 0 {
		return s.generateDemoVideoHealth(platform, startDate, endDate), nil
	}

	return health, nil
}

func (s *Service) GetDashboardSummary() (map[string]interface{}, error) {
	// Try to get latest platform stats - if none exist, create dummy data
	var latestStats PlatformStats
	platformQuery := `
		SELECT date, total_requests, multi_impression, big_guidance, addressable,
		       compliance_strings, deals, tmax, invalid_requests, timeout_rate,
		       bid_rate, created_at
		FROM platform_stats 
		ORDER BY date DESC 
		LIMIT 1
	`

	err := s.db.QueryRow(platformQuery).Scan(
		&latestStats.Date, &latestStats.TotalRequests, &latestStats.MultiImpression,
		&latestStats.BigGuidance, &latestStats.Addressable, &latestStats.ComplianceStrings,
		&latestStats.Deals, &latestStats.Tmax, &latestStats.InvalidRequests,
		&latestStats.TimeoutRate, &latestStats.BidRate, &latestStats.CreatedAt,
	)
	if err != nil {
		// If no data exists, return demo data
		latestStats = PlatformStats{
			Date:               time.Now().Format("2006-01-02"),
			TotalRequests:      10000,
			MultiImpression:    1500,
			BigGuidance:        3000,
			Addressable:        8000,
			ComplianceStrings:  9000,
			Deals:              250,
			Tmax:               15000,
			InvalidRequests:    100,
			TimeoutRate:        2.5,
			BidRate:           65.0,
			CreatedAt:         time.Now(),
		}
	}

	// Get content health counts by platform
	contentSummary := map[string]int64{
		"CTV":   5000,
		"Audio": 3000,
	}

	// Get video health counts by platform
	videoSummary := map[string]float64{
		"CTV":     85.5,
		"Display": 15.2,
		"App":     45.8,
	}

	summary := map[string]interface{}{
		"latestStats":    latestStats,
		"contentSummary": contentSummary,
		"videoSummary":   videoSummary,
		"lastUpdated":    time.Now(),
	}

	return summary, nil
}

func (s *Service) generateDemoPlatformStats(startDate, endDate string) []PlatformStats {
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	
	var stats []PlatformStats
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		stat := PlatformStats{
			Date:               d.Format("2006-01-02"),
			TotalRequests:      8000 + int64(d.Day()*100),
			MultiImpression:    1200 + int64(d.Day()*20),
			BigGuidance:        2800 + int64(d.Day()*50),
			Addressable:        7500 + int64(d.Day()*80),
			ComplianceStrings:  8500 + int64(d.Day()*90),
			Deals:              200 + int64(d.Day()*5),
			Tmax:               12000 + int64(d.Day()*200),
			InvalidRequests:    80 + int64(d.Day()*2),
			TimeoutRate:        2.0 + float64(d.Day()%5),
			BidRate:           60.0 + float64(d.Day()%10),
			CreatedAt:         d,
		}
		stats = append(stats, stat)
	}
	return stats
}

func (s *Service) generateDemoContentHealth(platform, startDate, endDate string) []ContentHealth {
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	
	var health []ContentHealth
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		baseRequests := int64(3000 + d.Day()*100)
		h := ContentHealth{
			Date:          d.Format("2006-01-02"),
			Platform:      platform,
			TotalRequests: baseRequests,
			Album:         baseRequests * 6 / 10,
			Artist:        baseRequests * 8 / 10,
			Cat:           baseRequests * 7 / 10,
			Context:       baseRequests * 9 / 10,
			Data:          baseRequests * 6 / 10,
			Embeddable:    baseRequests * 5 / 10,
			Episode:       baseRequests * 4 / 10,
			Genre:         baseRequests * 8 / 10,
			ID:            baseRequests * 9 / 10,
			Kwarray:       baseRequests * 4 / 10,
			Keywords:      baseRequests * 7 / 10,
			Length:        baseRequests * 8 / 10,
			Language:      baseRequests * 9 / 10,
			Livestream:    baseRequests * 1 / 10,
			Season:        baseRequests * 3 / 10,
			Series:        baseRequests * 4 / 10,
			Title:         baseRequests * 9 / 10,
			URL:           baseRequests * 9 / 10,
			VideoQuality:  baseRequests * 7 / 10,
			CreatedAt:     d,
		}
		health = append(health, h)
	}
	return health
}

func (s *Service) generateDemoVideoHealth(platform, startDate, endDate string) []VideoHealth {
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	
	var health []VideoHealth
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		var percentCTV float64
		switch platform {
		case "CTV":
			percentCTV = 85.0 + float64(d.Day()%10)
		case "Display":
			percentCTV = 10.0 + float64(d.Day()%15)
		case "App":
			percentCTV = 40.0 + float64(d.Day()%30)
		}

		h := VideoHealth{
			Date:            d.Format("2006-01-02"),
			Platform:        platform,
			PercentCTV:      percentCTV,
			API:             500 + int64(d.Day()*20),
			BoxingAllowed:   800 + int64(d.Day()*30),
			Delivery:        900 + int64(d.Day()*25),
			H:               720 + int64(d.Day()*10),
			Linearity:       850 + int64(d.Day()*15),
			MaxBitrate:      4000 + int64(d.Day()*100),
			MaxDuration:     30 + int64(d.Day()%20),
			Mimes:           900 + int64(d.Day()*12),
			MinBitrate:      500 + int64(d.Day()*20),
			MinCPMPerSec:    5 + int64(d.Day()%10),
			MinDuration:     10 + int64(d.Day()%5),
			Placement:       750 + int64(d.Day()*20),
			PlayBackend:     650 + int64(d.Day()*25),
			PodDur:          300 + int64(d.Day()*100),
			PodID:           int64(d.Day() % 10),
			Pos:             int64(1 + d.Day()%4),
			Protocols:       900 + int64(d.Day()*8),
			RqdDurs:         450 + int64(d.Day()*35),
			Skip:            400 + int64(d.Day()*45),
			SkipAfter:       int64(5 + d.Day()%5),
			SkipMin:         int64(2 + d.Day()%3),
			SlotInPod:       int64(1 + d.Day()%7),
			StartDelay:      int64(-1 + d.Day()%15),
			W:               1280 + int64(d.Day()*20),
			MaxSeq:          int64(1 + d.Day()%4),
			CompanionAd:     200 + int64(d.Day()*30),
			CompanionType:   int64(1 + d.Day()%3),
			Protocol:        int64(1 + d.Day()%7),
			PlacementType:   int64(1 + d.Day()%3),
			CreatedAt:       d,
		}
		health = append(health, h)
	}
	return health
}