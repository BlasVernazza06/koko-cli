import { defineConfig } from 'drizzle-kit';

export default defineConfig({
  schema: './schema.ts',
  out: './drizzle',
  dialect: '[[ if eq .Database "mysql" ]]mysql2[[ else ]]postgresql[[ end ]]',
  dbCredentials: {
    url: process.env.DATABASE_URL || '',
  },
});
