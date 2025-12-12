import { NextRequest, NextResponse } from 'next/server';

// Validate required environment variables at module load
const NEXTAUTH_URL = process.env.NEXTAUTH_URL;
if (!NEXTAUTH_URL) {
  throw new Error('NEXTAUTH_URL environment variable is required but not set');
}

/**
 * GET /api/v1/auth/github/callback
 * Handles GitHub OAuth callback
 */
export async function GET(request: NextRequest) {
  try {
    const searchParams = request.nextUrl.searchParams;
    const code = searchParams.get('code');
    const error = searchParams.get('error');

    if (error) {
      return NextResponse.redirect(`${NEXTAUTH_URL}/login?error=${error}`);
    }

    if (!code) {
      return NextResponse.redirect(`${NEXTAUTH_URL}/login?error=no_code`);
    }

    // Exchange code for access token
    const tokenResponse = await fetch('https://github.com/login/oauth/access_token', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
      body: JSON.stringify({
        client_id: process.env.GITHUB_CLIENT_ID,
        client_secret: process.env.GITHUB_CLIENT_SECRET,
        code,
      }),
    });

    // Check if token exchange was successful
    if (!tokenResponse.ok) {
      const errorText = await tokenResponse.text();
      console.error('GitHub token exchange failed:', {
        status: tokenResponse.status,
        statusText: tokenResponse.statusText,
        body: errorText,
      });
      return NextResponse.redirect(`${NEXTAUTH_URL}/login?error=token_exchange_failed`);
    }

    const tokenData = await tokenResponse.json();

    if (tokenData.error || !tokenData.access_token) {
      console.error('GitHub token exchange error:', tokenData);
      return NextResponse.redirect(`${NEXTAUTH_URL}/login?error=token_exchange_failed`);
    }

    // Get user data from GitHub
    const userResponse = await fetch('https://api.github.com/user', {
      headers: {
        'Authorization': `Bearer ${tokenData.access_token}`,
        'Accept': 'application/json',
      },
    });

    if (!userResponse.ok) {
      console.error('Failed to fetch GitHub user:', {
        status: userResponse.status,
        statusText: userResponse.statusText,
      });
      return NextResponse.redirect(`${NEXTAUTH_URL}/login?error=user_fetch_failed`);
    }

    const userData = await userResponse.json();

    // TODO: Save user to database and create session
    // For now, we'll redirect to dashboard with token in URL (not secure for production)
    // In production, you should:
    // 1. Save user to database
    // 2. Create a session/JWT
    // 3. Set secure HTTP-only cookie
    
    const response = NextResponse.redirect(`${NEXTAUTH_URL}/dashboard`);
    
    // Set auth token in cookie (simplified for MVP)
    response.cookies.set('devpulse_token', tokenData.access_token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      maxAge: 60 * 60 * 24 * 7, // 7 days
    });

    return response;
  } catch (error) {
    console.error('GitHub OAuth callback error:', error);
    return NextResponse.redirect(`${NEXTAUTH_URL}/login?error=callback_failed`);
  }
}
