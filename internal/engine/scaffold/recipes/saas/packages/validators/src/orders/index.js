import { z } from 'zod';
export const createOrderItemSchema = z.object({
    productId: z.string().min(1, 'El ID del producto es requerido'),
    quantity: z.number().int().positive('La cantidad debe ser mayor a 0'),
    price: z.number().positive('El precio debe ser mayor a 0'),
});
export const createOrderSchema = z.object({
    items: z
        .array(createOrderItemSchema)
        .min(1, 'La orden debe tener al menos 1 producto'),
});
//# sourceMappingURL=index.js.map