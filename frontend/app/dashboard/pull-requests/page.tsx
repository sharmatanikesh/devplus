'use client';

import { useState, useEffect } from 'react';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { GitPullRequest, GitCommit, FileCode, Loader2, Sparkles } from 'lucide-react';
import Link from 'next/link';
import type { PullRequest } from '@/lib/types';
import { formatRelativeTime } from '@/lib/utils/format';

export default function PullRequestsPage() {
  const [pullRequests, setPullRequests] = useState<PullRequest[]>([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState('open');

  useEffect(() => {
    fetchPullRequests();
  }, [filter]);

  const fetchPullRequests = async () => {
    try {
      setLoading(true);
      // TODO: Fetch from all repos or selected repo
      const response = await fetch(`/api/v1/repos/1/prs?state=${filter}`);
      const data = await response.json();
      if (data.success) {
        setPullRequests(data.data.items);
      }
    } catch (error) {
      console.error('Failed to fetch PRs:', error);
    } finally {
      setLoading(false);
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
          <h1 className="text-3xl font-bold tracking-tight">Pull Requests</h1>
          <p className="text-muted-foreground">
            Review and manage pull requests across your repositories
          </p>
        </div>
        <Select value={filter} onValueChange={setFilter}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Filter PRs" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="open">Open</SelectItem>
            <SelectItem value="closed">Closed</SelectItem>
            <SelectItem value="merged">Merged</SelectItem>
            <SelectItem value="all">All</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* PR List */}
      <div className="space-y-3">
        {pullRequests.map((pr) => (
          <Card key={pr.id} className="hover:shadow-md transition-shadow">
            <CardContent className="p-6">
              <div className="flex gap-4">
                {/* PR Icon */}
                <div className="shrink-0 mt-1">
                  <div className={`h-10 w-10 rounded-full flex items-center justify-center ${
                    pr.state === 'open' ? 'bg-green-100' : 
                    pr.state === 'merged' ? 'bg-purple-100' : 'bg-red-100'
                  }`}>
                    <GitPullRequest className={`h-5 w-5 ${
                      pr.state === 'open' ? 'text-green-600' : 
                      pr.state === 'merged' ? 'text-purple-600' : 'text-red-600'
                    }`} />
                  </div>
                </div>

                {/* PR Details */}
                <div className="flex-1 space-y-3">
                  {/* Title and badges */}
                  <div className="flex items-start justify-between gap-4">
                    <div className="space-y-1">
                      <Link
                        href={`/pull-requests/${pr.id}`}
                        className="text-lg font-semibold hover:underline inline-block"
                      >
                        {pr.title}
                      </Link>
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <span>#{pr.number}</span>
                        <span>•</span>
                        <span>opened {formatRelativeTime(pr.createdAt)}</span>
                        <span>•</span>
                        <div className="flex items-center gap-1">
                          <Avatar className="h-5 w-5">
                            <AvatarImage src={pr.author.avatarUrl} />
                            <AvatarFallback>{pr.author.username[0]}</AvatarFallback>
                          </Avatar>
                          <span>{pr.author.username}</span>
                        </div>
                      </div>
                    </div>
                    <div className="flex gap-2">
                      <PRStatusBadge status={pr.state} />
                      {pr.reviewStatus && <ReviewStatusBadge status={pr.reviewStatus} />}
                      {pr.aiAnalysisStatus === 'completed' && (
                        <Badge className="bg-purple-100 text-purple-800">
                          <Sparkles className="mr-1 h-3 w-3" />
                          AI Reviewed
                        </Badge>
                      )}
                    </div>
                  </div>

                  {/* AI Summary */}
                  {pr.aiSummary && (
                    <div className="bg-purple-50 border border-purple-200 rounded-lg p-3">
                      <div className="flex items-start gap-2">
                        <Sparkles className="h-4 w-4 text-purple-600 mt-0.5" />
                        <div className="flex-1">
                          <p className="text-sm font-medium text-purple-900">AI Summary</p>
                          <p className="text-sm text-purple-700 mt-1">{pr.aiSummary}</p>
                        </div>
                      </div>
                    </div>
                  )}

                  {/* Metadata */}
                  <div className="flex items-center gap-6 text-sm text-muted-foreground">
                    <div className="flex items-center gap-1">
                      <GitCommit className="h-4 w-4" />
                      {pr.commits} commits
                    </div>
                    <div className="flex items-center gap-1">
                      <FileCode className="h-4 w-4" />
                      {pr.filesChanged} files
                    </div>
                    <div className="flex items-center gap-2">
                      <span className="text-green-600">+{pr.additions}</span>
                      <span className="text-red-600">-{pr.deletions}</span>
                    </div>
                    <div>
                      {pr.baseBranch} ← {pr.headBranch}
                    </div>
                  </div>

                  {/* Actions */}
                  <div className="flex gap-2">
                    <Button size="sm" variant="outline" asChild>
                      <Link href={`/pull-requests/${pr.id}`}>
                        View Details
                      </Link>
                    </Button>
                    <Button size="sm" variant="outline" asChild>
                      <a href={pr.url} target="_blank" rel="noopener noreferrer">
                        Open in GitHub
                      </a>
                    </Button>
                    {pr.aiAnalysisStatus !== 'completed' && (
                      <Button size="sm" variant="outline">
                        <Sparkles className="mr-2 h-4 w-4" />
                        Analyze
                      </Button>
                    )}
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        ))}

        {pullRequests.length === 0 && (
          <Card>
            <CardContent className="flex flex-col items-center justify-center py-12">
              <GitPullRequest className="h-12 w-12 text-muted-foreground mb-4" />
              <h3 className="text-lg font-semibold mb-2">No pull requests found</h3>
              <p className="text-muted-foreground text-center max-w-md">
                There are no {filter} pull requests at the moment.
              </p>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}

function PRStatusBadge({ status }: { status: string }) {
  const config = {
    open: { color: 'bg-green-100 text-green-800', label: 'Open' },
    closed: { color: 'bg-red-100 text-red-800', label: 'Closed' },
    merged: { color: 'bg-purple-100 text-purple-800', label: 'Merged' },
  };
  
  const statusConfig = config[status as keyof typeof config];
  if (!statusConfig) return null;

  return (
    <Badge className={statusConfig.color}>
      {statusConfig.label}
    </Badge>
  );
}

function ReviewStatusBadge({ status }: { status: string }) {
  const config = {
    pending: { color: 'bg-yellow-100 text-yellow-800', label: 'Pending Review' },
    approved: { color: 'bg-green-100 text-green-800', label: 'Approved' },
    changes_requested: { color: 'bg-orange-100 text-orange-800', label: 'Changes Requested' },
    blocked: { color: 'bg-red-100 text-red-800', label: 'Blocked' },
  };
  
  const statusConfig = config[status as keyof typeof config];
  if (!statusConfig) return null;

  return (
    <Badge className={statusConfig.color}>
      {statusConfig.label}
    </Badge>
  );
}
