#!/bin/bash

# Volt Deployment Webhook Handler
# Simple webhook server untuk menerima deployment trigger dari GitHub Actions
# Install: chmod +x deploy-webhook.sh && ./deploy-webhook.sh

PORT=${WEBHOOK_PORT:-9000}
TOKEN=${WEBHOOK_TOKEN:-"your-secret-webhook-token"}

echo "🚀 Starting Volt Deployment Webhook Server on port $PORT"
echo "📡 Webhook URL: http://your-server:$PORT/deploy"
echo "🔑 Token: $TOKEN"
echo ""

# Function to handle deployment
deploy_volt() {
    local environment=$1
    local backend_image=$2
    local frontend_image=$3
    
    echo "🚀 Deploying Volt to $environment..."
    echo "Backend: $backend_image"
    echo "Frontend: $frontend_image"
    
    cd /opt/volt || exit 1
    
    if [ "$environment" == "production" ]; then
        # Production deployment
        source .env.production 2>/dev/null || echo "Warning: .env.production not found"
        
        echo "Pulling production images..."
        docker pull "$backend_image" || echo "Failed to pull backend image"
        docker pull "$frontend_image" || echo "Failed to pull frontend image"
        
        echo "Updating production services..."
        docker service update --image "$backend_image" volt-prod_backend || \
            docker stack deploy -c docker-compose.prod.yml volt-prod
        
        docker service update --image "$frontend_image" volt-prod_frontend || \
            docker stack deploy -c docker-compose.prod.yml volt-prod
        
        # Cleanup old images
        docker image prune -f
        
        echo "✅ Production deployment completed!"
        
    elif [ "$environment" == "staging" ]; then
        # Staging deployment
        source .env.staging 2>/dev/null || echo "Warning: .env.staging not found"
        
        echo "Pulling staging images..."
        docker pull "$backend_image" || echo "Failed to pull backend image"
        docker pull "$frontend_image" || echo "Failed to pull frontend image"
        
        echo "Updating staging services..."
        docker service update --image "$backend_image" volt-staging_backend || \
            docker stack deploy -c docker-compose.staging.yml volt-staging
        
        docker service update --image "$frontend_image" volt-staging_frontend || \
            docker stack deploy -c docker-compose.staging.yml volt-staging
        
        echo "✅ Staging deployment completed!"
        
    else
        echo "❌ Unknown environment: $environment"
        return 1
    fi
}

# Simple webhook server using netcat
while true; do
    echo "🔄 Waiting for deployment requests..."
    
    # Listen for HTTP requests
    response=$(echo -e "HTTP/1.1 200 OK\r\nContent-Type: application/json\r\nConnection: close\r\n\r\n{\"status\":\"received\"}" | nc -l -p $PORT)
    
    # Parse the request
    if echo "$response" | grep -q "POST /deploy"; then
        echo "📨 Deployment request received"
        
        # Extract authorization header
        auth_header=$(echo "$response" | grep -i "authorization:" | sed 's/.*Bearer //' | tr -d '\r\n ')
        
        if [ "$auth_header" == "$TOKEN" ]; then
            echo "✅ Token validated"
            
            # Extract JSON payload (this is a simple parser, might need improvement)
            json_payload=$(echo "$response" | tail -n 1)
            
            # Parse JSON manually (simple approach)
            environment=$(echo "$json_payload" | grep -o '"environment":"[^"]*"' | cut -d'"' -f4)
            backend_image=$(echo "$json_payload" | grep -o '"backend":"[^"]*"' | cut -d'"' -f4)
            frontend_image=$(echo "$json_payload" | grep -o '"frontend":"[^"]*"' | cut -d'"' -f4)
            
            if [ -n "$environment" ] && [ -n "$backend_image" ] && [ -n "$frontend_image" ]; then
                echo "Environment: $environment"
                echo "Backend: $backend_image" 
                echo "Frontend: $frontend_image"
                
                # Run deployment in background
                deploy_volt "$environment" "$backend_image" "$frontend_image" &
                
            else
                echo "❌ Invalid payload format"
            fi
            
        else
            echo "❌ Invalid token"
        fi
        
    else
        echo "📨 Non-deployment request received"
    fi
    
    sleep 1
done