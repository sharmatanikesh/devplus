'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Code2, Star, TrendingUp, Award, GitCommit, GitPullRequest, PieChart as PieChartIcon } from 'lucide-react';
import { motion } from 'framer-motion';
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from 'recharts';
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
        <div className="flex flex-col items-center gap-4">
          {/* Custom Loading Pulse */}
          <div className="relative flex items-center justify-center">
            <div className="absolute inset-0 bg-sky-500/20 rounded-full blur-xl animate-pulse" />
            <div className="h-10 w-10 text-sky-500 animate-spin">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 12a9 9 0 1 1-6.219-8.56" /></svg>
            </div>
          </div>
          <p className="text-muted-foreground animate-pulse">Loading your personalized metrics...</p>
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
  // Prepare data for Recharts PieChart
  const pieData = topLanguages.map(l => ({ name: l.language, value: l.percent, fill: getColorForLanguage(l.language) }));

  const container = {
    hidden: { opacity: 0 },
    show: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1
      }
    }
  };

  const item = {
    hidden: { opacity: 0, y: 20 },
    show: { opacity: 1, y: 0 }
  };

  return (
    <motion.div
      variants={container}
      initial="hidden"
      animate="show"
      className="space-y-8"
    >
      {/* Header */}
      <motion.div variants={item} className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight bg-clip-text text-transparent bg-linear-to-r from-foreground to-foreground/70">
            Your GitHub Metrics
          </h1>
          <p className="text-muted-foreground">
            Personalized insights into your coding activity and contributions
          </p>
        </div>
        <Select value={timeRange} onValueChange={setTimeRange}>
          <SelectTrigger className="w-[180px] border-sky-100 dark:border-sky-900/50 focus:ring-sky-500/20">
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
      </motion.div>

      {/* Overview Cards */}
      <motion.div variants={item} className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <MetricCard
          title="Total Commits"
          value={metrics.total_commits.toLocaleString()}
          description="All time commits"
          icon={<GitCommit className="h-4 w-4 text-sky-500" />}
          trend="up"
          delay={0}
        />
        <MetricCard
          title="Pull Requests"
          value={metrics.pull_request_count.toLocaleString()}
          description="PRs created"
          icon={<GitPullRequest className="h-4 w-4 text-sky-500" />}
          trend="up"
          delay={0.1}
        />
        <MetricCard
          title="Total Contributions"
          value={metrics.total_contributions.toLocaleString()}
          description="Commits + PRs"
          icon={<TrendingUp className="h-4 w-4 text-sky-500" />}
          trend="up"
          delay={0.2}
        />
        <MetricCard
          title="Issues Opened"
          value={metrics.issues_opened.toLocaleString()}
          description="Issues created"
          icon={<Award className="h-4 w-4 text-sky-500" />}
          delay={0.3}
        />
      </motion.div>

      {/* Tabs for Different Views */}
      <motion.div variants={item}>
        <Tabs defaultValue="languages" className="space-y-6">
          <TabsList className="bg-muted/50 p-1 border border-border/50">
            <TabsTrigger value="languages" className="data-[state=active]:bg-background data-[state=active]:text-sky-600 data-[state=active]:shadow-sm">Languages</TabsTrigger>
            <TabsTrigger value="top-repos" className="data-[state=active]:bg-background data-[state=active]:text-sky-600 data-[state=active]:shadow-sm">Top Repos</TabsTrigger>
          </TabsList>

          {/* Languages Tab */}
          <TabsContent value="languages" className="space-y-4">
            <div className="grid md:grid-cols-2 gap-6">
              {/* Chart */}
              <Card className="border-none shadow-lg bg-card/40 backdrop-blur-sm ring-1 ring-border/50">
                <CardHeader>
                  <CardTitle className="text-lg flex items-center gap-2">
                    <PieChartIcon className="h-5 w-5 text-sky-500" />
                    Language Distribution
                  </CardTitle>
                </CardHeader>
                <CardContent className="h-[300px] flex items-center justify-center">
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie
                        data={pieData}
                        cx="50%"
                        cy="50%"
                        innerRadius={60}
                        outerRadius={80}
                        paddingAngle={5}
                        dataKey="value"
                      >
                        {pieData.map((entry, index) => (
                          <Cell key={`cell-${index}`} fill={entry.fill} stroke="transparent" />
                        ))}
                      </Pie>
                      <Tooltip
                        contentStyle={{ borderRadius: '8px', border: 'none', background: 'rgba(255,255,255,0.9)', boxShadow: '0 4px 12px rgba(0,0,0,0.1)' }}
                        formatter={(value: number, name: string) => [`${value.toFixed(1)}%`, name]}
                      />
                    </PieChart>
                  </ResponsiveContainer>
                </CardContent>
              </Card>

              {/* List */}
              <Card className="border-none shadow-lg bg-card/40 backdrop-blur-sm ring-1 ring-border/50">
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Code2 className="h-5 w-5 text-sky-500" />
                    Details
                  </CardTitle>
                  <CardDescription>Breakdown by usage percentage</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-6">
                    {topLanguages.map((lang, index) => (
                      <div key={lang.language} className="group">
                        <div className="flex items-center justify-between mb-2">
                          <div className="flex items-center gap-3">
                            <span className={`flex items-center justify-center w-6 h-6 rounded-full text-xs font-bold ${index === 0 ? 'bg-sky-100 text-sky-700' : 'bg-muted text-muted-foreground'}`}>
                              {index + 1}
                            </span>
                            <span className="font-medium group-hover:text-sky-500 transition-colors">{lang.language}</span>
                          </div>
                          <span className="font-mono text-sm font-bold text-muted-foreground group-hover:text-foreground transition-colors">
                            {lang.percent.toFixed(1)}%
                          </span>
                        </div>
                        <div className="w-full bg-muted/50 rounded-full h-2 overflow-hidden">
                          <motion.div
                            initial={{ width: 0 }}
                            animate={{ width: `${lang.percent}%` }}
                            transition={{ duration: 1, delay: 0.5 + (index * 0.1) }}
                            className="h-full rounded-full"
                            style={{ backgroundColor: getColorForLanguage(lang.language) }}
                          />
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          {/* Top Repos Tab */}
          <TabsContent value="top-repos" className="space-y-4">
            <Card className="border-none shadow-lg bg-card/40 backdrop-blur-sm ring-1 ring-border/50">
              <CardHeader>
                <CardTitle>Your Top Repositories</CardTitle>
                <CardDescription>Most impactful repositories</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid gap-4 md:grid-cols-2">
                  {metrics.top_repos.slice(0, 6).map((repo, index) => (
                    <RepoTiltCard key={`${repo.owner}/${repo.name}`} repo={repo} index={index} />
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </motion.div>
    </motion.div>
  );
}

function MetricCard({
  title,
  value,
  description,
  icon,
  trend,
  delay = 0
}: any) {
  const [rotate, setRotate] = useState({ x: 0, y: 0 });

  const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
    const card = e.currentTarget;
    const box = card.getBoundingClientRect();
    const x = e.clientX - box.left;
    const y = e.clientY - box.top;
    const centerX = box.width / 2;
    const centerY = box.height / 2;
    const rotateX = (centerY - y) / 15; // Reduced intensity
    const rotateY = (x - centerX) / 15;

    setRotate({ x: rotateX, y: rotateY });
  };

  const handleMouseLeave = () => {
    setRotate({ x: 0, y: 0 });
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4, delay }}
      style={{ perspective: 1000 }}
    >
      <Card
        className="overflow-hidden border-none shadow-md hover:shadow-xl transition-all duration-200 ease-out bg-card/50 backdrop-blur-sm group hover:scale-[1.02]"
        onMouseMove={handleMouseMove}
        onMouseLeave={handleMouseLeave}
        style={{
          transform: `rotateX(${rotate.x}deg) rotateY(${rotate.y}deg)`,
          transition: "transform 0.1s ease-out"
        }}
      >
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2 relative z-10">
          <CardTitle className="text-sm font-medium text-muted-foreground group-hover:text-foreground transition-colors">{title}</CardTitle>
          <div className="p-2 bg-background/50 rounded-lg group-hover:bg-sky-50 dark:group-hover:bg-sky-900/20 transition-colors">
            {icon}
          </div>
        </CardHeader>
        <CardContent className="relative z-10">
          <div className="text-2xl font-bold tracking-tight">{value}</div>
          <p className="text-xs text-muted-foreground mt-1 flex items-center gap-1">
            {trend === 'up' && <TrendingUp className="h-3 w-3 text-green-500" />}
            <span className={trend === 'up' ? 'text-green-600 dark:text-green-400' : ''}>{description}</span>
          </p>
        </CardContent>
      </Card>
    </motion.div>
  );
}

function RepoTiltCard({ repo, index }: { repo: TopRepository, index: number }) {
  const [rotate, setRotate] = useState({ x: 0, y: 0 });

  const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
    const card = e.currentTarget;
    const box = card.getBoundingClientRect();
    const x = e.clientX - box.left;
    const y = e.clientY - box.top;
    const centerX = box.width / 2;
    const centerY = box.height / 2;
    const rotateX = (centerY - y) / 45; // Reduced intensity 
    const rotateY = (x - centerX) / 45;

    setRotate({ x: rotateX, y: rotateY });
  };

  const handleMouseLeave = () => {
    setRotate({ x: 0, y: 0 });
  };

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.95 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ delay: index * 0.05 }}
      style={{ perspective: 1000 }}
    >
      <div
        className="p-4 rounded-xl border bg-background/50 hover:bg-background hover:border-sky-200 dark:hover:border-sky-800 transition-all duration-200 group h-full shadow-sm"
        onMouseMove={handleMouseMove}
        onMouseLeave={handleMouseLeave}
        style={{
          transform: `rotateX(${rotate.x}deg) rotateY(${rotate.y}deg)`,
          transition: "transform 0.1s ease-out"
        }}
      >
        <div className="flex items-start justify-between mb-2">
          <div className="flex items-center gap-2">
            <span className="text-lg font-bold text-muted-foreground/30 group-hover:text-sky-500/50 transition-colors">#{index + 1}</span>
            <h4 className="font-semibold group-hover:text-sky-600 transition-colors">
              {repo.name}
            </h4>
          </div>
          {repo.language && (
            <Badge variant="secondary" className="bg-sky-50 text-sky-700 dark:bg-sky-900/30 dark:text-sky-300">
              {repo.language}
            </Badge>
          )}
        </div>
        <p className="text-sm text-muted-foreground line-clamp-2 mb-4 h-10">
          {repo.description || "No description provided."}
        </p>

        <div className="flex items-center gap-4 text-xs text-muted-foreground">
          <div className="flex items-center gap-1">
            <Star className="h-3 w-3 text-amber-400" />
            {repo.stars}
          </div>
          <div className="flex items-center gap-1">
            <GitCommit className="h-3 w-3" />
            {repo.commits}
          </div>
        </div>
      </div>
    </motion.div>
  );
}

function getColorForLanguage(language: string) {
  const colors: Record<string, string> = {
    'TypeScript': '#3178c6',
    'JavaScript': '#f1e05a',
    'Python': '#3572A5',
    'Java': '#b07219',
    'Go': '#00ADD8',
    'Rust': '#dea584',
    'HTML': '#e34c26',
    'CSS': '#563d7c',
  };
  return colors[language] || '#94a3b8'; // default slate-400
}
