import { z } from 'zod';

export const createTenantSchema = z.object({
  name: z.string().min(3, 'El nombre de la tienda debe tener al menos 3 caracteres'),
});

export type CreateTenantDto = z.infer<typeof createTenantSchema>;
