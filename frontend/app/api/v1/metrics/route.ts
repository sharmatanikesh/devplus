import { NextRequest, NextResponse } from 'next/server';
import type { EngineeringMetrics } from '@/lib/types';

/**
 * GET /api/v1/metrics
 * Get engineering metrics dashboard data
 */
export async function GET(request: NextRequest) {
  try {
    const searchParams = request.nextUrl.searchParams;
    const from = searchParams.get('from') || new Date(Date.now() - 30 * 86400000).toISOString();
    const to = searchParams.get('to') || new Date().toISOString();
    const repos = searchParams.get('repos')?.split(',');

    console.log('Fetching metrics:', { from, to, repos });

    // TODO: Fetch from database and calculate metrics
    const mockMetrics: EngineeringMetrics = {
      period: {
        from,
        to,
      },
      repositories: repos,
      metrics: {
        prLeadTime: {
          average: 18.5, // hours
          median: 12.3,
          p95: 48.2,
          trend: -12.5, // -12.5% improvement
        },
        reviewLatency: {
          average: 4.2,
          median: 2.8,
          p95: 12.5,
          trend: -8.3,
        },
        mergeRate: {
          percentage: 85.5,
          total: 145,
          merged: 124,
          trend: 3.2,
        },
        throughput: {
          prsPerWeek: 28.5,
          commitsPerWeek: 185.2,
          trend: 15.8,
        },
        codeChurn: {
          average: 245.5,
          trend: -5.2,
        },
      },
      topContributors: [
        {
          username: 'developer1',
          avatarUrl: 'https://github.com/identicons/developer1.png',
          prsCreated: 42,
          prsReviewed: 58,
          commitsCount: 156,
        },
        {
          username: 'developer2',
          avatarUrl: 'https://github.com/identicons/developer2.png',
          prsCreated: 38,
          prsReviewed: 45,
          commitsCount: 142,
        },
      ],
      hotspots: [
        {
          file: 'src/core/analyzer.ts',
          changes: 45,
          authors: 6,
          lastModified: new Date(Date.now() - 86400000).toISOString(),
        },
        {
          file: 'src/api/routes.ts',
          changes: 32,
          authors: 4,
          lastModified: new Date(Date.now() - 172800000).toISOString(),
        },
      ],
    };

    return NextResponse.json({
      success: true,
      data: mockMetrics,
    });
  } catch (error) {
    console.error('Metrics error:', error);
    return NextResponse.json(
      {
        success: false,
        error: {
          code: 'METRICS_ERROR',
          message: 'Failed to fetch metrics',
        },
      },
      { status: 500 }
    );
  }
}
