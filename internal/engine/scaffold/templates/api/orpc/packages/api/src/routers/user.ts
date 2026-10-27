import { z } from 'zod';
import { publicProcedure } from '../orpc';

const UserSchema = z.object({
  id: z.string(),
  name: z.string(),
  email: z.string().email(),
});

export const userRouter = {
  list: publicProcedure
    .route({
      method: 'GET',
      path: '/users',
      summary: 'List users',
    })
    .output(z.array(UserSchema))
    .handler(async () => {
      return [
        { id: '1', name: 'Alice', email: 'alice@example.com' },
        { id: '2', name: 'Bob', email: 'bob@example.com' },
      ];
    }),
  byId: publicProcedure
    .route({
      method: 'GET',
      path: '/users/{id}',
      summary: 'Get user by ID',
    })
    .input(z.object({ id: z.string() }))
    .output(UserSchema)
    .handler(async ({ input }) => {
      return {
        id: input.id,
        name: `User ${input.id}`,
        email: `user${input.id}@example.com`,
      };
    }),
};
