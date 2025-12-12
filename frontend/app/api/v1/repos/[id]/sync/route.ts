import { NextRequest, NextResponse } from 'next/server';

/**
 * POST /api/v1/repos/:id/sync
 * Manually trigger repository sync
 */
export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id } = await params;

    // TODO: Trigger Kestra job or backend sync
    console.log(`Triggering sync for repository: ${id}`);

    return NextResponse.json({
      success: true,
      data: {
        id,
        status: 'syncing',
        message: 'Sync started',
      },
    });
  } catch (error) {
    console.error('Repo sync error:', error);
    return NextResponse.json(
      {
        success: false,
        error: {
          code: 'SYNC_ERROR',
          message: 'Failed to trigger sync',
        },
      },
      { status: 500 }
    );
  }
}
