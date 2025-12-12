import { NextRequest, NextResponse } from 'next/server';
import type { Repository, PaginatedResponse } from '@/lib/types';

/**
 * GET /api/v1/repos
 * List connected repositories
 */
export async function GET(request: NextRequest) {
  try {
    const searchParams = request.nextUrl.searchParams;
    const page = parseInt(searchParams.get('page') || '1');
    const limit = parseInt(searchParams.get('limit') || '20');

    // TODO: Fetch from database
    // For MVP, return mock data
    const mockRepos: Repository[] = [
      {
        id: '1',
        githubId: 123456,
        name: 'devpulse-frontend',
        fullName: 'yourorg/devpulse-frontend',
        owner: 'yourorg',
        description: 'DevPulse frontend application',
        url: 'https://github.com/yourorg/devpulse-frontend',
        defaultBranch: 'main',
        isPrivate: false,
        language: 'TypeScript',
        stars: 42,
        forks: 5,
        syncStatus: 'synced',
        lastSyncedAt: new Date().toISOString(),
        webhookConfigured: true,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
    ];

    const response: PaginatedResponse<Repository> = {
      items: mockRepos,
      pagination: {
        page,
        limit,
        total: mockRepos.length,
        totalPages: 1,
      },
    };

    return NextResponse.json({
      success: true,
      data: response,
    });
  } catch (error) {
    console.error('List repos error:', error);
    return NextResponse.json(
      {
        success: false,
        error: {
          code: 'REPOS_LIST_ERROR',
          message: 'Failed to fetch repositories',
        },
      },
      { status: 500 }
    );
  }
}
