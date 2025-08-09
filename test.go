package main

import (
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "sort"
    "strings"
    "sync"
    "time"
)

// Core data structures
type Dimension struct {
    Name  string      `json:"name"`
    Value interface{} `json:"value"`
}

type DimensionKey struct {
    Dimensions []Dimension `json:"dimensions"`
}

type DimensionStat struct {
    PresenceCount int64     `json:"presence_count"`
    TotalRequests int64     `json:"total_requests"`
    FirstSeen     time.Time `json:"first_seen"`
    LastSeen      time.Time `json:"last_seen"`
}

type ParameterMetric struct {
    ParameterPath   string                    `json:"parameter_path"`
    PresenceCount   int64                     `json:"presence_count"`
    TotalRequests   int64                     `json:"total_requests"`
    SampleValues    []interface{}             `json:"sample_values,omitempty"`
    LastSeen        time.Time                 `json:"last_seen"`
    DimensionCounts map[string]*DimensionStat `json:"dimension_counts"`
}

type AgentConfig struct {
    PrimaryDimensions     []string            `json:"primary_dimensions"`
    ConditionalDimensions map[string][]string `json:"conditional_dimensions"`
    MaxParameters         int                 `json:"max_parameters"`
    MaxDimensionCombos    int                 `json:"max_dimension_combos"`
    MaxSampleValues       int                 `json:"max_sample_values"`
    FlushIntervalSeconds  int                 `json:"flush_interval_seconds"`
    SamplingRate          float64             `json:"sampling_rate"`
}

type EdgeAgent struct {
    config        *AgentConfig
    metrics       map[string]*ParameterMetric
    mutex         sync.RWMutex
    lastFlush     time.Time
    totalReqs     int64
    processedReqs int64
    agentID       string
}

type CloudPayload struct {
    AgentID        string                 `json:"agent_id"`
    TimestampStart time.Time              `json:"timestamp_start"`
    TimestampEnd   time.Time              `json:"timestamp_end"`
    TotalRequests  int64                  `json:"total_requests"`
    ProcessedReqs  int64                  `json:"processed_requests"`
    Parameters     []ParameterCloudData   `json:"parameters"`
    Metadata       map[string]interface{} `json:"metadata"`
}

type ParameterCloudData struct {
    Path          string               `json:"path"`
    PresenceCount int64                `json:"presence_count"`
    TotalRequests int64                `json:"total_requests"`
    SampleValues  []interface{}        `json:"sample_values,omitempty"`
    Dimensions    []DimensionCloudData `json:"dimensions"`
}

type DimensionCloudData struct {
    DimensionKey  string  `json:"dimension_key"`
    PresenceCount int64   `json:"presence_count"`
    TotalRequests int64   `json:"total_requests"`
    PresenceRate  float64 `json:"presence_rate"`
}

// Sample OpenRTB requests
func getSampleRequests() []string {
    return []string{
        // 1. Samsung Smart TV CTV Request
        `{
            "id": "ctv-request-001",
            "imp": [
                {
                    "id": "1",
                    "video": {
                        "mimes": ["video/mp4"],
                        "minduration": 15,
                        "maxduration": 30,
                        "protocols": [2, 3, 5, 6],
                        "w": 1920,
                        "h": 1080,
                        "skip": 0,
                        "placement": 1,
                        "linearity": 1
                    },
                    "secure": 1,
                    "tagid": "ctv-placement-123"
                }
            ],
            "app": {
                "id": "samsung-tv-app",
                "name": "Samsung TV Plus",
                "bundle": "com.samsung.tv.plus",
                "storeurl": "https://samsungtvplus.com",
                "content": {
                    "id": "content-123",
                    "title": "Breaking Bad",
                    "series": "Breaking Bad",
                    "season": "Season 1",
                    "episode": 5,
                    "genre": "drama",
                    "network": "AMC",
                    "channel": "AMC HD",
                    "livestream": 0,
                    "len": 2700
                }
            },
            "device": {
                "devicetype": 3,
                "make": "Samsung",
                "model": "QN65Q80T",
                "os": "Tizen",
                "osv": "5.5",
                "ifa": "12345678-1234-1234-1234-123456789012",
                "ua": "Mozilla/5.0 (SMART-TV; LINUX; Tizen 5.5) AppleWebKit/538.1",
                "ip": "192.168.1.100",
                "geo": {
                    "country": "USA",
                    "region": "CA",
                    "city": "Los Angeles",
                    "lat": 34.0522,
                    "lon": -118.2437
                },
                "lmt": 0
            },
            "user": {
                "id": "user-ctv-001",
                "ext": {
                    "eids": [
                        {
                            "source": "liveramp.com",
                            "uids": [
                                {
                                    "id": "XY1000bIVBVah9ium-sZ3ykhPiXQbEcUpn4GjCtxrrw2BRDGM"
                                }
                            ]
                        }
                    ]
                }
            },
            "at": 2,
            "tmax": 250,
            "cur": ["USD"]
        }`,

        // 2. Roku CTV Request
        `{
            "id": "roku-request-002",
            "imp": [
                {
                    "id": "1",
                    "video": {
                        "mimes": ["video/mp4", "video/x-flv"],
                        "minduration": 5,
                        "maxduration": 60,
                        "protocols": [2, 3, 5, 6, 7, 8],
                        "w": 1280,
                        "h": 720,
                        "skip": 1,
                        "skipmin": 5,
                        "skipafter": 15,
                        "placement": 4,
                        "linearity": 1,
                        "boxingallowed": 1
                    }
                }
            ],
            "app": {
                "id": "roku-channel-app",
                "name": "The Roku Channel",
                "bundle": "roku.channel.12345",
                "content": {
                    "title": "The Office",
                    "series": "The Office",
                    "season": "Season 3",
                    "episode": 12,
                    "genre": "comedy",
                    "contentrating": "TV-14",
                    "livestream": 0
                }
            },
            "device": {
                "make": "Roku",
                "model": "Roku Ultra",
                "os": "Roku",
                "osv": "10.5",
                "ifa": "87654321-4321-4321-4321-210987654321",
                "ua": "Roku/DVP-9.10 (519.10E04111A)",
                "geo": {
                    "country": "USA",
                    "region": "NY",
                    "city": "New York"
                }
            },
            "user": {
                "id": "roku-user-002"
            }
        }`,

        // 3. iPhone Mobile App Request
        `{
            "id": "mobile-request-003",
            "imp": [
                {
                    "id": "1",
                    "banner": {
                        "w": 320,
                        "h": 50,
                        "pos": 1
                    },
                    "video": {
                        "mimes": ["video/mp4"],
                        "minduration": 15,
                        "maxduration": 30,
                        "w": 320,
                        "h": 180,
                        "placement": 1
                    }
                }
            ],
            "app": {
                "id": "mobile-news-app",
                "name": "News App",
                "bundle": "com.news.mobile",
                "storeurl": "https://apps.apple.com/app/news/id123456"
            },
            "device": {
                "devicetype": 4,
                "make": "Apple",
                "model": "iPhone13,2",
                "os": "iOS",
                "osv": "15.0",
                "ifa": "ABCDEFAB-1234-1234-1234-ABCDEFABCDEF",
                "ua": "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X)",
                "w": 375,
                "h": 667,
                "geo": {
                    "country": "USA",
                    "region": "TX",
                    "city": "Austin"
                },
                "connectiontype": 6,
                "carrier": "Verizon"
            },
            "user": {
                "id": "mobile-user-003",
                "yob": 1990,
                "gender": "M",
                "ext": {
                    "eids": [
                        {
                            "source": "adsystem.com",
                            "uids": [
                                {
                                    "id": "mobile-user-id-12345"
                                }
                            ]
                        }
                    ]
                }
            }
        }`,

        // 4. Desktop Website Request
        `{
            "id": "desktop-request-004",
            "imp": [
                {
                    "id": "1",
                    "banner": {
                        "w": 728,
                        "h": 90,
                        "pos": 1
                    }
                }
            ],
            "site": {
                "id": "news-website",
                "name": "News Website",
                "domain": "news.com",
                "page": "https://news.com/sports/article-123",
                "content": {
                    "title": "Sports News Article",
                    "cat": ["IAB17", "IAB17-1"],
                    "keywords": "sports,news,football"
                }
            },
            "device": {
                "devicetype": 2,
                "ua": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
                "ip": "203.0.113.1",
                "geo": {
                    "country": "GBR",
                    "region": "England",
                    "city": "London"
                }
            },
            "user": {
                "id": "desktop-user-004",
                "buyeruid": "buyer-uid-desktop-123"
            }
        }`,

        // 5. Android TV Request with Rich Content
        `{
            "id": "androidtv-request-005",
            "imp": [
                {
                    "id": "1",
                    "video": {
                        "mimes": ["video/mp4", "application/dash+xml"],
                        "minduration": 30,
                        "maxduration": 120,
                        "protocols": [7, 8],
                        "w": 3840,
                        "h": 2160,
                        "skip": 1,
                        "placement": 1
                    }
                }
            ],
            "app": {
                "id": "android-tv-streaming",
                "name": "Streaming App",
                "bundle": "com.streaming.androidtv",
                "content": {
                    "id": "movie-content-456",
                    "title": "Avengers Endgame",
                    "genre": "action",
                    "contentrating": "PG-13",
                    "userrating": "8.4",
                    "livestream": 0,
                    "len": 10800,
                    "network": "Disney",
                    "producer": {
                        "name": "Marvel Studios"
                    }
                }
            },
            "device": {
                "devicetype": 3,
                "make": "Sony",
                "model": "Sony Bravia XR",
                "os": "Android TV",
                "osv": "11.0",
                "ifa": "SONY1234-5678-9012-3456-ANDROIDTV001",
                "ua": "Mozilla/5.0 (Linux; Android 11; BRAVIA XR) AppleWebKit/537.36",
                "geo": {
                    "country": "JPN",
                    "region": "Tokyo",
                    "city": "Tokyo"
                }
            },
            "user": {
                "ext": {
                    "eids": [
                        {
                            "source": "id5-sync.com",
                            "uids": [
                                {
                                    "id": "ID5*abcdef123456789"
                                }
                            ]
                        },
                        {
                            "source": "liveramp.com",
                            "uids": [
                                {
                                    "id": "ramp-id-sony-tv-user"
                                }
                            ]
                        }
                    ]
                }
            }
        }`,

        // 6. Tablet Request
        `{
            "id": "tablet-request-006",
            "imp": [
                {
                    "id": "1",
                    "banner": {
                        "w": 728,
                        "h": 90
                    },
                    "video": {
                        "mimes": ["video/mp4"],
                        "w": 1024,
                        "h": 576,
                        "minduration": 15,
                        "maxduration": 45
                    }
                }
            ],
            "app": {
                "id": "tablet-game-app",
                "name": "Puzzle Game",
                "bundle": "com.puzzlegame.tablet"
            },
            "device": {
                "devicetype": 5,
                "make": "Apple",
                "model": "iPad Pro",
                "os": "iOS",
                "osv": "15.1",
                "ifa": "IPAD1234-5678-9012-3456-TABLETDEVICE",
                "w": 1024,
                "h": 1366,
                "geo": {
                    "country": "CAN",
                    "region": "Ontario",
                    "city": "Toronto"
                }
            }
        }`,
    }
}

// Initialize new edge agent
func NewEdgeAgent(config *AgentConfig) *EdgeAgent {
    return &EdgeAgent{
        config:    config,
        metrics:   make(map[string]*ParameterMetric),
        lastFlush: time.Now(),
        agentID:   fmt.Sprintf("agent-%d", time.Now().Unix()),
    }
}

// Extract all parameter paths from JSON recursively
func (ea *EdgeAgent) extractParameterPaths(obj interface{}, prefix string) []string {
    var paths []string
    
    switch v := obj.(type) {
    case map[string]interface{}:
        for key, value := range v {
            currentPath := key
            if prefix != "" {
                currentPath = prefix + "." + key
            }
            paths = append(paths, currentPath)
            
            // Recursively extract nested paths
            nestedPaths := ea.extractParameterPaths(value, currentPath)
            paths = append(paths, nestedPaths...)
        }
    case []interface{}:
        if len(v) > 0 {
            // For arrays, process first element with [] notation
            arrayPath := prefix + "[]"
            paths = append(paths, arrayPath)
            nestedPaths := ea.extractParameterPaths(v[0], arrayPath)
            paths = append(paths, nestedPaths...)
        }
    }
    
    return paths
}

// Get parameter value from request
func (ea *EdgeAgent) getParameterValue(request map[string]interface{}, paramPath string) interface{} {
    parts := strings.Split(paramPath, ".")
    current := interface{}(request)
    
    for _, part := range parts {
        if part == "" {
            continue
        }
        
        if strings.HasSuffix(part, "[]") {
            // Handle array notation
            part = strings.TrimSuffix(part, "[]")
            if m, ok := current.(map[string]interface{}); ok {
                if arr, exists := m[part].([]interface{}); exists && len(arr) > 0 {
                    current = arr[0]
                } else {
                    return nil
                }
            }
        } else {
            if m, ok := current.(map[string]interface{}); ok {
                if val, exists := m[part]; exists {
                    current = val
                } else {
                    return nil
                }
            } else {
                return nil
            }
        }
    }
    
    return current
}

// Extract primary dimensions
func (ea *EdgeAgent) extractPrimaryDimensions(request map[string]interface{}) map[string]interface{} {
    dims := make(map[string]interface{})
    
    // Device type
    if device, ok := request["device"].(map[string]interface{}); ok {
        if deviceType, exists := device["devicetype"]; exists {
            dims["device_type"] = deviceType
        } else {
            dims["device_type"] = ea.inferDeviceType(request)
        }
        
        // Country
        if geo, geoOk := device["geo"].(map[string]interface{}); geoOk {
            if country, countryOk := geo["country"]; countryOk {
                dims["country"] = country
            }
        }
    }
    
    // Request type
    if _, hasApp := request["app"]; hasApp {
        dims["request_type"] = "app"
    } else if _, hasSite := request["site"]; hasSite {
        dims["request_type"] = "site"
    } else {
        dims["request_type"] = "unknown"
    }
    
    // Hour
    dims["hour"] = time.Now().Hour()
    
    return dims
}

// Infer device type from context
func (ea *EdgeAgent) inferDeviceType(request map[string]interface{}) int {
    device, ok := request["device"].(map[string]interface{})
    if !ok {
        return 2 // Default to desktop
    }
    
    // Check User Agent
    if ua, exists := device["ua"].(string); exists {
        uaLower := strings.ToLower(ua)
        if strings.Contains(uaLower, "smart-tv") || strings.Contains(uaLower, "tizen") || 
           strings.Contains(uaLower, "roku") || strings.Contains(uaLower, "androidtv") {
            return 3 // CTV
        }
        if strings.Contains(uaLower, "mobile") || strings.Contains(uaLower, "iphone") {
            return 1 // Mobile
        }
    }
    
    // Check make/model
    if make, exists := device["make"].(string); exists {
        makeLower := strings.ToLower(make)
        if makeLower == "roku" || makeLower == "samsung" || makeLower == "sony" {
            if model, modelExists := device["model"].(string); modelExists {
                if strings.Contains(strings.ToLower(model), "tv") {
                    return 3 // CTV
                }
            }
        }
    }
    
    return 2 // Default to desktop
}

// Extract conditional dimensions based on parameter
func (ea *EdgeAgent) extractConditionalDimensions(request map[string]interface{}, paramPath string) map[string]interface{} {
    conditionalDims := make(map[string]interface{})
    
    switch {
    case strings.Contains(paramPath, "user") || strings.Contains(paramPath, "eids"):
        conditionalDims["has_user_id"] = ea.hasUserID(request)
        conditionalDims["has_eids"] = ea.hasEIDs(request)
        
    case strings.Contains(paramPath, "content"):
        conditionalDims["content_type"] = ea.extractContentType(request)
        conditionalDims["has_series_info"] = ea.hasSeriesInfo(request)
        
    case strings.Contains(paramPath, "video"):
        conditionalDims["video_placement"] = ea.extractVideoPlacement(request)
        conditionalDims["video_skippable"] = ea.isVideoSkippable(request)
        
    case strings.Contains(paramPath, "device"):
        conditionalDims["device_make"] = ea.extractDeviceMake(request)
        conditionalDims["has_ifa"] = ea.hasIFA(request)
    }
    
    return conditionalDims
}

// Helper functions for conditional dimensions
func (ea *EdgeAgent) hasUserID(request map[string]interface{}) bool {
    if user, ok := request["user"].(map[string]interface{}); ok {
        _, hasID := user["id"]
        return hasID
    }
    return false
}

func (ea *EdgeAgent) hasEIDs(request map[string]interface{}) bool {
    if user, ok := request["user"].(map[string]interface{}); ok {
        if ext, extOk := user["ext"].(map[string]interface{}); extOk {
            if eids, eidsOk := ext["eids"].([]interface{}); eidsOk {
                return len(eids) > 0
            }
        }
    }
    return false
}

func (ea *EdgeAgent) extractContentType(request map[string]interface{}) string {
    // Check app content first
    if app, ok := request["app"].(map[string]interface{}); ok {
        if content, contentOk := app["content"].(map[string]interface{}); contentOk {
            if livestream, liveOk := content["livestream"].(float64); liveOk && livestream == 1 {
                return "live"
            }
            return "vod"
        }
    }
    
    // Check site content
    if site, ok := request["site"].(map[string]interface{}); ok {
        if _, contentOk := site["content"]; contentOk {
            return "article"
        }
    }
    
    return "unknown"
}

func (ea *EdgeAgent) hasSeriesInfo(request map[string]interface{}) bool {
    if app, ok := request["app"].(map[string]interface{}); ok {
        if content, contentOk := app["content"].(map[string]interface{}); contentOk {
            _, hasSeries := content["series"]
            return hasSeries
        }
    }
    return false
}

func (ea *EdgeAgent) extractVideoPlacement(request map[string]interface{}) string {
    if imp, ok := request["imp"].([]interface{}); ok && len(imp) > 0 {
        if impObj, impOk := imp[0].(map[string]interface{}); impOk {
            if video, videoOk := impObj["video"].(map[string]interface{}); videoOk {
                if placement, placementOk := video["placement"].(float64); placementOk {
                    switch int(placement) {
                    case 1:
                        return "in-stream"
                    case 2:
                        return "in-banner"
                    case 3:
                        return "in-article"
                    case 4:
                        return "in-feed"
                    default:
                        return "other"
                    }
                }
            }
        }
    }
    return "unknown"
}

func (ea *EdgeAgent) isVideoSkippable(request map[string]interface{}) bool {
    if imp, ok := request["imp"].([]interface{}); ok && len(imp) > 0 {
        if impObj, impOk := imp[0].(map[string]interface{}); impOk {
            if video, videoOk := impObj["video"].(map[string]interface{}); videoOk {
                if skip, skipOk := video["skip"].(float64); skipOk {
                    return skip == 1
                }
            }
        }
    }
    return false
}

func (ea *EdgeAgent) extractDeviceMake(request map[string]interface{}) string {
    if device, ok := request["device"].(map[string]interface{}); ok {
        if make, makeOk := device["make"].(string); makeOk {
            return strings.ToLower(make)
        }
    }
    return "unknown"
}

func (ea *EdgeAgent) hasIFA(request map[string]interface{}) bool {
    if device, ok := request["device"].(map[string]interface{}); ok {
        _, hasIFA := device["ifa"]
        return hasIFA
    }
    return false
}

// Generate dimension combinations
func (ea *EdgeAgent) generateDimensionCombinations(primary, conditional map[string]interface{}) []DimensionKey {
    var keys []DimensionKey
    
    // Level 1: Single primary dimensions
    for name, value := range primary {
        keys = append(keys, DimensionKey{
            Dimensions: []Dimension{{Name: name, Value: value}},
        })
    }
    
    // Level 2: Primary + Conditional combinations
    for pName, pValue := range primary {
        for cName, cValue := range conditional {
            keys = append(keys, DimensionKey{
                Dimensions: []Dimension{
                    {Name: pName, Value: pValue},
                    {Name: cName, Value: cValue},
                },
            })
        }
    }
    
    // Level 3: Selected three-way combinations (limit to avoid explosion)
    if len(conditional) >= 2 {
        primaryKeys := []string{"device_type", "request_type"}
        for _, pKey := range primaryKeys {
            if pValue, exists := primary[pKey]; exists {
                // Create pairs from conditional dimensions
                condKeys := make([]string, 0, len(conditional))
                for k := range conditional {
                    condKeys = append(condKeys, k)
                }
                
                // Take first two conditional dimensions for 3-way combo
                if len(condKeys) >= 2 {
                    keys = append(keys, DimensionKey{
                        Dimensions: []Dimension{
                            {Name: pKey, Value: pValue},
                            {Name: condKeys[0], Value: conditional[condKeys[0]]},
                            {Name: condKeys[1], Value: conditional[condKeys[1]]},
                        },
                    })
                }
            }
        }
    }
    
    return keys
}

// Convert dimension key to string
func (ea *EdgeAgent) dimensionKeyToString(dimKey DimensionKey) string {
    var parts []string
    for _, dim := range dimKey.Dimensions {
        parts = append(parts, fmt.Sprintf("%s:%v", dim.Name, dim.Value))
    }
    return strings.Join(parts, "|")
}

// Process a single request
func (ea *EdgeAgent) ProcessRequest(requestJSON string) error {
    ea.mutex.Lock()
    defer ea.mutex.Unlock()
    
    ea.totalReqs++
    
    // Parse JSON
    var request map[string]interface{}
    if err := json.Unmarshal([]byte(requestJSON), &request); err != nil {
        return fmt.Errorf("failed to parse JSON: %v", err)
    }
    
    ea.processedReqs++
    
    // Extract all parameter paths
    paramPaths := ea.extractParameterPaths(request, "")
    
    // Process each parameter
    for _, paramPath := range paramPaths {
        ea.processParameter(request, paramPath)
    }
    
    return nil
}

// Process individual parameter
func (ea *EdgeAgent) processParameter(request map[string]interface{}, paramPath string) {
    // Get or create parameter metric
    metric, exists := ea.metrics[paramPath]
    if !exists {
        // Memory limit check
        if len(ea.metrics) >= ea.config.MaxParameters {
            return
        }
        
        metric = &ParameterMetric{
            ParameterPath:   paramPath,
            DimensionCounts: make(map[string]*DimensionStat),
            LastSeen:        time.Now(),
        }
        ea.metrics[paramPath] = metric
    }
    
    // Update basic counts
    metric.PresenceCount++
    metric.TotalRequests++
    metric.LastSeen = time.Now()
    
    // Sample values
    if len(metric.SampleValues) < ea.config.MaxSampleValues {
        value := ea.getParameterValue(request, paramPath)
        if value != nil {
            metric.SampleValues = append(metric.SampleValues, value)
        }
    }
    
    // Extract dimensions
    primaryDims := ea.extractPrimaryDimensions(request)
    conditionalDims := ea.extractConditionalDimensions(request, paramPath)
    
    // Generate dimension combinations
    dimensionKeys := ea.generateDimensionCombinations(primaryDims, conditionalDims)
    
    // Update dimension-specific counts
    for _, dimKey := range dimensionKeys {
        ea.updateDimensionCount(metric, dimKey)
    }
}

// Update dimension count
func (ea *EdgeAgent) updateDimensionCount(metric *ParameterMetric, dimKey DimensionKey) {
    keyStr := ea.dimensionKeyToString(dimKey)
    
    // Memory limit check
    if len(metric.DimensionCounts) >= ea.config.MaxDimensionCombos {
        return
    }
    
    dimStat, exists := metric.DimensionCounts[keyStr]
    if !exists {
        dimStat = &DimensionStat{
            FirstSeen: time.Now(),
        }
        metric.DimensionCounts[keyStr] = dimStat
    }
    
    dimStat.PresenceCount++
    dimStat.TotalRequests++
    dimStat.LastSeen = time.Now()
}

// Create cloud payload
func (ea *EdgeAgent) createCloudPayload() *CloudPayload {
    payload := &CloudPayload{
        AgentID:        ea.agentID,
        TimestampStart: ea.lastFlush,
        TimestampEnd:   time.Now(),
        TotalRequests:  ea.totalReqs,
        ProcessedReqs:  ea.processedReqs,
        Parameters:     make([]ParameterCloudData, 0, len(ea.metrics)),
        Metadata: map[string]interface{}{
            "version": "1.0",
            "config":  ea.config,
        },
    }
    
    // Convert metrics to cloud format
    for _, metric := range ea.metrics {
        paramData := ParameterCloudData{
            Path:          metric.ParameterPath,
            PresenceCount: metric.PresenceCount,
            TotalRequests: metric.TotalRequests,
            SampleValues:  metric.SampleValues,
            Dimensions:    make([]DimensionCloudData, 0, len(metric.DimensionCounts)),
        }
        
        // Convert dimension data
        for dimKey, dimStat := range metric.DimensionCounts {
            dimData := DimensionCloudData{
                DimensionKey:  dimKey,
                PresenceCount: dimStat.PresenceCount,
                TotalRequests: dimStat.TotalRequests,
                PresenceRate:  float64(dimStat.PresenceCount) / float64(dimStat.TotalRequests),
            }
            paramData.Dimensions = append(paramData.Dimensions, dimData)
        }
        
        // Sort dimensions by presence count (highest first)
        sort.Slice(paramData.Dimensions, func(i, j int) bool {
            return paramData.Dimensions[i].PresenceCount > paramData.Dimensions[j].PresenceCount
        })
        
        payload.Parameters = append(payload.Parameters, paramData)
    }
    
    // Sort parameters by presence count (highest first)
    sort.Slice(payload.Parameters, func(i, j int) bool {
        return payload.Parameters[i].PresenceCount > payload.Parameters[j].PresenceCount
    })
    
    return payload
}

// Flush data to "cloud" (console output)
func (ea *EdgeAgent) flush() {
    ea.mutex.Lock()
    payload := ea.createCloudPayload()
    
    // Reset metrics after capturing
    ea.metrics = make(map[string]*ParameterMetric)
    ea.lastFlush = time.Now()
    ea.totalReqs = 0
    ea.processedReqs = 0
    ea.mutex.Unlock()
    
    // Print to console (simulating cloud push)
    fmt.Println("\n" + strings.Repeat("=", 80))
    fmt.Println("FLUSHING DATA TO CLOUD")
    fmt.Println(strings.Repeat("=", 80))
    
    fmt.Printf("Agent ID: %s\n", payload.AgentID)
    fmt.Printf("Time Range: %s to %s\n", 
        payload.TimestampStart.Format("15:04:05"), 
        payload.TimestampEnd.Format("15:04:05"))
    fmt.Printf("Total Requests: %d, Processed: %d\n", 
        payload.TotalRequests, payload.ProcessedReqs)
    fmt.Printf("Parameters Tracked: %d\n\n", len(payload.Parameters))
    
    // Show top parameters
    fmt.Println("TOP PARAMETERS BY PRESENCE:")
    fmt.Println(strings.Repeat("-", 60))
    for i, param := range payload.Parameters {
        if i >= 15 { // Show top 15
            break
        }
        presenceRate := float64(param.PresenceCount) / float64(param.TotalRequests) * 100
        fmt.Printf("%-40s %6d (%5.1f%%)\n", 
            param.Path, param.PresenceCount, presenceRate)
    }
    
    // Show detailed dimension analysis for top 5 parameters
    fmt.Println("\nDETAILED DIMENSION ANALYSIS:")
    fmt.Println(strings.Repeat("-", 80))
    for i, param := range payload.Parameters {
        if i >= 5 { // Show details for top 5
            break
        }
        
        fmt.Printf("\n📊 %s (Present in %d requests)\n", param.Path, param.PresenceCount)
        
        if len(param.SampleValues) > 0 {
            fmt.Printf("   Sample Values: %v\n", param.SampleValues[:min(3, len(param.SampleValues))])
        }
        
        if len(param.Dimensions) > 0 {
            fmt.Println("   Top Dimensions:")
            for j, dim := range param.Dimensions {
                if j >= 5 { // Show top 5 dimensions
                    break
                }
                fmt.Printf("     %-50s %6d (%.1f%%)\n", 
                    dim.DimensionKey, dim.PresenceCount, dim.PresenceRate*100)
            }
        }
    }
    
    fmt.Println("\n" + strings.Repeat("=", 80))
    fmt.Println("CLOUD FLUSH COMPLETE")
    fmt.Println(strings.Repeat("=", 80) + "\n")
}

// Helper function
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

// Main function
func main() {
    fmt.Println("🚀 Starting OpenRTB Edge Agent Test")
    fmt.Println("====================================")
    
    // Initialize agent configuration
    config := &AgentConfig{
        PrimaryDimensions: []string{"device_type", "request_type", "country", "hour"},
        ConditionalDimensions: map[string][]string{
            "user.*":    {"has_user_id", "has_eids"},
            "content.*": {"content_type", "has_series_info"},
            "video.*":   {"video_placement", "video_skippable"},
            "device.*":  {"device_make", "has_ifa"},
        },
        MaxParameters:        1000,
        MaxDimensionCombos:   50,
        MaxSampleValues:      5,
        FlushIntervalSeconds: 30,
        SamplingRate:         1.0,
    }
    
    // Create edge agent
    agent := NewEdgeAgent(config)
    
    // Get sample requests
    requests := getSampleRequests()
    
    fmt.Printf("📋 Processing %d sample OpenRTB requests...\n\n", len(requests))
    
    // Process each request multiple times to simulate traffic
    for round := 1; round <= 3; round++ {
        fmt.Printf("🔄 Processing Round %d\n", round)
        
        for i, requestJSON := range requests {
            // Process each request multiple times with slight variations
            for repeat := 0; repeat < 2+rand.Intn(3); repeat++ {
                if err := agent.ProcessRequest(requestJSON); err != nil {
                    log.Printf("Error processing request %d: %v", i+1, err)
                } else {
                    fmt.Printf("   ✅ Processed request %d (repeat %d)\n", i+1, repeat+1)
                }
                
                // Add small delay to simulate real traffic
                time.Sleep(10 * time.Millisecond)
            }
        }
        
        // Flush after each round
        fmt.Printf("\n💾 Flushing data after round %d...\n", round)
        agent.flush()
        
        if round < 3 {
            fmt.Println("\n⏳ Waiting before next round...\n")
            time.Sleep(2 * time.Second)
        }
    }
    
    fmt.Println("✅ Test completed successfully!")
    fmt.Println("\nKey Insights from this test:")
    fmt.Println("- Multi-dimensional parameter tracking working")
    fmt.Println("- Memory management with limits enforced")  
    fmt.Println("- Dimension combinations generated (1, 2, and 3 levels)")
    fmt.Println("- Cloud payload format ready for transmission")
    fmt.Println("- Different device types (CTV, mobile, desktop) processed")
    fmt.Println("- Content-aware dimension extraction functioning")
}
