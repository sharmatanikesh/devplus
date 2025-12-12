-- Create users table if not exists
CREATE TABLE IF NOT EXISTS public.users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    github_id BIGINT UNIQUE NOT NULL,
    username TEXT NOT NULL,
    email TEXT,
    avatar_url TEXT,
    access_token TEXT,
    refresh_token TEXT
);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON public.users(deleted_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_github_id ON public.users(github_id);

-- Create auth_states table
CREATE TABLE IF NOT EXISTS public.auth_states (
    id BIGSERIAL PRIMARY KEY,
    state TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create sessions table
CREATE TABLE IF NOT EXISTS public.sessions (
    id TEXT PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON public.sessions(user_id);

-- Create repositories table
CREATE TABLE IF NOT EXISTS public.repositories (
  id bigserial not null,
  created_at timestamp with time zone null,
  updated_at timestamp with time zone null,
  deleted_at timestamp with time zone null,
  github_repo_id bigint null,
  name text null,
  owner text null,
  url text null,
  installation_id bigint null,
  constraint repositories_pkey primary key (id)
) TABLESPACE pg_default;

CREATE UNIQUE INDEX IF NOT EXISTS idx_repositories_github_repo_id on public.repositories using btree (github_repo_id) TABLESPACE pg_default;
CREATE INDEX IF NOT EXISTS idx_repositories_deleted_at on public.repositories using btree (deleted_at) TABLESPACE pg_default;

-- Create pull_requests table
CREATE TABLE IF NOT EXISTS public.pull_requests (
  id bigserial not null,
  created_at timestamp with time zone null,
  updated_at timestamp with time zone null,
  deleted_at timestamp with time zone null,
  github_pr_id bigint null,
  number bigint null,
  title text null,
  state text null,
  repo_id bigint null,
  author_id bigint null,
  author_name text null,
  ai_summary text null,
  a_idecision text null,
  constraint pull_requests_pkey primary key (id),
  constraint fk_pull_requests_repository foreign KEY (repo_id) references repositories (id)
) TABLESPACE pg_default;

CREATE UNIQUE INDEX IF NOT EXISTS idx_pull_requests_github_pr_id on public.pull_requests using btree (github_pr_id) TABLESPACE pg_default;
CREATE INDEX IF NOT EXISTS idx_pull_requests_deleted_at on public.pull_requests using btree (deleted_at) TABLESPACE pg_default;

-- Create metrics table
CREATE TABLE IF NOT EXISTS public.metrics (
  id bigserial not null,
  created_at timestamp with time zone null,
  updated_at timestamp with time zone null,
  deleted_at timestamp with time zone null,
  repo_id bigint null,
  type text null,
  value numeric null,
  date timestamp with time zone null,
  constraint metrics_pkey primary key (id),
  constraint fk_metrics_repository foreign KEY (repo_id) references repositories (id)
) TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS idx_metrics_deleted_at on public.metrics using btree (deleted_at) TABLESPACE pg_default;

-- Create commits table
CREATE TABLE IF NOT EXISTS public.commits (
  id bigserial not null,
  created_at timestamp with time zone null,
  updated_at timestamp with time zone null,
  deleted_at timestamp with time zone null,
  sha text null,
  message text null,
  author_name text null,
  repo_id bigint null,
  constraint commits_pkey primary key (id),
  constraint fk_commits_repository foreign KEY (repo_id) references repositories (id)
) TABLESPACE pg_default;

CREATE UNIQUE INDEX IF NOT EXISTS idx_commits_sha on public.commits using btree (sha) TABLESPACE pg_default;
CREATE INDEX IF NOT EXISTS idx_commits_deleted_at on public.commits using btree (deleted_at) TABLESPACE pg_default;
