# OpenRTB Insights Reporting App

A comprehensive OpenRTB 2.6 bid request/response analytics platform that provides operational visibility through real-time dashboards. Built for DevOps and ad-ops teams handling up to 10,000 requests per second.

## Features

- **Real-time Analytics Dashboard** - Overview of key metrics and system health
- **Platform Statistics** - Detailed bid request analytics with timeout and bid rates
- **Content Health Monitoring** - Track content field availability across platforms (CTV/Audio)
- **Video Health Analytics** - Monitor video properties, protocols, and placement metrics
- **Role-based Access Control** - Multi-tier user permissions (Viewer/Analyst/Admin)
- **Export Functionality** - CSV export for all data tables
- **Responsive Design** - Mobile-friendly interface with dark/light mode support

## Tech Stack

### Backend
- **Go 1.21** with Gin framework
- **DuckDB** for high-performance analytics
- **JWT Authentication** with HttpOnly cookies
- **Rate limiting** and security middleware
- **Docker** containerization

### Frontend
- **React 18** with TypeScript
- **Vite** for fast development and builds
- **TanStack Query** for data fetching and caching
- **TanStack Table** for data tables with sorting/filtering
- **Recharts** for interactive charts and visualizations
- **Tailwind CSS** for styling

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Or: Go 1.21+ and Node.js 18+ for local development

### Using Docker (Recommended)

1. **Clone and start the application:**
```bash
git clone <repository-url>
cd openrtb-insights
docker-compose up -d
```

2. **Access the application:**
- Frontend: http://localhost
- Backend API: http://localhost:8080
- Health checks: http://localhost/health and http://localhost:8080/health

3. **Login with demo accounts:**
- Admin: `admin` / `admin123`
- Analyst: `analyst` / `admin123`
- Viewer: `viewer` / `admin123`

### Local Development

1. **Backend setup:**
```bash
cd backend
go mod download
go run scripts/seed-data.go  # Populate sample data
go run cmd/server/main.go
```

2. **Frontend setup:**
```bash
cd frontend
npm install
npm run dev
```

3. **Access:**
- Frontend: http://localhost:3000
- Backend: http://localhost:8080

## API Documentation

### Authentication Endpoints
- `POST /api/auth/login` - User login
- `POST /api/auth/refresh` - Refresh JWT token
- `POST /api/auth/logout` - User logout
- `GET /api/auth/me` - Get current user info

### Reports Endpoints (Protected)
- `GET /api/reports/dashboard` - Dashboard summary data
- `GET /api/reports/platform?start=YYYY-MM-DD&end=YYYY-MM-DD` - Platform statistics
- `GET /api/reports/content?platform={CTV|Audio}&start=YYYY-MM-DD&end=YYYY-MM-DD` - Content health
- `GET /api/reports/video?platform={CTV|Display|App}&start=YYYY-MM-DD&end=YYYY-MM-DD` - Video health

## Configuration

### Environment Variables

**Backend (.env):**
```bash
DB_PATH=./analytics.db
JWT_SECRET=your-jwt-secret-key
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h
PORT=8080
CORS_ORIGINS=http://localhost:3000
LOG_LEVEL=info
RATE_LIMIT=100
```

**Frontend (.env):**
```bash
VITE_API_BASE_URL=http://localhost:8080/api
VITE_APP_NAME=OpenRTB Insights
```

## Data Models

### Platform Stats
- Total requests, multi-impression counts
- Bid rates, timeout rates
- Deal counts, compliance metrics
- Invalid request tracking

### Content Health
- Platform-specific content field tracking
- Title, series, episode availability
- Genre, language, duration metrics
- Live stream vs. on-demand content

### Video Health
- CTV percentage by platform
- Placement and protocol metrics
- Skip behavior and duration tracking
- Bitrate and quality metrics

## Development

### Project Structure
```
├── backend/
│   ├── cmd/server/          # Main application
│   ├── internal/
│   │   ├── auth/            # Authentication logic
│   │   ├── database/        # Database connection & migrations
│   │   ├── reports/         # Business logic for reports
│   │   └── config/          # Configuration management
│   └── scripts/             # Utility scripts
├── frontend/
│   ├── src/
│   │   ├── api/             # API client
│   │   ├── components/      # React components
│   │   ├── hooks/           # Custom React hooks
│   │   ├── pages/           # Page components
│   │   └── types/           # TypeScript types
└── docker-compose.yml      # Production deployment
```

### Available Scripts

**Backend:**
```bash
go run cmd/server/main.go     # Start server
go run scripts/seed-data.go   # Generate sample data
go test ./...                 # Run tests
```

**Frontend:**
```bash
npm run dev                   # Development server
npm run build                # Production build
npm run lint                 # ESLint
npm run preview              # Preview build
```

**Docker:**
```bash
docker-compose up -d          # Production deployment
docker-compose -f docker-compose.dev.yml up  # Development
docker-compose logs -f        # View logs
docker-compose down           # Stop services
```

## Sample Data

The application includes a data seeding script that generates 30 days of realistic sample data:

```bash
# Backend container
go run scripts/seed-data.go

# Or via Docker
docker-compose exec backend go run scripts/seed-data.go
```

## Security Features

- **JWT Authentication** with automatic token refresh
- **HttpOnly cookies** for secure token storage
- **Rate limiting** (100 requests/minute per IP)
- **CORS protection** with configurable origins
- **Input validation** and SQL injection prevention
- **Security headers** (XSS, CSRF, Content-Type)
- **Non-root container users** for production

## Performance Optimizations

- **DuckDB** for fast analytical queries
- **React Query** with 30-second cache TTL
- **Data pagination** with configurable page sizes
- **Lazy loading** and code splitting
- **Nginx compression** and static asset caching
- **Database indexing** on date and platform fields

## Deployment

### Production Deployment
```bash
# Clone repository
git clone <repository-url>
cd openrtb-insights

# Configure environment variables
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env

# Deploy with Docker
docker-compose up -d

# Check health
curl http://localhost/health
curl http://localhost:8080/health
```

### CI/CD Integration
The application includes health checks and can be integrated with:
- GitHub Actions
- GitLab CI/CD
- Jenkins
- Kubernetes deployments

## Monitoring

- **Health check endpoints** for load balancers
- **Structured logging** with configurable levels
- **Request metrics** and error tracking
- **Database connection monitoring**
- **Real-time dashboard updates** every 30 seconds

## Troubleshooting

### Common Issues

1. **Database connection errors:**
   - Check file permissions for DuckDB file
   - Ensure data directory is writable

2. **Authentication issues:**
   - Verify JWT_SECRET configuration
   - Check browser cookies are enabled

3. **CORS errors:**
   - Update CORS_ORIGINS environment variable
   - Ensure frontend URL matches CORS configuration

4. **Port conflicts:**
   - Modify ports in docker-compose.yml
   - Update VITE_API_BASE_URL in frontend

### Log Collection
```bash
# View application logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Export logs
docker-compose logs backend > backend.log
docker-compose logs frontend > frontend.log
```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make changes and test thoroughly
4. Submit a pull request with detailed description

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions:
- Create an issue in the GitHub repository
- Review the troubleshooting section above
- Check application logs for error details