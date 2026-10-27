import { jsonb, pgTable, text, timestamp } from 'drizzle-orm/pg-core';

import { user } from './auth';
import { tenant } from './tenants';

export const auditLog = pgTable('audit_log', {
  id: text('id').primaryKey(),
  action: text('action').notNull(),
  userId: text('user_id')
    .notNull()
    .references(() => user.id),
  tenantId: text('tenant_id')
    .notNull()
    .references(() => tenant.id),
  details: jsonb('details'),
  createdAt: timestamp('created_at').notNull(),
});
