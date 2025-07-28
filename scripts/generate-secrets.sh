#!/bin/bash

# Volt - Secure Password Generator
# This script generates cryptographically secure passwords for production deployment

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

# Generate a complex password with special characters
generate_password() {
    local length="$1"
    local prefix="$2"
    
    if [[ -z "$length" ]]; then
        length=32
    fi
    
    # Generate random password with uppercase, lowercase, numbers, and special chars
    local password=$(openssl rand -base64 48 | tr -d "=+/" | cut -c1-${length})
    
    # Add some special characters and numbers for complexity
    local special_chars="!@#$%^&*()_+-=[]{}|;:,.<>?"
    local random_special=$(echo $special_chars | fold -w1 | shuf | head -n3 | tr -d '\n')
    local random_numbers=$(shuf -i 1000-9999 -n 1)
    
    if [[ -n "$prefix" ]]; then
        echo "${prefix}_${password}${random_special}${random_numbers}"
    else
        echo "${password}${random_special}${random_numbers}"
    fi
}

# Generate username with random suffix
generate_username() {
    local base="$1"
    local random_suffix=$(openssl rand -hex 4)
    echo "${base}_${random_suffix}"
}

# Generate JWT secret (longer and more complex)
generate_jwt_secret() {
    local base64_part=$(openssl rand -base64 64 | tr -d "=+/\n")
    local hex_part=$(openssl rand -hex 16)
    echo "VoltJWT2024_${base64_part}_${hex_part}"
}

# Create .env.prod with generated secrets
create_env_file() {
    local env_file=".env.prod"
    
    log "Generating secure credentials..."
    
    # Generate credentials
    local mongo_username=$(generate_username "voltadmin")
    local mongo_password=$(generate_password 40 "VoltMongo2024")
    local redis_password=$(generate_password 35 "VoltRedis2024")
    local jwt_secret=$(generate_jwt_secret)
    
    # Create environment file
    cat > "$env_file" << EOF
# Production Environment Variables
# Generated on $(date)
# IMPORTANT: Keep these credentials secure and never commit to version control!

# Docker Registry
DOCKER_REGISTRY=bimapangestu
VERSION=latest

# Database Configuration
MONGO_ROOT_USERNAME=${mongo_username}
MONGO_ROOT_PASSWORD=${mongo_password}
MONGO_DATABASE=volt

# Redis Configuration
REDIS_PASSWORD=${redis_password}

# Backend Configuration
JWT_SECRET=${jwt_secret}
FRONTEND_URL=https://yourdomain.com
WS_ORIGIN=wss://yourdomain.com

# Frontend Configuration
API_URL=https://api.yourdomain.com
WS_URL=wss://api.yourdomain.com

# Optional: Monitoring
ENABLE_MONITORING=true
EOF
    
    # Set restrictive permissions
    chmod 600 "$env_file"
    
    log "Secure environment file created: $env_file"
    warn "IMPORTANT: Update the domain URLs before deployment!"
    
    echo ""
    info "Generated credentials:"
    echo "MongoDB Username: ${mongo_username}"
    echo "MongoDB Password: ${mongo_password:0:20}... (${#mongo_password} chars)"
    echo "Redis Password: ${redis_password:0:20}... (${#redis_password} chars)"
    echo "JWT Secret: ${jwt_secret:0:30}... (${#jwt_secret} chars)"
    echo ""
    warn "Store these credentials in a secure password manager!"
}

# Update docker-compose.prod.yml with placeholders
update_compose_file() {
    local compose_file="docker-compose.prod.yml"
    
    if [[ ! -f "$compose_file" ]]; then
        error "Compose file not found: $compose_file"
    fi
    
    log "Updating compose file to use environment variables..."
    
    # Backup original file
    cp "$compose_file" "${compose_file}.backup"
    
    # Update MongoDB defaults to use env vars
    sed -i.tmp 's/MONGO_ROOT_USERNAME:-[^}]*/MONGO_ROOT_USERNAME:-voltadmin_changeme/g' "$compose_file"
    sed -i.tmp 's/MONGO_ROOT_PASSWORD:-[^}]*/MONGO_ROOT_PASSWORD:-changeme_mongo_password/g' "$compose_file"
    sed -i.tmp 's/REDIS_PASSWORD:-[^}]*/REDIS_PASSWORD:-changeme_redis_password/g' "$compose_file"
    
    # Remove temp file
    rm -f "${compose_file}.tmp"
    
    log "Compose file updated with environment variable placeholders"
}

# Generate Docker secrets
generate_docker_secrets() {
    log "Generating Docker secrets..."
    
    # Check if running in swarm mode
    if ! docker info --format '{{.Swarm.LocalNodeState}}' | grep -q active; then
        warn "Docker Swarm not initialized. Secrets will be created when swarm is initialized."
        return
    fi
    
    # Generate and create secrets
    local mongo_password=$(generate_password 40 "VoltMongo2024")
    local redis_password=$(generate_password 35 "VoltRedis2024") 
    local jwt_secret=$(generate_jwt_secret)
    
    # Create secrets if they don't exist
    if ! docker secret ls --format '{{.Name}}' | grep -q '^volt_mongo_password$'; then
        echo "$mongo_password" | docker secret create volt_mongo_password -
        log "Created Docker secret: volt_mongo_password"
    fi
    
    if ! docker secret ls --format '{{.Name}}' | grep -q '^volt_redis_password$'; then
        echo "$redis_password" | docker secret create volt_redis_password -
        log "Created Docker secret: volt_redis_password"
    fi
    
    if ! docker secret ls --format '{{.Name}}' | grep -q '^volt_jwt_secret$'; then
        echo "$jwt_secret" | docker secret create volt_jwt_secret -
        log "Created Docker secret: volt_jwt_secret"
    fi
    
    info "Docker secrets created successfully"
}

# Show security recommendations
show_security_tips() {
    echo ""
    log "Security Recommendations:"
    echo ""
    echo "1. 🔐 Change all default passwords immediately"
    echo "2. 🚫 Never commit .env.prod to version control"
    echo "3. 🔄 Rotate passwords regularly (every 90 days)"
    echo "4. 📝 Store credentials in a secure password manager"
    echo "5. 🌐 Use strong TLS/SSL certificates in production"
    echo "6. 🛡️ Enable MongoDB authentication and SSL"
    echo "7. 🔍 Monitor access logs regularly"
    echo "8. 🚨 Set up alerting for failed authentication attempts"
    echo "9. 📊 Use Docker secrets in production"
    echo "10. 🔒 Restrict network access with firewalls"
    echo ""
    warn "Add .env.prod to .gitignore if not already present!"
}

# Main script
main() {
    local command="$1"
    
    case "$command" in
        "generate"|"gen")
            create_env_file
            show_security_tips
            ;;
        "compose")
            update_compose_file
            ;;
        "secrets")
            generate_docker_secrets
            ;;
        "all")
            create_env_file
            update_compose_file
            generate_docker_secrets
            show_security_tips
            ;;
        "password")
            local length="${2:-32}"
            local prefix="$3"
            generate_password "$length" "$prefix"
            ;;
        "jwt")
            generate_jwt_secret
            ;;
        *)
            echo "Usage: $0 {generate|compose|secrets|all|password|jwt}"
            echo ""
            echo "Commands:"
            echo "  generate              Generate .env.prod with secure credentials"
            echo "  compose               Update compose file with env var placeholders"
            echo "  secrets               Create Docker secrets (requires swarm mode)"
            echo "  all                   Run all above commands"
            echo "  password [length] [prefix]  Generate a single password"
            echo "  jwt                   Generate JWT secret only"
            echo ""
            echo "Examples:"
            echo "  $0 generate"
            echo "  $0 password 50 MyApp"
            echo "  $0 all"
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"