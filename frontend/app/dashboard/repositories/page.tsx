'use client';

import { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { RefreshCw, GitBranch, Star, GitFork, Clock, Check, AlertCircle, Loader2 } from 'lucide-react';
import Link from 'next/link';
import type { Repository } from '@/lib/types';
import { formatRelativeTime } from '@/lib/utils/format';

export default function RepositoriesPage() {
  const [repositories, setRepositories] = useState<Repository[]>([]);
  const [loading, setLoading] = useState(true);
  const [syncing, setSyncing] = useState<Record<string, boolean>>({});

  useEffect(() => {
    fetchRepositories();
  }, []);

  const fetchRepositories = async () => {
    try {
      setLoading(true);
      const response = await fetch('/api/v1/repos');
      const data = await response.json();
      if (data.success) {
        setRepositories(data.data.items);
      }
    } catch (error) {
      console.error('Failed to fetch repositories:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSync = async (repoId: string) => {
    setSyncing({ ...syncing, [repoId]: true });
    try {
      await fetch(`/api/v1/repos/${repoId}/sync`, { method: 'POST' });
      // Wait a moment then refresh
      setTimeout(() => {
        fetchRepositories();
        setSyncing({ ...syncing, [repoId]: false });
      }, 2000);
    } catch (error) {
      console.error('Sync failed:', error);
      setSyncing({ ...syncing, [repoId]: false });
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-96">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Repositories</h1>
          <p className="text-muted-foreground">
            Manage your connected GitHub repositories
          </p>
        </div>
        <Button onClick={() => fetchRepositories()}>
          <RefreshCw className="mr-2 h-4 w-4" />
          Refresh
        </Button>
      </div>

      {/* Repository List */}
      <div className="grid gap-4">
        {repositories.map((repo) => (
          <Card key={repo.id} className="hover:shadow-md transition-shadow">
            <CardHeader>
              <div className="flex items-start justify-between">
                <div className="space-y-1">
                  <div className="flex items-center gap-2">
                    <GitBranch className="h-5 w-5 text-muted-foreground" />
                    <CardTitle className="text-xl">
                      <Link
                        href={`/repositories/${repo.id}`}
                        className="hover:underline"
                      >
                        {repo.fullName}
                      </Link>
                    </CardTitle>
                    {repo.isPrivate && (
                      <Badge variant="outline">Private</Badge>
                    )}
                  </div>
                  <CardDescription>{repo.description || 'No description'}</CardDescription>
                </div>
                <SyncStatusBadge
                  status={repo.syncStatus}
                  syncing={syncing[repo.id]}
                />
              </div>
            </CardHeader>
            <CardContent>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-6 text-sm text-muted-foreground">
                  {repo.language && (
                    <div className="flex items-center gap-1">
                      <div className="h-3 w-3 rounded-full bg-blue-500" />
                      {repo.language}
                    </div>
                  )}
                  <div className="flex items-center gap-1">
                    <Star className="h-4 w-4" />
                    {repo.stars}
                  </div>
                  <div className="flex items-center gap-1">
                    <GitFork className="h-4 w-4" />
                    {repo.forks}
                  </div>
                  {repo.lastSyncedAt && (
                    <div className="flex items-center gap-1">
                      <Clock className="h-4 w-4" />
                      Synced {formatRelativeTime(repo.lastSyncedAt)}
                    </div>
                  )}
                </div>
                <div className="flex gap-2">
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => handleSync(repo.id)}
                    disabled={syncing[repo.id]}
                  >
                    {syncing[repo.id] ? (
                      <>
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        Syncing...
                      </>
                    ) : (
                      <>
                        <RefreshCw className="mr-2 h-4 w-4" />
                        Sync
                      </>
                    )}
                  </Button>
                  <Button size="sm" asChild>
                    <Link href={`/repositories/${repo.id}/prs`}>
                      View PRs
                    </Link>
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        ))}

        {repositories.length === 0 && (
          <Card>
            <CardContent className="flex flex-col items-center justify-center py-12">
              <GitBranch className="h-12 w-12 text-muted-foreground mb-4" />
              <h3 className="text-lg font-semibold mb-2">No repositories connected</h3>
              <p className="text-muted-foreground text-center max-w-md">
                Connect your GitHub repositories to start using DevPulse&apos;s AI-powered PR reviews and metrics.
              </p>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}

function SyncStatusBadge({ status, syncing }: { status: string; syncing?: boolean }) {
  if (syncing) {
    return (
      <Badge className="bg-blue-100 text-blue-800">
        <Loader2 className="mr-1 h-3 w-3 animate-spin" />
        Syncing
      </Badge>
    );
  }

  switch (status) {
    case 'synced':
      return (
        <Badge className="bg-green-100 text-green-800">
          <Check className="mr-1 h-3 w-3" />
          Synced
        </Badge>
      );
    case 'error':
      return (
        <Badge className="bg-red-100 text-red-800">
          <AlertCircle className="mr-1 h-3 w-3" />
          Error
        </Badge>
      );
    case 'pending':
      return (
        <Badge className="bg-gray-100 text-gray-800">
          <Clock className="mr-1 h-3 w-3" />
          Pending
        </Badge>
      );
    default:
      return null;
  }
}
