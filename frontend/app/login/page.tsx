'use client';

import { useState } from 'react';
import { useSearchParams } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Github, Zap, BarChart3, GitPullRequest, Package, AlertCircle } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { createClient } from '@/lib/supabase/client';

export default function LoginPage() {
  const [isLoading, setIsLoading] = useState(false);
  const searchParams = useSearchParams();
  const error = searchParams.get('error');
  const next = searchParams.get('next');

  const handleGithubLogin = async () => {
    setIsLoading(true);
    try {
      const supabase = createClient();
      
      // Construct the callback URL with next parameter if exists
      const redirectTo = next 
        ? `${window.location.origin}/api/v1/auth/callback?next=${encodeURIComponent(next)}`
        : `${window.location.origin}/api/v1/auth/callback`;

      const { data, error } = await supabase.auth.signInWithOAuth({
        provider: 'github',
        options: {
          redirectTo,
        },
      });

      if (error) {
        console.error('Login error:', error);
        setIsLoading(false);
      }
      // Note: If successful, the browser will redirect to GitHub OAuth page
      // The user will be redirected back to /api/v1/auth/callback after authorization
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
            {/* Error message */}
            {error && (
              <div className="bg-red-50 border border-red-200 rounded-lg p-3 flex items-start gap-2">
                <AlertCircle className="h-5 w-5 text-red-600 shrink-0 mt-0.5" />
                <div className="flex-1">
                  <p className="text-sm font-medium text-red-900">Authentication failed</p>
                  <p className="text-sm text-red-700 mt-1">
                    Please try again or contact support if the issue persists.
                  </p>
                </div>
              </div>
            )}

            {/* Redirect notice */}
            {next && (
              <Badge variant="outline" className="w-full justify-center py-2">
                You&apos;ll be redirected after login
              </Badge>
            )}
            
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
