import React, { useState, useEffect } from 'react';
import {
  StyleSheet,
  Text,
  View,
  TextInput,
  TouchableOpacity,
  FlatList,
  SafeAreaView,
  ActivityIndicator,
  StatusBar,
} from 'react-native';
import { StatusBar as ExpoStatusBar } from 'expo-status-bar';
import { api } from './src/services/api';

interface Task {
  id: string;
  title: string;
  description?: string;
  completed: boolean;
}

export default function App() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [title, setTitle] = useState('');
  const [loading, setLoading] = useState(true);
  const [apiStatus, setApiStatus] = useState<'checking' | 'online' | 'offline'>('checking');

  const checkHealth = async () => {
    try {
      await api.get('/health');
      setApiStatus('online');
    } catch {
      setApiStatus('offline');
    }
  };

  const fetchTasks = async () => {
    try {
      setLoading(true);
      const res = await api.get('/tasks');
      setTasks(res.data);
    } catch {
      setTasks([
        { id: '1', title: 'Iniciar Expo en simulador o dispositivo', completed: true },
        { id: '2', title: 'Conectar endpoints de Express y PostgreSQL', completed: false },
      ]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    checkHealth();
    fetchTasks();
  }, []);

  const handleAddTask = async () => {
    if (!title.trim()) return;
    try {
      const res = await api.post('/tasks', { title, completed: false });
      setTasks([res.data, ...tasks]);
      setTitle('');
    } catch {
      const localTask: Task = {
        id: Date.now().toString(),
        title,
        completed: false,
      };
      setTasks([localTask, ...tasks]);
      setTitle('');
    }
  };

  const handleToggle = async (task: Task) => {
    const updated = !task.completed;
    try {
      await api.patch(`/tasks/${task.id}`, { completed: updated });
    } catch {
      // optimistic
    }
    setTasks(tasks.map((t) => (t.id === task.id ? { ...t, completed: updated } : t)));
  };

  const handleDelete = async (id: string) => {
    try {
      await api.delete(`/tasks/${id}`);
    } catch {
      // optimistic
    }
    setTasks(tasks.filter((t) => t.id !== id));
  };

  return (
    <SafeAreaView style={styles.container}>
      <ExpoStatusBar style="light" />
      <View style={styles.header}>
        <View>
          <Text style={styles.title}>[[.ProjectName]]</Text>
          <Text style={styles.subtitle}>Expo React Native + Express API</Text>
        </View>
        <View
          style={[
            styles.statusBadge,
            { backgroundColor: apiStatus === 'online' ? '#059669' : '#dc2626' },
          ]}
        >
          <Text style={styles.statusText}>{apiStatus.toUpperCase()}</Text>
        </View>
      </View>

      <View style={styles.form}>
        <TextInput
          style={styles.input}
          placeholder="Escribe una nueva tarea..."
          placeholderTextColor="#71717a"
          value={title}
          onChangeText={setTitle}
        />
        <TouchableOpacity style={styles.addButton} onPress={handleAddTask}>
          <Text style={styles.addButtonText}>Agregar</Text>
        </TouchableOpacity>
      </View>

      {loading ? (
        <ActivityIndicator size="large" color="#6366f1" style={{ marginTop: 40 }} />
      ) : (
        <FlatList
          data={tasks}
          keyExtractor={(item) => item.id}
          contentContainerStyle={styles.list}
          renderItem={({ item }) => (
            <View style={styles.taskCard}>
              <TouchableOpacity
                onPress={() => handleToggle(item)}
                style={styles.checkboxContainer}
              >
                <View
                  style={[
                    styles.checkbox,
                    item.completed && styles.checkboxCompleted,
                  ]}
                />
                <Text
                  style={[
                    styles.taskText,
                    item.completed && styles.taskTextCompleted,
                  ]}
                >
                  {item.title}
                </Text>
              </TouchableOpacity>
              <TouchableOpacity onPress={() => handleDelete(item.id)}>
                <Text style={styles.deleteText}>✕</Text>
              </TouchableOpacity>
            </View>
          )}
        />
      )}
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#09090b',
    paddingTop: StatusBar.currentHeight ? StatusBar.currentHeight + 10 : 20,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 20,
    paddingBottom: 15,
    borderBottomWidth: 1,
    borderBottomColor: '#27272a',
  },
  title: {
    fontSize: 22,
    fontWeight: 'bold',
    color: '#fafafa',
  },
  subtitle: {
    fontSize: 13,
    color: '#a1a1aa',
    marginTop: 2,
  },
  statusBadge: {
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 6,
  },
  statusText: {
    color: '#ffffff',
    fontSize: 10,
    fontWeight: 'bold',
  },
  form: {
    flexDirection: 'row',
    padding: 20,
    gap: 10,
  },
  input: {
    flex: 1,
    backgroundColor: '#18181b',
    borderWidth: 1,
    borderColor: '#27272a',
    borderRadius: 10,
    paddingHorizontal: 14,
    paddingVertical: 10,
    color: '#fafafa',
    fontSize: 14,
  },
  addButton: {
    backgroundColor: '#4f46e5',
    paddingHorizontal: 16,
    justifyContent: 'center',
    borderRadius: 10,
  },
  addButtonText: {
    color: '#ffffff',
    fontWeight: '600',
    fontSize: 14,
  },
  list: {
    paddingHorizontal: 20,
    gap: 10,
  },
  taskCard: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    backgroundColor: '#18181b',
    borderWidth: 1,
    borderColor: '#27272a',
    padding: 14,
    borderRadius: 10,
    marginBottom: 8,
  },
  checkboxContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    flex: 1,
    gap: 10,
  },
  checkbox: {
    width: 20,
    height: 20,
    borderRadius: 4,
    borderWidth: 2,
    borderColor: '#6366f1',
  },
  checkboxCompleted: {
    backgroundColor: '#6366f1',
  },
  taskText: {
    color: '#fafafa',
    fontSize: 14,
    flex: 1,
  },
  taskTextCompleted: {
    textDecorationLine: 'line-through',
    color: '#71717a',
  },
  deleteText: {
    color: '#ef4444',
    fontSize: 16,
    paddingHorizontal: 6,
  },
});
