import { z } from 'zod';

export const createProductSchema = z.object({
  name: z.string().min(3, 'El nombre debe tener al menos 3 caracteres'),
  description: z.string().optional(),
  price: z.number().positive('El precio debe ser mayor a cero'), // Guardamos en centavos en DB, pero el front envía decimales? O centavos?
  stock: z.number().int().min(0, 'El stock no puede ser negativo'),
});

export type CreateProductDto = z.infer<typeof createProductSchema>;

export const updateProductSchema = createProductSchema.partial(); // Hace todos los campos opcionales para parcheo

export type UpdateProductDto = z.infer<typeof updateProductSchema>;
