'use client';

import { useEffect, useState, useRef } from 'react';
import { useParams } from 'next/navigation';
import ReactMarkdown from 'react-markdown';
import { apiClient } from '@/lib/api-client';
import { PullRequest } from '@/lib/types';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { ArrowLeft, Bot, CheckCircle, AlertTriangle, AlertOctagon, RefreshCw, GitPullRequest } from 'lucide-react';
import Link from 'next/link';
import { PR_STATUS_COLORS } from '@/lib/constants';
import { toast } from 'sonner';
import { motion } from 'framer-motion';

const POLL_INTERVAL = 5000; // 5 seconds
const MAX_POLL_DURATION = 120000; // 2 minutes

export default function PullRequestDetailsPage() {
    const params = useParams();
    const repositoryId = params?.repositoryId as string;
    const number = params?.number as string;

    const [pr, setPr] = useState<PullRequest | null>(null);
    const [loading, setLoading] = useState(true);
    const [analyzing, setAnalyzing] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const pollIntervalRef = useRef<NodeJS.Timeout | null>(null);
    const pollStartTimeRef = useRef<number | null>(null);
    const previousSummaryRef = useRef<string | undefined>(undefined);

    async function fetchPR() {
        if (!repositoryId || !number) return;
        try {
            setLoading(true);
            const res = await apiClient.pullRequests.get(repositoryId, parseInt(number));
            const newPr = res.data as PullRequest;
            setPr(newPr);

            // Check if analysis just completed
            if (analyzing && newPr.ai_summary && newPr.ai_summary !== previousSummaryRef.current) {
                toast.success('AI Analysis Complete!', {
                    description: 'The code review is ready to view.',
                });
                stopPolling();
            }

            previousSummaryRef.current = newPr.ai_summary;
        } catch (err: any) {
            console.error('Failed to fetch PR:', err);
            setError('Failed to load Pull Request details.');
        } finally {
            setLoading(false);
        }
    }

    const startPolling = () => {
        if (pollIntervalRef.current) return; // Already polling

        pollStartTimeRef.current = Date.now();

        pollIntervalRef.current = setInterval(() => {
            const elapsed = Date.now() - (pollStartTimeRef.current || 0);

            if (elapsed >= MAX_POLL_DURATION) {
                toast.warning('Analysis Taking Longer Than Expected', {
                    description: 'The analysis is still processing. Please check back later.',
                });
                stopPolling();
                return;
            }

            fetchPR();
        }, POLL_INTERVAL);
    };

    const stopPolling = () => {
        if (pollIntervalRef.current) {
            clearInterval(pollIntervalRef.current);
            pollIntervalRef.current = null;
            pollStartTimeRef.current = null;
            setAnalyzing(false);
        }
    };

    useEffect(() => {
        fetchPR();

        // Cleanup on unmount
        return () => {
            stopPolling();
        };
    }, [repositoryId, number]);

    // Reset analyzing state when PR data loads
    useEffect(() => {
        if (pr && !pollIntervalRef.current) {
            // If we have PR data and we're not actively polling, stop analyzing state
            setAnalyzing(false);
        }
    }, [pr]);

    const handleReanalyze = async () => {
        if (!pr) return;
        try {
            setAnalyzing(true);
            previousSummaryRef.current = pr.ai_summary; // Store current summary

            await apiClient.pullRequests.analyze(repositoryId, parseInt(number));

            toast.info('Analysis Started', {
                description: 'AI is reviewing your code. This may take a minute...',
            });

            // Use EventSource for real-time updates via SSE
            const baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
            const streamUrl = baseUrl.endsWith('/api')
                ? `${baseUrl}/v1/repos/${repositoryId}/prs/${number}/analyze/stream`
                : `${baseUrl}/api/v1/repos/${repositoryId}/prs/${number}/analyze/stream`;
            const eventSource = new EventSource(streamUrl, {
                withCredentials: true
            });

            eventSource.onmessage = (event) => {
                const data = JSON.parse(event.data);

                if (data.status === 'completed' && data.ai_summary) {
                    // Update PR with the new summary
                    setPr(prev => prev ? {
                        ...prev,
                        ai_summary: data.ai_summary,
                        ai_decision: data.ai_decision
                    } : null);
                    setAnalyzing(false);
                    eventSource.close();
                    stopPolling();
                    toast.success('AI Analysis Complete!', {
                        description: 'The code review is ready to view.',
                    });
                } else if (data.status === 'error') {
                    toast.error('Analysis Failed', {
                        description: data.message || 'Please try again later.',
                    });
                    setAnalyzing(false);
                    eventSource.close();
                    stopPolling();
                }
            };

            eventSource.onerror = (error) => {
                console.error('SSE error:', error);
                // Fallback to polling if SSE fails
                toast.info('Using fallback polling', {
                    description: 'Checking for updates every few seconds...',
                });
                startPolling();
                eventSource.close();
            };

        } catch (err) {
            console.error('Failed to trigger analysis:', err);
            toast.error('Failed to Start Analysis', {
                description: 'Please try again later.',
            });
            setAnalyzing(false);
        }
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
                    <p className="text-muted-foreground animate-pulse">Loading pull request...</p>
                </div>
            </div>
        )
    }

    if (error || !pr) {
        return (
            <div className="flex flex-col items-center justify-center min-h-[400px] space-y-4">
                <p className="text-red-500 font-medium">{error || 'Pull Request not found'}</p>
                <Button variant="outline" asChild className="hover:bg-sky-50 dark:hover:bg-sky-900/20">
                    <Link href={`/repositories/${repositoryId}`}>Back to Repository</Link>
                </Button>
            </div>
        )
    }

    const getDecisionBadge = (decision?: string) => {
        switch (decision?.toLowerCase()) {
            case 'approve':
            case 'approved':
                return <Badge className="bg-green-500 hover:bg-green-600 shadow-md shadow-green-500/20"><CheckCircle className="w-3 h-3 mr-1" /> Approved</Badge>;
            case 'request_changes':
                return <Badge className="bg-orange-500 hover:bg-orange-600 shadow-md shadow-orange-500/20"><AlertTriangle className="w-3 h-3 mr-1" /> Changes Requested</Badge>;
            case 'block':
            case 'blocked':
                return <Badge variant="destructive" className="shadow-md shadow-red-500/20"><AlertOctagon className="w-3 h-3 mr-1" /> Blocked</Badge>;
            default:
                return <Badge variant="secondary" className="bg-muted text-muted-foreground">Pending Review</Badge>;
        }
    };

    const hasExistingAnalysis = pr.ai_summary && pr.ai_summary.trim() !== '';

    return (
        <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="space-y-6 max-w-5xl mx-auto pb-10"
        >
            {/* Header */}
            <div className="space-y-4">
                <Link href={`/repositories/${repositoryId}`} className="text-sm text-muted-foreground hover:text-sky-500 flex items-center gap-1 w-fit transition-colors">
                    <ArrowLeft className="h-4 w-4" /> Back to Repository
                </Link>

                <div className="flex flex-col md:flex-row md:items-start justify-between gap-4">
                    <div className="space-y-1">
                        <div className="flex items-center gap-3">
                            <h1 className="text-2xl font-bold tracking-tight">#{pr.number} {pr.title}</h1>
                            <Badge variant="outline" className={PR_STATUS_COLORS[pr.state]}>
                                {pr.state}
                            </Badge>
                        </div>
                        <p className="text-muted-foreground">
                            Opened by <span className="font-medium text-foreground">{pr.author?.username || pr.author_name || 'unknown'}</span> on {new Date(pr.createdAt || pr.created_at || new Date()).toLocaleDateString()}
                        </p>
                    </div>

                    <div className="flex gap-2">
                        <Button variant="outline" onClick={() => window.open(pr.url, '_blank')} className="hover:bg-sky-50 dark:hover:bg-sky-900/20 border-sky-200 dark:border-sky-800">
                            View on GitHub
                        </Button>
                        <Button
                            onClick={handleReanalyze}
                            disabled={analyzing}
                            className="bg-sky-500 hover:bg-sky-600 text-white shadow-lg shadow-sky-500/20"
                        >
                            <RefreshCw className={`mr-2 h-4 w-4 ${analyzing ? 'animate-spin' : ''}`} />
                            {analyzing ? 'Analyzing...' : hasExistingAnalysis ? 'Re-analyze' : 'Analyze'}
                        </Button>
                    </div>
                </div>
            </div>

            <div className="grid gap-6 md:grid-cols-3">
                {/* Main Content - AI Analysis */}
                <div className="md:col-span-2 space-y-6">
                    <Card className="border-none shadow-lg bg-card/40 backdrop-blur-md ring-1 ring-border/50">
                        <CardHeader className="bg-sky-50/50 dark:bg-sky-900/10 pb-4 border-b border-sky-100 dark:border-sky-800/50">
                            <div className="flex items-center justify-between">
                                <div className="flex items-center gap-2">
                                    <Bot className={`h-5 w-5 text-sky-500 ${analyzing ? 'animate-pulse' : ''}`} />
                                    <CardTitle>AI Review Summary</CardTitle>
                                </div>
                                {getDecisionBadge(pr.ai_decision)}
                            </div>
                        </CardHeader>
                        <CardContent className="pt-6 relative">
                            {analyzing ? (
                                <div className="rounded-lg p-8 relative overflow-hidden flex flex-col items-center justify-center min-h-[200px]">
                                    <div className="absolute top-0 right-0 p-4 opacity-10">
                                        <Bot className="h-20 w-20 text-sky-600 animate-pulse" />
                                    </div>
                                    <div className="relative z-10 text-center">
                                        <div className="flex items-center justify-center gap-3 mb-3">
                                            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-sky-500" />
                                        </div>
                                        <p className="text-lg font-medium">Analyzing code...</p>
                                        <p className="text-sm text-muted-foreground mt-2">
                                            Our AI is reviewing your pull request.<br />This may take a minute or two.
                                        </p>
                                    </div>
                                </div>
                            ) : pr.ai_summary ? (
                                <div className="text-sm leading-relaxed prose dark:prose-invert max-w-none">
                                    <ReactMarkdown>{pr.ai_summary}</ReactMarkdown>
                                </div>
                            ) : (
                                <div className="text-center py-12 bg-muted/20 rounded-lg border-2 border-dashed border-sky-200 dark:border-sky-800">
                                    <Bot className="h-12 w-12 mx-auto mb-4 text-muted-foreground opacity-40" />
                                    <p className="font-medium">No AI analysis available yet</p>
                                    <p className="text-sm text-muted-foreground mt-1">Click &quot;Analyze&quot; to start a new review</p>
                                </div>
                            )}
                        </CardContent>
                    </Card>

                    {/* Add more tabs or sections for File Changes later */}
                </div>

                {/* Sidebar - Meta Info */}
                <div className="space-y-6">
                    <Card className="border-none shadow-md bg-card/40 backdrop-blur-sm ring-1 ring-border/50">
                        <CardHeader>
                            <CardTitle className="text-sm font-medium">Details</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div className="flex justify-between text-sm">
                                <span className="text-muted-foreground">Repository ID</span>
                                <span className="font-mono text-xs text-sky-600 dark:text-sky-400">{repositoryId.substring(0, 8)}...</span>
                            </div>
                            <div className="flex justify-between text-sm">
                                <span className="text-muted-foreground">Branch</span>
                                <span className="font-mono text-xs flex items-center gap-1"><GitPullRequest className="h-3 w-3" /> main</span>
                            </div>
                            {/* Add more meta info if available */}
                        </CardContent>
                    </Card>
                </div>
            </div>
        </motion.div>
    )
}
