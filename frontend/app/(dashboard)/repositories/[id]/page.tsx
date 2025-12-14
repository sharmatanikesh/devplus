'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import ReactMarkdown from 'react-markdown';
import { apiClient } from '@/lib/api-client';
import { Repository, PullRequest } from '@/lib/types';
import { Button } from '@/components/ui/button';
import { ArrowLeft, GitFork, Star, Globe, Loader2, GitPullRequest, Sparkles, RefreshCw } from 'lucide-react';
import Link from 'next/link';
import { Badge } from '@/components/ui/badge';
import { motion } from 'framer-motion';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { cn } from "@/lib/utils";

export default function RepositoryDetailsPage() {
    const params = useParams();
    const id = params?.id as string;

    const [repo, setRepo] = useState<Repository | null>(null);
    const [pullRequests, setPullRequests] = useState<PullRequest[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const [analyzing, setAnalyzing] = useState(false);
    const [isSyncingPRs, setIsSyncingPRs] = useState(false);
    const [pollingForAnalysis, setPollingForAnalysis] = useState(false);

    useEffect(() => {
        async function fetchData() {
            if (!id) return;
            try {
                setLoading(true);
                const repoRes = await apiClient.repos.get(id);
                const repoData = repoRes.data as Repository;
                setRepo(repoData);

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
            setPollingForAnalysis(true);

            await apiClient.repos.analyze(id);

            const baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
            const streamUrl = baseUrl.endsWith('/api')
                ? `${baseUrl}/v1/repos/${id}/analyze/stream`
                : `${baseUrl}/api/v1/repos/${id}/analyze/stream`;
            const eventSource = new EventSource(streamUrl, {
                withCredentials: true
            });

            eventSource.onmessage = (event) => {
                const data = JSON.parse(event.data);

                if (data.status === 'completed' && data.ai_summary) {
                    setRepo(prev => prev ? { ...prev, ai_summary: data.ai_summary } : null);
                    setAnalyzing(false);
                    setPollingForAnalysis(false);
                    eventSource.close();
                } else if (data.status === 'error') {
                    alert('Analysis failed: ' + data.message);
                    setAnalyzing(false);
                    setPollingForAnalysis(false);
                    eventSource.close();
                }
            };

            eventSource.onerror = (error) => {
                console.error('SSE error:', error);
                startPolling();
                eventSource.close();
            };

            setTimeout(() => {
                if (pollingForAnalysis) {
                    eventSource.close();
                    setAnalyzing(false);
                    setPollingForAnalysis(false);
                }
            }, 120000);

        } catch (error) {
            console.error("Failed to trigger analysis:", error);
            alert("Failed to start analysis");
            setAnalyzing(false);
            setPollingForAnalysis(false);
        }
    };

    const startPolling = () => {
        const pollInterval = setInterval(async () => {
            try {
                const repoRes = await apiClient.repos.get(id);
                const updatedRepo = repoRes.data as Repository;

                if (updatedRepo.ai_summary && updatedRepo.ai_summary !== repo?.ai_summary) {
                    setRepo(updatedRepo);
                    setAnalyzing(false);
                    setPollingForAnalysis(false);
                    clearInterval(pollInterval);
                }
            } catch (error) {
                console.error('Polling error:', error);
            }
        }, 3000);

        setTimeout(() => {
            clearInterval(pollInterval);
        }, 120000);
    };

    if (loading) {
        return (
            <div className="flex items-center justify-center min-h-[400px]">
                <div className="flex flex-col items-center gap-4">
                    <div className="relative flex items-center justify-center">
                        <div className="absolute inset-0 bg-sky-500/20 rounded-full blur-xl animate-pulse" />
                        <div className="h-10 w-10 text-sky-500 animate-spin">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 12a9 9 0 1 1-6.219-8.56" /></svg>
                        </div>
                    </div>
                    <p className="text-muted-foreground animate-pulse">Loading repository...</p>
                </div>
            </div>
        )
    }

    if (error || !repo) {
        return (
            <div className="flex flex-col items-center justify-center min-h-[400px] space-y-4">
                <p className="text-red-500 font-medium">{error || 'Repository not found'}</p>
                <Button variant="outline" asChild className="hover:bg-sky-50 dark:hover:bg-sky-900/20">
                    <Link href="/repositories">Back to Repositories</Link>
                </Button>
            </div>
        )
    }

    return (
        <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="space-y-6"
        >
            {/* Header */}
            <div className="space-y-4">
                <Link href="/repositories" className="text-sm text-muted-foreground hover:text-sky-500 flex items-center gap-1 w-fit transition-colors">
                    <ArrowLeft className="h-4 w-4" /> Back
                </Link>
                <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                    <div>
                        <h1 className="text-3xl font-bold tracking-tight text-sky-600 dark:text-sky-400">{repo.fullName || repo.name}</h1>
                        {repo.description && (
                            <p className="text-muted-foreground mt-1 max-w-2xl text-lg">{repo.description}</p>
                        )}
                    </div>
                    <div className="flex gap-2">
                        <Button
                            variant="default"
                            onClick={handleAnalyze}
                            disabled={analyzing}
                            className="bg-sky-500 hover:bg-sky-600 text-white shadow-lg shadow-sky-500/20"
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
                        <Button variant="outline" onClick={() => window.open(repo.url, '_blank')} className="hover:bg-sky-50 dark:hover:bg-sky-900/20 border-sky-200 dark:border-sky-800">
                            <Globe className="mr-2 h-4 w-4" /> View on GitHub
                        </Button>
                    </div>
                </div>

                {/* Stats */}
                <div className="flex items-center gap-4 text-sm text-muted-foreground pt-2">
                    {repo.language && (
                        <Badge variant="secondary" className="font-normal text-sm bg-sky-100 text-sky-700 dark:bg-sky-900/30 dark:text-sky-300">
                            ● {repo.language}
                        </Badge>
                    )}
                    <span className="flex items-center gap-1"><Star className="h-4 w-4 text-amber-400" /> {repo.stars} stars</span>
                    <span className="flex items-center gap-1 text-muted-foreground"><GitFork className="h-4 w-4" /> {repo.forks} forks</span>
                </div>
            </div>

            {/* AI Summary Section */}
            {analyzing && (
                <Card className="border-sky-200 dark:border-sky-800 bg-sky-50/50 dark:bg-sky-900/10 overflow-hidden relative">
                    <CardContent className="p-6 relative z-10">
                        <div className="absolute top-0 right-0 p-4 opacity-10">
                            <Sparkles className="h-24 w-24 text-sky-500 animate-pulse" />
                        </div>
                        <h3 className="text-lg font-semibold mb-3 text-sky-700 dark:text-sky-400 flex items-center gap-2">
                            <Loader2 className="h-5 w-5 animate-spin" />
                            Analyzing Repository...
                        </h3>
                        <div className="text-sm text-muted-foreground">
                            Our AI is analyzing your codebase architecture, patterns, and best practices. This may take a minute or two.
                        </div>
                    </CardContent>
                </Card>
            )}

            {repo.ai_summary && !analyzing && (
                <motion.div
                    initial={{ opacity: 0, scale: 0.98 }}
                    animate={{ opacity: 1, scale: 1 }}
                >
                    <Card className="border-none shadow-lg bg-card/40 backdrop-blur-md ring-1 ring-border/50 overflow-hidden relative group">
                        <div className="absolute inset-0 bg-linear-to-r from-sky-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500" />
                        <CardContent className="p-6 relative z-10">
                            <div className="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
                                <Sparkles className="h-24 w-24 text-sky-500" />
                            </div>
                            <h3 className="text-lg font-semibold mb-3 text-sky-600 dark:text-sky-400 flex items-center gap-2">
                                <Sparkles className="h-5 w-5" />
                                AI Repository Analysis
                            </h3>
                            <div className="text-sm leading-relaxed text-muted-foreground prose dark:prose-invert max-w-none">
                                <ReactMarkdown>{repo.ai_summary}</ReactMarkdown>
                            </div>
                        </CardContent>
                    </Card>
                </motion.div>
            )}

            {/* PR Section */}
            <div className="space-y-4 pt-4">
                <div className="flex items-center justify-between">
                    <h2 className="text-xl font-semibold flex items-center gap-2">
                        <GitPullRequest className="h-5 w-5 text-sky-500" /> Pull Requests
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
                                    const prsRes = await apiClient.repos.getPullRequests(repo.owner, repo.name);
                                    setPullRequests(prsRes.data as PullRequest[] || []);
                                } catch (error) {
                                    console.error("Failed to sync PRs:", error);
                                } finally {
                                    setIsSyncingPRs(false);
                                }
                            }}
                            disabled={isSyncingPRs}
                            className="hover:bg-sky-50 dark:hover:bg-sky-900/20"
                        >
                            <RefreshCw className={`mr-2 h-4 w-4 ${isSyncingPRs ? 'animate-spin' : ''}`} />
                            {isSyncingPRs ? 'Syncing...' : 'Sync PRs'}
                        </Button>
                        <Badge variant="secondary" className="bg-muted text-muted-foreground">{pullRequests.length}</Badge>
                    </div>
                </div>

                {pullRequests.length === 0 ? (
                    <div className="rounded-xl border border-dashed border-sky-200 dark:border-sky-800 p-8 text-center text-muted-foreground bg-sky-50/30 dark:bg-sky-900/10">
                        No pull requests synced yet.
                        <div className="mt-2 text-xs">
                            Sync data to see Pull Requests.
                        </div>
                    </div>
                ) : (
                    <div className="grid gap-3">
                        {pullRequests.map((pr, i) => (
                            <motion.div
                                initial={{ opacity: 0, y: 10 }}
                                animate={{ opacity: 1, y: 0 }}
                                transition={{ delay: i * 0.05 }}
                                key={pr.id}
                            >
                                <Link href={`/pull-requests/${pr.repo_id}/${pr.number}`} className="block group">
                                    <div className="p-4 rounded-xl border border-border/50 bg-card/40 hover:bg-card/80 transition-all duration-200 hover:shadow-md hover:scale-[1.01] hover:border-sky-200 dark:hover:border-sky-800 flex justify-between items-start">
                                        <div className="space-y-1">
                                            <div className="font-medium group-hover:text-sky-600 transition-colors flex items-center gap-2">
                                                #{pr.number} {pr.title}
                                            </div>
                                            <div className="text-sm text-muted-foreground flex items-center gap-2">
                                                <span className={cn("inline-block w-2 h-2 rounded-full",
                                                    pr.state === 'open' ? 'bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.4)]' :
                                                        pr.state === 'merged' ? 'bg-purple-500 shadow-[0_0_8px_rgba(168,85,247,0.4)]' :
                                                            'bg-red-500'
                                                )} />
                                                {pr.state} • by <span className="text-foreground font-medium">{pr.author?.username || 'unknown'}</span>
                                            </div>
                                        </div>
                                        <div className="text-xs text-muted-foreground font-mono bg-muted/50 px-2 py-1 rounded">
                                            {new Date(pr.createdAt).toLocaleDateString()}
                                        </div>
                                    </div>
                                </Link>
                            </motion.div>
                        ))}
                    </div>
                )}
            </div>
        </motion.div>
    )
}
