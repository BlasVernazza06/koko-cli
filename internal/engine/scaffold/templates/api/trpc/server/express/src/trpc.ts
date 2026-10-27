import * as trpcExpress from '@trpc/server/adapters/express';
import { appRouter, createTRPCContext } from '@repo/api';

export const trpcMiddleware = trpcExpress.createExpressMiddleware({
  router: appRouter,
  createContext: () => createTRPCContext({ session: null }),
});
