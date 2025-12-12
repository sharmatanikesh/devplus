import { useState, useEffect, useCallback } from 'react';
import { apiClient } from '@/lib/api-client';

export interface Repository {
    id: string; // Updated to string for UUID
    github_repo_id: number;
    name: string;
    owner: string;
    url: string;
    updated_at: string;
}

export function useRepositories() {
    const [repositories, setRepositories] = useState<Repository[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const refetch = useCallback(async () => {
        setIsLoading(true);
        try {
            const response = await apiClient.repos.list();
            if (response.success && response.data) {
                setRepositories(response.data as Repository[]);
                setError(null);
            } else {
                setError(response.error?.message || 'Failed to fetch repositories');
            }
        } catch (err) {
            setError('An unexpected error occurred');
        } finally {
            setIsLoading(false);
        }
    }, []);

    useEffect(() => {
        refetch();
    }, [refetch]);

    return { repositories, isLoading, error, refetch };
}
