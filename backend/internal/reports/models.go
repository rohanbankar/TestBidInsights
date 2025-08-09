package reports

import "time"

type PlatformStats struct {
	Date               string  `json:"date" db:"date"`
	TotalRequests      int64   `json:"totalRequests" db:"total_requests"`
	MultiImpression    int64   `json:"multiImpression" db:"multi_impression"`
	BigGuidance        int64   `json:"bigGuidance" db:"big_guidance"`
	Addressable        int64   `json:"addressable" db:"addressable"`
	ComplianceStrings  int64   `json:"complianceStrings" db:"compliance_strings"`
	Deals              int64   `json:"deals" db:"deals"`
	Tmax               int64   `json:"tmax" db:"tmax"`
	InvalidRequests    int64   `json:"invalidRequests" db:"invalid_requests"`
	TimeoutRate        float64 `json:"timeoutRate" db:"timeout_rate"`
	BidRate            float64 `json:"bidRate" db:"bid_rate"`
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`
}

type ContentHealth struct {
	Date          string    `json:"date" db:"date"`
	Platform      string    `json:"platform" db:"platform"`
	TotalRequests int64     `json:"totalRequests" db:"total_requests"`
	Album         int64     `json:"album" db:"album"`
	Artist        int64     `json:"artist" db:"artist"`
	Cat           int64     `json:"cat" db:"cat"`
	Context       int64     `json:"context" db:"context"`
	Data          int64     `json:"data" db:"data"`
	Embeddable    int64     `json:"embeddable" db:"embeddable"`
	Episode       int64     `json:"episode" db:"episode"`
	Genre         int64     `json:"genre" db:"genre"`
	ID            int64     `json:"id" db:"id"`
	Kwarray       int64     `json:"kwarray" db:"kwarray"`
	Keywords      int64     `json:"keywords" db:"keywords"`
	Length        int64     `json:"length" db:"length"`
	Language      int64     `json:"language" db:"language"`
	Livestream    int64     `json:"livestream" db:"livestream"`
	Season        int64     `json:"season" db:"season"`
	Series        int64     `json:"series" db:"series"`
	Title         int64     `json:"title" db:"title"`
	URL           int64     `json:"url" db:"url"`
	VideoQuality  int64     `json:"videoquality" db:"videoquality"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
}

type VideoHealth struct {
	Date            string    `json:"date" db:"date"`
	Platform        string    `json:"platform" db:"platform"`
	PercentCTV      float64   `json:"percentCtv" db:"percent_ctv"`
	API             int64     `json:"api" db:"api"`
	BoxingAllowed   int64     `json:"boxingAllowed" db:"boxing_allowed"`
	Delivery        int64     `json:"delivery" db:"delivery"`
	H               int64     `json:"h" db:"h"`
	Linearity       int64     `json:"linearity" db:"linearity"`
	MaxBitrate      int64     `json:"maxBitrate" db:"max_bitrate"`
	MaxDuration     int64     `json:"maxDuration" db:"max_duration"`
	Mimes           int64     `json:"mimes" db:"mimes"`
	MinBitrate      int64     `json:"minBitrate" db:"min_bitrate"`
	MinCPMPerSec    int64     `json:"minCpmPerSec" db:"min_cpm_per_sec"`
	MinDuration     int64     `json:"minDuration" db:"min_duration"`
	Placement       int64     `json:"placement" db:"placement"`
	PlayBackend     int64     `json:"playBackend" db:"play_backend"`
	PodDur          int64     `json:"podDur" db:"pod_dur"`
	PodID           int64     `json:"podId" db:"pod_id"`
	Pos             int64     `json:"pos" db:"pos"`
	Protocols       int64     `json:"protocols" db:"protocols"`
	RqdDurs         int64     `json:"rqdDurs" db:"rqd_durs"`
	Skip            int64     `json:"skip" db:"skip"`
	SkipAfter       int64     `json:"skipAfter" db:"skip_after"`
	SkipMin         int64     `json:"skipMin" db:"skip_min"`
	SlotInPod       int64     `json:"slotInPod" db:"slot_in_pod"`
	StartDelay      int64     `json:"startDelay" db:"start_delay"`
	W               int64     `json:"w" db:"w"`
	MaxSeq          int64     `json:"maxSeq" db:"max_seq"`
	CompanionAd     int64     `json:"companionAd" db:"companion_ad"`
	CompanionType   int64     `json:"companionType" db:"companion_type"`
	Protocol        int64     `json:"protocol" db:"protocol"`
	PlacementType   int64     `json:"placementType" db:"placement_type"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
}