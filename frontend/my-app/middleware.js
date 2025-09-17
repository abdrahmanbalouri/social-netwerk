import { NextResponse } from "next/server";

export function middleware(req) {
  const url = req.nextUrl.clone();

  const hasSession = req.cookies.get("session")?.value;

  const protectedPaths = ["/dashboard", "/profile", "/home"];

  if (!hasSession && protectedPaths.some(path => url.pathname.startsWith(path))) {
    url.pathname = "/login";
    return NextResponse.redirect(url);
  }

  return NextResponse.next();
}

// Matcher li ghadi tapply middleware 3la had paths
export const config = {
  matcher: ["/dashboard/:path*", "/profile/:path*", "/home/:path*"],
};
