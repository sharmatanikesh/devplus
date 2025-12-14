# DevPulse Backend

This is the backend service for **DevPulse**, an AI-powered engineering intelligence platform. It is built with **Go** and uses **PostgreSQL** for data persistence and **Kestra** for AI workflow orchestration.

## Features

- 🔐 **GitHub OAuth Authentication** - Secure user authentication via GitHub
- 📊 **Repository Management** - Sync and manage GitHub repositories
- 🔄 **Pull Request Tracking** - Fetch and track pull requests with status updates
- 📈 **Engineering Metrics** - Calculate and track development metrics
- 🤖 **AI-Powered Analysis** - Integration with Kestra workflows for intelligent PR reviews and release notes
- 🔔 **Webhook Support** - GitHub webhook handling for real-time updates

## Tech Stack

- **Language**: Go (Golang)
- **Framework**: Gorilla Mux (Routing)
- **ORM**: GORM
- **Database**: PostgreSQL
- **Authentication**: GitHub OAuth 2.0
- **AI Workflows**: Kestra

## Prerequisites

- [Go](https://go.dev/dl/) (1.20+)
- [PostgreSQL](https://www.postgresql.org/download/) (13+)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/) (for Kestra)
- [Make](https://www.gnu.org/software/make/)

## Getting Started

### 1. Kestra Setup

Kestra is required for AI-powered workflows. Navigate to the kestra directory and start it using Docker:

```bash
cd ../kestra
docker run -d \
  --name kestra \
  -p 8080:8080 \
  -v $(pwd)/storage:/app/storage \
  kestra/kestra:latest
```

Or build from the Dockerfile:

```bash
cd ../kestra
docker build -t kestra-local .
docker run -d --name kestra -p 8080:8080 kestra-local
```

Kestra will be available at http://localhost:8080

See the [Kestra README](../kestra/README.md) for more details.

### 2. Database Setup

Ensure you have a PostgreSQL database running:

```sql
CREATE DATABASE devplus;
```

The application will automatically run migrations on startup (using GORM AutoMigrate).

### 3. Environment Variables

Create a `.env` file in the `backend` directory:

```bash
cp .env.example .env
```

Edit the `.env` file with your configuration. See `.env.example` for all required variables.

### 4. Running the Server

Use `make` to run the application:

```bash
# Run locally
make run
```

The server will start at `http://localhost:8081`.

### 5. Other Commands

```bash
# Install dependencies
make tidy

# Build binary
make build
```

## Kestra Workflows

The `workflows/` directory contains Kestra workflow definitions:

- `ai-pull-request-analysis.yaml` - AI-powered PR review workflow
- `ai-release-risk.yaml` - Release risk assessment workflow
- `ai-repo-analysis.yaml` - Repository analysis workflow

These workflows need to be deployed to your Kestra instance. You can upload them via the Kestra UI at http://localhost:8080 or use the Kestra CLI.

## API Endpoints

### Authentication

- `GET /api/v1/auth/github/login` - Initiates GitHub OAuth flow
- `GET /api/v1/auth/github/callback` - Handle OAuth callback
- `POST /api/v1/auth/logout` - Logout user

### Repositories

- `GET /api/v1/repos` - List synced repositories
- `GET /api/v1/repos/{id}` - Get repository details
- `POST /api/v1/repos/sync` - Sync repositories from GitHub
- `POST /api/v1/repos/{id}/sync` - Sync pull requests for a repository

### Pull Requests

- `GET /api/v1/repos/{id}/prs` - Get pull requests for a repository
- `GET /api/v1/repos/{id}/prs/{number}` - Get specific PR details
- `POST /api/v1/repos/{id}/prs/{number}/analyze` - Trigger AI PR analysis

### Metrics

- `GET /api/v1/metrics` - Get engineering metrics for user

### Webhooks

- `POST /api/v1/webhook/github` - GitHub webhook receiver
- `POST /api/v1/webhook/ai` - AI workflow callback

## Folder Structure

```
backend/
├── cmd/server/           # Entry point (main.go)
├── internal/
│   ├── config/          # Configuration and env loading
│   ├── controllers/     # HTTP handlers
│   │   └── rest/        # REST API controllers
│   ├── models/          # Database models
│   ├── services/        # Business logic
│   │   ├── auth_service/    # Authentication logic
│   │   ├── github_service/  # GitHub integration
│   │   └── ai/             # AI service factory
│   ├── repositories/    # Data access layer
│   ├── middleware/      # HTTP middleware (auth, CORS, session)
│   ├── router/          # Route definitions
│   ├── db/             # Database connection
│   └── migrations/     # SQL migrations
├── workflows/          # Kestra workflow definitions
├── pkg/
│   ├── logger/         # Logging utilities
│   └── utils/          # Utility functions
├── Dockerfile          # Docker configuration
├── docker-compose.yaml # Docker Compose (Kestra)
├── Makefile           # Build commands
└── go.mod             # Go dependencies
```

## Environment Variables Reference

See `.env.example` for a complete list of required environment variables:

- **Server**: PORT, BACKEND_URL, FRONTEND_URL
- **Database**: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
- **GitHub OAuth**: GITHUB_CLIENT_ID, GITHUB_CLIENT_SECRET, GITHUB_REDIRECT_URI
- **Kestra**: KESTRA_URL, KESTRA_USERNAME, KESTRA_PASSWORD
- **Environment**: ENVIRONMENT (development/production)

## Contributing

When contributing to this project:

1. Ensure all tests pass
2. Follow Go best practices and idioms
3. Update documentation for new features
4. Add proper error handling and logging

## License

MIT
