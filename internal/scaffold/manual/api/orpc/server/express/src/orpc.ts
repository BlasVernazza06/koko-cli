import { RPCHandler } from '@orpc/server/fetch';
import { router, createORPCContext } from '@repo/api';
import type { Request, Response } from 'express';

const handler = new RPCHandler(router);

export async function orpcMiddleware(req: Request, res: Response) {
  const ctx = await createORPCContext({ headers: new Headers(req.headers as any) });
  const protocol = req.protocol;
  const host = req.get('host');
  const url = `${protocol}://${host}${req.originalUrl}`;
  const standardReq = new Request(url, {
    method: req.method,
    headers: new Headers(req.headers as any),
  });

  const result = await handler.handle(standardReq, {
    prefix: '/orpc',
    context: ctx,
  });

  if (result.matched && result.response) {
    res.status(result.response.status);
    result.response.headers.forEach((value, key) => res.setHeader(key, value));
    const text = await result.response.text();
    res.send(text);
  } else {
    res.status(404).send('Not Found');
  }
}
