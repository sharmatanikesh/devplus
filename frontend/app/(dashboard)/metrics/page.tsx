'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Code2, Star, TrendingUp, Award, GitCommit, GitPullRequest } from 'lucide-react';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { apiClient } from '@/lib/api-client';

interface LanguageStat {
  language: string;
  count: number;
  bytes: number;
  percent: number;
}

interface RepositoryContribution {
  repo_name: string;
  repo_owner: string;
  commit_count: number;
  added_lines: number;
  deleted_lines: number;
  pull_requests: number;
  last_commit_date: string;
}

interface Activity {
  type: string;
  repo_name: string;
  message: string;
  date: string;
  url: string;
  additions?: number;
  deletions?: number;
}

interface TopRepository {
  name: string;
  owner: string;
  stars: number;
  forks: number;
  commits: number;
  language: string;
  description: string;
}

interface PersonalMetrics {
  total_commits: number;
  total_contributions: number;
  pull_request_count: number;
  issues_opened: number;
  time_range: string;
  language_stats: LanguageStat[];
  contribution_by_repo: RepositoryContribution[];
  recent_activity: Activity[];
  top_repos: TopRepository[];
}

export default function MetricsPage() {
  const router = useRouter();
  const [metrics, setMetrics] = useState<PersonalMetrics | null>(null);
  const [loading, setLoading] = useState(true);
  const [timeRange, setTimeRange] = useState('90');

  useEffect(() => {
    const fetchMetrics = async () => {
      setLoading(true);
      try {
        const res = await apiClient.metrics.personal(timeRange);

        if (!res.success && res.error?.code === 'HTTP_401') {
          router.push('/login');
          return;
        }

        if (res.success && res.data) {
          setMetrics(res.data);
        }
      } catch (error) {
        console.error('Failed to fetch metrics', error);
      } finally {
        setLoading(false);
      }
    };
    fetchMetrics();
  }, [router, timeRange]);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">Loading your personalized metrics...</p>
        </div>
      </div>
    );
  }

  if (!metrics) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <p className="text-muted-foreground">Failed to load metrics. Please try again later.</p>
      </div>
    );
  }

  const topLanguages = metrics.language_stats.slice(0, 5);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Your GitHub Metrics</h1>
          <p className="text-muted-foreground">
            Personalized insights into your coding activity and contributions
          </p>
        </div>
        <Select value={timeRange} onValueChange={setTimeRange}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Select time range" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="7">Last 7 days</SelectItem>
            <SelectItem value="30">Last 30 days</SelectItem>
            <SelectItem value="90">Last 90 days</SelectItem>
            <SelectItem value="180">Last 6 months</SelectItem>
            <SelectItem value="365">Last year</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Overview Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <MetricCard
          title="Total Commits"
          value={metrics.total_commits.toLocaleString()}
          description="All time commits"
          icon={<GitCommit className="h-4 w-4 text-muted-foreground" />}
          trend="up"
        />
        <MetricCard
          title="Pull Requests"
          value={metrics.pull_request_count.toLocaleString()}
          description="PRs created"
          icon={<GitPullRequest className="h-4 w-4 text-muted-foreground" />}
          trend="up"
        />
        <MetricCard
          title="Total Contributions"
          value={metrics.total_contributions.toLocaleString()}
          description="Commits + PRs"
          icon={<TrendingUp className="h-4 w-4 text-muted-foreground" />}
          trend="up"
        />
        <MetricCard
          title="Issues Opened"
          value={metrics.issues_opened.toLocaleString()}
          description="Issues created"
          icon={<Award className="h-4 w-4 text-muted-foreground" />}
        />
      </div>

      {/* Tabs for Different Views */}
      <Tabs defaultValue="languages" className="space-y-4">
        <TabsList>
          <TabsTrigger value="languages">Languages</TabsTrigger>
          <TabsTrigger value="top-repos">Top Repos</TabsTrigger>
        </TabsList>

        {/* Languages Tab */}
        <TabsContent value="languages" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Code2 className="h-5 w-5" />
                Programming Languages
              </CardTitle>
              <CardDescription>Languages you&apos;ve used the most in your repositories</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {topLanguages.map((lang, index) => (
                  <div key={lang.language} className="space-y-2">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <span className="text-2xl font-bold text-muted-foreground">#{index + 1}</span>
                        <span className="font-medium">{lang.language}</span>
                      </div>
                      <span className="text-sm text-muted-foreground">
                        {lang.percent.toFixed(1)}%
                      </span>
                    </div>
                    <div className="w-full bg-secondary rounded-full h-2">
                      <div
                        className="bg-primary rounded-full h-2 transition-all"
                        style={{ width: `${lang.percent}%` }}
                      />
                    </div>
                    <p className="text-xs text-muted-foreground">
                      {(lang.bytes / 1024).toFixed(1)} KB written
                    </p>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Top Repos Tab */}
        <TabsContent value="top-repos" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Your Top Repositories</CardTitle>
              <CardDescription>Most impactful repositories based on stars and commits</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {metrics.top_repos.slice(0, 10).map((repo, index) => (
                  <div key={`${repo.owner}/${repo.name}`} className="border-b pb-4 last:border-0">
                    <div className="flex items-start justify-between">
                      <div className="space-y-1 flex-1">
                        <div className="flex items-center gap-2">
                          <span className="text-xl font-bold text-muted-foreground">#{index + 1}</span>
                          <h4 className="font-semibold">
                            {repo.owner}/{repo.name}
                          </h4>
                        </div>
                        {repo.description && (
                          <p className="text-sm text-muted-foreground line-clamp-2">{repo.description}</p>
                        )}
                        <div className="flex flex-wrap gap-3 text-sm">
                          {repo.language && (
                            <Badge variant="secondary">
                              <Code2 className="h-3 w-3 mr-1" />
                              {repo.language}
                            </Badge>
                          )}
                          <span className="flex items-center gap-1 text-muted-foreground">
                            <Star className="h-3 w-3" />
                            {repo.stars.toLocaleString()} stars
                          </span>
                          <span className="flex items-center gap-1 text-muted-foreground">
                            <GitCommit className="h-3 w-3" />
                            {repo.commits} commits
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}

function MetricCard({
  title,
  value,
  description,
  icon,
  trend,
}: {
  title: string;
  value: string;
  description: string;
  icon: React.ReactNode;
  trend?: 'up' | 'down';
}) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        {icon}
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{value}</div>
        <p className={`text-xs ${trend === 'up' ? 'text-green-600' : 'text-muted-foreground'}`}>
          {description}
        </p>
      </CardContent>
    </Card>
  );
}
