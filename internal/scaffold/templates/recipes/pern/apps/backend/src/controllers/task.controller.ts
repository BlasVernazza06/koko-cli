import { Router } from 'express';
import { PrismaClient } from '@prisma/client';
import { createTaskSchema, updateTaskSchema } from '../schemas/task.schema';
import { authMiddleware } from '../middlewares/auth.middleware';

const router = Router();
const prisma = new PrismaClient();

// All task routes require authentication
router.use(authMiddleware);

router.get('/', async (req, res) => {
  const userId = (req as any).userId;
  try {
    const tasks = await prisma.task.findMany({
      where: { userId },
      orderBy: { createdAt: 'desc' },
    });
    return res.json(tasks);
  } catch (error) {
    return res.status(500).json({ error: 'Error al obtener tareas' });
  }
});

router.post('/', async (req, res) => {
  const userId = (req as any).userId;
  const result = createTaskSchema.safeParse(req.body);
  if (!result.success) {
    return res.status(400).json({ error: result.error.errors[0]?.message || 'Datos inválidos' });
  }

  const { title, description } = result.data;

  try {
    const task = await prisma.task.create({
      data: {
        title,
        description,
        userId,
      },
    });
    return res.status(201).json(task);
  } catch (error) {
    return res.status(500).json({ error: 'Error al crear la tarea' });
  }
});

router.put('/:id', async (req, res) => {
  const userId = (req as any).userId;
  const taskId = parseInt(req.params.id);
  if (isNaN(taskId)) {
    return res.status(400).json({ error: 'ID de tarea inválido' });
  }

  const result = updateTaskSchema.safeParse(req.body);
  if (!result.success) {
    return res.status(400).json({ error: result.error.errors[0]?.message || 'Datos inválidos' });
  }

  try {
    const task = await prisma.task.findFirst({
      where: { id: taskId, userId },
    });

    if (!task) {
      return res.status(404).json({ error: 'Tarea no encontrada' });
    }

    const updatedTask = await prisma.task.update({
      where: { id: taskId },
      data: result.data,
    });

    return res.json(updatedTask);
  } catch (error) {
    return res.status(500).json({ error: 'Error al actualizar la tarea' });
  }
});

router.delete('/:id', async (req, res) => {
  const userId = (req as any).userId;
  const taskId = parseInt(req.params.id);
  if (isNaN(taskId)) {
    return res.status(400).json({ error: 'ID de tarea inválido' });
  }

  try {
    const task = await prisma.task.findFirst({
      where: { id: taskId, userId },
    });

    if (!task) {
      return res.status(404).json({ error: 'Tarea no encontrada' });
    }

    await prisma.task.delete({
      where: { id: taskId },
    });

    return res.json({ message: 'Tarea eliminada correctamente' });
  } catch (error) {
    return res.status(500).json({ error: 'Error al eliminar la tarea' });
  }
});

export default router;
