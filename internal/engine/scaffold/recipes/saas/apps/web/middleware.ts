import { NextRequest, NextResponse } from 'next/server';

export async function middleware(_request: NextRequest) {
  // 💡 Tutor Tip: Para "desactivar" un middleware sin borrar el archivo, 
  // simplemente retornamos NextResponse.next() al inicio.
  return NextResponse.next();

  /* COMENTADO TEMPORALMENTE PARA DESARROLLO
  const { pathname } = request.nextUrl;
  const sessionCookie = request.cookies.get('better-auth.session_token')?.value;

  // Redirección para la página inicial/raíz
  if (pathname === '/') {
    if (sessionCookie) {
      return NextResponse.redirect(new URL('/auth/login', request.url));
    }
    return NextResponse.redirect(new URL('/dashboard', request.url));
  }

  const protectedRoutes = ['/dashboard'];

  if (protectedRoutes.some((route) => pathname.startsWith(route))) {
    if (!sessionCookie) {
      const loginUrl = new URL('/auth/login', request.url);
      loginUrl.searchParams.set(
        'callbackUrl',
        request.nextUrl.pathname + request.nextUrl.search,
      );
      return NextResponse.redirect(loginUrl);
    }
  }

  return NextResponse.next();
  */
}

export const config = {
  matcher: ['/', '/dashboard', '/dashboard/:path*'],
};

