# DevPulse Backend

This is the backend service for **DevPulse**, an AI-powered engineering metrics and PR review platform. It is built with **Go** and uses **PostgreSQL** for data persistence.

## Tech Stack

- **Language**: Go (Golang)
- **Framework**: Gorilla Mux (Routing)
- **ORM**: GORM
- **Database**: PostgreSQL
- **Authentication**: GitHub OAuth 2.0

## Prerequisites

- [Go](https://go.dev/dl/) (1.20+)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Make](https://www.gnu.org/software/make/)

## Getting Started

### 1. Database Setup

Ensure you have a PostgreSQL database running.

```sql
CREATE DATABASE devplus;
```

The application will automatically run migrations on startup (using GORM AutoMigrate).

### 2. Environment Variables

Create a `.env` file in the `backend` directory:

```bash
cp .env.example .env
```

Configuration keys:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DB_HOST` | Database host | `localhost` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `DB_NAME` | Database name | `devplus` |
| `DB_PORT` | Database port | `5432` |
| `GITHUB_CLIENT_ID` | GitHub OAuth Client ID | - |
| `GITHUB_CLIENT_SECRET` | GitHub OAuth Client Secret | - |
| `GITHUB_REDIRECT_URI` | OAuth Redirect URI | `http://localhost:8080/api/v1/auth/github/callback` |
| `FRONTEND_URL` | Frontend Dashboard URL | `http://localhost:3000/dashboard` |

### 3. Running the Server

Use `make` to run the application:

```bash
# Run locally
make run
```

The server will start at `http://localhost:8080`.

### 4. Other Commands

```bash
# Install dependencies
make tidy

# Build binary (Mac ARM64)
make build-mac
```

## API Endpoints

### Authentication
- `GET /api/v1/auth/github/login`: Initiates GitHub OAuth flow.
- `GET /api/v1/auth/github/callback`: Handle OAuth callback.
- `POST /api/v1/auth/logout`: Logout user.

### Repositories
- `GET /api/v1/repos`: List synced repositories (DB cached).
- `GET /api/v1/repos/{id}`: Get details for a specific repository.
- `POST /api/v1/repos/sync`: Sync repositories from GitHub.
- `GET /api/v1/repos/{owner}/{repo}/pulls`: Get pull requests for a repository.

### Dashboard
- `GET /api/v1/dashboard/stats`: Get aggregated stats (PR counts, etc.).
- `GET /api/v1/dashboard/recent-prs`: Get recently active pull requests.

## Folder Structure

- `cmd/server`: Entry point (`main.go`).
- `internal/config`: Configuration and env loading.
- `internal/controllers`: HTTP handlers.
- `internal/models`: Database structs.
- `internal/services`: Business logic (Auth, GitHub integration).
- `internal/middleware`: Auth and CORS middleware.
- `internal/router`: Route definitions.
