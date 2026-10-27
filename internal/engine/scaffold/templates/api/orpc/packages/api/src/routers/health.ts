import { z } from 'zod';
import { publicProcedure } from '../orpc';

export const healthRouter = {
  check: publicProcedure
    .route({
      method: 'GET',
      path: '/health',
      summary: 'Health Check',
      description: 'Returns server health and uptime status',
    })
    .output(
      z.object({
        status: z.string(),
        timestamp: z.string(),
      })
    )
    .handler(async () => {
      return {
        status: 'ok',
        timestamp: new Date().toISOString(),
      };
    }),
};
