import {
  index,
  integer,
  pgEnum,
  pgTable,
  text,
  timestamp,
} from 'drizzle-orm/pg-core';

import { tenant } from './tenants';

export const orderStatusEnum = pgEnum('order_status', [
  'PENDING',
  'PAID',
  'SHIPPED',
  'DELIVERED',
  'CANCELLED',
]);

export const product = pgTable('product', {
  id: text('id').primaryKey(),
  name: text('name').notNull(),
  description: text('description'),
  price: integer('price').notNull(),
  stock: integer('stock').notNull(),
  tenantId: text('tenant_id')
    .notNull()
    .references(() => tenant.id),
  createdAt: timestamp('created_at', { withTimezone: true }).notNull(),
  updatedAt: timestamp('updated_at', { withTimezone: true }).notNull(),
});

export const order = pgTable(
  'order',
  {
    id: text('id').primaryKey(),
    total: integer('total').notNull(),
    status: orderStatusEnum('status').notNull().default('PENDING'),
    tenantId: text('tenant_id')
      .notNull()
      .references(() => tenant.id),
    createdAt: timestamp('created_at', { withTimezone: true }).notNull(),
  },
  (table) => {
    return {
      tenantIdIdx: index('order_tenant_id_idx').on(table.tenantId),
    };
  },
);

export const orderItem = pgTable('order_item', {
  id: text('id').primaryKey(),
  orderId: text('order_id')
    .notNull()
    .references(() => order.id),
  productId: text('product_id')
    .notNull()
    .references(() => product.id),
  quantity: integer('quantity').notNull(),
  price: integer('price').notNull(),
});

export const productSchema = {
  product,
};

export const orderSchema = {
  order,
  orderItem,
};
