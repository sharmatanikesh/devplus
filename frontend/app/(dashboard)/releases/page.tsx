'use client';

import { motion } from 'framer-motion';
import { useEffect, useState } from 'react';
import { apiClient } from '@/lib/api-client';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { cn } from '@/lib/utils';
import { Badge } from '@/components/ui/badge';
import { Checkbox } from '@/components/ui/checkbox';
import { Loader2, Package, AlertTriangle, CheckCircle2, Rocket, TrendingUp, RefreshCw, GitPullRequest } from 'lucide-react';
import { Repository, PullRequest } from '@/lib/types';
import ReactMarkdown from 'react-markdown';

export default function ReleasesPage() {
  const [repositories, setRepositories] = useState<Repository[]>([]);
  const [selectedRepo, setSelectedRepo] = useState<Repository | null>(null);
  const [pullRequests, setPullRequests] = useState<PullRequest[]>([]);
  const [selectedPRIds, setSelectedPRIds] = useState<Set<string>>(new Set());
  const [isLoading, setIsLoading] = useState(false);
  const [isLoadingPRs, setIsLoadingPRs] = useState(false);
  const [isSyncing, setIsSyncing] = useState(false);
  const [isCalculating, setIsCalculating] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Load repositories on mount
  useEffect(() => {
    loadRepositories();
  }, []);

  // Load PRs when repo is selected
  useEffect(() => {
    if (selectedRepo) {
      loadPullRequests(selectedRepo.id);
    } else {
      setPullRequests([]);
      setSelectedPRIds(new Set());
    }
  }, [selectedRepo?.id]);

  // Set up SSE for real-time updates when a repo is selected
  useEffect(() => {
    if (!selectedRepo?.id) return;

    const baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
    const sseUrl = baseUrl.endsWith('/api')
      ? `${baseUrl}/v1/repos/${selectedRepo.id}/analyze/stream`
      : `${baseUrl}/api/v1/repos/${selectedRepo.id}/analyze/stream`;

    console.log('[SSE] Connecting to:', sseUrl);

    const eventSource = new EventSource(sseUrl, { withCredentials: true });

    eventSource.onopen = () => {
      console.log('[SSE] Connection opened for repo:', selectedRepo.id);
    };

    eventSource.onmessage = (event) => {
      console.log('[SSE] Received message:', event.data);
      try {
        const data = JSON.parse(event.data);
        console.log('[SSE] Parsed data:', data);

        // Update repository with release risk data if present
        if (data.release_risk_score !== undefined || data.release_changelog) {
          console.log('[SSE] Updating release risk data:', {
            score: data.release_risk_score,
            hasChangelog: !!data.release_changelog,
            hasAnalysis: !!data.release_risk_analysis
          });

          // Update selected repo
          setSelectedRepo(prev => prev ? {
            ...prev,
            release_risk_score: data.release_risk_score ?? prev.release_risk_score,
            release_changelog: data.release_changelog ?? prev.release_changelog,
            release_risk_analysis: data.release_risk_analysis ?? prev.release_risk_analysis,
          } : null);

          // Also update in repositories list
          setRepositories(prev => prev.map(repo =>
            repo.id === selectedRepo.id ? {
              ...repo,
              release_risk_score: data.release_risk_score ?? repo.release_risk_score,
              release_changelog: data.release_changelog ?? repo.release_changelog,
              release_risk_analysis: data.release_risk_analysis ?? repo.release_risk_analysis,
            } : repo
          ));

          setIsCalculating(false);
        }
      } catch (err) {
        console.error('[SSE] Failed to parse SSE message:', err, 'Raw data:', event.data);
      }
    };

    eventSource.onerror = (error) => {
      console.error('[SSE] Connection error:', error);
      eventSource.close();
    };

    return () => {
      console.log('[SSE] Closing connection for repo:', selectedRepo.id);
      eventSource.close();
    };
  }, [selectedRepo?.id]);

  const loadRepositories = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await apiClient.repos.list();
      if (response.success && response.data) {
        setRepositories(response.data);
      }
    } catch (error) {
      console.error('Failed to load repositories:', error);
      setError('Failed to load repositories');
    } finally {
      setIsLoading(false);
    }
  };

  const loadPullRequests = async (repoId: string) => {
    setIsLoadingPRs(true);
    setError(null);
    try {
      const response = await apiClient.repos.getPullRequestsByRepoId(repoId);
      if (response.success && response.data) {
        setPullRequests(response.data);
      }
    } catch (error) {
      console.error('Failed to load pull requests:', error);
      setError('Failed to load pull requests');
    } finally {
      setIsLoadingPRs(false);
    }
  };

  const handleRepoSelect = (repoId: string) => {
    const repo = repositories.find(r => r.id === repoId);
    setSelectedRepo(repo || null);
    setSelectedPRIds(new Set());
    setError(null);
  };

  const handleSyncPRs = async () => {
    if (!selectedRepo) return;

    setIsSyncing(true);
    setError(null);
    try {
      const response = await apiClient.repos.syncPullRequests(selectedRepo.id);
      if (response.success) {
        // Reload PRs after sync
        await loadPullRequests(selectedRepo.id);
      } else {
        setError(response.error?.message || 'Failed to sync pull requests');
      }
    } catch (error) {
      console.error('Failed to sync pull requests:', error);
      setError('Failed to sync pull requests');
    } finally {
      setIsSyncing(false);
    }
  };

  const handlePRToggle = (prId: string) => {
    const newSelected = new Set(selectedPRIds);
    if (newSelected.has(prId)) {
      newSelected.delete(prId);
    } else {
      newSelected.add(prId);
    }
    setSelectedPRIds(newSelected);
  };

  const handleSelectAll = () => {
    if (selectedPRIds.size === pullRequests.length) {
      setSelectedPRIds(new Set());
    } else {
      setSelectedPRIds(new Set(pullRequests.map(pr => pr.id)));
    }
  };

  const handleNewRelease = async () => {
    if (!selectedRepo || selectedPRIds.size === 0) return;

    setIsCalculating(true);
    setError(null);
    try {
      const response = await apiClient.repos.calculateReleaseRisk(
        selectedRepo.id,
        Array.from(selectedPRIds)
      );
      if (response.success) {
        // SSE will handle updating the UI with results
      } else {
        setError(response.error?.message || 'Failed to create release');
        setIsCalculating(false);
      }
    } catch (error) {
      console.error('Failed to create release:', error);
      setError('Failed to create release');
      setIsCalculating(false);
    }
  };

  const getRiskBadge = (score?: number) => {
    if (score === undefined || score === 0) return null;

    if (score <= 30) {
      return (
        <Badge variant="outline" className="bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border-emerald-500/20">
          <CheckCircle2 className="w-3 h-3 mr-1" />
          Low Risk ({score}/100)
        </Badge>
      );
    } else if (score <= 60) {
      return (
        <Badge variant="outline" className="bg-amber-500/15 text-amber-700 dark:text-amber-400 border-amber-500/20">
          <AlertTriangle className="w-3 h-3 mr-1" />
          Medium Risk ({score}/100)
        </Badge>
      );
    } else if (score <= 85) {
      return (
        <Badge variant="outline" className="bg-orange-500/15 text-orange-700 dark:text-orange-400 border-orange-500/20">
          <AlertTriangle className="w-3 h-3 mr-1" />
          High Risk ({score}/100)
        </Badge>
      );
    } else {
      return (
        <Badge variant="destructive">
          <AlertTriangle className="w-3 h-3 mr-1" />
          Critical Risk ({score}/100)
        </Badge>
      );
    }
  };

  const getRiskColor = (score?: number) => {
    if (!score) return 'bg-muted';
    if (score <= 30) return 'bg-emerald-500';
    if (score <= 60) return 'bg-amber-500';
    if (score <= 85) return 'bg-orange-500';
    return 'bg-destructive';
  };

  const getStateBadge = (state: string) => {
    const variants = {
      open: 'bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border-emerald-500/20',
      closed: 'bg-muted text-muted-foreground border-border',
      merged: 'bg-purple-500/10 text-purple-700 dark:text-purple-400 border-purple-500/20',
    };
    return (
      <Badge variant="outline" className={variants[state as keyof typeof variants] || variants.closed}>
        {state}
      </Badge>
    );
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="space-y-6"
    >
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Release Management</h1>
          <p className="text-muted-foreground">
            Select pull requests and generate release changelog with risk analysis
          </p>
        </div>
      </div>

      {error && (
        <Card className="border-destructive">
          <CardContent className="pt-6">
            <p className="text-destructive">{error}</p>
          </CardContent>
        </Card>
      )}

      {/* Repository Selection */}
      <Card>
        <CardHeader>
          <CardTitle>Select Repository</CardTitle>
          <CardDescription>Choose a repository to manage releases</CardDescription>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="flex justify-center py-4">
              <Loader2 className="w-6 h-6 animate-spin" />
            </div>
          ) : (
            <Select value={selectedRepo?.id} onValueChange={handleRepoSelect}>
              <SelectTrigger>
                <SelectValue placeholder="Select a repository" />
              </SelectTrigger>
              <SelectContent>
                {repositories.map((repo) => (
                  <SelectItem key={repo.id} value={repo.id}>
                    <div className="flex items-center gap-2">
                      <Package className="w-4 h-4" />
                      {repo.owner}/{repo.name}
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
        </CardContent>
      </Card>

      {/* Last Release Summary */}
      {selectedRepo && (
        <Card className="bg-card border-border">
          <CardHeader className="pb-3">
            <CardTitle className="flex items-center gap-2 text-base">
              <RefreshCw className="w-4 h-4 text-primary" />
              Last Release Summary
            </CardTitle>
          </CardHeader>
          <CardContent>
            {(() => {
              // Try to parse the summary from the analysis JSON
              let summary = "";
              try {
                if (selectedRepo.release_risk_analysis) {
                  const analysis = JSON.parse(selectedRepo.release_risk_analysis);
                  summary = analysis.summary;
                }
              } catch (e) {
                // Fallback to changelog if parsing fails or valid JSON doesn't exist
                console.log("Failed to parse release analysis JSON", e);
              }

              const displayContent = summary || selectedRepo.release_changelog;

              return displayContent ? (
                <div className="prose prose-sm dark:prose-invert max-w-none text-muted-foreground">
                  <div className="line-clamp-3">
                    <ReactMarkdown>
                      {displayContent}
                    </ReactMarkdown>
                  </div>
                  <p className="text-xs text-muted-foreground mt-2">
                    Release Risk Score: <span className={cn("font-medium", getRiskColor(selectedRepo.release_risk_score).replace('bg-', 'text-'))}>{selectedRepo.release_risk_score ?? 'N/A'}/100</span>
                  </p>
                </div>
              ) : (
                <p className="text-sm text-muted-foreground">
                  No release summary available for this repository yet. Click &quot;New Release&quot; to generate one.
                </p>
              );
            })()}
          </CardContent>
        </Card>
      )}

      {/* Pull Requests List */}
      {selectedRepo && (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="flex items-center gap-2">
                  <GitPullRequest className="w-5 h-5" />
                  Pull Requests
                </CardTitle>
                <CardDescription>
                  Select PRs to include in the release
                </CardDescription>
              </div>
              <div className="flex gap-2">
                <Button
                  onClick={handleSyncPRs}
                  disabled={isSyncing}
                  variant="outline"
                >
                  {isSyncing ? (
                    <>
                      <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                      Syncing...
                    </>
                  ) : (
                    <>
                      <RefreshCw className="w-4 h-4 mr-2" />
                      Sync PRs
                    </>
                  )}
                </Button>
                <Button
                  onClick={handleNewRelease}
                  disabled={selectedPRIds.size === 0 || isCalculating}
                  className="bg-blue-600 hover:bg-blue-700 text-white border-blue-600 hover:border-blue-700 shadow-sm transition-all"
                >
                  {isCalculating ? (
                    <>
                      <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                      Analyzing...
                    </>
                  ) : (
                    <>
                      <Rocket className="w-4 h-4 mr-2" />
                      New Release ({selectedPRIds.size})
                    </>
                  )}
                </Button>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            {isLoadingPRs ? (
              <div className="flex justify-center py-8">
                <Loader2 className="w-8 h-8 animate-spin" />
              </div>
            ) : pullRequests.length === 0 ? (
              <div className="text-center py-8 space-y-2">
                <GitPullRequest className="w-12 h-12 mx-auto text-muted-foreground" />
                <p className="text-muted-foreground">No pull requests found</p>
                <p className="text-sm text-muted-foreground">
                  Click &quot;Sync PRs&quot; to fetch pull requests from GitHub
                </p>
              </div>
            ) : (
              <div className="space-y-4">
                <div className="flex items-center gap-2 pb-4 border-b">
                  <Checkbox
                    checked={selectedPRIds.size === pullRequests.length && pullRequests.length > 0}
                    onCheckedChange={handleSelectAll}
                    id="select-all"
                  />
                  <label htmlFor="select-all" className="text-sm font-medium cursor-pointer select-none">
                    Select All ({selectedPRIds.size}/{pullRequests.length})
                  </label>
                </div>
                <div className="flex flex-col gap-3">
                  {pullRequests.map((pr) => (
                    <div
                      key={pr.id}
                      className={`
                        group relative flex flex-col gap-2 p-4 rounded-lg border transition-all duration-200
                        ${selectedPRIds.has(pr.id)
                          ? 'bg-blue-50/50 dark:bg-blue-900/10 border-blue-200 dark:border-blue-800'
                          : 'bg-card hover:bg-muted/50 border-border hover:border-slate-300 dark:hover:border-slate-700'
                        }
                      `}
                    >
                      <div className="flex items-start gap-3">
                        <Checkbox
                          checked={selectedPRIds.has(pr.id)}
                          onCheckedChange={() => handlePRToggle(pr.id)}
                          className="mt-1"
                        />
                        <div className="flex-1 min-w-0">
                          {/* Row 1: Header Info */}
                          <div className="flex items-center justify-between gap-4 mb-1">
                            <div className="flex items-center gap-2 min-w-0">
                              <span className="font-mono text-xs text-muted-foreground bg-muted px-1.5 py-0.5 rounded">
                                #{pr.number}
                              </span>
                              <h3 className="font-semibold text-base truncate" title={pr.title}>
                                {pr.title}
                              </h3>
                              {getStateBadge(pr.state)}
                            </div>
                            <div className="flex items-center gap-2 text-xs text-muted-foreground whitespace-nowrap">
                              <span>by</span>
                              <span className="font-medium text-foreground">{pr.author_name || pr.author?.username || 'Unknown'}</span>
                            </div>
                          </div>

                          {/* Row 2: Summary */}
                          <p className="text-sm text-muted-foreground line-clamp-2 leading-relaxed">
                            {pr.ai_summary || pr.description || 'No summary available.'}
                          </p>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {/* Release Risk Display */}
      {selectedRepo && (selectedRepo.release_risk_score !== undefined && selectedRepo.release_risk_score > 0) && (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="flex items-center gap-2">
                  <TrendingUp className="w-5 h-5" />
                  Latest Release Risk Analysis
                </CardTitle>
                <CardDescription>
                  Results from the most recent release analysis
                </CardDescription>
              </div>
              {getRiskBadge(selectedRepo.release_risk_score)}
            </div>
          </CardHeader>
          <CardContent>
            {isCalculating ? (
              <div className="flex flex-col items-center justify-center py-12 space-y-4">
                <Loader2 className="w-12 h-12 animate-spin text-primary" />
                <p className="text-lg font-medium">Analyzing selected pull requests...</p>
                <p className="text-sm text-muted-foreground">
                  This may take a minute while our AI analyzes the PRs
                </p>
              </div>
            ) : (
              <div className="space-y-6">
                {/* Risk Score Visualization */}
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm font-medium">Risk Level</span>
                    <span className="text-2xl font-bold">
                      {selectedRepo.release_risk_score}/100
                    </span>
                  </div>
                  <div className="w-full bg-secondary rounded-full h-3 overflow-hidden">
                    <motion.div
                      initial={{ width: 0 }}
                      animate={{ width: `${selectedRepo.release_risk_score}%` }}
                      transition={{ duration: 1, ease: "easeOut" }}
                      className={`h-3 ${getRiskColor(selectedRepo.release_risk_score)} transition-colors`}
                    />
                  </div>
                </div>

                {/* Changelog */}
                {selectedRepo.release_changelog && (
                  <div>
                    <h3 className="text-lg font-semibold mb-3">Generated Changelog</h3>
                    <div className="prose prose-sm dark:prose-invert max-w-none bg-muted/50 rounded-lg p-4 border">
                      <ReactMarkdown>{selectedRepo.release_changelog}</ReactMarkdown>
                    </div>
                  </div>
                )}


              </div>
            )}
          </CardContent>
        </Card>
      )}

      {/* Info Card */}
      <Card className="bg-muted/30 border-primary/20">
        <CardContent className="pt-6">
          <div className="flex gap-3">
            <div className="p-2 bg-primary/10 rounded-full h-fit">
              <AlertTriangle className="w-5 h-5 text-primary flex-shrink-0" />
            </div>
            <div className="space-y-1">
              <p className="text-sm font-medium text-foreground">
                How Release Management Works
              </p>
              <ul className="text-sm text-muted-foreground space-y-1 list-disc list-inside">
                <li>Select a repository to view its pull requests</li>
                <li>Choose specific PRs to include in the release</li>
                <li>Click &quot;New Release&quot; to generate changelog and risk analysis</li>
                <li>AI analyzes selected PRs for breaking changes and complexity</li>
                <li>Generates a risk score from 0-100 and markdown changelog</li>
                <li>Results update in real-time via server-sent events</li>
              </ul>
            </div>
          </div>
        </CardContent>
      </Card>
    </motion.div>
  );
}
