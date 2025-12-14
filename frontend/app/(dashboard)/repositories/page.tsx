'use client';

import { motion } from 'framer-motion';

import { useEffect } from 'react';
import { useRepositories } from '@/hooks/use-repositories';
import { RepositoryList } from '@/components/repositories/repository-list';
import { Button } from '@/components/ui/button';
import { RefreshCw, Loader2, ListFilter } from 'lucide-react';
import { apiClient } from '@/lib/api-client';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { Checkbox } from '@/components/ui/checkbox';
import { Label } from '@/components/ui/label';
import { useAppDispatch, useAppSelector } from '@/lib/store/hooks';
import {
  initializeSelectedRepos,
  toggleRepoVisibility,
  selectAllRepos,
  deselectAllRepos,
} from '@/lib/store/slices/repositoryVisibilitySlice';
import { useState } from 'react';

export default function RepositoriesPage() {
  const { repositories, isLoading, error, refetch } = useRepositories();
  const [isSyncing, setIsSyncing] = useState(false);
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  const dispatch = useAppDispatch();
  const { selectedRepoIds, isInitialized } = useAppSelector(
    (state) => state.repositoryVisibility
  );

  // Initialize selectedRepoIds when repositories load
  useEffect(() => {
    if (repositories.length > 0 && !isInitialized) {
      dispatch(initializeSelectedRepos(repositories.map(repo => repo.id)));
    }
  }, [repositories, isInitialized, dispatch]);

  const handleToggleRepo = (repoId: string) => {
    dispatch(toggleRepoVisibility(repoId));
  };

  const handleToggleSelectAll = () => {
    if (selectedRepoIds.length === repositories.length) {
      dispatch(deselectAllRepos());
    } else {
      dispatch(selectAllRepos(repositories.map(repo => repo.id)));
    }
  };

  const filteredRepositories = repositories.filter(repo =>
    selectedRepoIds.includes(repo.id)
  );

  const handleSync = async () => {
    setIsSyncing(true);
    try {
      if (apiClient.repos.syncAll) {
        await apiClient.repos.syncAll();
      } else {
        // Fallback if type definition hasn't reloaded in IDE yet, though runtime should be fine
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
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="space-y-6"
    >
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight bg-clip-text text-transparent bg-linear-to-r from-foreground to-foreground/70">Repositories</h1>
          <p className="text-muted-foreground">
            Manage your synced repositories and configurations
          </p>
        </div>
        <div className="flex gap-2">
          <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
            <DialogTrigger asChild>
              <Button variant="outline" className="border-sky-100 hover:bg-sky-50 hover:text-sky-600 dark:border-sky-900/50 dark:hover:bg-sky-900/20">
                <ListFilter className="mr-2 h-4 w-4" />
                Edit List
              </Button>
            </DialogTrigger>
            <DialogContent className="max-w-md max-h-[80vh] flex flex-col">
              <DialogHeader>
                <DialogTitle>Edit Repository List</DialogTitle>
                <DialogDescription>
                  Select which repositories to display on your dashboard.
                </DialogDescription>
              </DialogHeader>
              <div className="flex-1 overflow-y-auto space-y-4 pr-2">
                <div className="flex items-center space-x-2 py-2 border-b">
                  <Checkbox
                    id="select-all"
                    checked={selectedRepoIds.length === repositories.length && repositories.length > 0}
                    onCheckedChange={handleToggleSelectAll}
                  />
                  <Label htmlFor="select-all" className="font-semibold cursor-pointer">
                    Select All ({selectedRepoIds.length}/{repositories.length})
                  </Label>
                </div>
                {repositories.map((repo) => (
                  <div key={repo.id} className="flex items-center space-x-2 py-2">
                    <Checkbox
                      id={repo.id}
                      checked={selectedRepoIds.includes(repo.id)}
                      onCheckedChange={() => handleToggleRepo(repo.id)}
                    />
                    <Label htmlFor={repo.id} className="flex-1 cursor-pointer">
                      <div className="font-medium">{repo.name}</div>
                      <div className="text-xs text-muted-foreground">{repo.owner}</div>
                    </Label>
                  </div>
                ))}
              </div>
            </DialogContent>
          </Dialog>
          <Button disabled={isLoading || isSyncing} onClick={handleSync} className="bg-sky-500 hover:bg-sky-600 text-white shadow-lg shadow-sky-500/20">
            {isSyncing ? (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            ) : (
              <RefreshCw className="mr-2 h-4 w-4" />
            )}
            {isSyncing ? 'Syncing...' : 'Sync Now'}
          </Button>
        </div>
      </div>

      {isLoading ? (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {[1, 2, 3, 4, 5, 6].map((i) => (
            <div key={i} className="h-[180px] rounded-xl border border-sky-100 dark:border-sky-900/50 bg-card/50 p-6 space-y-4 relative overflow-hidden">
              <div className="absolute inset-0 bg-linear-to-r from-transparent via-sky-500/5 to-transparent animate-shimmer" style={{ transform: 'translateX(-100%)' }} />
              <div className="h-6 w-1/2 bg-muted/50 rounded animate-pulse" />
              <div className="h-4 w-1/3 bg-muted/50 rounded animate-pulse" />
              <div className="mt-auto pt-4 h-8 w-full bg-muted/50 rounded animate-pulse" />
            </div>
          ))}
        </div>
      ) : (
        <RepositoryList repositories={filteredRepositories} />
      )}
    </motion.div>
  );
}
