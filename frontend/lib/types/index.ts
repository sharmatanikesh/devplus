// ==================== Auth Types ====================
export interface User {
  id: string;
  githubId: string;
  username: string;
  email: string;
  avatarUrl: string;
  accessToken?: string;
  createdAt: string;
  updatedAt: string;
}

export interface AuthResponse {
  user: User;
  token: string;
}

// ==================== Repository Types ====================
export interface Repository {
  id: string;
  githubId: number;
  name: string;
  fullName: string;
  owner: string;
  description?: string;
  url: string;
  defaultBranch: string;
  isPrivate: boolean;
  language?: string;
  stars: number;
  forks: number;
  syncStatus: "pending" | "syncing" | "synced" | "error";
  lastSyncedAt?: string;
  webhookConfigured: boolean;
  createdAt: string;
  updatedAt: string;
}

// ==================== Pull Request Types ====================
export interface PullRequest {
  id: string;
  githubId: number;
  repositoryId: string;
  number: number;
  title: string;
  description?: string;
  state: "open" | "closed" | "merged";
  author: {
    username: string;
    avatarUrl: string;
  };
  baseBranch: string;
  headBranch: string;
  url: string;
  filesChanged: number;
  additions: number;
  deletions: number;
  commits: number;
  reviewStatus?: "pending" | "approved" | "changes_requested" | "blocked";
  aiAnalysisStatus?: "pending" | "analyzing" | "completed" | "failed";
  aiSummary?: string;
  aiDecision?: "approve" | "request_changes" | "block";
  aiComments?: AIComment[];
  createdAt: string;
  updatedAt: string;
  mergedAt?: string;
  closedAt?: string;
}

export interface AIComment {
  id: string;
  file: string;
  line: number;
  body: string;
  severity: "info" | "warning" | "error";
  createdAt: string;
}

export interface PRAnalysisResult {
  prId: string;
  summary: string;
  decision: "approve" | "request_changes" | "block";
  confidence: number;
  issues: {
    file: string;
    line: number;
    type: string;
    severity: "info" | "warning" | "error";
    message: string;
  }[];
  suggestions: string[];
  impactScore: number;
  analyzedAt: string;
}

// ==================== Metrics Types ====================
export interface EngineeringMetrics {
  period: {
    from: string;
    to: string;
  };
  repositories?: string[];
  metrics: {
    prLeadTime: {
      average: number;
      median: number;
      p95: number;
      trend: number; // percentage change
    };
    reviewLatency: {
      average: number;
      median: number;
      p95: number;
      trend: number;
    };
    mergeRate: {
      percentage: number;
      total: number;
      merged: number;
      trend: number;
    };
    throughput: {
      prsPerWeek: number;
      commitsPerWeek: number;
      trend: number;
    };
    codeChurn: {
      average: number;
      trend: number;
    };
  };
  topContributors: {
    username: string;
    avatarUrl: string;
    prsCreated: number;
    prsReviewed: number;
    commitsCount: number;
  }[];
  hotspots: {
    file: string;
    changes: number;
    authors: number;
    lastModified: string;
  }[];
}

// ==================== Impact Analysis Types ====================
export interface ImpactAnalysis {
  prId: string;
  moduleCoupling: {
    affectedModules: string[];
    couplingScore: number;
    depth: number;
  };
  filesTouched: {
    path: string;
    changes: number;
    isHotspot: boolean;
    complexity: number;
  }[];
  hotspotFlags: {
    file: string;
    reason: string;
    severity: "low" | "medium" | "high";
    recommendation: string;
  }[];
  riskScore: number;
  recommendations: string[];
  analyzedAt: string;
}

// ==================== Release Notes Types ====================
export interface ReleaseNote {
  id: string;
  repositoryId: string;
  version: string;
  title: string;
  description: string;
  tagName: string;
  targetCommitish: string;
  isDraft: boolean;
  isPrerelease: boolean;
  changes: {
    type: "feature" | "bugfix" | "breaking" | "improvement" | "docs" | "chore";
    description: string;
    pr?: number;
    commit?: string;
  }[];
  createdAt: string;
  publishedAt?: string;
}

// ==================== Webhook Types ====================
export interface WebhookPayload {
  event: "pull_request" | "push" | "check_run" | "check_suite";
  action: string;
  repository: {
    id: number;
    name: string;
    full_name: string;
  };
  sender: {
    login: string;
    avatar_url: string;
  };
  [key: string]: unknown;
}

// ==================== API Response Types ====================
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
    details?: Record<string, unknown>;
  };
  meta?: {
    page?: number;
    limit?: number;
    total?: number;
  };
}

export interface PaginatedResponse<T> {
  items: T[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

// ==================== Job/Task Types ====================
export interface JobStatus {
  id: string;
  type: "pr_analysis" | "release_generation" | "repo_sync";
  status: "queued" | "running" | "completed" | "failed";
  progress?: number;
  result?: Record<string, unknown>;
  error?: string;
  createdAt: string;
  completedAt?: string;
}

// ==================== Health Check ====================
export interface HealthCheck {
  status: "healthy" | "degraded" | "unhealthy";
  version: string;
  timestamp: string;
  services: {
    database: "up" | "down";
    github: "up" | "down";
    kestra: "up" | "down";
    ai: "up" | "down";
  };
}
