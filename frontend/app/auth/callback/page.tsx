'use client';

import { useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { LoadingPulse } from '@/components/ui/loading-pulse';

export default function CallbackPage() {
    const router = useRouter();
    const searchParams = useSearchParams();

    useEffect(() => {
        const token = searchParams?.get('token');
        const error = searchParams?.get('error');

        if (error) {
            console.error('Authentication error:', error);
            router.push('/login?error=' + error);
            return;
        }

        if (token) {
            // Store token in localStorage
            localStorage.setItem('session_token', token);
            console.log('Token stored successfully');

            // Redirect to dashboard
            router.push('/dashboard');
        } else {
            console.error('No token received');
            router.push('/login?error=no_token');
        }
    }, [router, searchParams]);

    return (
        <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-background via-background to-muted/20">
            <div className="text-center space-y-4">
                <LoadingPulse />
                <div className="space-y-2">
                    <h2 className="text-2xl font-semibold">Completing Sign In...</h2>
                    <p className="text-muted-foreground">
                        Setting up your session
                    </p>
                </div>
            </div>
        </div>
    );
}
