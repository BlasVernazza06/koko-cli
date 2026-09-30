import { createORPCClient } from '@orpc/client';
import { createORPCReactQueryUtils } from '@orpc/react-query';
import type { AppRouter } from '@repo/api';

export const orpcClient = createORPCClient<AppRouter>({
  baseURL: import.meta.env.VITE_API_URL
    ? `${import.meta.env.VITE_API_URL}/orpc`
    : 'http://localhost:3001/orpc',
});

export const orpc = createORPCReactQueryUtils(orpcClient);
