export declare const userRelations: import("drizzle-orm").Relations<"user", {
    profiles: import("drizzle-orm").Many<"user_profile">;
    auditLogs: import("drizzle-orm").Many<"audit_log">;
}>;
export declare const tenantRelations: import("drizzle-orm").Relations<"tenant", {
    profiles: import("drizzle-orm").Many<"user_profile">;
    products: import("drizzle-orm").Many<"product">;
    orders: import("drizzle-orm").Many<"order">;
    auditLogs: import("drizzle-orm").Many<"audit_log">;
}>;
export declare const userProfileRelations: import("drizzle-orm").Relations<"user_profile", {
    user: import("drizzle-orm").One<"user", true>;
    tenant: import("drizzle-orm").One<"tenant", true>;
}>;
export declare const productRelations: import("drizzle-orm").Relations<"product", {
    tenant: import("drizzle-orm").One<"tenant", true>;
    orderItems: import("drizzle-orm").Many<"order_item">;
}>;
export declare const orderRelations: import("drizzle-orm").Relations<"order", {
    tenant: import("drizzle-orm").One<"tenant", true>;
    items: import("drizzle-orm").Many<"order_item">;
}>;
export declare const orderItemRelations: import("drizzle-orm").Relations<"order_item", {
    order: import("drizzle-orm").One<"order", true>;
    product: import("drizzle-orm").One<"product", true>;
}>;
export declare const auditLogRelations: import("drizzle-orm").Relations<"audit_log", {
    user: import("drizzle-orm").One<"user", true>;
    tenant: import("drizzle-orm").One<"tenant", true>;
}>;
