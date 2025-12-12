'use client';

import { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { TrendingUp, TrendingDown, Clock, GitPullRequest, GitMerge, Code, Loader2 } from 'lucide-react';
import type { EngineeringMetrics } from '@/lib/types';

export default function MetricsPage() {
  const [metrics, setMetrics] = useState<EngineeringMetrics | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchMetrics();
  }, []);

  const fetchMetrics = async () => {
    try {
      setLoading(true);
      const response = await fetch('/api/v1/metrics');
      const data = await response.json();
      if (data.success) {
        setMetrics(data.data);
      }
    } catch (error) {
      console.error('Failed to fetch metrics:', error);
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

  if (!metrics) {
    return <div>Failed to load metrics</div>;
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Engineering Metrics</h1>
        <p className="text-muted-foreground">
          Track performance and productivity across your team
        </p>
      </div>

      {/* Key Metrics */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <MetricCard
          title="PR Lead Time"
          value={`${metrics.metrics.prLeadTime.average.toFixed(1)}h`}
          description={`Median: ${metrics.metrics.prLeadTime.median.toFixed(1)}h`}
          trend={metrics.metrics.prLeadTime.trend}
          icon={<Clock className="h-4 w-4 text-muted-foreground" />}
        />
        <MetricCard
          title="Review Latency"
          value={`${metrics.metrics.reviewLatency.average.toFixed(1)}h`}
          description={`Median: ${metrics.metrics.reviewLatency.median.toFixed(1)}h`}
          trend={metrics.metrics.reviewLatency.trend}
          icon={<Clock className="h-4 w-4 text-muted-foreground" />}
        />
        <MetricCard
          title="Merge Rate"
          value={`${metrics.metrics.mergeRate.percentage.toFixed(1)}%`}
          description={`${metrics.metrics.mergeRate.merged}/${metrics.metrics.mergeRate.total} PRs`}
          trend={metrics.metrics.mergeRate.trend}
          icon={<GitMerge className="h-4 w-4 text-muted-foreground" />}
        />
        <MetricCard
          title="PR Throughput"
          value={`${metrics.metrics.throughput.prsPerWeek.toFixed(1)}/wk`}
          description={`${metrics.metrics.throughput.commitsPerWeek.toFixed(0)} commits/wk`}
          trend={metrics.metrics.throughput.trend}
          icon={<GitPullRequest className="h-4 w-4 text-muted-foreground" />}
        />
      </div>

      {/* Additional Insights */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Top Contributors */}
        <Card>
          <CardHeader>
            <CardTitle>Top Contributors</CardTitle>
            <CardDescription>Most active team members this period</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {metrics.topContributors.map((contributor, index) => (
              <div key={contributor.username} className="flex items-center justify-between pb-4 border-b last:border-0 last:pb-0">
                <div className="flex items-center gap-3">
                  <div className="text-lg font-bold text-muted-foreground w-6">
                    #{index + 1}
                  </div>
                  <Avatar className="h-10 w-10">
                    <AvatarImage src={contributor.avatarUrl} />
                    <AvatarFallback>{contributor.username[0]}</AvatarFallback>
                  </Avatar>
                  <div>
                    <p className="font-medium">{contributor.username}</p>
                    <p className="text-sm text-muted-foreground">
                      {contributor.commitsCount} commits
                    </p>
                  </div>
                </div>
                <div className="flex gap-4 text-sm">
                  <div className="text-center">
                    <div className="font-semibold text-green-600">{contributor.prsCreated}</div>
                    <div className="text-muted-foreground">Created</div>
                  </div>
                  <div className="text-center">
                    <div className="font-semibold text-blue-600">{contributor.prsReviewed}</div>
                    <div className="text-muted-foreground">Reviewed</div>
                  </div>
                </div>
              </div>
            ))}
          </CardContent>
        </Card>

        {/* Hotspots */}
        <Card>
          <CardHeader>
            <CardTitle>Code Hotspots</CardTitle>
            <CardDescription>Files with the most changes</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {metrics.hotspots.map((hotspot) => (
              <div key={hotspot.file} className="pb-4 border-b last:border-0 last:pb-0">
                <div className="flex items-start justify-between gap-4">
                  <div className="flex-1 space-y-1">
                    <div className="flex items-center gap-2">
                      <Code className="h-4 w-4 text-muted-foreground" />
                      <p className="font-medium text-sm font-mono">{hotspot.file}</p>
                    </div>
                    <p className="text-sm text-muted-foreground">
                      {hotspot.changes} changes • {hotspot.authors} authors
                    </p>
                  </div>
                  <Badge variant="outline" className="bg-orange-50 text-orange-700 border-orange-200">
                    High Activity
                  </Badge>
                </div>
              </div>
            ))}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

function MetricCard({
  title,
  value,
  description,
  trend,
  icon,
}: {
  title: string;
  value: string;
  description: string;
  trend: number;
  icon: React.ReactNode;
}) {
  const isPositive = trend < 0; // For time metrics, negative is good
  const trendValue = Math.abs(trend);

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        {icon}
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{value}</div>
        <div className="flex items-center gap-1 mt-1">
          <p className="text-xs text-muted-foreground">{description}</p>
        </div>
        <div className={`flex items-center gap-1 mt-2 text-xs ${isPositive ? 'text-green-600' : 'text-red-600'}`}>
          {isPositive ? (
            <TrendingDown className="h-3 w-3" />
          ) : (
            <TrendingUp className="h-3 w-3" />
          )}
          <span>{trendValue.toFixed(1)}% from last period</span>
        </div>
      </CardContent>
    </Card>
  );
}
