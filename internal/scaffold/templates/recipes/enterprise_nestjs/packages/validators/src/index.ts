import { z } from 'zod';

export const createTaskSchema = z.object({
  title: z.string().min(1, 'El título es requerido').max(120),
  description: z.string().max(500).optional(),
  status: z.enum(['pending', 'in_progress', 'completed']).optional().default('pending'),
  priority: z.enum(['low', 'medium', 'high', 'urgent']).optional().default('medium'),
});

export const updateTaskSchema = createTaskSchema.partial();

export const userAuthSchema = z.object({
  email: z.string().email('Email inválido'),
  password: z.string().min(8, 'La contraseña debe tener al menos 8 caracteres'),
  name: z.string().min(2).optional(),
});

export type CreateTaskInput = z.infer<typeof createTaskSchema>;
export type UpdateTaskInput = z.infer<typeof updateTaskSchema>;
export type UserAuthInput = z.infer<typeof userAuthSchema>;
