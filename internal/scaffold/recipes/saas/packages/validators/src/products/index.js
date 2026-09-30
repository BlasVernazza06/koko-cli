import { z } from 'zod';
export const createProductSchema = z.object({
    name: z.string().min(3, 'El nombre debe tener al menos 3 caracteres'),
    description: z.string().optional(),
    price: z.number().positive('El precio debe ser mayor a cero'),
    stock: z.number().int().min(0, 'El stock no puede ser negativo'),
});
export const updateProductSchema = createProductSchema.partial();
//# sourceMappingURL=index.js.map