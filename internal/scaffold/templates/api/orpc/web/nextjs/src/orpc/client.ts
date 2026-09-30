import { createORPCClient } from '@orpc/client';
import { createORPCReactQueryUtils } from '@orpc/react-query';
import type { AppRouter } from '@repo/api';

export const orpcClient = createORPCClient<AppRouter>({
  baseURL: typeof window !== 'undefined' ? '/api/orpc' : 'http://localhost:3000/api/orpc',
});

export const orpc = createORPCReactQueryUtils(orpcClient);
