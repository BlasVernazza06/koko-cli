import { drizzleAdapter } from 'better-auth/adapters/drizzle';
import { twoFactor } from 'better-auth/plugins';

import { authSchema, db as defaultDb } from '@repo/db';

export const getAuthConfig = (options: {
  db?: any;
  secret: string;
  baseURL: string;
  env: {
    googleClientId?: string;
    googleClientSecret?: string;
    githubClientId?: string;
    githubClientSecret?: string;
  };
}) => ({
  secret: options.secret,
  baseURL: options.baseURL,
  database: drizzleAdapter(options.db || defaultDb, {
    provider: 'pg',
    schema: authSchema,
  }),
  emailAndPassword: {
    enabled: true,
  },
  session: {
    expiresIn: 60 * 60 * 24 * 7, // 7 days
    updateAge: 60 * 60 * 24, // 1 day
    cookieCache: {
      enabled: true,
      maxAge: 5 * 60, // 5 minutes
    },
  },
  plugins: [
    twoFactor({
      issuer: 'StoreDash',
    }),
  ],
  socialProviders: {
    google: {
      clientId: options.env.googleClientId || '',
      clientSecret: options.env.googleClientSecret || '',
    },
    github: {
      clientId: options.env.githubClientId || '',
      clientSecret: options.env.githubClientSecret || '',
    },
  },
});
