import { NextRequest, NextResponse } from 'next/server';
import type { JobStatus } from '@/lib/types';

/**
 * POST /api/v1/repos/:id/prs/:pr_number/analyze
 * Trigger AI analysis for a pull request
 */
export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ id: string; pr_number: string }> }
) {
  try {
    const { id, pr_number } = await params;

    console.log(`Triggering AI analysis for PR #${pr_number} in repository: ${id}`);

    // TODO: Enqueue Kestra job for AI analysis
    const jobId = `job_${Date.now()}`;

    const job: JobStatus = {
      id: jobId,
      type: 'pr_analysis',
      status: 'queued',
      progress: 0,
      createdAt: new Date().toISOString(),
    };

    return NextResponse.json({
      success: true,
      data: job,
      message: 'AI analysis job queued successfully',
    });
  } catch (error) {
    console.error('PR analyze error:', error);
    return NextResponse.json(
      {
        success: false,
        error: {
          code: 'ANALYZE_ERROR',
          message: 'Failed to trigger AI analysis',
        },
      },
      { status: 500 }
    );
  }
}
