# Volt - Docker Swarm Deployment Guide

This guide explains how to deploy the Volt API client application using Docker Swarm with CI/CD pipeline via GitHub Actions.

## 🚀 Quick Start

### Prerequisites

- Docker Engine 20.10+
- Docker Compose v2+
- Git
- OpenSSL (for generating secrets)

### 1. Initial Setup

```bash
# Clone the repository
git clone https://github.com/bimapangestu/volt.git
cd volt

# Setup Docker Swarm
./scripts/setup-swarm.sh setup

# Generate secure credentials
./scripts/generate-secrets.sh generate
# Edit .env.prod with your domain URLs

# Deploy the application
./scripts/deploy.sh deploy
```

## 📋 Detailed Deployment Guide

### Step 1: Server Preparation

```bash
# Update system packages
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.21.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Verify installation
docker --version
docker-compose --version
```

### Step 2: Initialize Docker Swarm

```bash
# Option 1: Complete setup (recommended)
./scripts/setup-swarm.sh setup [your-server-ip]

# Option 2: Manual setup
./scripts/setup-swarm.sh init [your-server-ip]
./scripts/setup-swarm.sh networks
./scripts/setup-swarm.sh secrets

# Generate secure passwords (do this AFTER swarm init)
./scripts/generate-secrets.sh all
```

### Step 3: Generate Secure Credentials

Generate cryptographically secure passwords and configuration:

```bash
# Generate secure .env.prod file with strong passwords
./scripts/generate-secrets.sh generate

# Or generate all components (recommended)
./scripts/generate-secrets.sh all
```

This will create `.env.prod` with:
- **MongoDB**: Complex username/password (40+ chars)
- **Redis**: Secure password (35+ chars) 
- **JWT**: Cryptographically strong secret (100+ chars)
- **Proper file permissions** (600)

Edit the generated `.env.prod` to update domain URLs:

```bash
# Update these values for your domain
FRONTEND_URL=https://yourdomain.com
WS_ORIGIN=wss://yourdomain.com
API_URL=https://api.yourdomain.com
WS_URL=wss://api.yourdomain.com
```

### Step 4: Deploy Application

```bash
# Deploy the complete stack
./scripts/deploy.sh deploy

# Check deployment status
./scripts/deploy.sh status

# View service logs
./scripts/deploy.sh logs backend
./scripts/deploy.sh logs frontend
```

## 🔧 GitHub Actions CI/CD

### Setup Docker Hub Integration

1. **Create Docker Hub Repository**
   ```bash
   # Create repositories on Docker Hub:
   # - bimapangestu/volt-backend
   # - bimapangestu/volt-frontend
   ```

2. **Configure GitHub Secrets**
   
   Go to your GitHub repository → Settings → Secrets and variables → Actions
   
   Add the following secrets:
   ```
   DOCKER_USERNAME=your-dockerhub-username
   DOCKER_PASSWORD=your-dockerhub-password
   DOCKER_HOST=tcp://your-server:2376
   DOCKER_CERT_PATH=/path/to/docker/certs
   DOCKER_TLS_VERIFY=1
   ```

3. **GitHub Actions Workflow**
   
   The workflow (`.github/workflows/deploy.yml`) automatically:
   - Builds and tests both backend and frontend
   - Creates Docker images for multiple architectures (amd64, arm64)
   - Pushes images to Docker Hub
   - Performs security scanning with Trivy
   - Deploys to staging (develop branch) and production (main branch)
   - Cleans up old images

### Workflow Triggers

- **Pull Requests**: Build and test only
- **Develop Branch**: Deploy to staging
- **Main Branch**: Deploy to production
- **Tags (v*)**: Deploy versioned release

## 🏗️ Architecture Overview

### Services

1. **Frontend** (2 replicas)
   - SvelteKit application with Nginx
   - Port: 80, 443
   - Health check: `/health`

2. **Backend** (2 replicas)  
   - Go API with Gin framework
   - Load balanced by Nginx
   - Health check: `/health`

3. **Nginx Load Balancer**
   - Reverse proxy for backend
   - Rate limiting and security headers
   - WebSocket support

4. **MongoDB**
   - Persistent data storage
   - Authentication enabled
   - Volume: `mongo_data`

5. **Redis**
   - Session and cache storage
   - Password protected
   - Volume: `redis_data`

### Network Architecture

```
Internet
    ↓
[Load Balancer/CDN]
    ↓
[Docker Swarm Manager]
    ↓
[Frontend Services] → [Nginx LB] → [Backend Services]
                           ↓
                    [MongoDB] [Redis]
```

## 📊 Monitoring and Management

### Service Management

```bash
# Scale services
./scripts/deploy.sh scale backend 4
./scripts/deploy.sh scale frontend 3

# Update specific service
./scripts/deploy.sh update backend bimapangestu/volt-backend:v1.2.0

# View service logs
./scripts/deploy.sh logs backend

# Health checks
./scripts/deploy.sh health
```

### Stack Management

```bash
# View stack status
docker stack ls
docker stack services volt
docker stack ps volt

# Remove stack
./scripts/deploy.sh remove

# Backup data
./scripts/deploy.sh backup
```

### Node Management

```bash
# Add worker node
./scripts/setup-swarm.sh add-worker

# Add manager node  
./scripts/setup-swarm.sh add-manager

# Label nodes for constraints
./scripts/setup-swarm.sh label node-id env=production
./scripts/setup-swarm.sh label node-id type=database
```

## 🔒 Security Considerations

### Docker Images
- Multi-stage builds for minimal attack surface
- Non-root user execution
- Security scanning with Trivy
- Regular base image updates

### Network Security
- Overlay network isolation
- Rate limiting on API endpoints
- Security headers (CSP, HSTS, etc.)
- HTTPS enforcement

### Secrets Management
- Docker secrets for sensitive data
- Environment-specific configurations
- No secrets in images or repository

### Access Control
- Manager vs worker node separation
- Service-specific resource limits
- Health checks and restart policies

## 🚨 Troubleshooting

### Common Issues

1. **Services not starting**
   ```bash
   # Check service logs
   ./scripts/deploy.sh logs <service-name>
   
   # Check service constraints
   docker service inspect volt_<service-name>
   
   # Check node availability
   docker node ls
   ```

2. **Database connection issues**
   ```bash
   # Check MongoDB logs
   ./scripts/deploy.sh logs mongo
   
   # Verify environment variables
   docker service inspect volt_backend --format '{{.Spec.TaskTemplate.ContainerSpec.Env}}'
   ```

3. **Image pull failures**
   ```bash
   # Check Docker Hub connectivity
   docker pull bimapangestu/volt-backend:latest
   
   # Verify registry authentication
   docker login
   ```

### Rolling Back

```bash
# Rollback to previous version
docker service rollback volt_backend
docker service rollback volt_frontend

# Or deploy specific version
VERSION=v1.1.0 ./scripts/deploy.sh deploy
```

### Backup and Recovery

```bash
# Create backup
./scripts/deploy.sh backup

# Restore from backup (manual process)
# 1. Stop services
docker service scale volt_mongo=0

# 2. Restore volume
docker run --rm \
  -v volt_mongo_data:/data \
  -v ./backups/20240127_120000:/backup \
  alpine tar xzf /backup/mongo_data.tar.gz -C /data

# 3. Start services
docker service scale volt_mongo=1
```

## 📚 Additional Resources

- [Docker Swarm Documentation](https://docs.docker.com/engine/swarm/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Nginx Configuration Guide](https://nginx.org/en/docs/)
- [MongoDB Best Practices](https://docs.mongodb.com/manual/administration/production-notes/)

## 🤝 Contributing

1. Fork the repository
2. Create feature branch
3. Make changes and test locally
4. Submit pull request
5. GitHub Actions will automatically build and test

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.