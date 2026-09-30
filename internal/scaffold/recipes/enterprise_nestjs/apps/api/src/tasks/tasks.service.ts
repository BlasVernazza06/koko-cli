import { Injectable, NotFoundException } from '@nestjs/common';
import { PrismaService } from '../prisma/prisma.service';
import { CreateTaskDto, UpdateTaskDto } from './dto/create-task.dto';

@Injectable()
export class TasksService {
  constructor(private readonly prisma: PrismaService) {}

  // Fallback in-memory items if DB is offline
  private inMemoryTasks = [
    {
      id: 'task-1',
      title: 'Configurar NestJS con Prisma ORM en Koko',
      description: 'Arquitectura monorrepo escalable con Turborepo',
      status: 'completed',
      priority: 'high',
      createdAt: new Date(),
      updatedAt: new Date(),
    },
    {
      id: 'task-2',
      title: 'Explorar documentación Swagger en /api/docs',
      description: 'Documentación OpenAPI autogenerada',
      status: 'in_progress',
      priority: 'medium',
      createdAt: new Date(),
      updatedAt: new Date(),
    },
  ];

  async findAll() {
    try {
      return await this.prisma.task.findMany({
        orderBy: { createdAt: 'desc' },
      });
    } catch {
      return this.inMemoryTasks;
    }
  }

  async findOne(id: string) {
    try {
      return await this.prisma.task.findUnique({ where: { id } });
    } catch {
      return this.inMemoryTasks.find((t) => t.id === id) || null;
    }
  }

  async create(dto: CreateTaskDto) {
    try {
      return await this.prisma.task.create({
        data: {
          title: dto.title,
          description: dto.description || null,
          status: dto.status || 'pending',
          priority: dto.priority || 'medium',
        },
      });
    } catch {
      const newTask = {
        id: `task-${Date.now()}`,
        title: dto.title,
        description: dto.description || '',
        status: dto.status || 'pending',
        priority: dto.priority || 'medium',
        createdAt: new Date(),
        updatedAt: new Date(),
      };
      this.inMemoryTasks.unshift(newTask as any);
      return newTask;
    }
  }

  async update(id: string, dto: UpdateTaskDto) {
    try {
      return await this.prisma.task.update({
        where: { id },
        data: dto,
      });
    } catch {
      const idx = this.inMemoryTasks.findIndex((t) => t.id === id);
      if (idx === -1) throw new NotFoundException('Task not found');
      this.inMemoryTasks[idx] = { ...this.inMemoryTasks[idx], ...dto, updatedAt: new Date() };
      return this.inMemoryTasks[idx];
    }
  }

  async remove(id: string) {
    try {
      return await this.prisma.task.delete({ where: { id } });
    } catch {
      this.inMemoryTasks = this.inMemoryTasks.filter((t) => t.id !== id);
      return { success: true, id };
    }
  }
}
