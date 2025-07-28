#!/bin/bash

# Volt - Docker Swarm Setup Script
# This script initializes Docker Swarm and sets up the environment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Check if Docker is installed and running
check_docker() {
    if ! command -v docker &> /dev/null; then
        error "Docker is not installed. Please install Docker first."
    fi
    
    if ! docker info &> /dev/null; then
        error "Docker is not running. Please start Docker daemon."
    fi
    
    log "Docker is installed and running"
}

# Initialize Docker Swarm
init_swarm() {
    local advertise_addr="$1"
    
    if docker info --format '{{.Swarm.LocalNodeState}}' | grep -q active; then
        log "Docker Swarm is already initialized"
        return
    fi
    
    log "Initializing Docker Swarm..."
    
    if [[ -n "$advertise_addr" ]]; then
        docker swarm init --advertise-addr "$advertise_addr"
    else
        docker swarm init
    fi
    
    log "Docker Swarm initialized successfully"
}

# Create overlay networks
create_networks() {
    log "Creating overlay networks..."
    
    # Check if network already exists
    if docker network ls --format '{{.Name}}' | grep -q '^volt_network$'; then
        log "Network 'volt_network' already exists"
    else
        docker network create \
            --driver overlay \
            --attachable \
            volt_network
        log "Network 'volt_network' created"
    fi
}

# Create secrets
create_secrets() {
    log "Creating Docker secrets..."
    
    # JWT Secret
    if docker secret ls --format '{{.Name}}' | grep -q '^jwt_secret$'; then
        log "Secret 'jwt_secret' already exists"
    else
        # Generate a random JWT secret
        openssl rand -base64 32 | docker secret create jwt_secret -
        log "Secret 'jwt_secret' created"
    fi
    
    # MongoDB root password
    if docker secret ls --format '{{.Name}}' | grep -q '^mongo_root_password$'; then
        log "Secret 'mongo_root_password' already exists"
    else
        openssl rand -base64 32 | docker secret create mongo_root_password -
        log "Secret 'mongo_root_password' created"
    fi
    
    # Redis password
    if docker secret ls --format '{{.Name}}' | grep -q '^redis_password$'; then
        log "Secret 'redis_password' already exists"
    else
        openssl rand -base64 32 | docker secret create redis_password -
        log "Secret 'redis_password' created"
    fi
}

# Setup environment file
setup_env() {
    local env_file=".env.prod"
    
    if [[ -f "$env_file" ]]; then
        log "Environment file '$env_file' already exists"
        return
    fi
    
    log "Creating environment file '$env_file'..."
    
    # Copy from example
    if [[ -f ".env.prod.example" ]]; then
        cp ".env.prod.example" "$env_file"
        log "Environment file created from example"
        warn "Please edit '$env_file' with your production values"
    else
        error "Example environment file '.env.prod.example' not found"
    fi
}

# Pull required images
pull_images() {
    log "Pulling required Docker images..."
    
    local images=(
        "mongo:8.0"
        "redis:alpine"
        "nginx:alpine"
        "bimapangestu/volt-backend:latest"
        "bimapangestu/volt-frontend:latest"
    )
    
    for image in "${images[@]}"; do
        info "Pulling $image..."
        docker pull "$image" || warn "Failed to pull $image"
    done
    
    log "Image pulling completed"
}

# Show swarm information
show_swarm_info() {
    log "Docker Swarm Information:"
    echo ""
    
    info "Swarm Status:"
    docker info --format 'Swarm: {{.Swarm.LocalNodeState}}'
    
    echo ""
    info "Node Information:"
    docker node ls
    
    echo ""
    info "Available Networks:"
    docker network ls --filter driver=overlay
    
    echo ""
    info "Available Secrets:"
    docker secret ls
    
    echo ""
    info "Join Token (Worker):"
    docker swarm join-token worker | grep "docker swarm join"
    
    echo ""
    info "Join Token (Manager):"
    docker swarm join-token manager | grep "docker swarm join"
}

# Add worker node
add_worker() {
    log "Worker Join Token:"
    docker swarm join-token worker
}

# Add manager node
add_manager() {
    log "Manager Join Token:"
    docker swarm join-token manager
}

# Remove node
remove_node() {
    local node_id="$1"
    
    if [[ -z "$node_id" ]]; then
        error "Usage: $0 remove-node <node_id>"
    fi
    
    log "Removing node: $node_id"
    docker node rm "$node_id"
}

# Promote node to manager
promote_node() {
    local node_id="$1"
    
    if [[ -z "$node_id" ]]; then
        error "Usage: $0 promote <node_id>"
    fi
    
    log "Promoting node to manager: $node_id"
    docker node promote "$node_id"
}

# Demote node to worker
demote_node() {
    local node_id="$1"
    
    if [[ -z "$node_id" ]]; then
        error "Usage: $0 demote <node_id>"
    fi
    
    log "Demoting node to worker: $node_id"
    docker node demote "$node_id"
}

# Label node
label_node() {
    local node_id="$1"
    local label="$2"
    
    if [[ -z "$node_id" || -z "$label" ]]; then
        error "Usage: $0 label <node_id> <label_key>=<label_value>"
    fi
    
    log "Adding label to node: $node_id"
    docker node update --label-add "$label" "$node_id"
}

# Complete setup
setup_all() {
    local advertise_addr="$1"
    
    check_docker
    init_swarm "$advertise_addr"
    create_networks
    create_secrets
    setup_env
    pull_images
    show_swarm_info
    
    echo ""
    log "Docker Swarm setup completed successfully!"
    echo ""
    warn "Next steps:"
    echo "1. Edit .env.prod with your production configuration"
    echo "2. Run './scripts/deploy.sh deploy' to deploy the application"
    echo "3. Configure your load balancer/reverse proxy to point to this server"
}

# Main script
main() {
    local command="$1"
    
    case "$command" in
        "init")
            check_docker
            init_swarm "$2"
            ;;
        "setup")
            setup_all "$2"
            ;;
        "networks")
            create_networks
            ;;
        "secrets")
            create_secrets
            ;;
        "env")
            setup_env
            ;;
        "pull")
            pull_images
            ;;
        "info")
            show_swarm_info
            ;;
        "add-worker")
            add_worker
            ;;
        "add-manager")
            add_manager
            ;;
        "remove-node")
            remove_node "$2"
            ;;
        "promote")
            promote_node "$2"
            ;;
        "demote")
            demote_node "$2"
            ;;
        "label")
            label_node "$2" "$3"
            ;;
        *)
            echo "Usage: $0 {init|setup|networks|secrets|env|pull|info|add-worker|add-manager|remove-node|promote|demote|label}"
            echo ""
            echo "Commands:"
            echo "  init [advertise-addr]     Initialize Docker Swarm"
            echo "  setup [advertise-addr]    Complete setup (recommended)"
            echo "  networks                  Create overlay networks"
            echo "  secrets                   Create Docker secrets"
            echo "  env                       Setup environment file"
            echo "  pull                      Pull required images"
            echo "  info                      Show swarm information"
            echo "  add-worker                Show worker join token"
            echo "  add-manager               Show manager join token"
            echo "  remove-node <id>          Remove a node"
            echo "  promote <id>              Promote node to manager"
            echo "  demote <id>               Demote node to worker"
            echo "  label <id> <key>=<value>  Add label to node"
            echo ""
            echo "Examples:"
            echo "  $0 setup 192.168.1.100"
            echo "  $0 add-worker"
            echo "  $0 label node-id env=production"
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"