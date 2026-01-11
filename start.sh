#!/bin/bash

# AI of the World - Backend Start Script

echo "🚀 Starting AI of the World Backend..."
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

# Check if .env file exists
if [ ! -f .env ]; then
    echo "⚠️  .env file not found. Creating from .env.example..."
    cp .env.example .env
    echo "✅ .env file created. Please update it with your configuration."
fi

# Install dependencies
echo "📦 Installing dependencies..."
go mod download

# Run the server
echo ""
echo "🎯 Starting server..."
go run main.go
