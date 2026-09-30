import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Coffee, Database, Server, CheckCircle2, Circle, Plus, Trash2, Activity } from 'lucide-react';

interface Item {
  id: number;
  name: string;
  description?: string;
  completed: boolean;
  createdAt?: string;
}

export function App() {
  const [items, setItems] = useState<Item[]>([]);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [loading, setLoading] = useState(true);
  const [healthStatus, setHealthStatus] = useState<'checking' | 'online' | 'offline'>('checking');
  const [dbStatus, setDbStatus] = useState<string>('checking');

  const API_URL = '/api';

  const checkHealth = async () => {
    try {
      const res = await axios.get(`${API_URL}/health`);
      setHealthStatus('online');
      setDbStatus(res.data.database || 'connected');
    } catch {
      setHealthStatus('offline');
      setDbStatus('offline');
    }
  };

  const fetchItems = async () => {
    try {
      setLoading(true);
      const res = await axios.get(`${API_URL}/items`);
      setItems(res.data);
    } catch (err) {
      console.warn('Spring Boot API unavailable, fallback items displayed:', err);
      setItems([
        { id: 1, name: 'Configurar Spring Boot 3 con JPA en Koko', description: 'Entidades Hibernate y PostgreSQL', completed: true },
        { id: 2, name: 'Conectar cliente React con proxy Vite', description: 'Consumo tipado de endpoints REST', completed: false },
      ]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    checkHealth();
    fetchItems();
  }, []);

  const handleAddItem = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;

    try {
      const res = await axios.post(`${API_URL}/items`, {
        name,
        description,
        completed: false,
      });
      setItems([res.data, ...items]);
      setName('');
      setDescription('');
    } catch {
      const localItem: Item = {
        id: Date.now(),
        name,
        description,
        completed: false,
        createdAt: new Date().toISOString(),
      };
      setItems([localItem, ...items]);
      setName('');
      setDescription('');
    }
  };

  const handleToggle = async (item: Item) => {
    const updated = !item.completed;
    try {
      await axios.put(`${API_URL}/items/${item.id}`, {
        ...item,
        completed: updated,
      });
    } catch {
      // optimistic
    }
    setItems(items.map((i) => (i.id === item.id ? { ...i, completed: updated } : i)));
  };

  const handleDelete = async (id: number) => {
    try {
      await axios.delete(`${API_URL}/items/${id}`);
    } catch {
      // optimistic
    }
    setItems(items.filter((i) => i.id !== id));
  };

  return (
    <div className="min-h-screen bg-neutral-950 text-neutral-100 p-4 sm:p-8 font-sans">
      <div className="max-w-4xl mx-auto space-y-8">
        {/* Header */}
        <header className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 p-6 bg-neutral-900/80 border border-neutral-800 rounded-2xl backdrop-blur">
          <div className="space-y-1">
            <div className="flex items-center gap-2">
              <span className="p-2 bg-orange-500/10 text-orange-400 rounded-xl">
                <Coffee className="w-6 h-6" />
              </span>
              <h1 className="text-2xl font-bold tracking-tight text-white">[[.ProjectName]]</h1>
              <span className="text-xs font-semibold px-2.5 py-0.5 rounded-full bg-orange-500/10 text-orange-400 border border-orange-500/20">
                Spring Boot 3 + React
              </span>
            </div>
            <p className="text-sm text-neutral-400">
              Enterprise Java Spring Boot (Spring Data JPA + PostgreSQL) + React 19 SPA
            </p>
          </div>

          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-neutral-800/80 border border-neutral-700 text-xs">
              <Server className={`w-3.5 h-3.5 ${healthStatus === 'online' ? 'text-emerald-400' : 'text-rose-400'}`} />
              <span>Spring Boot: {healthStatus}</span>
            </div>
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-neutral-800/80 border border-neutral-700 text-xs">
              <Database className={`w-3.5 h-3.5 ${dbStatus === 'connected' ? 'text-emerald-400' : 'text-amber-400'}`} />
              <span>PostgreSQL: {dbStatus}</span>
            </div>
          </div>
        </header>

        {/* Form */}
        <div className="p-6 bg-neutral-900/50 border border-neutral-800 rounded-2xl">
          <h2 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
            <Plus className="w-4 h-4 text-orange-400" /> Crear Nuevo Registro JPA
          </h2>
          <form onSubmit={handleAddItem} className="space-y-4">
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <input
                type="text"
                placeholder="Nombre de la entidad..."
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="sm:col-span-2 px-4 py-2.5 rounded-xl bg-neutral-950 border border-neutral-800 text-sm focus:outline-none focus:border-orange-500 transition"
                required
              />
              <button
                type="submit"
                className="px-4 py-2.5 bg-orange-600 hover:bg-orange-500 font-semibold rounded-xl text-sm transition flex items-center justify-center gap-2 shadow-lg shadow-orange-600/20 text-white"
              >
                <Plus className="w-4 h-4" /> Guardar en BD
              </button>
            </div>
          </form>
        </div>

        {/* Items List */}
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold text-white flex items-center gap-2">
              <Activity className="w-4 h-4 text-orange-400" /> Entidades Spring Data JPA ({items.length})
            </h2>
            <button
              onClick={() => { checkHealth(); fetchItems(); }}
              className="text-xs text-neutral-400 hover:text-orange-400 transition"
            >
              ↻ Refrescar datos
            </button>
          </div>

          {loading ? (
            <div className="text-center py-12 text-neutral-500 text-sm">Consultando backend Spring Boot...</div>
          ) : items.length === 0 ? (
            <div className="text-center py-12 p-6 bg-neutral-900/30 border border-neutral-800/80 rounded-2xl text-neutral-400 text-sm">
              No hay entidades aún. ¡Crea la primera arriba!
            </div>
          ) : (
            <div className="grid gap-3">
              {items.map((item) => (
                <div
                  key={item.id}
                  className="flex items-center justify-between p-4 bg-neutral-900/60 border border-neutral-800/80 hover:border-neutral-700 rounded-xl transition group"
                >
                  <div className="flex items-center gap-3">
                    <button
                      onClick={() => handleToggle(item)}
                      className="text-neutral-400 hover:text-orange-400 transition"
                    >
                      {item.completed ? (
                        <CheckCircle2 className="w-5 h-5 text-emerald-400" />
                      ) : (
                        <Circle className="w-5 h-5" />
                      )}
                    </button>
                    <div>
                      <span className={`text-sm font-medium ${item.completed ? 'line-through text-neutral-500' : 'text-neutral-200'}`}>
                        {item.name}
                      </span>
                      {item.description && (
                        <p className="text-xs text-neutral-400 mt-0.5">{item.description}</p>
                      )}
                    </div>
                  </div>

                  <button
                    onClick={() => handleDelete(item.id)}
                    className="p-2 text-neutral-500 hover:text-rose-400 transition opacity-0 group-hover:opacity-100"
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
    </div>
  );
}

export default App;
