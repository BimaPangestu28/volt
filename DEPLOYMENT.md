# Volt Deployment Guide

Simple deployment menggunakan GitHub Actions + SSH + Docker Swarm.

## 🚀 Setup Server

### 1. Persiapan Server
```bash
# Install Docker & Docker Compose
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Initialize Docker Swarm
sudo docker swarm init
```

### 2. Setup Directory Structure
```bash
# Create project directory
sudo mkdir -p /opt/volt
cd /opt/volt

# Copy configuration files
sudo wget https://raw.githubusercontent.com/yourusername/volt/main/docker-compose.prod.yml
sudo wget https://raw.githubusercontent.com/yourusername/volt/main/docker-compose.staging.yml
sudo wget https://raw.githubusercontent.com/yourusername/volt/main/.env.production.example -O .env.production
sudo wget https://raw.githubusercontent.com/yourusername/volt/main/.env.staging.example -O .env.staging

# Edit environment variables
sudo nano .env.production
sudo nano .env.staging
```

### 3. Setup Environment Variables
Edit `/opt/volt/.env.production`:
```bash
# Database
MONGO_USERNAME=voltuser
MONGO_PASSWORD=your-secure-mongo-password-here
REDIS_PASSWORD=your-secure-redis-password-here
JWT_SECRET=your-super-secure-jwt-secret-key-minimum-64-chars

# Update dengan domain Anda
CORS_ORIGIN=https://volt.yourdomain.com
WS_ORIGIN=wss://api.yourdomain.com
```

## 🔑 Setup GitHub Secrets

Di GitHub repository, tambahkan secrets berikut:

### Docker Hub
- `DOCKER_USERNAME`: Username Docker Hub
- `DOCKER_PASSWORD`: Password Docker Hub

### Production Server
- `PROD_HOST`: IP address server production
- `PROD_USER`: Username SSH (biasanya `root` atau `ubuntu`)
- `PROD_PASSWORD`: Password SSH user
- `PROD_PORT`: SSH port (optional, default 22)

### Staging Server (optional)
- `STAGING_HOST`: IP address server staging
- `STAGING_USER`: Username SSH
- `STAGING_PASSWORD`: Password SSH user
- `STAGING_PORT`: SSH port (optional, default 22)

## 📦 Deployment Process

### Automatic Deployment
- **Push ke `main`** → Deploy ke production
- **Push ke `develop`** → Deploy ke staging
- **Pull Request** → Run tests only

### Manual Deployment
```bash
# Production
cd /opt/volt
source .env.production
docker stack deploy -c docker-compose.prod.yml volt-prod

# Staging
cd /opt/volt
source .env.staging
docker stack deploy -c docker-compose.staging.yml volt-staging
```

## 🔧 Management Commands

### Check Status
```bash
# Check services
docker service ls

# Check specific service logs
docker service logs volt-prod_backend
docker service logs volt-prod_frontend

# Check container status
docker ps
```

### Scale Services
```bash
# Scale backend to 3 replicas
docker service scale volt-prod_backend=3

# Scale frontend to 2 replicas
docker service scale volt-prod_frontend=2
```

### Update Services
```bash
# Update backend image
docker service update --image bimapangestu/volt-backend:latest volt-prod_backend

# Update frontend image
docker service update --image bimapangestu/volt-frontend:latest volt-prod_frontend
```

### Rollback
```bash
# Rollback backend to previous version
docker service rollback volt-prod_backend

# Rollback frontend to previous version
docker service rollback volt-prod_frontend
```

### Cleanup
```bash
# Remove unused images
docker image prune -f

# Remove unused volumes
docker volume prune -f

# Remove unused networks
docker network prune -f
```

## 🌐 Nginx Configuration (Optional)

Jika menggunakan external load balancer/reverse proxy, hapus service `nginx` dari docker-compose files.

Contoh Nginx config (`/etc/nginx/sites-available/volt`):
```nginx
server {
    listen 80;
    server_name volt.yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name volt.yourdomain.com;
    
    ssl_certificate /path/to/ssl/cert.pem;
    ssl_certificate_key /path/to/ssl/private.key;
    
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    location /webhook {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 🔍 Monitoring

### Health Checks
```bash
# Check application health
curl http://localhost:8080/health

# Check service status
docker service inspect volt-prod_backend --format='{{.UpdateStatus.State}}'
```

### Logs
```bash
# Follow backend logs
docker service logs -f volt-prod_backend

# Follow frontend logs
docker service logs -f volt-prod_frontend

# Check last 100 lines
docker service logs --tail 100 volt-prod_backend
```

## 🚨 Troubleshooting

### Service Won't Start
```bash
# Check service events
docker service ps volt-prod_backend --no-trunc

# Check detailed error logs
docker service logs volt-prod_backend
```

### Database Connection Issues
```bash
# Check MongoDB service
docker service logs volt-prod_mongo

# Test MongoDB connection
docker exec -it $(docker ps -qf name=volt-prod_mongo) mongo -u voltuser -p
```

### Image Pull Issues
```bash
# Login to Docker Hub manually
docker login

# Pull images manually
docker pull bimapangestu/volt-backend:latest
docker pull bimapangestu/volt-frontend:latest
```

## 📋 Maintenance

### Backup Database
```bash
# Create backup
docker exec $(docker ps -qf name=volt-prod_mongo) mongodump --out /backup --username voltuser --password your-password

# Copy backup to host
docker cp $(docker ps -qf name=volt-prod_mongo):/backup ./backup-$(date +%Y%m%d)
```

### Update Stack
```bash
# Pull latest images
docker pull bimapangestu/volt-backend:latest
docker pull bimapangestu/volt-frontend:latest

# Update stack (zero-downtime)
docker stack deploy -c docker-compose.prod.yml volt-prod
```

---

## 🎯 Quick Commands Cheatsheet

```bash
# Deploy production
docker stack deploy -c docker-compose.prod.yml volt-prod

# Check services
docker service ls

# Follow logs
docker service logs -f volt-prod_backend

# Scale backend
docker service scale volt-prod_backend=3

# Update backend
docker service update --image bimapangestu/volt-backend:latest volt-prod_backend

# Rollback
docker service rollback volt-prod_backend

# Remove stack
docker stack rm volt-prod
```