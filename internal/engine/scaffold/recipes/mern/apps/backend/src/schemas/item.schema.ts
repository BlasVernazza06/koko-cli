import { z } from 'zod';

export const createItemSchema = z.object({
  title: z.string().min(1, 'El título es requerido').max(120),
  description: z.string().max(500).optional().default(''),
  status: z.enum(['pending', 'in_progress', 'completed']).optional().default('pending'),
});

export const updateItemSchema = createItemSchema.partial();
