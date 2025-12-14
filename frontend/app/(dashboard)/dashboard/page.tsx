'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { GitPullRequest, Package, TrendingUp, Clock, CheckCircle, AlertCircle, ShieldAlert } from 'lucide-react';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { apiClient } from '@/lib/api-client';
import { cn } from "@/lib/utils";

import { Area, AreaChart, CartesianGrid, ResponsiveContainer, Tooltip as RechartsTooltip, XAxis, YAxis } from 'recharts';
import { motion } from 'framer-motion';
import { Activity, Zap, FileText } from 'lucide-react';
import ReactMarkdown from 'react-markdown';

const velocityData = [
  { name: 'Mon', speed: 40 },
  { name: 'Tue', speed: 55 },
  { name: 'Wed', speed: 45 },
  { name: 'Thu', speed: 70 },
  { name: 'Fri', speed: 65 },
  { name: 'Sat', speed: 30 },
  { name: 'Sun', speed: 45 },
];

export default function DashboardPage() {
  const router = useRouter();
  const [stats, setStats] = useState<any>(null);
  const [recentPRs, setRecentPRs] = useState<any[]>([]);
  const [latestReleaseRepo, setLatestReleaseRepo] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Verify session by fetching a protected resource
    const checkAuthAndFetch = async () => {
      try {
        const [statsRes, activeRes] = await Promise.all([
          apiClient.dashboard.stats(),
          apiClient.dashboard.recentActivity(),
          apiClient.repos.list()
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

        // Find repo with latest release info
        // The third item in the Promise.all result is the repos list
        const reposRes = await apiClient.repos.list(); // Re-calling just to be safe with types, or cast from array above
        if (reposRes.success && reposRes.data && reposRes.data.length > 0) {
          // Prefer one with a changelog
          const repoWithChangelog = reposRes.data.find((r: any) => r.release_changelog);
          setLatestReleaseRepo(repoWithChangelog || reposRes.data[0]);
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
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="flex flex-col items-center gap-4">
          <div className="relative flex items-center justify-center">
            <div className="absolute inset-0 bg-sky-500/20 rounded-full blur-xl animate-pulse" />
            <div className="h-10 w-10 text-sky-500 animate-spin">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 12a9 9 0 1 1-6.219-8.56" /></svg>
            </div>
          </div>
          <p className="text-muted-foreground animate-pulse">Loading dashboard...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
      >
        <h1 className="text-3xl font-bold tracking-tight bg-clip-text text-transparent bg-linear-to-r from-foreground to-foreground/70">Dashboard</h1>
        <p className="text-muted-foreground">Overview of your engineering capabilities.</p>
      </motion.div>

      {/* Overview Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <TiltCard
          title="Active PRs"
          value={stats?.open_prs?.toString() || "0"}
          description="Currently active"
          icon={GitPullRequest}
          color="text-sky-500"
        />
        <TiltCard
          title="Total PRs"
          value={stats?.total_prs?.toString() || "0"}
          description="Lifetime"
          icon={Activity}
          color="text-emerald-500"
        />
        <TiltCard
          title="System Status"
          value="Normal"
          description="All systems go"
          icon={CheckCircle}
          color="text-green-500"
        />
        <TiltCard
          title="Active Repos"
          value={stats?.active_repos?.toString() || "0"}
          description="Tracked"
          icon={Package}
          color="text-amber-500"
        />
      </div>

      {/* AI Analysis Section */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.2 }}
          className="col-span-4"
        >
          <Card className="h-full border-none shadow-lg bg-card/50 backdrop-blur-sm ring-1 ring-border/50">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <FileText className="h-5 w-5 text-primary" />
                Last Release Summary
              </CardTitle>
              <CardDescription>
                {latestReleaseRepo ? `Latest updates from ${latestReleaseRepo.name}` : 'No release data available'}
              </CardDescription>
            </CardHeader>
            <CardContent>
              {(() => {
                let summary = "";
                try {
                  if (latestReleaseRepo?.release_risk_analysis) {
                    const analysis = JSON.parse(latestReleaseRepo.release_risk_analysis);
                    summary = analysis.summary;
                  }
                } catch (e) {
                  console.log("Failed to parse release analysis JSON", e);
                }

                const displayContent = summary || latestReleaseRepo?.release_changelog;

                return displayContent ? (
                  <div className="h-[250px] overflow-y-auto pr-2 custom-scrollbar">
                    <div className="prose prose-sm dark:prose-invert max-w-none text-muted-foreground">
                      <ReactMarkdown>
                        {displayContent}
                      </ReactMarkdown>
                    </div>
                  </div>
                ) : (
                  <div className="h-[250px] flex items-center justify-center text-muted-foreground text-sm">
                    No release changelogs found. Run a release analysis in the Releases page.
                  </div>
                );
              })()}
            </CardContent>
          </Card>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.3 }}
          className="col-span-3"
        >
          <Card className="h-full border-none shadow-lg bg-card/50 backdrop-blur-sm ring-1 ring-border/50">
            <CardHeader>
              <CardTitle>Recent Activity</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-6">
                {recentPRs.length > 0 ? recentPRs.slice(0, 3).map((pr, i) => (
                  <div key={pr.id} className="flex items-center pb-4 border-b last:border-0 last:pb-0">
                    <div className="space-y-1">
                      <p className="text-sm font-medium leading-none truncate max-w-[200px]">{pr.title}</p>
                      <p className="text-xs text-muted-foreground">{new Date(pr.created_at).toLocaleDateString()} in {pr.repository?.name}</p>
                    </div>
                    <div className="ml-auto">
                      <Badge variant="secondary" className={cn(
                        "capitalize",
                        pr.state === 'merged' && "bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-300",
                        pr.state === 'closed' && "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300",
                        pr.state === 'open' && "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300"
                      )}>
                        {pr.state}
                      </Badge>
                    </div>
                  </div>
                )) : (
                  <div className="flex items-center">
                    <div className="space-y-1">
                      <p className="text-sm font-medium leading-none">System Ready</p>
                      <p className="text-xs text-muted-foreground">Waiting for activity...</p>
                    </div>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </motion.div>
      </div>

    </div>
  )
}

function TiltCard({ title, value, description, icon: Icon, color }: any) {
  const [rotate, setRotate] = useState({ x: 0, y: 0 });

  const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
    const card = e.currentTarget;
    const box = card.getBoundingClientRect();
    const x = e.clientX - box.left;
    const y = e.clientY - box.top;
    const centerX = box.width / 2;
    const centerY = box.height / 2;
    const rotateX = (centerY - y) / 10; // Stronger tilt
    const rotateY = (x - centerX) / 10;

    setRotate({ x: rotateX, y: rotateY });
  };

  const handleMouseLeave = () => {
    setRotate({ x: 0, y: 0 });
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
      style={{ perspective: 1000 }}
    >
      <Card
        className="border-none shadow-md hover:shadow-xl transition-all duration-200 ease-out bg-card/40 backdrop-blur-md ring-1 ring-border/50"
        onMouseMove={handleMouseMove}
        onMouseLeave={handleMouseLeave}
        style={{
          transform: `rotateX(${rotate.x}deg) rotateY(${rotate.y}deg)`,
          transition: "transform 0.1s ease-out"
        }}
      >
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">
            {title}
          </CardTitle>
          <Icon className={cn("h-4 w-4 text-muted-foreground", color)} />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{value}</div>
          <p className="text-xs text-muted-foreground">
            {description}
          </p>
        </CardContent>
      </Card>
    </motion.div>
  )
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
