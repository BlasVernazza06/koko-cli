import { betterAuth } from 'better-auth';
import { getAuthConfig } from '@repo/auth';
import { db } from '@repo/db';

export const auth = betterAuth(
  getAuthConfig({
    db,
    secret: process.env.BETTER_AUTH_SECRET || 'secret-placeholder-minimum-32-chars-long-secret-key',
    baseURL: process.env.BETTER_AUTH_URL || 'http://localhost:3000',
    env: {
      googleClientId: process.env.GOOGLE_CLIENT_ID,
      googleClientSecret: process.env.GOOGLE_CLIENT_SECRET,
      githubClientId: process.env.GITHUB_CLIENT_ID,
      githubClientSecret: process.env.GITHUB_CLIENT_SECRET,
    },
  })
);
