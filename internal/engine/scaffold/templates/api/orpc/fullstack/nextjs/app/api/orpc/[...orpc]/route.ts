import { RPCHandler } from '@orpc/server/fetch';
import { router, createORPCContext } from '@repo/api';

const handler = new RPCHandler(router);

export async function GET(req: Request) {
  const ctx = await createORPCContext({ headers: req.headers });
  const result = await handler.handle(req, {
    prefix: '/api/orpc',
    context: ctx,
  });
  return result.response ?? new Response('Not Found', { status: 404 });
}

export async function POST(req: Request) {
  const ctx = await createORPCContext({ headers: req.headers });
  const result = await handler.handle(req, {
    prefix: '/api/orpc',
    context: ctx,
  });
  return result.response ?? new Response('Not Found', { status: 404 });
}
