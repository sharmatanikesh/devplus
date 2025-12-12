# DevPulse Frontend

AI-Powered Engineering Intelligence platform for automated PR reviews, release notes generation, and engineering metrics.

## Features

- 🤖 **Automated PR Reviews** - AI analyzes pull requests and provides intelligent feedback
- 📦 **Smart Release Notes** - Generate comprehensive changelogs automatically
- 📊 **Engineering Metrics** - Track PR lead time, review latency, and team performance
- 🔍 **Architecture Impact Analysis** - Identify hotspots and analyze code coupling
- 🔗 **GitHub Integration** - Seamless OAuth and webhook integration

## Tech Stack

- **Framework**: Next.js 16 (App Router)
- **UI Library**: ShadCN UI + Tailwind CSS
- **Language**: TypeScript
- **State Management**: React Hooks
- **API**: RESTful API routes

## Getting Started

### Prerequisites

- Node.js 20+
- npm or yarn
- GitHub OAuth App credentials

### Installation

1. Clone the repository and install dependencies:

```bash
npm install
```

2. Set up environment variables:

```bash
cp .env.example .env.local
```

Edit `.env.local` and add your GitHub OAuth credentials:

```env
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
GITHUB_REDIRECT_URI=http://localhost:3000/api/v1/auth/github/callback
NEXTAUTH_SECRET=your_random_secret_key
```

3. Run the development server:

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) to see the application.

## Project Structure

```
app/
├── api/v1/                 # API Routes
│   ├── auth/              # Authentication endpoints
│   ├── repos/             # Repository management
│   ├── metrics/           # Engineering metrics
│   ├── impact/            # Impact analysis
│   ├── webhook/           # GitHub webhooks
│   └── health/            # Health check
├── dashboard/             # Dashboard pages
│   ├── layout.tsx         # Dashboard layout with sidebar
│   ├── page.tsx           # Main dashboard
│   ├── repositories/      # Repositories management
│   ├── pull-requests/     # PR list and details
│   ├── metrics/           # Metrics dashboard
│   └── releases/          # Release notes
├── login/                 # Authentication page
└── globals.css           # Global styles

components/
└── ui/                    # ShadCN UI components

lib/
├── types/                 # TypeScript type definitions
├── utils/                 # Utility functions
├── constants.ts           # App constants
└── api-client.ts          # API client wrapper
```

## API Endpoints

### Authentication

- `POST /api/v1/auth/github/connect` - Start GitHub OAuth
- `GET /api/v1/auth/github/callback` - OAuth callback

### Repositories

- `GET /api/v1/repos` - List connected repositories
- `POST /api/v1/repos/:id/sync` - Trigger repository sync
- `GET /api/v1/repos/:id/prs` - List pull requests
- `GET /api/v1/repos/:id/prs/:pr_number` - Get PR details
- `POST /api/v1/repos/:id/prs/:pr_number/analyze` - Trigger AI analysis
- `POST /api/v1/repos/:id/release` - Generate release notes

### Metrics & Analysis

- `GET /api/v1/metrics` - Get engineering metrics
- `GET /api/v1/impact/:pr_id` - Get impact analysis

### System

- `POST /api/v1/webhook/github` - GitHub webhook receiver
- `GET /api/v1/health` - Health check

## Development

### Adding New Pages

Pages follow Next.js App Router conventions. Create files in `app/dashboard/`:

```tsx
// app/dashboard/new-page/page.tsx
export default function NewPage() {
  return <div>New Page</div>;
}
```

### Adding UI Components

Use ShadCN CLI to add new components:

```bash
npx shadcn@latest add [component-name]
```

### Type Definitions

All types are defined in `lib/types/index.ts`. Import and use:

```tsx
import type { Repository, PullRequest } from "@/lib/types";
```

## Production Deployment

1. Build the application:

```bash
npm run build
```

2. Start production server:

```bash
npm start
```

## Environment Variables

| Variable               | Description                   | Required |
| ---------------------- | ----------------------------- | -------- |
| `GITHUB_CLIENT_ID`     | GitHub OAuth App Client ID    | Yes      |
| `GITHUB_CLIENT_SECRET` | GitHub OAuth App Secret       | Yes      |
| `GITHUB_REDIRECT_URI`  | OAuth callback URL            | Yes      |
| `NEXTAUTH_SECRET`      | Secret for session encryption | Yes      |
| `NEXTAUTH_URL`         | Application base URL          | Yes      |
| `API_BASE_URL`         | Backend API URL (if separate) | No       |

## Contributing

This is a hackathon MVP project. For production use, consider:

1. Implementing proper authentication with NextAuth.js
2. Adding database integration (PostgreSQL, MongoDB, etc.)
3. Implementing proper error handling and logging
4. Adding comprehensive testing
5. Setting up CI/CD pipelines
6. Implementing rate limiting and security measures

## License

MIT
