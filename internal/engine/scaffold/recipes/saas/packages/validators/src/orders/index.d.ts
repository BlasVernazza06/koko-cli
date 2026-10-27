import { z } from 'zod';
export declare const createOrderItemSchema: z.ZodObject<{
    productId: z.ZodString;
    quantity: z.ZodNumber;
    price: z.ZodNumber;
}, z.core.$strip>;
export declare const createOrderSchema: z.ZodObject<{
    items: z.ZodArray<z.ZodObject<{
        productId: z.ZodString;
        quantity: z.ZodNumber;
        price: z.ZodNumber;
    }, z.core.$strip>>;
}, z.core.$strip>;
export type CreateOrderDto = z.infer<typeof createOrderSchema>;
export type CreateOrderItemDto = z.infer<typeof createOrderItemSchema>;
