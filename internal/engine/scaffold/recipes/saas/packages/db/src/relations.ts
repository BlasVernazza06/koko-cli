import { relations } from 'drizzle-orm';
import * as schema from './schema/index';

export const userRelations = relations(schema.user, ({ many }) => ({
  profiles: many(schema.userProfile),
  auditLogs: many(schema.auditLog),
}));

export const tenantRelations = relations(schema.tenant, ({ many }) => ({
  profiles: many(schema.userProfile),
  products: many(schema.product),
  orders: many(schema.order),
  auditLogs: many(schema.auditLog),
}));

export const userProfileRelations = relations(schema.userProfile, ({ one }) => ({
  user: one(schema.user, {
    fields: [schema.userProfile.userId],
    references: [schema.user.id],
  }),
  tenant: one(schema.tenant, {
    fields: [schema.userProfile.tenantId],
    references: [schema.tenant.id],
  }),
}));

export const productRelations = relations(schema.product, ({ one, many }) => ({
  tenant: one(schema.tenant, {
    fields: [schema.product.tenantId],
    references: [schema.tenant.id],
  }),
  orderItems: many(schema.orderItem),
}));

export const orderRelations = relations(schema.order, ({ one, many }) => ({
  tenant: one(schema.tenant, {
    fields: [schema.order.tenantId],
    references: [schema.tenant.id],
  }),
  items: many(schema.orderItem),
}));

export const orderItemRelations = relations(schema.orderItem, ({ one }) => ({
  order: one(schema.order, {
    fields: [schema.orderItem.orderId],
    references: [schema.order.id],
  }),
  product: one(schema.product, {
    fields: [schema.orderItem.productId],
    references: [schema.product.id],
  }),
}));

export const auditLogRelations = relations(schema.auditLog, ({ one }) => ({
  user: one(schema.user, {
    fields: [schema.auditLog.userId],
    references: [schema.user.id],
  }),
  tenant: one(schema.tenant, {
    fields: [schema.auditLog.tenantId],
    references: [schema.tenant.id],
  }),
}));
