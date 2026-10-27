import { z } from 'zod';
import { router, publicProcedure } from '../trpc';

export const userRouter = router({
  all: publicProcedure.query(() => {
    return [
      { id: '1', name: 'Alice', role: 'admin' },
      { id: '2', name: 'Bob', role: 'member' },
    ];
  }),
  byId: publicProcedure
    .input(z.object({ id: z.string() }))
    .query(({ input }) => {
      return {
        id: input.id,
        name: `User ${input.id}`,
        role: 'member',
      };
    }),
});
