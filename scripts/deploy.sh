#!/bin/bash

# Volt - Docker Swarm Deployment Script
# This script deploys the Volt application using Docker Swarm

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
STACK_NAME="volt"
COMPOSE_FILE="docker-compose.prod.yml"
ENV_FILE=".env.prod"

# Functions
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}"
    exit 1
}

info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] INFO: $1${NC}"
}

# Check if Docker Swarm is initialized
check_swarm() {
    if ! docker info --format '{{.Swarm.LocalNodeState}}' | grep -q active; then
        error "Docker Swarm is not initialized. Run 'docker swarm init' first."
    fi
    log "Docker Swarm is active"
}

# Check if required files exist
check_files() {
    if [[ ! -f "$COMPOSE_FILE" ]]; then
        error "Compose file '$COMPOSE_FILE' not found"
    fi
    
    if [[ ! -f "$ENV_FILE" ]]; then
        error "Environment file '$ENV_FILE' not found. Copy from .env.prod.example and configure."
    fi
    
    log "Required files found"
}

# Load environment variables
load_env() {
    if [[ -f "$ENV_FILE" ]]; then
        set -a
        source "$ENV_FILE"
        set +a
        log "Environment variables loaded from $ENV_FILE"
    fi
}

# Pull latest images
pull_images() {
    log "Pulling latest images..."
    
    local backend_image="${DOCKER_REGISTRY:-bimapangestu}/volt-backend:${VERSION:-latest}"
    local frontend_image="${DOCKER_REGISTRY:-bimapangestu}/volt-frontend:${VERSION:-latest}"
    
    docker pull "$backend_image" || warn "Failed to pull backend image"
    docker pull "$frontend_image" || warn "Failed to pull frontend image"
    docker pull mongo:7.0 || warn "Failed to pull MongoDB image"
    docker pull redis:7.2-alpine || warn "Failed to pull Redis image"
    docker pull nginx:alpine || warn "Failed to pull Nginx image"
    
    log "Images pulled successfully"
}

# Deploy or update the stack
deploy_stack() {
    log "Deploying stack '$STACK_NAME'..."
    
    docker stack deploy \
        --compose-file "$COMPOSE_FILE" \
        --with-registry-auth \
        "$STACK_NAME"
    
    log "Stack '$STACK_NAME' deployed successfully"
}

# Check stack status
check_status() {
    log "Checking stack status..."
    
    echo ""
    info "Stack services:"
    docker stack services "$STACK_NAME"
    
    echo ""
    info "Stack tasks:"
    docker stack ps "$STACK_NAME" --no-trunc
    
    echo ""
    info "Waiting for services to be ready..."
    
    # Wait for services to be ready
    local max_attempts=30
    local attempt=1
    
    while [[ $attempt -le $max_attempts ]]; do
        local running_services=$(docker stack services "$STACK_NAME" --format "{{.Replicas}}" | grep -c "^[1-9]")
        local total_services=$(docker stack services "$STACK_NAME" --format "{{.Name}}" | wc -l)
        
        if [[ $running_services -eq $total_services ]]; then
            log "All services are running!"
            break
        fi
        
        info "Attempt $attempt/$max_attempts: $running_services/$total_services services running"
        sleep 10
        ((attempt++))
    done
    
    if [[ $attempt -gt $max_attempts ]]; then
        warn "Not all services are running after $max_attempts attempts"
    fi
}

# Health check
health_check() {
    log "Performing health checks..."
    
    # Wait a bit for services to fully start
    sleep 30
    
    # Check backend health
    info "Checking backend health..."
    if curl -f -s http://localhost:8080/health > /dev/null; then
        log "Backend health check passed"
    else
        warn "Backend health check failed"
    fi
    
    # Check frontend health
    info "Checking frontend health..."
    if curl -f -s http://localhost:80/health > /dev/null; then
        log "Frontend health check passed"
    else
        warn "Frontend health check failed"
    fi
}

# Remove stack
remove_stack() {
    log "Removing stack '$STACK_NAME'..."
    
    docker stack rm "$STACK_NAME"
    
    log "Stack '$STACK_NAME' removed successfully"
    
    # Wait for stack to be completely removed
    info "Waiting for stack to be completely removed..."
    while docker stack ls --format "{{.Name}}" | grep -q "^$STACK_NAME$"; do
        sleep 5
    done
    
    log "Stack removal completed"
}

# Show logs
show_logs() {
    local service_name="$1"
    
    if [[ -z "$service_name" ]]; then
        info "Available services:"
        docker stack services "$STACK_NAME" --format "{{.Name}}"
        return
    fi
    
    log "Showing logs for service: $service_name"
    docker service logs --follow "${STACK_NAME}_${service_name}"
}

# Scale service
scale_service() {
    local service_name="$1"
    local replicas="$2"
    
    if [[ -z "$service_name" || -z "$replicas" ]]; then
        error "Usage: $0 scale <service_name> <replicas>"
    fi
    
    log "Scaling service '$service_name' to $replicas replicas..."
    
    docker service scale "${STACK_NAME}_${service_name}=$replicas"
    
    log "Service scaled successfully"
}

# Update service
update_service() {
    local service_name="$1"
    local image="$2"
    
    if [[ -z "$service_name" ]]; then
        error "Usage: $0 update <service_name> [image]"
    fi
    
    log "Updating service '$service_name'..."
    
    if [[ -n "$image" ]]; then
        docker service update --image "$image" "${STACK_NAME}_${service_name}"
    else
        docker service update --force "${STACK_NAME}_${service_name}"
    fi
    
    log "Service updated successfully"
}

# Backup volumes
backup_volumes() {
    log "Creating backup of volumes..."
    
    local backup_dir="./backups/$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$backup_dir"
    
    # Backup MongoDB data
    info "Backing up MongoDB data..."
    docker run --rm \
        --volume "${STACK_NAME}_mongo_data:/data" \
        --volume "$backup_dir:/backup" \
        alpine tar czf /backup/mongo_data.tar.gz -C /data .
    
    # Backup Redis data
    info "Backing up Redis data..."
    docker run --rm \
        --volume "${STACK_NAME}_redis_data:/data" \
        --volume "$backup_dir:/backup" \
        alpine tar czf /backup/redis_data.tar.gz -C /data .
    
    log "Backup completed in: $backup_dir"
}

# Main script
main() {
    local command="$1"
    
    case "$command" in
        "deploy")
            check_swarm
            check_files
            load_env
            pull_images
            deploy_stack
            check_status
            health_check
            ;;
        "remove"|"down")
            remove_stack
            ;;
        "status")
            check_status
            ;;
        "logs")
            show_logs "$2"
            ;;
        "scale")
            scale_service "$2" "$3"
            ;;
        "update")
            update_service "$2" "$3"
            ;;
        "backup")
            backup_volumes
            ;;
        "health")
            health_check
            ;;
        *)
            echo "Usage: $0 {deploy|remove|status|logs|scale|update|backup|health}"
            echo ""
            echo "Commands:"
            echo "  deploy                    Deploy the stack"
            echo "  remove                    Remove the stack"
            echo "  status                    Show stack status"
            echo "  logs <service>            Show logs for a service"
            echo "  scale <service> <count>   Scale a service"
            echo "  update <service> [image]  Update a service"
            echo "  backup                    Backup volumes"
            echo "  health                    Perform health checks"
            echo ""
            echo "Examples:"
            echo "  $0 deploy"
            echo "  $0 logs backend"
            echo "  $0 scale backend 3"
            echo "  $0 update frontend bimapangestu/volt-frontend:v1.2.0"
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"