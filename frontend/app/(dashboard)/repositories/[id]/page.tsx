'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import { apiClient } from '@/lib/api-client';
import { Repository, PullRequest } from '@/lib/types';
import { Button } from '@/components/ui/button';
import { ArrowLeft, GitFork, Star, Globe, Loader2, GitPullRequest, Sparkles, RefreshCw } from 'lucide-react';
import Link from 'next/link';
import { Badge } from '@/components/ui/badge';

export default function RepositoryDetailsPage() {
    const params = useParams();
    const id = params?.id as string;

    const [repo, setRepo] = useState<Repository | null>(null);
    const [pullRequests, setPullRequests] = useState<PullRequest[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const [analyzing, setAnalyzing] = useState(false);
    const [isSyncingPRs, setIsSyncingPRs] = useState(false);

    useEffect(() => {
        async function fetchData() {
            if (!id) return;
            try {
                setLoading(true);
                // 1. Fetch Repo
                // Use repos.get as defined in api-client.ts
                const repoRes = await apiClient.repos.get(id);
                const repoData = repoRes.data as Repository;
                setRepo(repoData);

                // 2. Fetch PRs
                if (repoData && repoData.owner && repoData.name) {
                    const prsRes = await apiClient.repos.getPullRequests(repoData.owner, repoData.name);
                    setPullRequests(prsRes.data as PullRequest[] || []);
                }

            } catch (err: any) {
                console.error('Failed to fetch details:', err);
                setError('Failed to load repository details.');
            } finally {
                setLoading(false);
            }
        }
        fetchData();
    }, [id]);

    const handleAnalyze = async () => {
        if (!id) return;
        try {
            setAnalyzing(true);
            await apiClient.repos.analyze(id);
            // Optionally show toast
            alert("Analysis triggered! It will appear here shortly.");
        } catch (error) {
            console.error("Failed to trigger analysis:", error);
            alert("Failed to start analysis");
        } finally {
            setAnalyzing(false);
        }
    };

    if (loading) {
        return (
            <div className="flex h-[400px] items-center justify-center">
                <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
        )
    }

    if (error || !repo) {
        return (
            <div className="flex flex-col items-center justify-center min-h-[400px] space-y-4">
                <p className="text-red-500">{error || 'Repository not found'}</p>
                <Button variant="outline" asChild>
                    <Link href="/repositories">Back to Repositories</Link>
                </Button>
            </div>
        )
    }

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="space-y-2">
                <Link href="/repositories" className="text-sm text-muted-foreground hover:text-foreground flex items-center gap-1 w-fit">
                    <ArrowLeft className="h-4 w-4" /> Back
                </Link>
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-3xl font-bold">{repo.fullName || repo.name}</h1>
                        {repo.description && (
                            <p className="text-muted-foreground mt-1 max-w-2xl">{repo.description}</p>
                        )}
                    </div>
                    <div className="flex gap-2">
                        <Button
                            variant="default"
                            onClick={handleAnalyze}
                            disabled={analyzing}
                            className="bg-purple-600 hover:bg-purple-700 text-white"
                        >
                            {analyzing ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    Analyzing...
                                </>
                            ) : (
                                <>
                                    <Sparkles className="mr-2 h-4 w-4" />
                                    Analyze Codebase
                                </>
                            )}
                        </Button>
                        <Button variant="outline" onClick={() => window.open(repo.url, '_blank')}>
                            <Globe className="mr-2 h-4 w-4" /> View on GitHub
                        </Button>
                    </div>
                </div>

                {/* Stats */}
                <div className="flex items-center gap-4 text-sm text-muted-foreground pt-2">
                    {repo.language && (
                        <Badge variant="secondary" className="font-normal text-sm">
                            ● {repo.language}
                        </Badge>
                    )}
                    <span className="flex items-center gap-1"><Star className="h-4 w-4 text-yellow-500" /> {repo.stars} stars</span>
                    <span className="flex items-center gap-1"><GitFork className="h-4 w-4" /> {repo.forks} forks</span>
                </div>
            </div>

            {/* AI Summary Section */}
            {repo.aiSummary && (
                <div className="bg-muted/30 rounded-lg p-6 border border-purple-200 dark:border-purple-900/50 relative overflow-hidden">
                    <div className="absolute top-0 right-0 p-4 opacity-10">
                        <Sparkles className="h-24 w-24 text-purple-600" />
                    </div>
                    <h3 className="text-lg font-semibold mb-3 text-purple-700 dark:text-purple-400 flex items-center gap-2">
                        <Sparkles className="h-5 w-5" />
                        AI Repository Analysis
                    </h3>
                    <div className="prose dark:prose-invert max-w-none text-sm text-muted-foreground whitespace-pre-wrap">
                        {repo.aiSummary}
                    </div>
                </div>
            )}

            {/* PR Section */}
            <div className="space-y-4 pt-4">
                <div className="flex items-center justify-between">
                    <h2 className="text-xl font-semibold flex items-center gap-2">
                        <GitPullRequest className="h-5 w-5" /> Pull Requests
                    </h2>
                    <div className="flex items-center gap-2">
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={async () => {
                                setIsSyncingPRs(true);
                                try {
                                    if (apiClient.repos.syncOne) {
                                        await apiClient.repos.syncOne(id);
                                    } else {
                                        await (apiClient.repos as any).syncOne(id);
                                    }
                                    // Refresh PRs
                                    const prsRes = await apiClient.repos.getPullRequests(repo.owner, repo.name);
                                    setPullRequests(prsRes.data as PullRequest[] || []);
                                } catch (error) {
                                    console.error("Failed to sync PRs:", error);
                                } finally {
                                    setIsSyncingPRs(false);
                                }
                            }}
                            disabled={isSyncingPRs}
                        >
                            <RefreshCw className={`mr-2 h-4 w-4 ${isSyncingPRs ? 'animate-spin' : ''}`} />
                            {isSyncingPRs ? 'Syncing...' : 'Sync PRs'}
                        </Button>
                        <Badge variant="outline">{pullRequests.length}</Badge>
                    </div>
                </div>

                {pullRequests.length === 0 ? (
                    <div className="rounded-lg border border-dashed p-8 text-center text-muted-foreground bg-muted/20">
                        No pull requests synced yet.
                        <div className="mt-2 text-xs">
                            Sync data to see Pull Requests.
                        </div>
                    </div>
                ) : (
                    <div className="rounded-md border divide-y">
                        {pullRequests.map(pr => (
                            <Link key={pr.id} href={`/pull-requests/${pr.repositoryId}/${pr.number}`} className="block group">
                                <div className="p-4 flex justify-between items-start hover:bg-muted/50 transition-colors">
                                    <div className="space-y-1">
                                        <div className="font-medium group-hover:text-primary transition-colors">
                                            #{pr.number} {pr.title}
                                        </div>
                                        <div className="text-sm text-muted-foreground flex items-center gap-2">
                                            <span className={`inline-block w-2 h-2 rounded-full ${pr.state === 'open' ? 'bg-green-500' : pr.state === 'merged' ? 'bg-purple-500' : 'bg-red-500'}`} />
                                            {pr.state} • by {pr.author?.username || 'unknown'}
                                        </div>
                                    </div>
                                    <div className="text-xs text-muted-foreground">
                                        {new Date(pr.createdAt).toLocaleDateString()}
                                    </div>
                                </div>
                            </Link>
                        ))}
                    </div>
                )}
            </div>

        </div>
    )
}
