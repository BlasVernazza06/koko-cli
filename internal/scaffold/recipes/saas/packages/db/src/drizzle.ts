import { neon } from '@neondatabase/serverless';
import { type NeonHttpDatabase, drizzle } from 'drizzle-orm/neon-http';

import * as schemaExports from './schema/index';
import * as relationsExports from './relations';

const combinedSchema = { ...schemaExports, ...relationsExports };

export type Database = NeonHttpDatabase<typeof combinedSchema>;

const createDb = () => {
  const url = process.env.DATABASE_URL || 'postgresql://placeholder:placeholder@localhost:5432/placeholder';
  const sql = neon(url);
  return drizzle(sql, { schema: combinedSchema }) as Database;
};

export const db = createDb();

export { combinedSchema as schema };
