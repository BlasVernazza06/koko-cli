import { auth } from '@/lib/auth';
import { toNextRouteHandler } from 'better-auth/next-js';

export const { GET, POST } = toNextRouteHandler(auth);
