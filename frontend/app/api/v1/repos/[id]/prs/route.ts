import { NextRequest, NextResponse } from 'next/server';
import type { PullRequest, PaginatedResponse } from '@/lib/types';

/**
 * GET /api/v1/repos/:id/prs
 * List pull requests for a repository
 */
export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id } = await params;
    const searchParams = request.nextUrl.searchParams;
    const page = parseInt(searchParams.get('page') || '1');
    const limit = parseInt(searchParams.get('limit') || '20');
    const state = searchParams.get('state') || 'open'; // open, closed, merged, all

    console.log(`Fetching PRs for repository: ${id}, state: ${state}`);

    // TODO: Fetch from database
    // For MVP, return mock data
    const mockPRs: PullRequest[] = [
      {
        id: '1',
        githubId: 789,
        repositoryId: id,
        number: 42,
        title: 'Add new feature: AI-powered code review',
        description: 'This PR implements the core AI review functionality',
        state: 'open',
        author: {
          username: 'developer1',
          avatarUrl: 'https://github.com/identicons/developer1.png',
        },
        baseBranch: 'main',
        headBranch: 'feature/ai-review',
        url: 'https://github.com/yourorg/repo/pull/42',
        filesChanged: 15,
        additions: 450,
        deletions: 120,
        commits: 8,
        reviewStatus: 'pending',
        aiAnalysisStatus: 'completed',
        aiSummary: 'This PR introduces AI-powered code review features. Code quality is good with no major issues found.',
        aiDecision: 'approve',
        createdAt: new Date(Date.now() - 86400000 * 2).toISOString(),
        updatedAt: new Date(Date.now() - 3600000).toISOString(),
      },
    ];

    const filteredPRs = state === 'all' ? mockPRs : mockPRs.filter(pr => pr.state === state);

    const response: PaginatedResponse<PullRequest> = {
      items: filteredPRs,
      pagination: {
        page,
        limit,
        total: filteredPRs.length,
        totalPages: Math.ceil(filteredPRs.length / limit),
      },
    };

    return NextResponse.json({
      success: true,
      data: response,
    });
  } catch (error) {
    console.error('List PRs error:', error);
    return NextResponse.json(
      {
        success: false,
        error: {
          code: 'PRS_LIST_ERROR',
          message: 'Failed to fetch pull requests',
        },
      },
      { status: 500 }
    );
  }
}
