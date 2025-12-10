import { NextRequest, NextResponse } from 'next/server';

/**
 * GET /api/v1/auth/github/connect
 * Initiates GitHub OAuth flow
 */
export async function GET(request: NextRequest) {
  try {
    const clientId = process.env.GITHUB_CLIENT_ID;
    const redirectUri = process.env.GITHUB_REDIRECT_URI || `${process.env.NEXTAUTH_URL}/api/v1/auth/github/callback`;
    
    if (!clientId) {
      return NextResponse.json(
        { 
          success: false, 
          error: { 
            code: 'MISSING_CONFIG', 
            message: 'GitHub OAuth not configured' 
          } 
        },
        { status: 500 }
      );
    }

    // Build GitHub OAuth URL
    const githubAuthUrl = new URL('https://github.com/login/oauth/authorize');
    githubAuthUrl.searchParams.append('client_id', clientId);
    githubAuthUrl.searchParams.append('redirect_uri', redirectUri);
    githubAuthUrl.searchParams.append('scope', 'read:user user:email repo');
    
    // Redirect to GitHub
    return NextResponse.redirect(githubAuthUrl.toString());
  } catch (error) {
    console.error('GitHub OAuth connect error:', error);
    return NextResponse.json(
      { 
        success: false, 
        error: { 
          code: 'OAUTH_ERROR', 
          message: 'Failed to initiate GitHub OAuth' 
        } 
      },
      { status: 500 }
    );
  }
}
