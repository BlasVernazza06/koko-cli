import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Database, Server, CheckCircle2, Circle, Plus, Trash2, Layers, Activity } from 'lucide-react';

interface Item {
  _id: string;
  title: string;
  description?: string;
  status: 'pending' | 'in_progress' | 'completed';
  createdAt: string;
}

export function App() {
  const [items, setItems] = useState<Item[]>([]);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [loading, setLoading] = useState(true);
  const [serverStatus, setServerStatus] = useState<'checking' | 'online' | 'offline'>('checking');
  const [dbStatus, setDbStatus] = useState<string>('checking');

  const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

  const checkHealth = async () => {
    try {
      const res = await axios.get(`${API_URL}/health`);
      setServerStatus('online');
      setDbStatus(res.data.database || 'connected');
    } catch {
      setServerStatus('offline');
      setDbStatus('disconnected');
    }
  };

  const fetchItems = async () => {
    try {
      setLoading(true);
      const res = await axios.get(`${API_URL}/items`);
      setItems(res.data);
    } catch (err) {
      console.warn('Backend unavailable, using local items preview:', err);
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
    if (!title.trim()) return;

    try {
      const res = await axios.post(`${API_URL}/items`, { title, description, status: 'pending' });
      setItems([res.data, ...items]);
      setTitle('');
      setDescription('');
    } catch {
      // Offline fallback
      const localItem: Item = {
        _id: Date.now().toString(),
        title,
        description,
        status: 'pending',
        createdAt: new Date().toISOString(),
      };
      setItems([localItem, ...items]);
      setTitle('');
      setDescription('');
    }
  };

  const handleToggleStatus = async (item: Item) => {
    const nextStatus = item.status === 'completed' ? 'pending' : 'completed';
    try {
      await axios.put(`${API_URL}/items/${item._id}`, { status: nextStatus });
    } catch {
      // optimistic update
    }
    setItems(items.map((i) => (i._id === item._id ? { ...i, status: nextStatus } : i)));
  };

  const handleDeleteItem = async (id: string) => {
    try {
      await axios.delete(`${API_URL}/items/${id}`);
    } catch {
      // optimistic
    }
    setItems(items.filter((i) => i._id !== id));
  };

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 p-4 sm:p-8 font-sans">
      <div className="max-w-4xl mx-auto space-y-8">
        {/* Header Banner */}
        <header className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 p-6 bg-slate-900/80 border border-slate-800 rounded-2xl backdrop-blur">
          <div className="space-y-1">
            <div className="flex items-center gap-2">
              <span className="p-2 bg-emerald-500/10 text-emerald-400 rounded-xl">
                <Layers className="w-6 h-6" />
              </span>
              <h1 className="text-2xl font-bold tracking-tight text-white">[[.ProjectName]]</h1>
              <span className="text-xs font-semibold px-2.5 py-0.5 rounded-full bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                MERN Stack
              </span>
            </div>
            <p className="text-sm text-slate-400">
              React 19 SPA + Node Express (TypeScript) + MongoDB & Mongoose
            </p>
          </div>

          {/* Service status badges */}
          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-slate-800/80 border border-slate-700 text-xs">
              <Server className={`w-3.5 h-3.5 ${serverStatus === 'online' ? 'text-emerald-400' : 'text-rose-400'}`} />
              <span>API: {serverStatus}</span>
            </div>
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-slate-800/80 border border-slate-700 text-xs">
              <Database className={`w-3.5 h-3.5 ${dbStatus === 'connected' ? 'text-emerald-400' : 'text-amber-400'}`} />
              <span>MongoDB: {dbStatus}</span>
            </div>
          </div>
        </header>

        {/* Create Item Card */}
        <div className="p-6 bg-slate-900/50 border border-slate-800 rounded-2xl">
          <h2 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
            <Plus className="w-4 h-4 text-emerald-400" /> Crear Nuevo Registro
          </h2>
          <form onSubmit={handleAddItem} className="space-y-4">
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <input
                type="text"
                placeholder="Título del item o tarea..."
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="sm:col-span-2 px-4 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-sm focus:outline-none focus:border-emerald-500 transition"
                required
              />
              <button
                type="submit"
                className="px-4 py-2.5 bg-emerald-600 hover:bg-emerald-500 font-semibold rounded-xl text-sm transition flex items-center justify-center gap-2 shadow-lg shadow-emerald-600/20 text-white"
              >
                <Plus className="w-4 h-4" /> Agregar Item
              </button>
            </div>
          </form>
        </div>

        {/* Item List */}
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold text-white flex items-center gap-2">
              <Activity className="w-4 h-4 text-emerald-400" /> Registros en Base de Datos ({items.length})
            </h2>
            <button
              onClick={() => { checkHealth(); fetchItems(); }}
              className="text-xs text-slate-400 hover:text-emerald-400 transition"
            >
              ↻ Refrescar datos
            </button>
          </div>

          {loading ? (
            <div className="text-center py-12 text-slate-500 text-sm">Cargando registros...</div>
          ) : items.length === 0 ? (
            <div className="text-center py-12 p-6 bg-slate-900/30 border border-slate-800/80 rounded-2xl text-slate-400 text-sm">
              No hay items creados aún. ¡Agrega el primero arriba!
            </div>
          ) : (
            <div className="grid gap-3">
              {items.map((item) => (
                <div
                  key={item._id}
                  className="flex items-center justify-between p-4 bg-slate-900/60 border border-slate-800/80 hover:border-slate-700 rounded-xl transition group"
                >
                  <div className="flex items-center gap-3">
                    <button
                      onClick={() => handleToggleStatus(item)}
                      className="text-slate-400 hover:text-emerald-400 transition"
                    >
                      {item.status === 'completed' ? (
                        <CheckCircle2 className="w-5 h-5 text-emerald-400" />
                      ) : (
                        <Circle className="w-5 h-5" />
                      )}
                    </button>
                    <div>
                      <span className={`text-sm font-medium ${item.status === 'completed' ? 'line-through text-slate-500' : 'text-slate-200'}`}>
                        {item.title}
                      </span>
                      {item.description && (
                        <p className="text-xs text-slate-400 mt-0.5">{item.description}</p>
                      )}
                    </div>
                  </div>

                  <button
                    onClick={() => handleDeleteItem(item._id)}
                    className="p-2 text-slate-500 hover:text-rose-400 transition opacity-0 group-hover:opacity-100"
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