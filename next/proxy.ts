import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

const PUBLIC_PATHS = [
  '/login',
  '/forgot-username',
  '/forgot-password',
  '/reset-password',
]

export default function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl
  const sessionId = request.cookies.get('session_id')?.value

  const isPublicPath =
    PUBLIC_PATHS.includes(pathname) || pathname.startsWith('/oauth')

  if (!isPublicPath && !sessionId) {
    return NextResponse.redirect(new URL('/login', request.url))
  }

  return NextResponse.next()
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - public files (images, etc)
     */
    '/((?!api|_next/static|_next/image|favicon.ico|.*\\.png$|.*\\.jpg$|.*\\.jpeg$|.*\\.gif$|.*\\.svg$).*)',
  ],
}
