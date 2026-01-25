# Docker Compose Deployment Guide

## Production-Ready Features

### 🔧 **Infrastructure Setup**
- **Network Isolation**: Custom bridge network with static IPs for reliable service discovery
- **Persistent Volumes**: Separate volumes for each database ensuring data durability
- **Health Checks**: All services include health checks with configurable intervals
- **Resource Management**: CPU and memory limits/reservations for stability

### 🚀 **Service Dependencies**
- Proper startup order with `depends_on` conditions
- Services wait for database health before starting
- Microservices wait for each other before initialization

### 2. **Build and Run**
```bash
# Build all services
docker-compose build

# Start all services
docker-compose up -d

# Check status
docker-compose ps
docker-compose logs -f
```

### 3. **Verify Services**
```bash
# Check GraphQL gateway
curl http://localhost:8080/health

# Test GraphQL endpoint
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"{ accounts { id name } }"}'
```

## Service URLs (Internal)

| Service | URL | Port |
|---------|-----|------|
| Account | `grpc://account:50051` | 50051 |
| Catalog | `grpc://catalog:50052` | 50052 |
| Order | `grpc://order:50053` | 50053 |
| GraphQL | `http://graphql:8080` | 8080 |

## Database Access (Local Development)

| Database | Host | Port | User |
|----------|------|------|------|
| Account | localhost | 5432 | postgres |
| Catalog | localhost | 5433 | postgres |
| Order | localhost | 5434 | postgres |

```bash
psql -h localhost -p 5432 -U postgres -d account_db
psql -h localhost -p 5433 -U postgres -d catalog_db
psql -h localhost -p 5434 -U postgres -d order_db
```

## Common Commands

```bash
# View logs
docker-compose logs -f [service_name]

# Restart a service
docker-compose restart account

# Scale a service
docker-compose up -d --scale catalog=3

# Stop all services
docker-compose down

# Remove volumes (WARNING: deletes data)
docker-compose down -v

# Run a one-off command
docker-compose exec account sh

# Check service health
docker-compose ps
```

## Environment Variables

Copy `.env.example` to `.env` and customize:

```bash
cp .env.example .env
nano .env
```

| Variable | Default | Purpose |
|----------|---------|---------|
| `ENVIRONMENT` | production | Environment mode |
| `LOG_LEVEL` | info | Logging level |
| `GRACEFUL_SHUTDOWN_TIMEOUT` | 30 | Shutdown grace period (seconds) |

### Out of memory
Increase resource limits in `docker-compose.yaml`:
```yaml
deploy:
  resources:
    limits:
      memory: 1024M
```

## Backup & Recovery

```bash
# Backup all database volumes
docker-compose exec account_db pg_dump -U postgres account_db > account_backup.sql

# Restore
docker-compose exec -T account_db psql -U postgres account_db < account_backup.sql
```
