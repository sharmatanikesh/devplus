'use client';
// Force HMR Update

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { GitPullRequest, Package, TrendingUp, Clock, CheckCircle, AlertCircle } from 'lucide-react';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { apiClient } from '@/lib/api-client';

export default function DashboardPage() {
  const router = useRouter();
  const [stats, setStats] = useState<any>(null);
  const [recentPRs, setRecentPRs] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Verify session by fetching a protected resource
    const checkAuthAndFetch = async () => {
      try {
        const [statsRes, activeRes] = await Promise.all([
          apiClient.dashboard.stats(),
          apiClient.dashboard.recentActivity()
        ]);

        if (!statsRes.success && statsRes.error?.code === 'HTTP_401') {
          router.push('/login');
          return;
        }

        if (statsRes.success) {
          setStats(statsRes.data);
        }
        if (activeRes.success) {
          setRecentPRs(activeRes.data || []);
        }
      } catch (error) {
        console.error("Failed to fetch dashboard data", error);
      } finally {
        setLoading(false);
      }
    };
    checkAuthAndFetch();
  }, [router]);

  if (loading) {
    return <div>Loading...</div>; // TODO: Add better skeleton
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground">
          Overview of your engineering metrics and recent activity
        </p>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <StatsCard
          title="Total PRs"
          value={stats?.total_prs?.toString() || "0"}
          description="All time"
          icon={<GitPullRequest className="h-4 w-4 text-muted-foreground" />}
        />
        <StatsCard
          title="Open PRs"
          value={stats?.open_prs?.toString() || "0"}
          description="Currently active"
          icon={<Clock className="h-4 w-4 text-muted-foreground" />}
        />
        <StatsCard
          title="Merged PRs"
          value={stats?.merged_prs?.toString() || "0"}
          description="Successfully merged"
          icon={<TrendingUp className="h-4 w-4 text-muted-foreground" />}
          trendPositive
        />
        <StatsCard
          title="Active Repos"
          value={stats?.active_repos?.toString() || "0"}
          description="Tracked repositories"
          icon={<Package className="h-4 w-4 text-muted-foreground" />}
        />
      </div>

      {/* Recent Activity */}
      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Recent Pull Requests</CardTitle>
            <CardDescription>Latest PR reviews and activity</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {recentPRs.length === 0 ? (
              <p className="text-sm text-muted-foreground">No recent pull requests found.</p>
            ) : (
              recentPRs.map((pr) => (
                <PRItem
                  key={pr.id}
                  number={pr.number}
                  title={pr.title || "Untitled"}
                  repo={pr.repository?.name || "Unknown Repo"}
                  status={pr.state === 'open' ? 'pending' : (pr.state === 'closed' ? 'approved' : 'changes_requested')} // Mapping github state to UI
                  time={new Date(pr.created_at).toLocaleDateString()}
                />
              ))
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>System Health</CardTitle>
            <CardDescription>Service status and uptime</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <ServiceStatus name="Database" status="up" />
            <ServiceStatus name="GitHub Integration" status="up" />
            <ServiceStatus name="Kestra Workflows" status="up" />
            <ServiceStatus name="AI Analysis" status="up" />
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

function StatsCard({
  title,
  value,
  description,
  icon,
  trend,
  trendPositive,
}: {
  title: string;
  value: string;
  description: string;
  icon: React.ReactNode;
  trend?: 'up' | 'down';
  trendPositive?: boolean;
}) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        {icon}
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{value}</div>
        <p className={`text-xs ${trendPositive ? 'text-green-600' : 'text-muted-foreground'}`}>
          {description}
        </p>
      </CardContent>
    </Card>
  );
}

function PRItem({
  number,
  title,
  repo,
  status,
  time,
}: {
  number: number;
  title: string;
  repo: string;
  status: 'approved' | 'pending' | 'changes_requested';
  time: string;
}) {
  const statusConfig = {
    approved: { color: 'bg-green-100 text-green-800', label: 'Approved' },
    pending: { color: 'bg-yellow-100 text-yellow-800', label: 'Pending' },
    changes_requested: { color: 'bg-orange-100 text-orange-800', label: 'Changes Requested' },
  };

  const config = statusConfig[status];

  return (
    <div className="flex items-start gap-3 pb-4 border-b last:border-0 last:pb-0">
      <div className="flex-1 space-y-1">
        <p className="text-sm font-medium leading-none">
          #{number} {title}
        </p>
        <p className="text-sm text-muted-foreground">
          {repo} • {time}
        </p>
      </div>
      <Badge className={config.color} variant="secondary">
        {config.label}
      </Badge>
    </div>
  );
}

function ServiceStatus({ name, status }: { name: string; status: 'up' | 'down' }) {
  return (
    <div className="flex items-center justify-between pb-4 border-b last:border-0 last:pb-0">
      <span className="text-sm font-medium">{name}</span>
      <div className="flex items-center gap-2">
        {status === 'up' ? (
          <>
            <CheckCircle className="h-4 w-4 text-green-600" />
            <span className="text-sm text-green-600">Operational</span>
          </>
        ) : (
          <>
            <AlertCircle className="h-4 w-4 text-red-600" />
            <span className="text-sm text-red-600">Down</span>
          </>
        )}
      </div>
    </div>
  );
}
