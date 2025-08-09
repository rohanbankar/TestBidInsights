#!/bin/bash

# OpenRTB Insights Backend Setup Script

set -e

echo "🚀 Setting up OpenRTB Insights Backend..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
MIN_VERSION="1.21"

if [ "$(printf '%s\n' "$MIN_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$MIN_VERSION" ]; then
    echo "❌ Go version $GO_VERSION is not supported. Please upgrade to Go $MIN_VERSION or later."
    exit 1
fi

echo "✅ Go version $GO_VERSION is supported"

# Create necessary directories
echo "📁 Creating directories..."
mkdir -p data
mkdir -p logs

# Copy environment file if it doesn't exist
if [ ! -f .env ]; then
    echo "📄 Creating .env file..."
    cat > .env << EOF
# Database configuration
DB_PATH=./data/analytics.db

# JWT configuration
JWT_SECRET=$(openssl rand -hex 32)
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h

# Server configuration
PORT=8080
CORS_ORIGINS=http://localhost:3000,http://localhost:5173

# Logging
LOG_LEVEL=info

# Rate limiting (requests per minute per IP)
RATE_LIMIT=100
EOF
    echo "✅ Created .env file with secure JWT secret"
else
    echo "✅ .env file already exists"
fi

# Download dependencies
echo "📦 Downloading Go dependencies..."
go mod download
go mod verify

echo "✅ Dependencies downloaded successfully"

# Build the application
echo "🔨 Building the application..."
go build -o bin/server ./cmd/server
echo "✅ Application built successfully"

# Generate sample data
echo "🎲 Generating sample data..."
go run scripts/seed-data.go
echo "✅ Sample data generated successfully"

# Set proper permissions
chmod +x bin/server
chmod 755 data/

echo ""
echo "🎉 Backend setup completed successfully!"
echo ""
echo "Next steps:"
echo "  1. Review the .env file and update configurations as needed"
echo "  2. Start the server: ./bin/server"
echo "  3. Or use Go directly: go run cmd/server/main.go"
echo ""
echo "The server will be available at: http://localhost:8080"
echo "Health check: http://localhost:8080/health"
echo ""
echo "Default user accounts:"
echo "  • Admin:   admin / admin123"
echo "  • Analyst: analyst / admin123"
echo "  • Viewer:  viewer / admin123"
echo ""