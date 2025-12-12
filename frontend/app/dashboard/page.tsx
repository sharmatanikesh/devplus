'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { GitPullRequest, Package, TrendingUp, Clock, CheckCircle, AlertCircle } from 'lucide-react';

export default function DashboardPage() {
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
          title="Open PRs"
          value="12"
          description="+2 from last week"
          icon={<GitPullRequest className="h-4 w-4 text-muted-foreground" />}
          trend="up"
        />
        <StatsCard
          title="Avg Review Time"
          value="4.2h"
          description="-8% from last week"
          icon={<Clock className="h-4 w-4 text-muted-foreground" />}
          trend="down"
          trendPositive
        />
        <StatsCard
          title="Merge Rate"
          value="85.5%"
          description="+3.2% from last week"
          icon={<TrendingUp className="h-4 w-4 text-muted-foreground" />}
          trend="up"
          trendPositive
        />
        <StatsCard
          title="Active Repos"
          value="8"
          description="2 syncing now"
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
            <PRItem
              number={42}
              title="Add new feature: AI-powered code review"
              repo="devpulse-frontend"
              status="approved"
              time="2 hours ago"
            />
            <PRItem
              number={41}
              title="Fix: Update authentication flow"
              repo="devpulse-backend"
              status="pending"
              time="5 hours ago"
            />
            <PRItem
              number={40}
              title="Refactor: Improve metrics calculation"
              repo="devpulse-backend"
              status="changes_requested"
              time="1 day ago"
            />
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
