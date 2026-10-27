import { z } from 'zod';

export const createTaskSchema = z.object({
  title: z.string().min(1, 'El título es requerido').max(120),
  description: z.string().max(500).optional(),
  completed: z.boolean().optional().default(false),
});

export const updateTaskSchema = createTaskSchema.partial();

export type CreateTaskInput = z.infer<typeof createTaskSchema>;
export type UpdateTaskInput = z.infer<typeof updateTaskSchema>;
