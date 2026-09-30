import { pgEnum, pgTable, text, timestamp } from 'drizzle-orm/pg-core';
import { user } from './auth';
export const roleEnum = pgEnum('role', ['ADMIN', 'STAFF', 'VIEWER']);
export const tenant = pgTable('tenant', {
    id: text('id').primaryKey(),
    name: text('name').notNull(),
    createdAt: timestamp('created_at', { withTimezone: true }).notNull(),
    updatedAt: timestamp('updated_at', { withTimezone: true }).notNull(),
});
export const userProfile = pgTable('user_profile', {
    id: text('id').primaryKey(),
    userId: text('user_id')
        .notNull()
        .references(() => user.id),
    tenantId: text('tenant_id')
        .notNull()
        .references(() => tenant.id),
    role: roleEnum('role').notNull().default('STAFF'),
    createdAt: timestamp('created_at', { withTimezone: true }).notNull(),
});
//# sourceMappingURL=tenants.js.map