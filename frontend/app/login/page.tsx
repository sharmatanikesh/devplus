'use client';

import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Github, Zap, BarChart3, GitPullRequest, Package, ArrowRight, ShieldCheck } from 'lucide-react';
import { apiClient } from '@/lib/api-client';
import { LoadingPulse } from '@/components/ui/loading-pulse';

export default function LoginPage() {
  const [isLoading, setIsLoading] = useState(false);

  const handleGithubLogin = async () => {
    setIsLoading(true);
    try {
      // Redirect to GitHub OAuth
      apiClient.auth.login();
    } catch (error) {
      console.error('Login error:', error);
      setIsLoading(false);
    }
  };

  // Typing animation effect
  const [typedText, setTypedText] = useState('');
  const [isDeleting, setIsDeleting] = useState(false);
  const fullText = "Engineering intelligence reimagined for high-performance teams.";

  // 3D Tilt Effect State
  const [rotate, setRotate] = useState({ x: 0, y: 0 });

  const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
    const card = e.currentTarget;
    const box = card.getBoundingClientRect();
    const x = e.clientX - box.left;
    const y = e.clientY - box.top;
    const centerX = box.width / 2;
    const centerY = box.height / 2;
    const rotateX = (centerY - y) / 20;
    const rotateY = (x - centerX) / 20;

    setRotate({ x: rotateX, y: rotateY });
  };

  const handleMouseLeave = () => {
    setRotate({ x: 0, y: 0 });
  };

  useEffect(() => {
    let timeout: NodeJS.Timeout;

    const animate = () => {
      setTypedText(current => {
        if (isDeleting) {
          if (current === '') {
            setIsDeleting(false);
            return '';
          }
          return current.slice(0, -1);
        } else {
          if (current === fullText) {
            // This is a side effect inside state updater, but for this simple animation it works if we handle the pause outside
            return current;
          }
          return fullText.slice(0, current.length + 1);
        }
      });
    };

    // Control timing based on state
    if (!isDeleting && typedText === fullText) {
      timeout = setTimeout(() => setIsDeleting(true), 1500); // Wait 1.5s before deleting
    } else if (isDeleting && typedText === '') {
      timeout = setTimeout(animate, 500); // Wait 0.5s before typing
    } else {
      timeout = setTimeout(animate, isDeleting ? 20 : 50); // Fast delete, normal type
    }

    return () => clearTimeout(timeout);
  }, [typedText, isDeleting]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-background p-4 relative overflow-hidden">
      {/* Background Decorative Elements */}
      <div className="absolute top-0 left-0 w-full h-full overflow-hidden -z-10 pointer-events-none">
        <div className="absolute top-[-10%] left-[-10%] w-[50%] h-[50%] rounded-full bg-sky-500/10 blur-[120px]" />
        <div className="absolute bottom-[-10%] right-[-10%] w-[50%] h-[50%] rounded-full bg-blue-600/10 blur-[120px]" />
      </div>

      <div className="w-full max-w-6xl grid lg:grid-cols-2 gap-12 lg:gap-24 items-center">
        {/* Left Side - Branding & Features */}
        <div className="space-y-12 order-2 lg:order-1">
          <div className="space-y-6">
            <div className="flex items-center gap-3">
              <LoadingPulse size={40} className="shadow-2xl shadow-sky-500/50" />
              <h1 className="text-5xl font-extrabold tracking-tight text-transparent bg-clip-text bg-linear-to-r from-foreground to-foreground/70 pb-1">
                Dev<span className="text-sky-500">Pulse</span>
              </h1>
            </div>
            <div className="h-20"> {/* Fixed height to prevent layout shift */}
              <p className="text-2xl text-muted-foreground font-light leading-relaxed">
                {typedText}<span className="animate-pulse text-sky-500">|</span>
              </p>
            </div>
          </div>

          <div className="grid gap-8">
            <FeatureItem
              icon={<GitPullRequest className="h-6 w-6" />}
              title="Automated PR Analysis"
              description="Get instant, intelligent feedback on every pull request to maintain code quality."
            />
            {/* Removed Smart Release Cycles as requested */}
            <FeatureItem
              icon={<BarChart3 className="h-6 w-6" />}
              title="Performance Metrics"
              description="Visualize team velocity, review latency, and impact with real-time dashboards."
            />
            <FeatureItem
              icon={<Zap className="h-6 w-6" />}
              title="Architecture Insights"
              description="Detect coupling hotspots and architectural drift before they become debt."
            />
          </div>

          {/* Removed Trusted by 1000+ engineering teams as requested */}
        </div>

        {/* Right Side - Login Card */}
        <div
          className="order-1 lg:order-2 perspective-1000"
          onMouseMove={handleMouseMove}
          onMouseLeave={handleMouseLeave}
          style={{ perspective: '1000px' }}
        >
          <Card
            className="border-none shadow-2xl bg-card/40 backdrop-blur-xl ring-1 ring-white/10 dark:ring-white/5 relative overflow-hidden transition-transform duration-200 ease-out"
            style={{
              transform: `rotateX(${rotate.x}deg) rotateY(${rotate.y}deg)`,
            }}
          >
            {/* Subtle top gradient border effect */}
            <div className="absolute top-0 left-0 w-full h-1 bg-linear-to-r from-transparent via-sky-500/50 to-transparent opacity-50" />

            <CardHeader className="space-y-2 pb-8 text-center">
              <CardTitle className="text-3xl font-bold">Welcome back</CardTitle>
              <CardDescription className="text-base">
                Sign in to your dashboard to continue
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <Button
                className="w-full h-14 text-base font-medium relative group overflow-hidden transition-all hover:shadow-lg hover:shadow-sky-500/25 bg-foreground text-background hover:bg-foreground/90"
                size="lg"
                onClick={handleGithubLogin}
                disabled={isLoading}
              >
                <div className="absolute inset-0 w-full h-full bg-linear-to-r from-transparent via-white/10 to-transparent translate-x-[-100%] group-hover:translate-x-[100%] transition-transform duration-1000" />
                <Github className="mr-2 h-5 w-5" />
                {isLoading ? 'Connecting...' : 'Continue with GitHub'}
                {!isLoading && <ArrowRight className="ml-2 h-4 w-4 opacity-0 group-hover:opacity-100 transition-opacity translate-x-[-5px] group-hover:translate-x-0" />}
              </Button>

              {/* Removed Secure Authentication section as requested */}

              <div className="text-center text-sm text-muted-foreground px-8 pt-4">
                By continuing, you acknowledge that you have read and understood our{' '}
                <a href="/terms" className="hover:text-sky-500 transition-colors font-medium">Terms of Service</a>
                {' '}and{' '}
                <a href="/privacy" className="hover:text-sky-500 transition-colors font-medium">Privacy Policy</a>.
              </div>
            </CardContent>
          </Card>
        </div>
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
    <div className="flex gap-4 items-start group p-4 rounded-2xl hover:bg-white/50 dark:hover:bg-white/5 transition-all duration-300 border border-transparent hover:border-sky-100/50 dark:hover:border-white/5 hover:scale-105">
      <div className="shrink-0 h-12 w-12 bg-linear-to-br from-sky-500/10 to-blue-500/5 rounded-xl flex items-center justify-center text-sky-600 dark:text-sky-400 ring-1 ring-sky-500/20 group-hover:ring-sky-500/40 transition-all shadow-sm group-hover:shadow-md group-hover:shadow-sky-500/10">
        {icon}
      </div>
      <div className="space-y-1">
        <h3 className="font-semibold text-lg text-foreground group-hover:text-sky-600 dark:group-hover:text-sky-400 transition-colors">{title}</h3>
        <p className="text-base text-muted-foreground leading-snug">{description}</p>
      </div>
    </div>
  );
}
