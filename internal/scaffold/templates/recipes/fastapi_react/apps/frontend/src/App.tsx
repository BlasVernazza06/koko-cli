import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Terminal, Cpu, CheckCircle2, Circle, Plus, Trash2, BookOpen, ExternalLink } from 'lucide-react';

interface Item {
  id: number;
  title: string;
  description?: string;
  status: string;
  created_at: string;
}

export function App() {
  const [items, setItems] = useState<Item[]>([]);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [loading, setLoading] = useState(true);
  const [apiHealth, setApiHealth] = useState<'checking' | 'online' | 'offline'>('checking');

  const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000';

  const checkHealth = async () => {
    try {
      await axios.get(`${API_URL}/api/health`);
      setApiHealth('online');
    } catch {
      setApiHealth('offline');
    }
  };

  const fetchItems = async () => {
    try {
      setLoading(true);
      const res = await axios.get(`${API_URL}/api/items`);
      setItems(res.data);
    } catch (err) {
      console.warn('FastAPI backend unavailable, using fallback:', err);
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
      const res = await axios.post(`${API_URL}/api/items`, {
        title,
        description,
        status: 'pending',
      });
      setItems([res.data, ...items]);
      setTitle('');
      setDescription('');
    } catch {
      // Local fallback
      const localItem: Item = {
        id: Date.now(),
        title,
        description,
        status: 'pending',
        created_at: new Date().toISOString(),
      };
      setItems([localItem, ...items]);
      setTitle('');
      setDescription('');
    }
  };

  const handleToggleStatus = async (item: Item) => {
    const nextStatus = item.status === 'completed' ? 'pending' : 'completed';
    try {
      await axios.put(`${API_URL}/api/items/${item.id}`, { status: nextStatus });
    } catch {
      // optimistic
    }
    setItems(items.map((i) => (i.id === item.id ? { ...i, status: nextStatus } : i)));
  };

  const handleDeleteItem = async (id: number) => {
    try {
      await axios.delete(`${API_URL}/api/items/${id}`);
    } catch {
      // optimistic
    }
    setItems(items.filter((i) => i.id !== id));
  };

  return (
    <div className="min-h-screen bg-zinc-950 text-zinc-100 p-4 sm:p-8 font-sans">
      <div className="max-w-4xl mx-auto space-y-8">
        {/* Header */}
        <header className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 p-6 bg-zinc-900/80 border border-zinc-800 rounded-2xl backdrop-blur">
          <div className="space-y-1">
            <div className="flex items-center gap-2">
              <span className="p-2 bg-amber-500/10 text-amber-400 rounded-xl">
                <Terminal className="w-6 h-6" />
              </span>
              <h1 className="text-2xl font-bold tracking-tight text-white">[[.ProjectName]]</h1>
              <span className="text-xs font-semibold px-2.5 py-0.5 rounded-full bg-amber-500/10 text-amber-400 border border-amber-500/20">
                FastAPI + React SPA
              </span>
            </div>
            <p className="text-sm text-zinc-400">
              Python FastAPI Backend (Pydantic v2) + React 19 Single Page Application
            </p>
          </div>

          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-zinc-800/80 border border-zinc-700 text-xs">
              <Cpu className={`w-3.5 h-3.5 ${apiHealth === 'online' ? 'text-emerald-400' : 'text-rose-400'}`} />
              <span>FastAPI: {apiHealth}</span>
            </div>
            <a
              href="http://localhost:8000/docs"
              target="_blank"
              rel="noreferrer"
              className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-amber-500/10 hover:bg-amber-500/20 text-amber-400 border border-amber-500/30 text-xs font-medium transition"
            >
              <BookOpen className="w-3.5 h-3.5" />
              <span>Swagger Docs</span>
              <ExternalLink className="w-3 h-3" />
            </a>
          </div>
        </header>

        {/* Create Item Card */}
        <div className="p-6 bg-zinc-900/50 border border-zinc-800 rounded-2xl">
          <h2 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
            <Plus className="w-4 h-4 text-amber-400" /> Crear Nuevo Registro Pydantic
          </h2>
          <form onSubmit={handleAddItem} className="space-y-4">
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <input
                type="text"
                placeholder="Título del item o entidad..."
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="sm:col-span-2 px-4 py-2.5 rounded-xl bg-zinc-950 border border-zinc-800 text-sm focus:outline-none focus:border-amber-500 transition"
                required
              />
              <button
                type="submit"
                className="px-4 py-2.5 bg-amber-500 hover:bg-amber-400 font-semibold rounded-xl text-sm transition flex items-center justify-center gap-2 shadow-lg shadow-amber-500/20 text-zinc-950"
              >
                <Plus className="w-4 h-4" /> Agregar Item
              </button>
            </div>
          </form>
        </div>

        {/* Item List */}
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold text-white">Registros Activos ({items.length})</h2>
            <button
              onClick={() => { checkHealth(); fetchItems(); }}
              className="text-xs text-zinc-400 hover:text-amber-400 transition"
            >
              ↻ Refrescar datos
            </button>
          </div>

          {loading ? (
            <div className="text-center py-12 text-zinc-500 text-sm">Cargando registros desde FastAPI...</div>
          ) : items.length === 0 ? (
            <div className="text-center py-12 p-6 bg-zinc-900/30 border border-zinc-800/80 rounded-2xl text-zinc-400 text-sm">
              No hay items aún. Crea uno nuevo arriba.
            </div>
          ) : (
            <div className="grid gap-3">
              {items.map((item) => (
                <div
                  key={item.id}
                  className="flex items-center justify-between p-4 bg-zinc-900/60 border border-zinc-800/80 hover:border-zinc-700 rounded-xl transition group"
                >
                  <div className="flex items-center gap-3">
                    <button
                      onClick={() => handleToggleStatus(item)}
                      className="text-zinc-400 hover:text-amber-400 transition"
                    >
                      {item.status === 'completed' ? (
                        <CheckCircle2 className="w-5 h-5 text-emerald-400" />
                      ) : (
                        <Circle className="w-5 h-5" />
                      )}
                    </button>
                    <div>
                      <span className={`text-sm font-medium ${item.status === 'completed' ? 'line-through text-zinc-500' : 'text-zinc-200'}`}>
                        {item.title}
                      </span>
                      {item.description && (
                        <p className="text-xs text-zinc-400 mt-0.5">{item.description}</p>
                      )}
                    </div>
                  </div>

                  <button
                    onClick={() => handleDeleteItem(item.id)}
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
    </div>
  );
}

export default App;