import type { inferAsyncReturnType } from '@trpc/server';

export interface CreateContextOptions {
  headers?: Headers;
  session?: unknown;
}

export async function createTRPCContext(opts?: CreateContextOptions) {
  return {
    headers: opts?.headers,
    session: opts?.session,
  };
}

export type Context = inferAsyncReturnType<typeof createTRPCContext>;
