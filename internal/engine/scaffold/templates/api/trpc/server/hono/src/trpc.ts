import { fetchRequestHandler } from '@trpc/server/adapters/fetch';
import { appRouter, createTRPCContext } from '@repo/api';
import type { Context as HonoContext } from 'hono';

export async function handleTRPC(c: HonoContext) {
  return fetchRequestHandler({
    endpoint: '/trpc',
    req: c.req.raw,
    router: appRouter,
    createContext: () => createTRPCContext({ headers: c.req.raw.headers }),
  });
}
