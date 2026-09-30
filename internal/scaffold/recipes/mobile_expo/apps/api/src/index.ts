import express from 'express';
import cors from 'cors';
import dotenv from 'dotenv';
import { prisma } from '@repo/db';
import { createTaskSchema, updateTaskSchema } from '@repo/validators';

dotenv.config();

const app = express();
const port = process.env.PORT || 8080;

app.use(cors());
app.use(express.json());

// In-memory fallback
let fallbackTasks = [
  { id: '1', title: 'Explorar plantilla Expo + Express en Koko CLI', description: 'React Native con TypeScript y Prisma', completed: true },
  { id: '2', title: 'Conectar dispositivo móvil con Expo Go', description: 'Escanear QR o ejecutar en emulador', completed: false },
];

// Health endpoint
app.get('/api/health', async (req, res) => {
  let dbStatus = 'connected';
  try {
    await prisma.$queryRaw`SELECT 1`;
  } catch {
    dbStatus = 'disconnected';
  }

  res.json({
    status: 'ok',
    backend: 'Node.js Express + Prisma ORM',
    database: dbStatus,
    timestamp: new Date().toISOString(),
  });
});

// Tasks CRUD
app.get('/api/tasks', async (req, res) => {
  try {
    const tasks = await prisma.task.findMany({
      orderBy: { createdAt: 'desc' },
    });
    return res.json(tasks);
  } catch {
    return res.json(fallbackTasks);
  }
});

app.post('/api/tasks', async (req, res) => {
  try {
    const validated = createTaskSchema.parse(req.body);
    try {
      const task = await prisma.task.create({ data: validated });
      return res.status(201).json(task);
    } catch {
      const newTask = { id: Date.now().toString(), ...validated, completed: false };
      fallbackTasks.unshift(newTask as any);
      return res.status(201).json(newTask);
    }
  } catch (error: any) {
    return res.status(400).json({ error: 'Validation error', details: error.errors });
  }
});

app.patch('/api/tasks/:id', async (req, res) => {
  try {
    const validated = updateTaskSchema.parse(req.body);
    try {
      const task = await prisma.task.update({
        where: { id: req.params.id },
        data: validated,
      });
      return res.json(task);
    } catch {
      const idx = fallbackTasks.findIndex((t) => t.id === req.params.id);
      if (idx !== -1) {
        fallbackTasks[idx] = { ...fallbackTasks[idx], ...validated };
        return res.json(fallbackTasks[idx]);
      }
      return res.status(404).json({ error: 'Task not found' });
    }
  } catch (error: any) {
    return res.status(400).json({ error: 'Validation error', details: error.errors });
  }
});

app.delete('/api/tasks/:id', async (req, res) => {
  try {
    await prisma.task.delete({ where: { id: req.params.id } });
  } catch {
    fallbackTasks = fallbackTasks.filter((t) => t.id !== req.params.id);
  }
  res.json({ message: 'Task deleted successfully', id: req.params.id });
});

app.listen(port, () => {
  console.log(`📱 [Mobile Expo API] Backend running on http://localhost:${port}`);
});
