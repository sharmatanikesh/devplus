'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Github, Zap, BarChart3, GitPullRequest, Package } from 'lucide-react';

export default function LoginPage() {
  const [isLoading, setIsLoading] = useState(false);

  const handleGithubLogin = async () => {
    setIsLoading(true);
    try {
      // Redirect to GitHub OAuth
      window.location.href = '/api/v1/auth/github/connect';
    } catch (error) {
      console.error('Login error:', error);
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-linear-to-br from-slate-50 to-slate-100 dark:from-slate-950 dark:to-slate-900 p-4">
      <div className="w-full max-w-6xl grid lg:grid-cols-2 gap-8 items-center">
        {/* Left Side - Branding & Features */}
        <div className="space-y-6">
          <div className="space-y-2">
            <div className="flex items-center gap-2">
              <div className="h-10 w-10 bg-primary rounded-lg flex items-center justify-center">
                <Zap className="h-6 w-6 text-primary-foreground" />
              </div>
              <h1 className="text-4xl font-bold tracking-tight">DevPulse</h1>
            </div>
            <p className="text-xl text-muted-foreground">
              AI-Powered Engineering Intelligence
            </p>
          </div>

          <div className="space-y-4 pt-4">
            <FeatureItem
              icon={<GitPullRequest className="h-5 w-5" />}
              title="Automated PR Reviews"
              description="AI analyzes your pull requests and provides intelligent feedback"
            />
            <FeatureItem
              icon={<Package className="h-5 w-5" />}
              title="Smart Release Notes"
              description="Generate comprehensive changelogs from your commits automatically"
            />
            <FeatureItem
              icon={<BarChart3 className="h-5 w-5" />}
              title="Engineering Metrics"
              description="Track PR lead time, review latency, and team performance"
            />
            <FeatureItem
              icon={<Zap className="h-5 w-5" />}
              title="Architecture Impact"
              description="Identify hotspots and analyze code coupling in real-time"
            />
          </div>
        </div>

        {/* Right Side - Login Card */}
        <Card className="shadow-2xl">
          <CardHeader className="space-y-1">
            <CardTitle className="text-2xl">Welcome back</CardTitle>
            <CardDescription>
              Sign in with your GitHub account to get started
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <Button
              className="w-full h-12 text-base"
              size="lg"
              onClick={handleGithubLogin}
              disabled={isLoading}
            >
              <Github className="mr-2 h-5 w-5" />
              {isLoading ? 'Connecting...' : 'Continue with GitHub'}
            </Button>

            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <span className="w-full border-t" />
              </div>
              <div className="relative flex justify-center text-xs uppercase">
                <span className="bg-background px-2 text-muted-foreground">
                  Secure OAuth Authentication
                </span>
              </div>
            </div>

            <div className="text-center text-sm text-muted-foreground">
              By signing in, you agree to our Terms of Service and Privacy Policy
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

function FeatureItem({
  icon,
  title,
  description,
}: {
  icon: React.ReactNode;
  title: string;
  description: string;
}) {
  return (
    <div className="flex gap-3">
      <div className="shrink-0 h-10 w-10 bg-primary/10 rounded-lg flex items-center justify-center text-primary">
        {icon}
      </div>
      <div className="space-y-1">
        <h3 className="font-semibold">{title}</h3>
        <p className="text-sm text-muted-foreground">{description}</p>
      </div>
    </div>
  );
}
