# DevPlus Frontend

AI-Powered Engineering Intelligence platform for automated PR reviews, release notes generation, and engineering metrics.

## Features

- 🔐 **GitHub OAuth Authentication** - Secure login with GitHub
- 📦 **Repository Management** - View and manage connected GitHub repositories
- 🔄 **Pull Request Tracking** - Monitor pull requests with real-time sync
- 📊 **Engineering Metrics Dashboard** - Track PR lead time, review latency, and team performance
- 📝 **Release Notes Generation** - AI-powered release notes and changelogs
- 🤖 **AI PR Analysis** - Intelligent pull request reviews powered by Kestra workflows

## Tech Stack

- **Framework**: Next.js 16 (App Router)
- **UI Library**: ShadCN UI + Tailwind CSS
- **Language**: TypeScript
- **State Management**: Redux Toolkit
- **Styling**: Tailwind CSS
- **API Client**: Axios

## Getting Started

### Prerequisites

- Node.js 20+
- npm or yarn
- Backend API running (see [Backend README](../backend/README.md))

### Installation

1. Clone the repository and install dependencies:

```bash
npm install
```

2. Set up environment variables:

```bash
cp .env.example .env.local
```

Edit `.env.local` with your configuration. See `.env.example` for all required variables.

3. Run the development server:

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) to see the application.

## Project Structure

```
app/
├── (dashboard)/           # Dashboard routes (protected)
│   ├── layout.tsx        # Dashboard layout with sidebar
│   ├── dashboard/        # Main dashboard page
│   ├── repositories/     # Repositories management
│   ├── pull-requests/    # PR details pages
│   ├── metrics/          # Metrics dashboard
│   └── releases/         # Release notes
├── login/                # Authentication page
├── privacy/              # Privacy policy page
├── terms/                # Terms of service page
├── layout.tsx            # Root layout
├── page.tsx              # Landing page
└── globals.css           # Global styles

components/
├── dashboard/            # Dashboard-specific components
├── repositories/         # Repository components
└── ui/                   # ShadCN UI components

lib/
├── api-client.ts         # API client wrapper
├── constants.ts          # App constants
├── utils.ts              # Utility functions
├── store/                # Redux store configuration
│   ├── index.ts          # Store setup
│   ├── hooks.ts          # Redux hooks
│   ├── StoreProvider.tsx # Redux provider
│   └── slices/           # Redux slices
├── types/                # TypeScript type definitions
└── utils/                # Utility functions
    ├── auth.ts           # Auth utilities
    └── format.ts         # Formatting utilities
```

## API Integration

The frontend communicates with the backend API using axios. The base URL is configured via environment variables.

### API Client Usage

```tsx
import { apiClient } from "@/lib/api-client";

// Example: Fetch repositories
const response = await apiClient.repos.list();
if (response.success && response.data) {
  setRepositories(response.data);
}
```

## Available Scripts

```bash
# Development
npm run dev          # Start development server

# Production
npm run build        # Build for production
npm start            # Start production server

# Code Quality
npm run lint         # Run ESLint
npm run type-check   # Run TypeScript type checking

# Format
npm run format       # Format code with Prettier
```

## Development

### Adding New Pages

Pages follow Next.js App Router conventions. Create files in `app/(dashboard)/`:

```tsx
// app/(dashboard)/new-page/page.tsx
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

## Environment Variables

See `.env.example` for all required environment variables:

- `NEXT_PUBLIC_API_URL` - Backend API base URL
- Additional variables as needed for production deployment

## Contributing

When contributing to this project:

1. Follow the existing code structure and conventions
2. Use TypeScript types for all new code
3. Add proper error handling
4. Update documentation for new features
5. Test your changes thoroughly

## License

MIT
