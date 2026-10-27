import { cookies } from 'next/headers';

const API_URL =
  process.env.API_URL ||
  process.env.NEXT_PUBLIC_API_URL ||
  'http://localhost:3001';

export interface SessionUser {
  id: string;
  name: string;
  email: string;
  image: string | null;
}

export async function getSession() {
  const cookieStore = await cookies();
  const allCookies = cookieStore.getAll();

  if (allCookies.length === 0) return null;

  const cookieHeader = allCookies.map((c) => `${c.name}=${c.value}`).join('; ');

  try {
    const res = await fetch(`${API_URL}/api/auth/get-session`, {
      headers: {
        cookie: cookieHeader,
      },
      cache: 'no-store',
    });

    if (!res.ok) return null;

    const data = await res.json();
    return data as {
      session: {
        id: string;
        userId: string;
        expiresAt: string;
        token: string;
      };
      user: SessionUser;
    } | null;
  } catch (err) {
    console.error('[auth-session] Error:', err);
    return null;
  }
}
