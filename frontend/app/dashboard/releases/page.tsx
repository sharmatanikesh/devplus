'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Package, Calendar, Tag, GitCommit } from 'lucide-react';

export default function ReleasesPage() {
  // TODO: Fetch releases from API
  const mockReleases = [
    {
      id: '1',
      version: 'v1.2.0',
      title: 'Enhanced AI Analysis',
      tagName: 'v1.2.0',
      publishedAt: '2024-12-08T10:00:00Z',
      isDraft: false,
      isPrerelease: false,
      changes: [
        { type: 'feature', description: 'Added support for multi-file PR analysis' },
        { type: 'improvement', description: 'Improved metrics calculation performance' },
        { type: 'bugfix', description: 'Fixed webhook signature validation' },
      ],
    },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Release Notes</h1>
          <p className="text-muted-foreground">
            Auto-generated changelogs and version history
          </p>
        </div>
        <Button>
          <Package className="mr-2 h-4 w-4" />
          Generate Release
        </Button>
      </div>

      {/* Releases List */}
      <div className="space-y-4">
        {mockReleases.map((release) => (
          <Card key={release.id}>
            <CardHeader>
              <div className="flex items-start justify-between">
                <div className="space-y-1">
                  <div className="flex items-center gap-2">
                    <Tag className="h-5 w-5 text-muted-foreground" />
                    <CardTitle className="text-xl">{release.title}</CardTitle>
                    <Badge>{release.version}</Badge>
                    {release.isDraft && <Badge variant="outline">Draft</Badge>}
                    {release.isPrerelease && <Badge variant="outline">Pre-release</Badge>}
                  </div>
                  <CardDescription className="flex items-center gap-2">
                    <Calendar className="h-4 w-4" />
                    Released on {new Date(release.publishedAt).toLocaleDateString()}
                  </CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* Changes */}
              <div className="space-y-2">
                <h3 className="font-semibold text-sm">What&apos;s Changed</h3>
                <div className="space-y-2">
                  {release.changes.map((change, index) => (
                    <div key={index} className="flex items-start gap-2 text-sm">
                      <GitCommit className="h-4 w-4 text-muted-foreground mt-0.5" />
                      <div className="flex-1">
                        <ChangeTypeBadge type={change.type} />
                        <span className="ml-2">{change.description}</span>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}

function ChangeTypeBadge({ type }: { type: string }) {
  const config = {
    feature: { color: 'bg-green-100 text-green-800', label: 'Feature' },
    bugfix: { color: 'bg-red-100 text-red-800', label: 'Bug Fix' },
    improvement: { color: 'bg-blue-100 text-blue-800', label: 'Improvement' },
    breaking: { color: 'bg-orange-100 text-orange-800', label: 'Breaking' },
    docs: { color: 'bg-purple-100 text-purple-800', label: 'Docs' },
    chore: { color: 'bg-gray-100 text-gray-800', label: 'Chore' },
  };

  const typeConfig = config[type as keyof typeof config] || config.chore;

  return (
    <Badge variant="secondary" className={`${typeConfig.color} text-xs`}>
      {typeConfig.label}
    </Badge>
  );
}
