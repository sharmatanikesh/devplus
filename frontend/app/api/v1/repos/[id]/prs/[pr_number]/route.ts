import { NextRequest, NextResponse } from 'next/server';
import type { PullRequest } from '@/lib/types';

/**
 * GET /api/v1/repos/:id/prs/:pr_number
 * Get detailed PR information including AI analysis
 */
export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ id: string; pr_number: string }> }
) {
  try {
    const { id, pr_number } = await params;

    console.log(`Fetching PR #${pr_number} for repository: ${id}`);

    // TODO: Fetch from database
    const mockPR: PullRequest = {
      id: '1',
      githubId: 789,
      repositoryId: id,
      number: parseInt(pr_number),
      title: 'Add new feature: AI-powered code review',
      description: 'This PR implements the core AI review functionality with comprehensive testing',
      state: 'open',
      author: {
        username: 'developer1',
        avatarUrl: 'https://github.com/identicons/developer1.png',
      },
      baseBranch: 'main',
      headBranch: 'feature/ai-review',
      url: `https://github.com/yourorg/repo/pull/${pr_number}`,
      filesChanged: 15,
      additions: 450,
      deletions: 120,
      commits: 8,
      reviewStatus: 'pending',
      aiAnalysisStatus: 'completed',
      aiSummary: 'This PR introduces AI-powered code review features. Code quality is good with no major issues found.',
      aiDecision: 'approve',
      aiComments: [
        {
          id: '1',
          file: 'src/analyzer.ts',
          line: 45,
          body: 'Consider adding error handling for this async operation',
          severity: 'warning',
          createdAt: new Date().toISOString(),
        },
      ],
      createdAt: new Date(Date.now() - 86400000 * 2).toISOString(),
      updatedAt: new Date(Date.now() - 3600000).toISOString(),
    };

    return NextResponse.json({
      success: true,
      data: mockPR,
    });
  } catch (error) {
    console.error('Get PR error:', error);
    return NextResponse.json(
      {
        success: false,
        error: {
          code: 'PR_GET_ERROR',
          message: 'Failed to fetch pull request',
        },
      },
      { status: 500 }
    );
  }
}
