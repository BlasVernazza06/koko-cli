import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Card,
  Flex,
  Box,
  Text,
  Heading,
  TextField,
  Button,
  Grid,
  Badge,
  IconButton
} from '@radix-ui/themes';
import { LogOut, CheckCircle, Circle, Trash2, Plus, AlertCircle, FileText } from 'lucide-react';
import api from '../lib/api';

interface Task {
  id: number;
  title: string;
  description?: string;
  status: 'PENDING' | 'COMPLETED';
  createdAt: string;
}

export default function DashboardPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [newTitle, setNewTitle] = useState('');
  const [newDesc, setNewDesc] = useState('');
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  // Retrieve user from localStorage
  const userString = localStorage.getItem('user');
  const user = userString ? JSON.parse(userString) : { name: 'Usuario', email: '' };

  // 1. Fetch Tasks
  const { data: tasks = [], isLoading } = useQuery<Task[]>({
    queryKey: ['tasks'],
    queryFn: async () => {
      const res = await api.get('/tasks');
      return res.data;
    }
  });

  // 2. Create Task Mutation
  const createTaskMutation = useMutation({
    mutationFn: async (taskData: { title: string; description?: string }) => {
      const res = await api.post('/tasks', taskData);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      setNewTitle('');
      setNewDesc('');
      setErrorMsg(null);
    },
    onError: (err: any) => {
      setErrorMsg(err.response?.data?.error || 'Error al crear la tarea');
    }
  });

  // 3. Toggle Status Mutation
  const toggleStatusMutation = useMutation({
    mutationFn: async (task: Task) => {
      const newStatus = task.status === 'PENDING' ? 'COMPLETED' : 'PENDING';
      const res = await api.put(`/tasks/${task.id}`, { status: newStatus });
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
    }
  });

  // 4. Delete Task Mutation
  const deleteTaskMutation = useMutation({
    mutationFn: async (id: number) => {
      const res = await api.delete(`/tasks/${id}`);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
    }
  });

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    navigate('/login');
  };

  const handleCreate = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newTitle.trim()) return;
    createTaskMutation.mutate({ title: newTitle, description: newDesc });
  };

  const completedCount = tasks.filter(t => t.status === 'COMPLETED').length;
  const pendingCount = tasks.length - completedCount;

  return (
    <Box style={{ background: 'var(--gray-2)', minHeight: '100vh', paddingBottom: '40px' }}>
      {/* Navbar */}
      <Box style={{ background: 'var(--gray-1)', borderBottom: '1px solid var(--gray-5)', padding: '16px 24px' }}>
        <Flex justify="between" align="center" style={{ maxWidth: '1200px', margin: '0 auto' }}>
          <Flex align="center" gap="2">
            <Box style={{ width: '32px', height: '32px', borderRadius: '8px', background: 'var(--accent-9)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'white' }}>
              <FileText size={18} />
            </Box>
            <Heading size="5" weight="bold">PERN Task Manager</Heading>
          </Flex>

          <Flex align="center" gap="4">
            <Box style={{ textAlign: 'right' }}>
              <Text size="2" weight="bold" display="block">{user.name}</Text>
              <Text size="1" color="gray" display="block">{user.email}</Text>
            </Box>
            <Button variant="ghost" color="red" style={{ cursor: 'pointer' }} onClick={handleLogout}>
              <LogOut size={16} />
              Salir
            </Button>
          </Flex>
        </Flex>
      </Box>

      {/* Main Container */}
      <Box style={{ maxWidth: '1200px', margin: '32px auto 0', padding: '0 16px' }}>
        <Grid columns={{ initial: '1', md: '3' }} gap="6">
          
          {/* Left Column: Stats & Creation Form */}
          <Flex direction="column" gap="6" style={{ gridColumn: 'span 1' }}>
            {/* Stats Card */}
            <Card style={{ borderRadius: '16px' }}>
              <Heading size="4" style={{ marginBottom: '16px' }}>Progreso de Tareas</Heading>
              <Flex justify="between" gap="4">
                <Box style={{ background: 'var(--blue-2)', padding: '16px', borderRadius: '12px', flex: 1, textAlign: 'center' }}>
                  <Text size="5" weight="bold" color="blue" display="block">{pendingCount}</Text>
                  <Text size="1" color="gray" weight="medium">Pendientes</Text>
                </Box>
                <Box style={{ background: 'var(--green-2)', padding: '16px', borderRadius: '12px', flex: 1, textAlign: 'center' }}>
                  <Text size="5" weight="bold" color="green" display="block">{completedCount}</Text>
                  <Text size="1" color="gray" weight="medium">Completadas</Text>
                </Box>
              </Flex>
            </Card>

            {/* Creation Card */}
            <Card style={{ borderRadius: '16px' }}>
              <Heading size="4" style={{ marginBottom: '16px' }}>Nueva Tarea</Heading>
              
              {errorMsg && (
                <Flex gap="2" align="center" style={{ padding: '8px 12px', background: 'var(--red-3)', borderRadius: '8px', color: 'var(--red-11)', marginBottom: '12px' }}>
                  <AlertCircle size={14} />
                  <Text size="1" weight="medium">{errorMsg}</Text>
                </Flex>
              )}

              <form onSubmit={handleCreate}>
                <Flex direction="column" gap="3">
                  <Flex direction="column" gap="1">
                    <Text as="label" size="1" weight="bold" color="gray" htmlFor="task-title">Título</Text>
                    <TextField.Root
                      id="task-title"
                      placeholder="Comprar víveres, etc."
                      value={newTitle}
                      onChange={(e) => setNewTitle(e.target.value)}
                      required
                    />
                  </Flex>

                  <Flex direction="column" gap="1">
                    <Text as="label" size="1" weight="bold" color="gray" htmlFor="task-desc">Descripción (Opcional)</Text>
                    <TextField.Root
                      id="task-desc"
                      placeholder="Detalles adicionales..."
                      value={newDesc}
                      onChange={(e) => setNewDesc(e.target.value)}
                    />
                  </Flex>

                  <Button type="submit" style={{ cursor: 'pointer', marginTop: '8px' }} disabled={createTaskMutation.isPending}>
                    <Plus size={16} />
                    {createTaskMutation.isPending ? 'Creando...' : 'Agregar Tarea'}
                  </Button>
                </Flex>
              </form>
            </Card>
          </Flex>

          {/* Right Column: Tasks List */}
          <Box style={{ gridColumn: 'span 2' }}>
            <Card style={{ borderRadius: '16px', minHeight: '400px' }}>
              <Heading size="4" style={{ marginBottom: '16px' }}>Mis Tareas</Heading>

              {isLoading ? (
                <Text size="2" color="gray" align="center" style={{ display: 'block', padding: '32px' }}>Cargando tareas...</Text>
              ) : tasks.length === 0 ? (
                <Flex direction="column" align="center" justify="center" gap="2" style={{ padding: '48px 0', opacity: 0.5 }}>
                  <CheckCircle size={48} color="var(--gray-9)" />
                  <Text size="2" weight="bold">¡Todo al día!</Text>
                  <Text size="1">No tienes tareas pendientes.</Text>
                </Flex>
              ) : (
                <Flex direction="column" gap="3">
                  {tasks.map((task) => (
                    <Card key={task.id} style={{ border: '1px solid var(--gray-4)', borderRadius: '12px' }}>
                      <Flex justify="between" align="center" gap="4">
                        <Flex gap="3" align="start" style={{ flex: 1 }}>
                          <IconButton
                            variant="ghost"
                            color={task.status === 'COMPLETED' ? 'green' : 'gray'}
                            style={{ cursor: 'pointer', marginTop: '2px' }}
                            onClick={() => toggleStatusMutation.mutate(task)}
                          >
                            {task.status === 'COMPLETED' ? (
                              <CheckCircle size={20} />
                            ) : (
                              <Circle size={20} />
                            )}
                          </IconButton>
                          <Box>
                            <Text
                              size="2"
                              weight="bold"
                              style={{
                                textDecoration: task.status === 'COMPLETED' ? 'line-through' : 'none',
                                color: task.status === 'COMPLETED' ? 'var(--gray-9)' : 'inherit'
                              }}
                            >
                              {task.title}
                            </Text>
                            {task.description && (
                              <Text
                                size="1"
                                color="gray"
                                display="block"
                                style={{
                                  marginTop: '2px',
                                  textDecoration: task.status === 'COMPLETED' ? 'line-through' : 'none'
                                }}
                              >
                                {task.description}
                              </Text>
                            )}
                          </Box>
                        </Flex>

                        <Flex align="center" gap="3">
                          <Badge color={task.status === 'COMPLETED' ? 'green' : 'blue'}>
                            {task.status === 'COMPLETED' ? 'Completado' : 'Pendiente'}
                          </Badge>
                          <IconButton
                            variant="ghost"
                            color="red"
                            style={{ cursor: 'pointer' }}
                            onClick={() => deleteTaskMutation.mutate(task.id)}
                          >
                            <Trash2 size={16} />
                          </IconButton>
                        </Flex>
                      </Flex>
                    </Card>
                  ))}
                </Flex>
              )}
            </Card>
          </Box>

        </Grid>
      </Box>
    </Box>
  );
}
