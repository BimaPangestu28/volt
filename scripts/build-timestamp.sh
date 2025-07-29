#!/bin/bash

# Volt Build with Timestamp
# Auto-generates version based on current timestamp

set -e

# Generate timestamp-based version
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
SHORT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
VERSION="v${TIMESTAMP}-${SHORT_COMMIT}"

echo "🚀 Building Volt with timestamp version: $VERSION"

# Build backend
echo "📦 Building backend..."
cd backend
docker build -t bimapangestu/volt-backend:$VERSION .
docker push bimapangestu/volt-backend:$VERSION
echo "✅ Backend pushed: bimapangestu/volt-backend:$VERSION"

# Build frontend
echo "📦 Building frontend..."
cd ../frontend
docker build -t bimapangestu/volt-frontend:$VERSION .
docker push bimapangestu/volt-frontend:$VERSION
echo "✅ Frontend pushed: bimapangestu/volt-frontend:$VERSION"

# Update docker-compose.prod.yml
echo "📝 Updating docker-compose.prod.yml..."
cd ..
sed -i.bak "s|bimapangestu/volt-backend:.*|bimapangestu/volt-backend:$VERSION|g" docker-compose.prod.yml
sed -i.bak "s|bimapangestu/volt-frontend:.*|bimapangestu/volt-frontend:$VERSION|g" docker-compose.prod.yml
rm docker-compose.prod.yml.bak

echo "🎉 Build complete with version: $VERSION"
echo "🌐 Deploy with: ./scripts/deploy.sh $VERSION"