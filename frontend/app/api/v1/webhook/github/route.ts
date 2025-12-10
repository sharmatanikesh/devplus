import { NextRequest, NextResponse } from 'next/server';

/**
 * POST /api/v1/webhook/github
 * Handles GitHub webhook events (PR, push, checks)
 */
export async function POST(request: NextRequest) {
  try {
    const event = request.headers.get('x-github-event');
    const signature = request.headers.get('x-hub-signature-256');
    const delivery = request.headers.get('x-github-delivery');

    // TODO: Verify webhook signature for security
    // const isValid = verifyWebhookSignature(await request.text(), signature);
    // if (!isValid) {
    //   return NextResponse.json({ error: 'Invalid signature' }, { status: 401 });
    // }

    const payload = await request.json();

    console.log(`Received GitHub webhook: ${event}`, {
      delivery,
      repository: payload.repository?.full_name,
      action: payload.action,
    });

    // TODO: Forward to backend/Kestra for processing
    // In a real implementation, you would:
    // 1. Verify the webhook signature
    // 2. Enqueue the event for processing (Kestra, Bull, etc.)
    // 3. Return quickly to avoid timeout

    // For MVP, we'll just log and acknowledge
    switch (event) {
      case 'pull_request':
        console.log('PR Event:', {
          action: payload.action,
          pr: payload.pull_request?.number,
          title: payload.pull_request?.title,
        });
        break;
      
      case 'push':
        console.log('Push Event:', {
          ref: payload.ref,
          commits: payload.commits?.length,
        });
        break;
      
      case 'check_run':
      case 'check_suite':
        console.log('Check Event:', {
          action: payload.action,
          conclusion: payload.check_run?.conclusion || payload.check_suite?.conclusion,
        });
        break;
    }

    return NextResponse.json({ 
      success: true, 
      message: 'Webhook received',
      event,
      delivery,
    });
  } catch (error) {
    console.error('Webhook processing error:', error);
    return NextResponse.json(
      { 
        success: false, 
        error: { 
          code: 'WEBHOOK_ERROR', 
          message: 'Failed to process webhook' 
        } 
      },
      { status: 500 }
    );
  }
}
