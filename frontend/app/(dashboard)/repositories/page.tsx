'use client';

import { useState } from 'react';
import { useRepositories } from '@/hooks/use-repositories';
import { RepositoryList } from '@/components/repositories/repository-list';
import { Button } from '@/components/ui/button';
import { RefreshCw, Loader2 } from 'lucide-react';
import { apiClient } from '@/lib/api-client';

export default function RepositoriesPage() {
  const { repositories, isLoading, error, refetch } = useRepositories();
  const [isSyncing, setIsSyncing] = useState(false);

  const handleSync = async () => {
    setIsSyncing(true);
    try {
      if (apiClient.repos.syncAll) {
        await apiClient.repos.syncAll();
      } else {
        // Fallback if type definition hasn't reloaded in IDE yet, though runtime should be fine
        // casting to any to avoid TS error during transition
        await (apiClient.repos as any).sync();
      }
      await refetch();
    } catch (error) {
      console.error('Failed to sync repositories:', error);
    } finally {
      setIsSyncing(false);
    }
  };

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[400px] space-y-4 text-center">
        <p className="text-red-500 font-medium">{error}</p>
        <Button variant="outline" onClick={() => window.location.reload()}>
          Try Again
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Repositories</h1>
          <p className="text-muted-foreground">
            Manage your synced repositories and configurations
          </p>
        </div>
        <Button disabled={isLoading || isSyncing} onClick={handleSync}>
          {isSyncing ? (
            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
          ) : (
            <RefreshCw className="mr-2 h-4 w-4" />
          )}
          {isSyncing ? 'Syncing...' : 'Sync Now'}
        </Button>
      </div>

      {isLoading ? (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {[1, 2, 3, 4, 5, 6].map((i) => (
            <div key={i} className="h-[180px] rounded-xl border bg-card text-card-foreground shadow-sm p-6 space-y-4">
              <div className="h-6 w-1/2 bg-muted animate-pulse rounded" />
              <div className="h-4 w-1/3 bg-muted animate-pulse rounded" />
              <div className="mt-auto pt-4 h-8 w-full bg-muted animate-pulse rounded" />
            </div>
          ))}
        </div>
      ) : (
        <RepositoryList repositories={repositories} />
      )}
    </div>
  );
}
