# ZKTeco API - Docker Deployment Guide

Complete guide to deploying ZKTeco API with Face Version 35+ support using Docker and Official ZKTeco SDK.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Download Official ZKTeco SDK](#download-official-zkteco-sdk)
3. [Quick Start](#quick-start)
4. [Detailed Setup](#detailed-setup)
5. [Usage Examples](#usage-examples)
6. [Management Commands](#management-commands)
7. [Troubleshooting](#troubleshooting)
8. [Production Deployment](#production-deployment)

## Prerequisites

### Required Software
- [x] Docker 20.10+ installed
- [x] Docker Compose v2.0+ installed
- [x] Official ZKTeco SDK for Linux (download below)
- [x] Git (to clone the repository)

### Check Docker Installation
```bash
docker --version
docker-compose --version

# Should output versions like:
# Docker version 24.0.7
# Docker Compose version 2.21.0
```

## Download Official ZKTeco SDK

### Step 1: Download SDK

Visit: **https://www.zkteco.com/en/support_download.html**

Search for: **`Linux_AccessControl_RealTime_SDK`**

**Available Downloads:**
- `Linux_AccessControl_RealTime_SDK_v6.60.zip`
- `Linux_AccessControl_RealTime_SDK_v6.80.zip`
- `Linux_AccessControl_RealTime_SDK_Latest.zip`

### Step 2: Verify SDK Contents

Extract SDK to: `scripts/sync_zk/face_sync_cpp/zk_sdk/`

**Expected structure:**
```
scripts/sync_zk/face_sync_cpp/
└── zk_sdk/
    ├── Include/
    │   ├── zk_face_sdk.h       ← Required
    │   ├── zk_device.h         ← Required
    │   ├── zk_commpro.h        ← Required
    │   └── ...
    ├── Lib/
    │   ├── libzk.so            ← Required (Linux)
    │   ├── libzk.so.6          ← Versioned library
    │   └── ...
    ├── Sample Code/
    │   ├── C++/
    │   ├── C/
    │   └── ...
    └── Documentation/
```

### Step 3: Verify Files
```bash
# Check SDK files exist
ls -la scripts/sync_zk/face_sync_cpp/zk_sdk/Include/zk_face_sdk.h
ls -la scripts/sync_zk/face_sync_cpp/zk_sdk/Lib/libzk.so

# Expected output:
# -rw-r--r-- ... zk_face_sdk.h
# -rw-r--r-- ... libzk.so
```

## Quick Start

### One-Command Setup
```bash
# 1. Clone repository (if needed)
git clone <repository-url>
cd zkteco-api-go

# 2. Extract SDK to correct location
unzip Linux_AccessControl_RealTime_SDK.zip \
  -d scripts/sync_zk/face_sync_cpp/zk_sdk/

# 3. Build and start all services
make -f Makefile.docker quickstart
```

**What this does:**
- ✅ Checks for SDK installation
- ✅ Builds Docker image with Face Sync support
- ✅ Starts PostgreSQL and API services
- ✅ Runs database migrations
- ✅ Tests API endpoints
- ✅ Shows status

### Access Services
- **API**: http://localhost:8090
- **Health Check**: http://localhost:8090/health
- **PgAdmin**: http://localhost:8080 (admin@zkteco.local / admin)

## Detailed Setup

### Method 1: Using Makefile (Recommended)

**Build image:**
```bash
make -f Makefile.docker build
```

**Start services:**
```bash
make -f Makefile.docker up
```

**View logs:**
```bash
make -f Makefile.docker logs
```

**Test face sync:**
```bash
make -f Makefile.docker sync-face-api
```

**Stop services:**
```bash
make -f Makefile.docker down
```

### Method 2: Using Docker Compose Directly

**Build image:**
```bash
docker build -f Dockerfile.face-sync -t zkteco-api-go:face-sync .
```

**Start services:**
```bash
docker-compose -f docker-compose.face-sync.yml up -d
```

**View logs:**
```bash
docker-compose -f docker-compose.face-sync.yml logs -f
```

**Stop services:**
```bash
docker-compose -f docker-compose.face-sync.yml down
```

### Method 3: Step-by-Step

**1. Set up environment:**
```bash
# Copy environment template
cp .env.docker .env

# Edit if needed
vim .env
```

**2. Build image:**
```bash
docker build -f Dockerfile.face-sync -t zkteco-api-go:face-sync .
```

**3. Start PostgreSQL:**
```bash
docker-compose -f docker-compose.face-sync.yml up -d postgres
```

**4. Wait for database ready:**
```bash
# Check health
docker-compose -f docker-compose.face-sync.yml ps

# Should show "healthy" for postgres
```

**5. Start API:**
```bash
docker-compose -f docker-compose.face-sync.yml up -d zkteco-api
```

## Usage Examples

### API Testing

**Health check:**
```bash
curl http://localhost:8090/health
```

**Get devices:**
```bash
curl http://localhost:8090/api/adms/devices
```

**Get attendance events:**
```bash
curl "http://localhost:8090/api/attendance/events?startDate=2024-01-01&endDate=2024-12-31"
```

**Sync face templates (via API):**
```bash
curl -X POST http://localhost:8090/api/adms/devices/SPK7253000015/sync-face-templates-v35
```

**Get face templates:**
```bash
curl http://localhost:8090/api/adms/devices/SPK7253000015/face-templates-v35
```

### Database Operations

**Access database shell:**
```bash
make -f Makefile.docker shell-db
```

**Query face templates:**
```bash
make -f Makefile.docker face-templates
```

**Backup database:**
```bash
make -f Makefile.docker db-backup
```

**Restore database:**
```bash
make -f Makefile.docker db-restore FILE=./backups/zkteco_db_20251121_120000.sql
```

**Reset database (WARNING: destroys all data):**
```bash
make -f Makefile.docker db-reset
```

### Face Sync Operations

**Trigger sync manually (from container):**
```bash
make -f Makefile.docker sync-face
```

**Trigger sync via API:**
```bash
make -f Makefile.docker sync-face-api
```

**Check face sync logs:**
```bash
make -f Makefile.docker logs-face
```

### Container Shell Access

**API container:**
```bash
make -f Makefile.docker shell
```

**Database container:**
```bash
make -f Makefile.docker shell-db
```

**Run arbitrary commands:**
```bash
# Inside API container
docker-compose -f docker-compose.face-sync.yml exec zkteco-api ls -la /app

# Inside database container
docker-compose -f docker-compose.face-sync.yml exec postgres psql -U postgres -d zkteco_db -c "SELECT * FROM devices;"
```

## Management Commands

### Service Management
```bash
# Start all services
make -f Makefile.docker up

# Stop all services
make -f Makefile.docker down

# Restart services
make -f Makefile.docker restart

# View service status
make -f Makefile.docker status

# View resource usage
make -f Makefile.docker stats
```

### Logs and Monitoring
```bash
# View all logs
make -f Makefile.docker logs

# View API logs only
make -f Makefile.docker logs-api

# View database logs only
make -f Makefile.docker logs-db

# View error logs only
make -f Makefile.docker logs-errors

# View face sync logs
make -f Makefile.docker logs-face
```

### Build and Update
```bash
# Rebuild image (after code changes)
make -f Makefile.docker build

# Full rebuild (no cache)
docker build -f Dockerfile.face-sync -t zkteco-api-go:face-sync . --no-cache

# Update services (pull latest)
docker-compose -f docker-compose.face-sync.yml pull
```

### Cleanup
```bash
# Remove stopped containers
make -f Makefile.docker clean

# Remove everything (containers, images, volumes)
make -f Makefile.docker clean-all

# Remove unused images
docker image prune -f

# Remove unused volumes
docker volume prune -f
```

### Backup and Restore
```bash
# Backup database
make -f Makefile.docker db-backup
# Creates: ./backups/zkteco_db_YYYYMMDD_HHMMSS.sql

# Restore database
make -f Makefile.docker db-restore FILE=./backups/zkteco_db_20251121_120000.sql

# Export Docker image
make -f Makefile.docker export-image FILE=zkteco.tar.gz

# Import Docker image
make -f Makefile.docker import-image FILE=zkteco.tar.gz
```

## Troubleshooting

### Common Issues

#### 1. SDK Not Found

**Error:**
```
ERROR: ZKTeco SDK not found!
```

**Solution:**
```bash
# Check SDK location
ls -la scripts/sync_zk/face_sync_cpp/zk_sdk/Include/zk_face_sdk.h

# If missing, download and extract SDK
unzip Linux_AccessControl_RealTime_SDK.zip \
  -d scripts/sync_zk/face_sync_cpp/zk_sdk/
```

#### 2. Build Fails

**Error:**
```
ERROR: After scaling service 'zkteco-api' from 1 to 1, this resulted in 2 total
```

**Solution:**
```bash
# Clean and rebuild
make -f Makefile.docker clean
make -f Makefile.docker build
```

#### 3. Database Connection Failed

**Error:**
```
dial error: dial tcp 172.20.0.2:5432: connect: connection refused
```

**Solution:**
```bash
# Check database status
make -f Makefile.docker status

# Wait for database to be healthy
docker-compose -f docker-compose.face-sync.yml up -d postgres
sleep 10
curl http://localhost:8090/health
```

#### 4. Face Templates Not Syncing

**Error:**
```
Face templates: 0 fetched
```

**Debug:**
```bash
# Check device connectivity
docker-compose -f docker-compose.face-sync.yml exec zkteco-api \
  ping -c 3 100.100.2.215

# Check C++ service can connect
docker-compose -f docker-compose.face-sync.yml exec zkteco-api \
  /app/face_sync_service 100.100.2.215 4370 0

# Check logs
make -f Makefile.docker logs-face
```

#### 5. API Returns 500 Error

**Debug:**
```bash
# Check API logs
make -f Makefile.docker logs-api

# Check database connection
docker-compose -f docker-compose.face-sync.yml exec zkteco-api \
  psql "postgresql://postgres:postgres@postgres:5432/zkteco_db" \
  -c "SELECT 1;"
```

#### 6. Permission Denied

**Error:**
```
permission denied while trying to connect to the Docker daemon socket
```

**Solution:**
```bash
# Add user to docker group
sudo usermod -aG docker $USER

# Log out and back in, or:
newgrp docker

# Verify
groups
```

### Debugging Commands

**Check all resources:**
```bash
make -f Makefile.docker info
```

**View container processes:**
```bash
docker ps -a
```

**Inspect container:**
```bash
docker inspect zkteco-api
```

**View container filesystem:**
```bash
docker exec -it zkteco-api ls -la /app
```

**Test network connectivity:**
```bash
# From API container
docker-compose -f docker-compose.face-sync.yml exec zkteco-api \
  nc -zv postgres 5432

# From database container
docker-compose -f docker-compose.face-sync.yml exec postgres \
  nc -zv zkteco-api 8090
```

## Production Deployment

### Environment Configuration

**1. Create production .env:**
```bash
cp .env.docker .env.production
```

**2. Edit production settings:**
```env
ENVIRONMENT=production
DB_PASSWORD=secure_random_password_here
ERPNEXT_API_KEY=your_production_api_key
ERPNEXT_SECRET=your_production_secret
API_SECRET_KEY=generate_random_secret
JWT_SECRET=generate_random_jwt_secret
```

**3. Use production compose file:**
```bash
# Create production override
cp docker-compose.face-sync.yml docker-compose.prod.yml

# Edit to disable PgAdmin, change ports, etc.
vim docker-compose.prod.yml

# Start with production config
docker-compose -f docker-compose.prod.yml -f docker-compose.face-sync.yml up -d
```

### Security Hardening

**1. Change default passwords:**
```env
POSTGRES_PASSWORD=<generate_secure_password>
PGADMIN_DEFAULT_PASSWORD=<generate_secure_password>
```

**2. Use secrets management:**
```yaml
# In docker-compose.prod.yml
services:
  postgres:
    environment:
      POSTGRES_PASSWORD_FILE: /run/secrets/db_password
    secrets:
      - db_password
```

**3. Restrict network access:**
```yaml
networks:
  zkteco-network:
    driver: bridge
    internal: true  # No external access
```

**4. Run as non-root:**
- Already configured in Dockerfile ✅
- All services run as non-root users

### Scaling

**Scale API instances:**
```bash
make -f Makefile.docker scale-api N=3
```

**Load balancer setup:**
```yaml
# docker-compose.prod.yml
nginx:
  image: nginx:alpine
  ports:
    - "80:80"
  volumes:
    - ./nginx.conf:/etc/nginx/nginx.conf
  depends_on:
    - zkteco-api
```

### Monitoring

**Health checks:**
```bash
# Monitor service health
watch -n 5 'curl -s http://localhost:8090/health | jq .'
```

**Logging:**
```bash
# Centralized logging with ELK stack
# Or use Docker logging drivers
docker-compose -f docker-compose.face-sync.yml up -d \
  --log-driver=json-file \
  --log-opt max-size=10m \
  --log-opt max-file=3
```

### Backup Strategy

**Automated daily backups:**
```bash
# Add to crontab
0 2 * * * cd /path/to/zkteco-api-go && make -f Makefile.docker db-backup
```

**Retain backups:**
```bash
# Keep only last 30 days
find ./backups -name "*.sql" -mtime +30 -delete
```

### Update Process

**1. Pull latest changes:**
```bash
git pull origin main
```

**2. Rebuild image:**
```bash
make -f Makefile.docker build
```

**3. Rolling update:**
```bash
docker-compose -f docker-compose.face-sync.yml up -d --no-deps zkteco-api
```

**4. Verify:**
```bash
make -f Makefile.docker test
```

## File Reference

### Docker Files
- `Dockerfile.face-sync` - Multi-stage build with Face Sync support
- `docker-compose.face-sync.yml` - Complete service stack
- `Makefile.docker` - Management commands
- `.env.docker` - Environment template

### C++ Service Files
- `scripts/sync_zk/face_sync_cpp/Dockerfile.cpp` - Standalone C++ service
- `scripts/sync_zk/face_sync_cpp/face_sync_service.cpp` - C++ source
- `scripts/sync_zk/face_sync_cpp/Makefile` - C++ build system

### Documentation
- `DOCKER.md` - This file (Docker deployment guide)
- `scripts/sync_zk/face_sync_cpp/README.md` - Face Sync documentation
- `scripts/sync_zk/face_sync_cpp/QUICKSTART.md` - Quick start guide

## Quick Reference

### Common Commands
```bash
# Quick start
make -f Makefile.docker quickstart

# Build
make -f Makefile.docker build

# Start/Stop
make -f Makefile.docker up
make -f Makefile.docker down

# Logs
make -f Makefile.docker logs
make -f Makefile.docker logs-api
make -f Makefile.docker logs-db

# Face sync
make -f Makefile.docker sync-face-api
make -f Makefile.docker face-templates

# Database
make -f Makefile.docker shell-db
make -f Makefile.docker db-backup
make -f Makefile.docker db-reset

# Cleanup
make -f Makefile.docker clean
make -f Makefile.docker clean-all

# Info
make -f Makefile.docker status
make -f Makefile.docker info
```

### URLs
- **API**: http://localhost:8090
- **Health**: http://localhost:8090/health
- **PgAdmin**: http://localhost:8080
- **API Docs**: http://localhost:8090/docs (if configured)

## Support

- **Full Documentation**: `scripts/sync_zk/face_sync_cpp/README.md`
- **Quick Start**: `scripts/sync_zk/face_sync_cpp/QUICKSTART.md`
- **ZKTeco SDK**: https://www.zkteco.com/en/support_download.html
- **Docker Docs**: https://docs.docker.com/

---

**🎉 You now have a complete Docker-based ZKTeco API with Face Version 35+ support!**
