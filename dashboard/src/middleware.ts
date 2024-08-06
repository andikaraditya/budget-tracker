import { NextResponse, NextRequest } from "next/server";
import { cookies } from "next/headers";

export function middleware(request: NextRequest) {
  const token = cookies().has("token");

  if (request.nextUrl.pathname.startsWith("/_next")) {
    return NextResponse.next();
  }

  const toAuth =
    request.nextUrl.pathname.startsWith("/login") ||
    request.nextUrl.pathname.startsWith("/register");

  if (!token && !toAuth) {
    return NextResponse.redirect(new URL("/login", request.url));
  } else if (token && toAuth) {
    return NextResponse.redirect(new URL("/home", request.url));
  }
  return NextResponse.next();
}
