# Deployment Guide

This guide covers different deployment options for the OpenRTB Insights application.

## Quick Start (Recommended)

Run the automated setup script:

```bash
chmod +x scripts/quick-start.sh
./scripts/quick-start.sh
```

This will detect your environment and choose the best deployment method automatically.

## Docker Deployment (Production)

### Prerequisites
- Docker 20.10+
- Docker Compose 2.0+

### Step 1: Clone and Configure

```bash
git clone <repository-url>
cd openrtb-insights

# Copy environment files
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
```

### Step 2: Configure Environment Variables

Edit `backend/.env`:
```bash
# Generate secure JWT secret
JWT_SECRET=$(openssl rand -hex 32)

# Set production database path
DB_PATH=/app/data/analytics.db

# Configure CORS for your domain
CORS_ORIGINS=https://your-domain.com,https://api.your-domain.com
```

Edit `frontend/.env`:
```bash
VITE_API_BASE_URL=https://api.your-domain.com/api
VITE_APP_NAME=OpenRTB Insights
```

### Step 3: Deploy

```bash
# Production deployment
docker-compose up -d

# View logs
docker-compose logs -f

# Check health
curl http://localhost/health
curl http://localhost:8080/health
```

## Local Development

### Prerequisites
- Go 1.21+
- Node.js 18+
- npm 9+

### Backend Setup

```bash
cd backend
./scripts/setup.sh

# Or manual setup:
go mod download
go run scripts/seed-data.go
go run cmd/server/main.go
```

### Frontend Setup

```bash
cd frontend
./scripts/setup.sh

# Or manual setup:
npm ci
npm run dev
```

## Kubernetes Deployment

### Prerequisites
- Kubernetes 1.20+
- kubectl configured
- Persistent storage class

### Step 1: Create Namespace

```yaml
# namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: openrtb-insights
```

```bash
kubectl apply -f namespace.yaml
```

### Step 2: Create ConfigMaps and Secrets

```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: openrtb-config
  namespace: openrtb-insights
data:
  JWT_EXPIRY: "15m"
  REFRESH_TOKEN_EXPIRY: "168h"
  PORT: "8080"
  LOG_LEVEL: "info"
  RATE_LIMIT: "100"
---
apiVersion: v1
kind: Secret
metadata:
  name: openrtb-secrets
  namespace: openrtb-insights
type: Opaque
data:
  JWT_SECRET: <base64-encoded-secret>
  DB_PATH: L2FwcC9kYXRhL2FuYWx5dGljcy5kYg==  # /app/data/analytics.db
```

### Step 3: Create Persistent Volume

```yaml
# pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: openrtb-data
  namespace: openrtb-insights
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: fast-ssd  # Adjust for your cluster
```

### Step 4: Deploy Backend

```yaml
# backend-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: openrtb-backend
  namespace: openrtb-insights
spec:
  replicas: 2
  selector:
    matchLabels:
      app: openrtb-backend
  template:
    metadata:
      labels:
        app: openrtb-backend
    spec:
      containers:
      - name: backend
        image: openrtb-insights-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_PATH
          value: "/app/data/analytics.db"
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: openrtb-secrets
              key: JWT_SECRET
        envFrom:
        - configMapRef:
            name: openrtb-config
        volumeMounts:
        - name: data-volume
          mountPath: /app/data
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
      volumes:
      - name: data-volume
        persistentVolumeClaim:
          claimName: openrtb-data
---
apiVersion: v1
kind: Service
metadata:
  name: openrtb-backend-service
  namespace: openrtb-insights
spec:
  selector:
    app: openrtb-backend
  ports:
  - port: 8080
    targetPort: 8080
  type: ClusterIP
```

### Step 5: Deploy Frontend

```yaml
# frontend-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: openrtb-frontend
  namespace: openrtb-insights
spec:
  replicas: 2
  selector:
    matchLabels:
      app: openrtb-frontend
  template:
    metadata:
      labels:
        app: openrtb-frontend
    spec:
      containers:
      - name: frontend
        image: openrtb-insights-frontend:latest
        ports:
        - containerPort: 80
        livenessProbe:
          httpGet:
            path: /health
            port: 80
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 80
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "64Mi"
            cpu: "50m"
          limits:
            memory: "128Mi"
            cpu: "100m"
---
apiVersion: v1
kind: Service
metadata:
  name: openrtb-frontend-service
  namespace: openrtb-insights
spec:
  selector:
    app: openrtb-frontend
  ports:
  - port: 80
    targetPort: 80
  type: ClusterIP
```

### Step 6: Create Ingress

```yaml
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: openrtb-ingress
  namespace: openrtb-insights
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  tls:
  - hosts:
    - your-domain.com
    - api.your-domain.com
    secretName: openrtb-tls
  rules:
  - host: your-domain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: openrtb-frontend-service
            port:
              number: 80
  - host: api.your-domain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: openrtb-backend-service
            port:
              number: 8080
```

### Deploy to Kubernetes

```bash
kubectl apply -f configmap.yaml
kubectl apply -f pvc.yaml
kubectl apply -f backend-deployment.yaml
kubectl apply -f frontend-deployment.yaml
kubectl apply -f ingress.yaml

# Check status
kubectl get pods -n openrtb-insights
kubectl get services -n openrtb-insights
kubectl get ingress -n openrtb-insights
```

## Environment-Specific Configurations

### Production Environment

```bash
# backend/.env
DB_PATH=/app/data/analytics.db
JWT_SECRET=<generated-with-openssl-rand-hex-32>
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h
PORT=8080
CORS_ORIGINS=https://your-domain.com
LOG_LEVEL=info
RATE_LIMIT=100

# frontend/.env
VITE_API_BASE_URL=https://api.your-domain.com/api
VITE_APP_NAME=OpenRTB Insights
```

### Staging Environment

```bash
# backend/.env
DB_PATH=./data/analytics-staging.db
JWT_SECRET=staging-jwt-secret
PORT=8080
CORS_ORIGINS=https://staging.your-domain.com
LOG_LEVEL=debug
RATE_LIMIT=200

# frontend/.env
VITE_API_BASE_URL=https://api-staging.your-domain.com/api
VITE_APP_NAME=OpenRTB Insights (Staging)
```

## Load Balancing and High Availability

### Nginx Configuration

```nginx
upstream backend {
    server backend1:8080;
    server backend2:8080;
    server backend3:8080;
}

upstream frontend {
    server frontend1:80;
    server frontend2:80;
}

server {
    listen 80;
    server_name api.your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.your-domain.com;

    ssl_certificate /etc/ssl/certs/your-domain.crt;
    ssl_certificate_key /etc/ssl/private/your-domain.key;

    location / {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/ssl/certs/your-domain.crt;
    ssl_certificate_key /etc/ssl/private/your-domain.key;

    location / {
        proxy_pass http://frontend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Database Scaling

For high-traffic deployments, consider:

1. **Read Replicas**: Set up DuckDB read replicas for analytics queries
2. **Connection Pooling**: Implement connection pooling for database access
3. **Caching**: Add Redis for caching frequent queries
4. **Partitioning**: Partition large tables by date ranges

## Monitoring and Observability

### Health Checks

Both services provide health endpoints:
- Backend: `GET /health`
- Frontend: `GET /health`

### Metrics Collection

Add Prometheus metrics collection:

```yaml
# monitoring/prometheus-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-config
data:
  prometheus.yml: |
    global:
      scrape_interval: 15s
    scrape_configs:
    - job_name: 'openrtb-backend'
      static_configs:
      - targets: ['openrtb-backend-service:8080']
```

### Log Aggregation

Use Fluentd or similar for log collection:

```yaml
# logging/fluentd-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: fluentd-config
data:
  fluent.conf: |
    <source>
      @type tail
      path /var/log/containers/*openrtb*.log
      pos_file /var/log/fluentd-containers.log.pos
      tag kubernetes.*
      format json
    </source>
```

## Backup and Recovery

### Database Backup

```bash
# Create backup script
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups"
DB_PATH="/app/data/analytics.db"

# Create backup
cp "$DB_PATH" "$BACKUP_DIR/analytics_backup_$DATE.db"

# Compress backup
gzip "$BACKUP_DIR/analytics_backup_$DATE.db"

# Clean old backups (keep last 30 days)
find $BACKUP_DIR -name "analytics_backup_*.db.gz" -mtime +30 -delete
```

### Automated Backups

```yaml
# backup-cronjob.yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: database-backup
  namespace: openrtb-insights
spec:
  schedule: "0 2 * * *"  # Daily at 2 AM
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: backup
            image: alpine:latest
            command: ["/bin/sh", "-c"]
            args:
            - |
              apk add --no-cache gzip
              DATE=$(date +%Y%m%d_%H%M%S)
              cp /app/data/analytics.db /backups/analytics_backup_$DATE.db
              gzip /backups/analytics_backup_$DATE.db
              find /backups -name "analytics_backup_*.db.gz" -mtime +30 -delete
            volumeMounts:
            - name: data-volume
              mountPath: /app/data
            - name: backup-volume
              mountPath: /backups
          volumes:
          - name: data-volume
            persistentVolumeClaim:
              claimName: openrtb-data
          - name: backup-volume
            persistentVolumeClaim:
              claimName: backup-storage
          restartPolicy: OnFailure
```

## Troubleshooting

### Common Issues

1. **Database Permission Errors**
   ```bash
   # Fix permissions
   docker-compose exec backend chown -R appuser:appgroup /app/data
   ```

2. **CORS Issues**
   ```bash
   # Update backend environment
   CORS_ORIGINS=https://your-actual-domain.com
   ```

3. **Memory Issues**
   ```bash
   # Increase container memory limits
   docker-compose up -d --scale backend=2
   ```

4. **Port Conflicts**
   ```bash
   # Change ports in docker-compose.yml
   ports:
     - "8081:8080"  # Backend
     - "3001:80"    # Frontend
   ```

### Debugging Commands

```bash
# View logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Check container status
docker-compose ps

# Access container shell
docker-compose exec backend sh
docker-compose exec frontend sh

# Check database
docker-compose exec backend sqlite3 /app/data/analytics.db ".tables"

# Test API endpoints
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

This deployment guide should cover most production scenarios. Adjust the configurations based on your specific infrastructure requirements.