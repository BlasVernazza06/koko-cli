import { z } from 'zod';

// 🛒 Validador para cada item de la orden
export const createOrderItemSchema = z.object({
  productId: z.string().min(1, 'El ID del producto es requerido'),
  quantity: z.number().int().positive('La cantidad debe ser mayor a 0'),
  price: z.number().positive('El precio debe ser mayor a 0'),
});

// 📦 Validador para la orden completa
export const createOrderSchema = z.object({
  items: z
    .array(createOrderItemSchema)
    .min(1, 'La orden debe tener al menos 1 producto'),
});

// 💡 Inferencia de tipos para usar como DTOs en el Backend
export type CreateOrderDto = z.infer<typeof createOrderSchema>;
export type CreateOrderItemDto = z.infer<typeof createOrderItemSchema>;
