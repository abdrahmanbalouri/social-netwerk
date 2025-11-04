import { NextRequest, NextResponse } from 'next/server';

export default async function middleware(request: NextRequest) {
  const pathname = request.nextUrl.pathname;
  const API ="http://localhost:8080/";



  try {
    const response = await fetch(`${API}api/me`, {
      headers: {
        Cookie: request.headers.get('cookie') || ""
      }
    });
    

    const ok  = response.ok;

    if (!ok) {
      if (pathname === "/login" || pathname === "/register") {
        return NextResponse.next();
      }
      return NextResponse.redirect(new URL("/login", request.url));
    }
    if (pathname === "/login" || pathname === "/register") {
      return NextResponse.redirect(new URL("/home", request.url));
    }

    return NextResponse.next();
  } catch (err) {
    console.error("Middleware Error:", err);
    if (pathname === "/login" || pathname === "/register") {
      return NextResponse.next();
    }
    return NextResponse.redirect(new URL("/login", request.url));
  }
}

export const config = {
  matcher: [
    '/home',
    '/login',
    '/register',
    '/chat/:path*',
    '/chat',
    '/profile/:path*',
    '/groups/:path*',
    '/games',
    '/Gallery',
    '/Follow/:path*',
    '/Events',
    '/Watch',
    
  ],
};