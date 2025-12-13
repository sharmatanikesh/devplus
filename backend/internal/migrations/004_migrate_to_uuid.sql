-- Enable pgcrypto for UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Drop tables in dependency order
DROP TABLE IF EXISTS public.commits;
DROP TABLE IF EXISTS public.metrics;
DROP TABLE IF EXISTS public.pull_requests;
DROP TABLE IF EXISTS public.repositories;
DROP TABLE IF EXISTS public.sessions;
DROP TABLE IF EXISTS public.auth_states;
DROP TABLE IF EXISTS public.users;

-- 1. Users Table
CREATE TABLE public.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    github_id BIGINT UNIQUE NOT NULL,
    username TEXT NOT NULL,
    email TEXT,
    avatar_url TEXT,
    access_token TEXT,
    refresh_token TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX idx_users_deleted_at ON public.users(deleted_at);
CREATE UNIQUE INDEX idx_users_github_id ON public.users(github_id);

-- 2. Auth States Table
CREATE TABLE public.auth_states (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    state TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 3. Sessions Table
CREATE TABLE public.sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX idx_sessions_user_id ON public.sessions(user_id);

-- 4. Repositories Table
CREATE TABLE public.repositories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    github_repo_id BIGINT UNIQUE,
    name TEXT,
    owner TEXT,
    url TEXT,
    installation_id BIGINT,
    user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX idx_repositories_github_repo_id ON public.repositories(github_repo_id);
CREATE INDEX idx_repositories_user_id ON public.repositories(user_id);
CREATE INDEX idx_repositories_deleted_at ON public.repositories(deleted_at);

-- 5. Pull Requests Table
CREATE TABLE public.pull_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    github_pr_id BIGINT UNIQUE,
    number BIGINT,
    title TEXT,
    state TEXT,
    repo_id UUID REFERENCES public.repositories(id) ON DELETE CASCADE,
    author_id BIGINT,
    author_name TEXT,
    ai_summary TEXT,
    ai_decision TEXT,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX idx_pull_requests_github_pr_id ON public.pull_requests(github_pr_id);
CREATE INDEX idx_pull_requests_deleted_at ON public.pull_requests(deleted_at);

-- 6. Metrics Table
CREATE TABLE public.metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_id UUID REFERENCES public.repositories(id) ON DELETE CASCADE,
    type TEXT,
    value NUMERIC,
    date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX idx_metrics_deleted_at ON public.metrics(deleted_at);

-- 7. Commits Table
CREATE TABLE public.commits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sha TEXT UNIQUE,
    message TEXT,
    author_name TEXT,
    repo_id UUID REFERENCES public.repositories(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX idx_commits_sha ON public.commits(sha);
CREATE INDEX idx_commits_deleted_at ON public.commits(deleted_at);
