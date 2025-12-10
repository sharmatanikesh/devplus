import { NextRequest, NextResponse } from 'next/server';
import type { JobStatus } from '@/lib/types';

/**
 * POST /api/v1/repos/:id/release
 * Generate changelog and create release
 */
export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id } = await params;
    const body = await request.json();

    const { version, tagName, targetCommitish, isDraft, isPrerelease } = body;

    console.log(`Generating release for repository: ${id}`, {
      version,
      tagName,
    });

    // TODO: Enqueue Kestra job for release generation
    const jobId = `job_${Date.now()}`;

    const job: JobStatus = {
      id: jobId,
      type: 'release_generation',
      status: 'queued',
      progress: 0,
      createdAt: new Date().toISOString(),
    };

    return NextResponse.json({
      success: true,
      data: job,
      message: 'Release generation job queued successfully',
    });
  } catch (error) {
    console.error('Release generation error:', error);
    return NextResponse.json(
      {
        success: false,
        error: {
          code: 'RELEASE_ERROR',
          message: 'Failed to generate release',
        },
      },
      { status: 500 }
    );
  }
}
