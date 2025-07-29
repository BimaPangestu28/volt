#!/bin/bash

# Volt Build and Push Script
# Usage: ./scripts/build-and-push.sh [version]
# If no version provided, auto-increment from latest git tag

set -e

# Function to get next version
get_next_version() {
    # Get latest git tag
    latest_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
    
    # Extract version numbers
    if [[ $latest_tag =~ ^v([0-9]+)\.([0-9]+)\.([0-9]+)$ ]]; then
        major=${BASH_REMATCH[1]}
        minor=${BASH_REMATCH[2]}
        patch=${BASH_REMATCH[3]}
        
        # Auto-increment patch version
        patch=$((patch + 1))
        echo "v${major}.${minor}.${patch}"
    else
        echo "v1.0.0"
    fi
}

# Get version
if [ -z "$1" ]; then
    VERSION=$(get_next_version)
    echo "🏷️  Auto-generated version: $VERSION"
    
    # Ask for confirmation
    read -p "Continue with version $VERSION? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "❌ Cancelled"
        exit 1
    fi
else
    VERSION=$1
fi
echo "🚀 Building and pushing Volt images with tag: $VERSION"

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

echo "🎉 Build complete!"

# Auto-commit and tag
echo "📝 Creating git tag..."
git add docker-compose.prod.yml
git commit -m "build: release $VERSION" || true
git tag $VERSION
git push origin main
git push origin $VERSION

echo ""
echo "✅ Release $VERSION completed!"
echo "🌐 Deploy with: ./scripts/deploy.sh $VERSION"