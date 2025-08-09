# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Structure

This repository contains two main components:

1. **Root Level**: Go-based OpenRTB request processing and analytics engine
   - `test.go`: Main OpenRTB edge agent implementation with parameter tracking and dimensional analysis
   
2. **Next.js Example Application** (`testing/rill-embedding-example/`):
   - React/Next.js 15 application demonstrating Rill dashboard embedding
   - Uses TypeScript, Tailwind CSS, and modern React patterns

## Development Commands

### Next.js Application (testing/rill-embedding-example/)
```bash
# Navigate to the Next.js app directory first
cd testing/rill-embedding-example/

# Development
npm run dev          # Start development server with turbopack
npm run build        # Build for production
npm run start        # Start production server
npm run lint         # Run ESLint
```

### Go Application (Root)
```bash
# Run the OpenRTB analytics test
go run test.go
```

## Architecture Overview

### OpenRTB Analytics Engine (test.go)
- **EdgeAgent**: Core component that processes OpenRTB bid requests and tracks parameter usage
- **Multi-dimensional Analysis**: Tracks parameters across device types, geography, request types, and content characteristics  
- **Memory Management**: Built-in limits for parameters (1000), dimension combinations (50), and sample values (5)
- **Sample Data**: Includes 6 different request types (CTV, mobile, desktop, etc.) with realistic OpenRTB payloads
- **Cloud Simulation**: Periodically flushes aggregated metrics to simulated cloud endpoint

Key data structures:
- `ParameterMetric`: Tracks presence, sample values, and dimensional breakdowns for each JSON path
- `DimensionKey`: Represents combinations of contextual dimensions (device type + country, etc.)
- `CloudPayload`: Aggregated data format for transmission to analytics backend

### Next.js Rill Embedding App
- **Component Architecture**: 
  - `IframeFetcher`: Handles API calls to generate Rill embed URLs
  - `RillFrame`: Renders the embedded Rill dashboard iframes
  - Page components demonstrate different embedding scenarios
- **API Integration**: `/api/get-iframe` endpoint for server-side iframe URL generation
- **Feature Demonstrations**: Navigation controls, row access policies, custom views, canvas dashboards

## Key Features Demonstrated

### OpenRTB Processing
- Recursive JSON parameter path extraction  
- Device type inference from user agent and device properties
- Content-aware dimension extraction (live vs VOD, series info)
- Video placement and skippability detection
- Geographic and temporal dimension tracking

### Rill Dashboard Embedding
- Multiple embed configurations (navigation enabled/disabled, pivot controls)
- Row-level security with custom attributes
- Canvas dashboard embedding
- Real-time iframe URL generation via API

## File Organization

- Root `test.go`: Complete OpenRTB analytics implementation
- `testing/rill-embedding-example/src/app/`: Next.js app router pages and components
- `testing/rill-embedding-example/src/app/api/`: API routes for iframe generation
- Individual page directories show specific embedding scenarios (navigation, row access policies, etc.)