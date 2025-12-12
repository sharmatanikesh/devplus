// API Configuration
// API Configuration
export const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';
export const APP_URL = process.env.NEXT_PUBLIC_APP_URL || 'http://localhost:3000';

// App Metadata
export const APP_NAME = 'DevPulse';
export const APP_DESCRIPTION = 'AI-powered PR review and engineering metrics platform';

// API Endpoints
export const API_ENDPOINTS = {
  // Auth
  AUTH_GITHUB_CONNECT: '/v1/auth/github/login',
  AUTH_GITHUB_CALLBACK: '/v1/auth/github/callback',

  // Repos
  REPOS_LIST: '/v1/repos',
  REPOS_DETAIL: (id: string) => `/v1/repos/${id}`,
  REPOS_SYNC: (id: string) => `/v1/repos/${id}/sync`,
  REPOS_PRS: (id: string) => `/v1/repos/${id}/prs`,

  // PRs
  PR_DETAIL: (repoId: string, prNumber: number) => `/v1/repos/${repoId}/prs/${prNumber}`,
  PR_ANALYZE: (repoId: string, prNumber: number) => `/v1/repos/${repoId}/prs/${prNumber}/analyze`,

  // Release
  RELEASE_CREATE: (repoId: string) => `/v1/repos/${repoId}/release`,

  // Metrics & Impact
  METRICS: '/v1/metrics',
  IMPACT: (prId: string) => `/v1/impact/${prId}`,

  // Webhook
  WEBHOOK_GITHUB: '/v1/webhook/github',

  // System
  HEALTH: '/v1/health',
} as const;

// Dashboard Navigation
export const NAVIGATION_ITEMS = [
  {
    title: 'Dashboard',
    href: '/dashboard',
    icon: 'LayoutDashboard',
  },
  {
    title: 'Repositories',
    href: '/repositories',
    icon: 'GitBranch',
  },
  {
    title: 'Pull Requests',
    href: '/pull-requests',
    icon: 'GitPullRequest',
  },
  {
    title: 'Metrics',
    href: '/metrics',
    icon: 'BarChart3',
  },
  {
    title: 'Release Notes',
    href: '/releases',
    icon: 'Package',
  },
] as const;

// PR Status Colors
export const PR_STATUS_COLORS = {
  open: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
  closed: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
  merged: 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200',
} as const;

// Review Status Colors
export const REVIEW_STATUS_COLORS = {
  pending: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200',
  approved: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
  changes_requested: 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200',
  blocked: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
} as const;

// Sync Status Colors
export const SYNC_STATUS_COLORS = {
  pending: 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200',
  syncing: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200',
  synced: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
  error: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
} as const;

// Date Formats
export const DATE_FORMATS = {
  FULL: 'PPP',
  SHORT: 'PP',
  TIME: 'p',
  DATETIME: 'PPp',
  RELATIVE: 'relative',
} as const;

// Pagination
export const DEFAULT_PAGE_SIZE = 20;
export const PAGE_SIZE_OPTIONS = [10, 20, 50, 100];

// Chart Colors
export const CHART_COLORS = {
  primary: 'hsl(var(--primary))',
  secondary: 'hsl(var(--secondary))',
  success: 'hsl(142, 76%, 36%)',
  warning: 'hsl(38, 92%, 50%)',
  error: 'hsl(0, 84%, 60%)',
  info: 'hsl(199, 89%, 48%)',
} as const;

// Local Storage Keys
export const STORAGE_KEYS = {
  AUTH_TOKEN: 'devpulse_auth_token',
  USER: 'devpulse_user',
  THEME: 'devpulse_theme',
} as const;
