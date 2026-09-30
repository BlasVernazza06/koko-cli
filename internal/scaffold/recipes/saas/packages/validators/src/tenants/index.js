import { z } from 'zod';
export const createTenantSchema = z.object({
    name: z.string().min(3, 'El nombre de la tienda debe tener al menos 3 caracteres'),
});
//# sourceMappingURL=index.js.map