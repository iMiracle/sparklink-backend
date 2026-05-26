# Sparklink Backend

API Server for Sparklink VPN Application

## Quick Start

```bash
# Install dependencies
go mod tidy

# Run
go run ./cmd/api
```

## Environment Variables

```bash
# Server
PORT=8080

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=sparklink

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-secret-key

# SMS (Optional)
SMS_API_KEY=
SMS_APP_SECRET=
SMS_ENDPOINT=
```

## API Endpoints

### Auth
- POST /api/v1/auth/sendcode - Send verification code
- POST /api/v1/auth/login - Login/Register
- POST /api/v1/auth/refresh - Refresh token
- POST /api/v1/auth/logout - Logout

### Nodes
- GET /api/v1/nodes/list - Get node list
- GET /api/v1/nodes/:id - Get node details
- POST /api/v1/nodes/ping - Update latency

### Rewards
- POST /api/v1/rewards/video - Claim video reward
- POST /api/v1/rewards/daily - Daily check-in
- GET /api/v1/rewards/invite - Get invite info

### Subscription
- GET /api/v1/subscription/plans - Get plan list
- POST /api/v1/subscription/create - Create subscription