'use client';

import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Building2, Server, Database, CheckCircle2, Circle, Plus, Trash2, ShieldCheck, Zap, BookOpen, ExternalLink, Activity } from 'lucide-react';

interface Task {
  id: string;
  title: string;
  description?: string;
  status: 'pending' | 'in_progress' | 'completed';
  priority: 'low' | 'medium' | 'high' | 'urgent';
  createdAt: string;
}

export default function EnterpriseDashboard() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [priority, setPriority] = useState<'low' | 'medium' | 'high' | 'urgent'>('medium');
  const [loading, setLoading] = useState(true);
  const [apiHealth, setApiHealth] = useState<'checking' | 'online' | 'offline'>('checking');
  const [dbHealth, setDbHealth] = useState<string>('checking');

  const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:4000/api';

  const checkHealth = async () => {
    try {
      const res = await axios.get(`${API_URL}/health`);
      setApiHealth('online');
      setDbHealth(res.data.database || 'connected');
    } catch {
      setApiHealth('offline');
      setDbHealth('disconnected');
    }
  };

  const fetchTasks = async () => {
    try {
      setLoading(true);
      const res = await axios.get(`${API_URL}/tasks`);
      setTasks(res.data);
    } catch (err) {
      console.warn('NestJS API not reachable, loading fallback preview:', err);
      setTasks([
        { id: '1', title: 'Configurar NestJS con Prisma ORM en Koko', status: 'completed', priority: 'high', createdAt: new Date().toISOString() },
        { id: '2', title: 'Explorar documentación Swagger en /api/docs', status: 'in_progress', priority: 'medium', createdAt: new Date().toISOString() },
      ]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    checkHealth();
    fetchTasks();
  }, []);

  const handleAddTask = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;

    try {
      const res = await axios.post(`${API_URL}/tasks`, { title, description, priority, status: 'pending' });
      setTasks([res.data, ...tasks]);
      setTitle('');
      setDescription('');
    } catch {
      const localTask: Task = {
        id: Date.now().toString(),
        title,
        description,
        status: 'pending',
        priority,
        createdAt: new Date().toISOString(),
      };
      setTasks([localTask, ...tasks]);
      setTitle('');
      setDescription('');
    }
  };

  const handleToggleTask = async (task: Task) => {
    const nextStatus = task.status === 'completed' ? 'pending' : 'completed';
    try {
      await axios.patch(`${API_URL}/tasks/${task.id}`, { status: nextStatus });
    } catch {
      // optimistic
    }
    setTasks(tasks.map((t) => (t.id === task.id ? { ...t, status: nextStatus } : t)));
  };

  const handleDeleteTask = async (id: string) => {
    try {
      await axios.delete(`${API_URL}/tasks/${id}`);
    } catch {
      // optimistic
    }
    setTasks(tasks.filter((t) => t.id !== id));
  };

  return (
    <main className="min-h-screen bg-zinc-950 text-zinc-100 p-4 sm:p-8">
      <div className="max-w-5xl mx-auto space-y-8">
        {/* Header Banner */}
        <header className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 p-6 bg-zinc-900/80 border border-zinc-800 rounded-2xl backdrop-blur">
          <div className="space-y-1">
            <div className="flex items-center gap-2">
              <span className="p-2 bg-red-500/10 text-red-400 rounded-xl">
                <Building2 className="w-6 h-6" />
              </span>
              <h1 className="text-2xl font-bold tracking-tight text-white">[[.ProjectName]]</h1>
              <span className="text-xs font-semibold px-2.5 py-0.5 rounded-full bg-red-500/10 text-red-400 border border-red-500/20">
                Enterprise NestJS + Next.js
              </span>
            </div>
            <p className="text-sm text-zinc-400">
              Next.js 15 App Router + NestJS 11 (Prisma PostgreSQL + Better-Auth + Stripe + Resend)
            </p>
          </div>

          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-zinc-800/80 border border-zinc-700 text-xs">
              <Server className={`w-3.5 h-3.5 ${apiHealth === 'online' ? 'text-emerald-400' : 'text-rose-400'}`} />
              <span>NestJS: {apiHealth}</span>
            </div>
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-zinc-800/80 border border-zinc-700 text-xs">
              <Database className={`w-3.5 h-3.5 ${dbHealth === 'connected' ? 'text-emerald-400' : 'text-amber-400'}`} />
              <span>Prisma DB: {dbHealth}</span>
            </div>
            <a
              href="http://localhost:4000/api/docs"
              target="_blank"
              rel="noreferrer"
              className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-red-500/10 hover:bg-red-500/20 text-red-400 border border-red-500/30 text-xs font-medium transition"
            >
              <BookOpen className="w-3.5 h-3.5" />
              <span>Swagger</span>
              <ExternalLink className="w-3 h-3" />
            </a>
          </div>
        </header>

        {/* Feature Highlights Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="p-5 bg-zinc-900/40 border border-zinc-800/80 rounded-xl space-y-2">
            <div className="flex items-center gap-2 text-violet-400 text-sm font-semibold">
              <ShieldCheck className="w-4 h-4" /> Tipado de Extremo a Extremo
            </div>
            <p className="text-xs text-zinc-400">
              Modelos Prisma compartidos en `@repo/db` y esquemas Zod en `@repo/validators`.
            </p>
          </div>
          <div className="p-5 bg-zinc-900/40 border border-zinc-800/80 rounded-xl space-y-2">
            <div className="flex items-center gap-2 text-emerald-400 text-sm font-semibold">
              <Zap className="w-4 h-4" /> Integraciones Listas
            </div>
            <p className="text-xs text-zinc-400">
              Helpers de Stripe (`lib/stripe.ts`) y Resend (`lib/email.ts`) preconfigurados.
            </p>
          </div>
          <div className="p-5 bg-zinc-900/40 border border-zinc-800/80 rounded-xl space-y-2">
            <div className="flex items-center gap-2 text-amber-400 text-sm font-semibold">
              <Activity className="w-4 h-4" /> Arquitectura Modular
            </div>
            <p className="text-xs text-zinc-400">
              NestJS con Dependency Injection, Global Validation Pipes y Swagger OpenAPI docs.
            </p>
          </div>
        </div>

        {/* Form Card */}
        <div className="p-6 bg-zinc-900/50 border border-zinc-800 rounded-2xl">
          <h2 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
            <Plus className="w-4 h-4 text-red-400" /> Crear Tarea Empresarial
          </h2>
          <form onSubmit={handleAddTask} className="space-y-4">
            <div className="grid grid-cols-1 sm:grid-cols-4 gap-4">
              <input
                type="text"
                placeholder="Título de la tarea..."
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="sm:col-span-2 px-4 py-2.5 rounded-xl bg-zinc-950 border border-zinc-800 text-sm focus:outline-none focus:border-red-500 transition"
                required
              />
              <select
                value={priority}
                onChange={(e) => setPriority(e.target.value as any)}
                className="px-4 py-2.5 rounded-xl bg-zinc-950 border border-zinc-800 text-sm focus:outline-none focus:border-red-500 text-zinc-300"
              >
                <option value="low">Prioridad Baja</option>
                <option value="medium">Prioridad Media</option>
                <option value="high">Prioridad Alta</option>
                <option value="urgent">Urgente</option>
              </select>
              <button
                type="submit"
                className="px-4 py-2.5 bg-red-600 hover:bg-red-500 font-semibold rounded-xl text-sm transition flex items-center justify-center gap-2 shadow-lg shadow-red-600/20 text-white"
              >
                <Plus className="w-4 h-4" /> Guardar Tarea
              </button>
            </div>
          </form>
        </div>

        {/* Task List */}
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold text-white">Tareas en Base de Datos ({tasks.length})</h2>
            <button
              onClick={() => { checkHealth(); fetchTasks(); }}
              className="text-xs text-zinc-400 hover:text-red-400 transition"
            >
              ↻ Refrescar datos
            </button>
          </div>

          {loading ? (
            <div className="text-center py-12 text-zinc-500 text-sm">Consultando API NestJS...</div>
          ) : tasks.length === 0 ? (
            <div className="text-center py-12 p-6 bg-zinc-900/30 border border-zinc-800/80 rounded-2xl text-zinc-400 text-sm">
              No hay tareas registradas. Agrega una arriba.
            </div>
          ) : (
            <div className="grid gap-3">
              {tasks.map((task) => (
                <div
                  key={task.id}
                  className="flex items-center justify-between p-4 bg-zinc-900/60 border border-zinc-800/80 hover:border-zinc-700 rounded-xl transition group"
                >
                  <div className="flex items-center gap-3">
                    <button
                      onClick={() => handleToggleTask(task)}
                      className="text-zinc-400 hover:text-red-400 transition"
                    >
                      {task.status === 'completed' ? (
                        <CheckCircle2 className="w-5 h-5 text-emerald-400" />
                      ) : (
                        <Circle className="w-5 h-5" />
                      )}
                    </button>
                    <div>
                      <div className="flex items-center gap-2">
                        <span className={`text-sm font-medium ${task.status === 'completed' ? 'line-through text-zinc-500' : 'text-zinc-200'}`}>
                          {task.title}
                        </span>
                        <span className={`text-[10px] px-2 py-0.5 rounded-full font-semibold ${
                          task.priority === 'urgent' ? 'bg-rose-500/10 text-rose-400 border border-rose-500/20' :
                          task.priority === 'high' ? 'bg-amber-500/10 text-amber-400 border border-amber-500/20' :
                          'bg-zinc-800 text-zinc-400'
                        }`}>
                          {task.priority}
                        </span>
                      </div>
                      {task.description && (
                        <p className="text-xs text-zinc-400 mt-0.5">{task.description}</p>
                      )}
                    </div>
                  </div>

                  <button
                    onClick={() => handleDeleteTask(task.id)}
                    className="p-2 text-zinc-500 hover:text-rose-400 transition opacity-0 group-hover:opacity-100"
                    title="Eliminar"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </main>
  );
}
