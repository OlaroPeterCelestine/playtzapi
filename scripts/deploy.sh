#!/bin/bash
# Railway deployment script for Playtz API

echo "🚀 Starting Playtz API deployment..."

# Run database migrations
echo "📊 Running database migrations..."
./main migrate 2>/dev/null || echo "Note: Migration command not available, migrations will run on startup"

# Start the server
echo "✅ Starting server..."
exec ./main



