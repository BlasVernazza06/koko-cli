import { RPCHandler } from '@orpc/server/fetch';
import { router, createORPCContext } from '@repo/api';
import type { Context as HonoContext } from 'hono';

const handler = new RPCHandler(router);

export async function handleORPC(c: HonoContext) {
  const ctx = await createORPCContext({ headers: c.req.raw.headers });
  const result = await handler.handle(c.req.raw, {
    prefix: '/orpc',
    context: ctx,
  });
  return result.response ?? c.text('Not Found', 404);
}
