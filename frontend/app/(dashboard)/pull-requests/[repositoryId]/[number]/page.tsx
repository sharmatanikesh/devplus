'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import { apiClient } from '@/lib/api-client';
import { PullRequest } from '@/lib/types';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { ArrowLeft, Bot, CheckCircle, AlertTriangle, AlertOctagon, RefreshCw } from 'lucide-react';
import Link from 'next/link';
import { PR_STATUS_COLORS } from '@/lib/constants';

export default function PullRequestDetailsPage() {
    const params = useParams();
    const repositoryId = params?.repositoryId as string;
    const number = params?.number as string;

    const [pr, setPr] = useState<PullRequest | null>(null);
    const [loading, setLoading] = useState(true);
    const [analyzing, setAnalyzing] = useState(false);
    const [error, setError] = useState<string | null>(null);

    async function fetchPR() {
        if (!repositoryId || !number) return;
        try {
            setLoading(true);
            const res = await apiClient.pullRequests.get(repositoryId, parseInt(number));
            setPr(res.data as PullRequest);
        } catch (err: any) {
            console.error('Failed to fetch PR:', err);
            setError('Failed to load Pull Request details.');
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        fetchPR();
    }, [repositoryId, number]);

    const handleReanalyze = async () => {
        if (!pr) return;
        try {
            setAnalyzing(true);
            // We need to implement the analyze endpoint call in apiClient or call via raw axios
            // Since it's not in apiClient explicitly as a method (checked previously, it was in constants but typically exposed)
            // Let's check api-client.ts again to be sure specific method exists or use generic post
            // The constants showed API_ENDPOINTS.PR_ANALYZE.
            // Let's use the generic client for now if specific one isn't there.
            // Actually apiClient.pullRequests.get is there. Let's assume we can add analyze or just use generic post.
            // For safety, I will use apiClient.post with the URL constructed from constants or manually.

            await apiClient.post(`/v1/repos/${repositoryId}/prs/${number}/analyze`);

            // Poll or just wait a bit and refresh? 
            // For now, simple user feedback.
            alert('Analysis triggered! Refresh in a few moments.');
            fetchPR();
        } catch (err) {
            console.error('Failed to trigger analysis:', err);
            alert('Failed to trigger analysis.');
        } finally {
            setAnalyzing(false);
        }
    };

    if (loading) {
        return (
            <div className="flex h-[400px] items-center justify-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary" />
            </div>
        )
    }

    if (error || !pr) {
        return (
            <div className="flex flex-col items-center justify-center min-h-[400px] space-y-4">
                <p className="text-red-500">{error || 'Pull Request not found'}</p>
                <Button variant="outline" asChild>
                    <Link href={`/repositories/${repositoryId}`}>Back to Repository</Link>
                </Button>
            </div>
        )
    }

    const getDecisionBadge = (decision?: string) => {
        switch (decision?.toLowerCase()) {
            case 'approve':
            case 'approved':
                return <Badge className="bg-green-500 hover:bg-green-600"><CheckCircle className="w-3 h-3 mr-1" /> Approved</Badge>;
            case 'request_changes':
                return <Badge className="bg-orange-500 hover:bg-orange-600"><AlertTriangle className="w-3 h-3 mr-1" /> Changes Requested</Badge>;
            case 'block':
            case 'blocked':
                return <Badge variant="destructive"><AlertOctagon className="w-3 h-3 mr-1" /> Blocked</Badge>;
            default:
                return <Badge variant="secondary">Pending Review</Badge>;
        }
    };

    return (
        <div className="space-y-6 max-w-5xl mx-auto pb-10">
            {/* Header */}
            <div className="space-y-4">
                <Link href={`/repositories/${repositoryId}`} className="text-sm text-muted-foreground hover:text-foreground flex items-center gap-1 w-fit">
                    <ArrowLeft className="h-4 w-4" /> Back to Repository
                </Link>

                <div className="flex items-start justify-between">
                    <div className="space-y-1">
                        <div className="flex items-center gap-3">
                            <h1 className="text-2xl font-bold">#{pr.number} {pr.title}</h1>
                            <Badge variant="outline" className={PR_STATUS_COLORS[pr.state]}>
                                {pr.state}
                            </Badge>
                        </div>
                        <p className="text-muted-foreground">
                            Opened by <span className="font-medium text-foreground">{pr.author?.username}</span> on {new Date(pr.createdAt).toLocaleDateString()}
                        </p>
                    </div>

                    <div className="flex gap-2">
                        <Button variant="outline" onClick={() => window.open(pr.url, '_blank')}>
                            View on GitHub
                        </Button>
                        <Button onClick={handleReanalyze} disabled={analyzing}>
                            <RefreshCw className={`mr-2 h-4 w-4 ${analyzing ? 'animate-spin' : ''}`} />
                            {analyzing ? 'Analyzing...' : 'Re-analyze'}
                        </Button>
                    </div>
                </div>
            </div>

            <div className="grid gap-6 md:grid-cols-3">
                {/* Main Content - AI Analysis */}
                <div className="md:col-span-2 space-y-6">
                    <Card className="border-primary/20 shadow-sm">
                        <CardHeader className="bg-primary/5 pb-4">
                            <div className="flex items-center justify-between">
                                <div className="flex items-center gap-2">
                                    <Bot className="h-5 w-5 text-primary" />
                                    <CardTitle>AI Review Summary</CardTitle>
                                </div>
                                {getDecisionBadge(pr.aiDecision)}
                            </div>
                        </CardHeader>
                        <CardContent className="pt-6">
                            {pr.aiSummary ? (
                                <div className="prose dark:prose-invert max-w-none">
                                    <div className="whitespace-pre-wrap text-sm leading-relaxed">
                                        {pr.aiSummary}
                                    </div>
                                </div>
                            ) : (
                                <div className="text-center py-8 text-muted-foreground">
                                    <Bot className="h-10 w-10 mx-auto mb-3 opacity-20" />
                                    <p>No AI analysis available yet.</p>
                                    <p className="text-xs mt-1">Click "Re-analyze" to start a new review.</p>
                                </div>
                            )}
                        </CardContent>
                    </Card>

                    {/* Add more tabs or sections for File Changes later */}
                </div>

                {/* Sidebar - Meta Info */}
                <div className="space-y-6">
                    <Card>
                        <CardHeader>
                            <CardTitle className="text-sm font-medium">Details</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div className="flex justify-between text-sm">
                                <span className="text-muted-foreground">Repository</span>
                                <span className="font-mono">{repositoryId.substring(0, 8)}...</span>
                            </div>
                            {/* Add more meta info if available */}
                        </CardContent>
                    </Card>
                </div>
            </div>
        </div>
    )
}
