'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
  LayoutDashboard,
  GitBranch,
  GitPullRequest,
  BarChart3,
  Package,
  LogOut,
  Zap,
  Heart,
  Github,
} from 'lucide-react';
import { NAVIGATION_ITEMS } from '@/lib/constants';
import { apiClient } from '@/lib/api-client';
import { User } from '@/lib/types';
import { motion, AnimatePresence } from 'framer-motion';

const iconMap = {
  LayoutDashboard,
  GitBranch,
  GitPullRequest,
  BarChart3,
  Package,
};

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    async function fetchUser() {
      try {
        const res = await apiClient.auth.me();
        if (res.success && res.data) {
          setUser(res.data);
        }
      } catch (error) {
        console.error('Failed to fetch user:', error);
      }
    }
    fetchUser();

    const handleScroll = () => {
      setScrolled(window.scrollY > 10);
    };
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  const handleLogout = async () => {
    try {
      await apiClient.auth.logout();
      router.push('/login');
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  return (
    <div className="min-h-screen bg-background relative flex flex-col">
      {/* Background Gradients */}
      <div className="fixed inset-0 w-full h-full overflow-hidden -z-10 pointer-events-none">
        <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] rounded-full bg-sky-500/5 blur-[120px]" />
        <div className="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] rounded-full bg-indigo-600/5 blur-[120px]" />
      </div>

      {/* Top Navigation Bar */}
      <header
        className={cn(
          "fixed top-0 left-0 right-0 z-50 transition-all duration-200 border-b border-transparent",
          scrolled ? "bg-background/80 backdrop-blur-md border-border/50 shadow-sm" : "bg-transparent"
        )}
      >
        <div className="container mx-auto px-4 h-16 flex items-center justify-between">
          {/* Logo */}
          <Link href="/dashboard" className="flex items-center gap-2">
            <motion.div
              initial={{ scale: 0.8, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              className="h-8 w-8 bg-sky-500 rounded-lg flex items-center justify-center shadow-lg shadow-sky-500/20"
            >
              <Zap className="h-5 w-5 text-white" />
            </motion.div>
            <span className="font-bold text-lg tracking-tight bg-clip-text text-transparent bg-linear-to-r from-foreground to-foreground/80">DevPulse</span>
          </Link>

          {/* Navigation */}
          <nav className="hidden md:flex items-center gap-1">
            {NAVIGATION_ITEMS.map((item) => {
              const Icon = iconMap[item.icon as keyof typeof iconMap];
              const isActive = pathname === item.href;

              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className="relative group"
                >
                  <div
                    className={cn(
                      'flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium transition-all duration-200',
                      isActive
                        ? 'text-sky-600 dark:text-sky-400'
                        : 'text-muted-foreground hover:text-foreground hover:bg-muted/50'
                    )}
                  >
                    {isActive && (
                      <motion.div
                        layoutId="activeNav"
                        className="absolute inset-0 bg-sky-100 dark:bg-sky-900/20 rounded-full"
                        initial={false}
                        transition={{ type: "spring", stiffness: 300, damping: 30 }}
                      />
                    )}
                    <Icon className={cn("h-4 w-4 relative z-10 transition-colors", isActive ? "text-sky-600 dark:text-sky-400" : "group-hover:text-sky-500")} />
                    <span className="relative z-10">{item.title}</span>
                  </div>
                </Link>
              );
            })}
          </nav>

          {/* User Menu */}
          <div className="flex items-center gap-4">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="relative h-9 w-9 rounded-full ring-2 ring-transparent hover:ring-sky-200 transition-all p-0 overflow-hidden">
                  <Avatar className="h-9 w-9">
                    <AvatarImage src={user?.avatarUrl || "https://github.com/github.png"} />
                    <AvatarFallback className="bg-sky-100 text-sky-700">{user?.username?.substring(0, 2).toUpperCase() || "DV"}</AvatarFallback>
                  </Avatar>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-56" sideOffset={8}>
                <DropdownMenuLabel>
                  <div className="flex flex-col space-y-1">
                    <p className="text-sm font-medium leading-none">{user?.username}</p>
                    <p className="text-xs leading-none text-muted-foreground">{user?.email}</p>
                  </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  className="text-red-600 focus:text-red-600 focus:bg-red-50 dark:focus:bg-red-950/20 cursor-pointer"
                  onClick={handleLogout}
                >
                  <LogOut className="mr-2 h-4 w-4" />
                  Sign out
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </header>

      {/* Main content */}
      <main className="flex-1 pt-16">
        <div className="container mx-auto p-6">
          <AnimatePresence mode="wait">
            <motion.div
              key={pathname}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              transition={{ duration: 0.3 }}
            >
              {children}
            </motion.div>
          </AnimatePresence>
        </div>
      </main>

      <footer className="py-6 border-t border-border/50 bg-background/50 backdrop-blur-sm">
        <div className="container mx-auto px-6 grid gap-4 md:grid-cols-3 items-center">
          {/* Empty left column for balance */}
          <div className="hidden md:block" />

          {/* Center Content */}
          <div className="flex flex-col md:flex-row items-center justify-center gap-2 text-sm text-muted-foreground whitespace-nowrap">
            <span className="font-bold text-sky-500 animate-pulse">DevPulse</span>
            <span className="hidden md:inline">•</span>
            <p className="flex items-center gap-1">
              Made with <Heart className="h-3 w-3 text-red-500 fill-red-500 animate-pulse" /> from team Shaitan
            </p>
          </div>

          {/* Right Content */}
          <div className="flex justify-center md:justify-end">
            <a
              href="https://github.com/sharmatanikesh/devplus"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 text-xs text-muted-foreground hover:text-sky-600 transition-colors"
            >
              <Github className="h-4 w-4" />
              See code
            </a>
          </div>
        </div>
      </footer>
    </div>
  );
}
