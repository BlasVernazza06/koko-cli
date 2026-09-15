import { healthRouter } from './routers/health';
import { userRouter } from './routers/user';

export const router = {
  health: healthRouter,
  user: userRouter,
};

export type AppRouter = typeof router;
