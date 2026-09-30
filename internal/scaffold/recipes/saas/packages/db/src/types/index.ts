import { type InferInsertModel, type InferSelectModel } from 'drizzle-orm';
import { user, tenant, userProfile, product, order, orderItem, roleEnum, orderStatusEnum } from '../schema/index';

// ==========================================
// Base Types (Direct from Schema)
// ==========================================

// Auth & Profiles & Enums
export type DbUser = InferSelectModel<typeof user>;
export type NewUser = InferInsertModel<typeof user>;

export type Role = (typeof roleEnum.enumValues)[number];

export type Tenant = InferSelectModel<typeof tenant>;
export type NewTenant = InferInsertModel<typeof tenant>;

export type UserProfile = InferSelectModel<typeof userProfile>;
export type NewUserProfile = InferInsertModel<typeof userProfile>;

// Business (Product)
export type Product = InferSelectModel<typeof product>;
export type NewProduct = InferInsertModel<typeof product>;

// Business (Order)
export type Order = InferSelectModel<typeof order>;
export type NewOrder = InferInsertModel<typeof order>;
export type OrderStatus = (typeof orderStatusEnum.enumValues)[number];

export type OrderItem = InferSelectModel<typeof orderItem>;
export type NewOrderItem = InferInsertModel<typeof orderItem>;

// ==========================================
// Composite Types (App Specific)
// ==========================================

export type UserProfileWithUser = UserProfile & {
  user: Pick<DbUser, 'name' | 'email' | 'image'>;
};

export type OrderWithItems = Order & {
  items: OrderItem[];
};
