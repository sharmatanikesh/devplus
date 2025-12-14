import Link from 'next/link';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { GitBranch, ExternalLink, Clock } from 'lucide-react';
import { Repository } from '@/hooks/use-repositories';
import { formatDistanceToNow } from 'date-fns';

interface RepositoryListProps {
    repositories: Repository[];
}

export function RepositoryList({ repositories }: RepositoryListProps) {
    if (repositories.length === 0) {
        return (
            <Card>
                <CardContent className="flex flex-col items-center justify-center py-10 space-y-4">
                    <div className="bg-muted p-4 rounded-full">
                        <GitBranch className="h-8 w-8 text-muted-foreground" />
                    </div>
                    <div className="text-center space-y-1">
                        <p className="text-lg font-medium">No repositories found</p>
                        <p className="text-sm text-muted-foreground">
                            We couldn't find any repositories linked to your account.
                        </p>
                    </div>
                </CardContent>
            </Card>
        );
    }

    return (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4">
            {repositories.map((repo) => (
                <Card key={repo.id} className="flex flex-col hover:shadow-xl hover:scale-[1.01] transition-all duration-200 border-border/50 bg-card/50 hover:border-sky-200 dark:hover:border-sky-800">
                    <CardHeader>
                        <div className="flex items-start justify-between gap-2">
                            <div className="space-y-1">
                                <CardTitle className="text-base font-semibold leading-none">
                                    <Link
                                        href={`/repositories/${repo.id}`}
                                        className="hover:underline hover:text-sky-600 transition-colors"
                                    >
                                        {repo.name}
                                    </Link>
                                </CardTitle>
                                <CardDescription className="text-xs">
                                    {repo.owner}
                                </CardDescription>
                            </div>
                            <Badge variant="outline" className="shrink-0 bg-sky-50 text-sky-700 dark:bg-sky-900/30 dark:text-sky-300 border-sky-100 dark:border-sky-800">
                                Public
                            </Badge>
                        </div>
                    </CardHeader>
                    <CardContent className="mt-auto pt-0">
                        <div className="flex items-center justify-between text-sm text-muted-foreground mt-4">
                            <div className="flex items-center gap-1 text-xs">
                                <Clock className="h-3 w-3" />
                                {repo.updated_at && (
                                    <span>
                                        Updated {formatDistanceToNow(new Date(repo.updated_at), { addSuffix: true })}
                                    </span>
                                )}
                            </div>
                            <Button variant="ghost" size="icon" className="h-8 w-8 hover:text-sky-600 hover:bg-sky-50 dark:hover:bg-sky-900/20" asChild>
                                <a href={repo.url} target="_blank" rel="noopener noreferrer">
                                    <ExternalLink className="h-4 w-4" />
                                    <span className="sr-only">View on GitHub</span>
                                </a>
                            </Button>
                        </div>
                    </CardContent>
                </Card>
            ))}
        </div>
    );
}
