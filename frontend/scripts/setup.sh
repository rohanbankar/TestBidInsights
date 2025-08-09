#!/bin/bash

# OpenRTB Insights Frontend Setup Script

set -e

echo "🚀 Setting up OpenRTB Insights Frontend..."

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 18 or later."
    exit 1
fi

# Check Node.js version
NODE_VERSION=$(node -v | sed 's/v//')
MIN_VERSION="18.0.0"

if [ "$(printf '%s\n' "$MIN_VERSION" "$NODE_VERSION" | sort -V | head -n1)" != "$MIN_VERSION" ]; then
    echo "❌ Node.js version $NODE_VERSION is not supported. Please upgrade to Node.js $MIN_VERSION or later."
    exit 1
fi

echo "✅ Node.js version $NODE_VERSION is supported"

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    echo "❌ npm is not installed. Please install npm."
    exit 1
fi

echo "✅ npm is available"

# Copy environment file if it doesn't exist
if [ ! -f .env ]; then
    echo "📄 Creating .env file..."
    cat > .env << EOF
VITE_API_BASE_URL=http://localhost:8080/api
VITE_APP_NAME=OpenRTB Insights
EOF
    echo "✅ Created .env file"
else
    echo "✅ .env file already exists"
fi

# Install dependencies
echo "📦 Installing dependencies..."
npm ci
echo "✅ Dependencies installed successfully"

# Build the application for production (optional)
echo "🔨 Building production version..."
npm run build
echo "✅ Production build completed successfully"

# Run type checking
echo "🔍 Running type check..."
npx tsc --noEmit
echo "✅ Type check passed"

# Run linting
echo "🧹 Running linter..."
npm run lint
echo "✅ Linting passed"

echo ""
echo "🎉 Frontend setup completed successfully!"
echo ""
echo "Next steps:"
echo "  1. Review the .env file and update API URL if needed"
echo "  2. Start development server: npm run dev"
echo "  3. Or serve production build: npm run preview"
echo ""
echo "Development server will be available at: http://localhost:3000"
echo "Make sure the backend is running at: http://localhost:8080"
echo ""
echo "Available commands:"
echo "  • npm run dev     - Start development server"
echo "  • npm run build   - Build for production"
echo "  • npm run preview - Preview production build"
echo "  • npm run lint    - Run ESLint"
echo ""