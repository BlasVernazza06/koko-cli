export interface CreateContextOptions {
  headers?: Headers;
  session?: unknown;
}

export async function createORPCContext(opts?: CreateContextOptions) {
  return {
    headers: opts?.headers,
    session: opts?.session,
  };
}

export type Context = Awaited<ReturnType<typeof createORPCContext>>;
