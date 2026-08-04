import { neon } from '@neondatabase/serverless';
import { drizzle } from 'drizzle-orm/neon-http';
import * as schemaExports from './schema/index';
import * as relationsExports from './relations';
const combinedSchema = { ...schemaExports, ...relationsExports };
const createDb = () => {
    const url = process.env.DATABASE_URL;
    if (!url) {
        return null;
    }
    const sql = neon(url);
    return drizzle(sql, { schema: combinedSchema });
};
export const db = createDb();
export { combinedSchema as schema };
//# sourceMappingURL=drizzle.js.map