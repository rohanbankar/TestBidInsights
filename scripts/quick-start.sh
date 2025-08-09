#!/bin/bash

# OpenRTB Insights Quick Start Script

set -e

echo "🚀 OpenRTB Insights - Quick Start"
echo "================================="
echo ""

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
echo "🔍 Checking prerequisites..."

if command_exists docker && command_exists docker-compose; then
    echo "✅ Docker and Docker Compose are available"
    USE_DOCKER=true
elif command_exists go && command_exists node && command_exists npm; then
    echo "✅ Go and Node.js are available"
    USE_DOCKER=false
else
    echo "❌ Missing prerequisites!"
    echo "Please install either:"
    echo "  Option 1: Docker and Docker Compose (recommended)"
    echo "  Option 2: Go 1.21+ and Node.js 18+"
    exit 1
fi

echo ""

if [ "$USE_DOCKER" = true ]; then
    echo "🐳 Using Docker deployment"
    echo "========================"
    
    # Check if docker-compose.yml exists
    if [ ! -f docker-compose.yml ]; then
        echo "❌ docker-compose.yml not found. Please run this script from the project root."
        exit 1
    fi
    
    echo "🏗️  Building and starting services..."
    docker-compose down --remove-orphans 2>/dev/null || true
    docker-compose up -d --build
    
    echo "⏳ Waiting for services to be ready..."
    sleep 10
    
    # Wait for backend health check
    echo "🔍 Checking backend health..."
    for i in {1..30}; do
        if curl -f -s http://localhost:8080/health >/dev/null 2>&1; then
            echo "✅ Backend is healthy"
            break
        fi
        if [ $i -eq 30 ]; then
            echo "❌ Backend health check failed"
            echo "Checking logs..."
            docker-compose logs backend
            exit 1
        fi
        sleep 2
    done
    
    # Wait for frontend health check
    echo "🔍 Checking frontend health..."
    for i in {1..20}; do
        if curl -f -s http://localhost/health >/dev/null 2>&1; then
            echo "✅ Frontend is healthy"
            break
        fi
        if [ $i -eq 20 ]; then
            echo "❌ Frontend health check failed"
            echo "Checking logs..."
            docker-compose logs frontend
            exit 1
        fi
        sleep 2
    done
    
    FRONTEND_URL="http://localhost"
    BACKEND_URL="http://localhost:8080"
    
else
    echo "💻 Using local development setup"
    echo "==============================="
    
    # Setup backend
    echo "🔧 Setting up backend..."
    cd backend
    if [ -f scripts/setup.sh ]; then
        chmod +x scripts/setup.sh
        ./scripts/setup.sh
    else
        echo "📦 Installing backend dependencies..."
        go mod download
        echo "🎲 Generating sample data..."
        go run scripts/seed-data.go
    fi
    
    # Start backend in background
    echo "🚀 Starting backend..."
    go run cmd/server/main.go &
    BACKEND_PID=$!
    
    cd ..
    
    # Setup frontend
    echo "🔧 Setting up frontend..."
    cd frontend
    if [ -f scripts/setup.sh ]; then
        chmod +x scripts/setup.sh
        ./scripts/setup.sh
    else
        echo "📦 Installing frontend dependencies..."
        npm ci
        echo "🔨 Building frontend..."
        npm run build
    fi
    
    # Start frontend in background
    echo "🚀 Starting frontend..."
    npm run dev -- --port 3000 &
    FRONTEND_PID=$!
    
    cd ..
    
    # Wait for services
    echo "⏳ Waiting for services to start..."
    sleep 5
    
    FRONTEND_URL="http://localhost:3000"
    BACKEND_URL="http://localhost:8080"
    
    # Cleanup function for local development
    cleanup() {
        echo ""
        echo "🛑 Shutting down services..."
        if [ ! -z "$BACKEND_PID" ]; then
            kill $BACKEND_PID 2>/dev/null || true
        fi
        if [ ! -z "$FRONTEND_PID" ]; then
            kill $FRONTEND_PID 2>/dev/null || true
        fi
        echo "✅ Services stopped"
        exit 0
    }
    
    trap cleanup SIGINT SIGTERM
fi

echo ""
echo "🎉 OpenRTB Insights is now running!"
echo "=================================="
echo ""
echo "📊 Dashboard:     $FRONTEND_URL"
echo "🔌 API:           $BACKEND_URL"
echo "❤️  Health Check: $BACKEND_URL/health"
echo ""
echo "👤 Demo Accounts:"
echo "   Admin:   admin / admin123"
echo "   Analyst: analyst / admin123" 
echo "   Viewer:  viewer / admin123"
echo ""

if command_exists open; then
    echo "🌐 Opening dashboard in browser..."
    sleep 2
    open "$FRONTEND_URL"
elif command_exists xdg-open; then
    echo "🌐 Opening dashboard in browser..."
    sleep 2
    xdg-open "$FRONTEND_URL"
else
    echo "🌐 Open $FRONTEND_URL in your browser to get started"
fi

if [ "$USE_DOCKER" = true ]; then
    echo ""
    echo "📝 Useful Docker commands:"
    echo "   View logs:     docker-compose logs -f"
    echo "   Stop services: docker-compose down"
    echo "   Restart:       docker-compose restart"
    echo ""
    echo "Press Ctrl+C to stop the services"
    
    # Keep script running for Docker
    while true; do
        sleep 1
    done
else
    echo ""
    echo "🔧 Services are running in the background"
    echo "Press Ctrl+C to stop all services"
    echo ""
    
    # Wait for interrupt
    while true; do
        sleep 1
    done
fi