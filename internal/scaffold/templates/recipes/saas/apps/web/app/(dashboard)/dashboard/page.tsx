import {
  ArrowDownRight,
  ArrowUpRight,
  Bell,
  Clock,
  Cpu,
  Search,
  Activity,
  TrendingUp,
  Users,
  ChevronRight,
} from 'lucide-react';

import { Badge } from '@repo/ui/components/ui/badge';
import { Button } from '@repo/ui/components/ui/button';

export default async function DashboardPage() {
  const stats = [
    {
      icon: Users,
      label: 'Usuarios Activos',
      value: '1,284',
      trend: '+12.5%',
      trendUp: true,
      color: 'text-emerald-600',
      bg: 'bg-emerald-500/10',
    },
    {
      icon: Activity,
      label: 'API Requests',
      value: '84.2K',
      trend: '+8.2%',
      trendUp: true,
      color: 'text-brand',
      bg: 'bg-brand/10',
    },
    {
      icon: Cpu,
      label: 'Uso del Servidor',
      value: '24%',
      trend: '-3.1%',
      trendUp: false,
      color: 'text-amber-600',
      bg: 'bg-amber-500/10',
    },
    {
      icon: TrendingUp,
      label: 'Conversión',
      value: '4.8%',
      trend: '+1.4%',
      trendUp: true,
      color: 'text-violet-600',
      bg: 'bg-violet-500/10',
    },
  ];

  const recentEvents = [
    {
      id: '1',
      event: 'Nuevo usuario registrado',
      user: 'alex@example.com',
      time: 'Hace 5m',
    },
    {
      id: '2',
      event: 'Suscripción Pro activada',
      user: 'maria@example.com',
      time: 'Hace 1h',
    },
    {
      id: '3',
      event: 'API Token creado',
      user: 'john@example.com',
      time: 'Hace 2h',
    },
    {
      id: '4',
      event: 'Límite de cuota alcanzado',
      user: 'company@example.com',
      time: 'Hace 4h',
    },
  ];

  return (
    <div className="h-full overflow-y-auto no-scrollbar p-6 lg:p-8 space-y-8 bg-zinc-50/50 dark:bg-zinc-950/20">
      {/* Header */}
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div>
          <h1 className="text-3xl font-black tracking-tight text-title">
            Panel de Control
          </h1>
          <p className="text-subtitle mt-1 font-medium">
            ¡Hola de nuevo! Aquí tienes un resumen de la plataforma.
          </p>
        </div>
        <div className="flex items-center gap-3">
          <div className="relative group hidden sm:block">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground group-focus-within:text-brand transition-colors" />
            <input
              type="text"
              placeholder="Buscar..."
              className="pl-9 pr-4 py-2 bg-white dark:bg-zinc-900 border border-border/60 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-brand/20 w-64 transition-all"
            />
          </div>
          <Button
            variant="outline"
            size="icon"
            className="rounded-xl border-border/60 bg-white dark:bg-zinc-900 relative"
          >
            <Bell className="h-5 w-5" />
            <span className="absolute top-2 right-2 h-2 w-2 bg-red-500 rounded-full border-2 border-white dark:border-zinc-900"></span>
          </Button>
        </div>
      </div>

      {/* Grid 1: Stats Rápidas */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat, i) => (
          <div
            key={i}
            className="bg-white dark:bg-card border border-border/50 rounded-2xl p-5 shadow-sm hover:shadow-md transition-all group"
          >
            <div className="flex justify-between items-start">
              <div className={`p-2.5 rounded-xl ${stat.bg} ${stat.color}`}>
                <stat.icon className="h-5 w-5" />
              </div>
              <Badge
                variant="secondary"
                className={`bg-transparent border-0 font-bold ${stat.trendUp ? 'text-emerald-600' : 'text-red-600'}`}
              >
                {stat.trendUp ? (
                  <ArrowUpRight className="h-3 w-3 mr-1" />
                ) : (
                  <ArrowDownRight className="h-3 w-3 mr-1" />
                )}
                {stat.trend}
              </Badge>
            </div>
            <div className="mt-4">
              <p className="text-xs font-bold text-muted-foreground uppercase tracking-wider">
                {stat.label}
              </p>
              <h3 className="text-2xl font-black mt-1 text-foreground">
                {stat.value}
              </h3>
            </div>
          </div>
        ))}
      </div>

      {/* Grid 2: Bento Mosaic Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 pb-4">
        {/* Visualización de Rendimiento */}
        <div className="lg:col-span-8 bg-white dark:bg-card border border-border/50 rounded-3xl p-6 shadow-sm overflow-hidden flex flex-col min-h-[400px]">
          <div className="flex justify-between items-center mb-6">
            <div>
              <h3 className="text-xl font-bold flex items-center gap-2">
                Rendimiento de la Plataforma
                <Badge
                  variant="outline"
                  className="text-[10px] bg-brand/5 text-brand border-brand/20"
                >
                  En vivo
                </Badge>
              </h3>
              <p className="text-xs text-muted-foreground font-medium mt-1">
                Visualización general de rendimiento y uso de API.
              </p>
            </div>
          </div>

          <div className="flex-1 bg-gradient-to-b from-brand/5 to-transparent rounded-2xl border border-dashed border-brand/20 p-8 flex flex-col items-center justify-center relative">
            <div className="absolute inset-x-0 bottom-0 h-32 bg-[radial-gradient(ellipse_at_center,_var(--tw-gradient-stops))] from-brand/10 via-transparent to-transparent opacity-50"></div>
            <TrendingUp className="h-12 w-12 text-brand opacity-20 mb-4 animate-bounce" />
            <div className="text-center">
              <p className="text-sm font-bold text-brand/60 uppercase tracking-widest">
                Esqueleto de Gráficos Listo
              </p>
              <p className="text-xs text-muted-foreground mt-2 max-w-[250px] mx-auto text-balance">
                Conecta aquí tu librería favorita de gráficos (Recharts, ChartJS, etc.) para monitorear el comportamiento de tu SaaS.
              </p>
            </div>

            <div className="absolute bottom-6 left-6 right-6 flex items-end justify-between h-20 gap-2 opacity-20">
              {[40, 70, 45, 90, 65, 80, 55, 95, 75, 60, 85, 50].map((h, i) => (
                <div
                  key={i}
                  className="flex-1 bg-brand rounded-t-lg transition-all duration-1000"
                  style={{ height: `${h}%` }}
                ></div>
              ))}
            </div>
          </div>
        </div>

        {/* Actividad Reciente */}
        <div className="lg:col-span-4 space-y-6">
          <div className="bg-white dark:bg-card border border-border/50 rounded-3xl p-6 shadow-sm flex flex-col h-full overflow-hidden justify-between">
            <div>
              <div className="flex justify-between items-center mb-6 px-1">
                <h3 className="text-lg font-bold">Eventos Recientes</h3>
              </div>

              <div className="space-y-5 overflow-y-auto no-scrollbar max-h-[300px]">
                {recentEvents.map((ev) => (
                  <div key={ev.id} className="flex gap-4 group cursor-pointer px-1">
                    <div className="h-10 w-10 rounded-full bg-brand/5 flex items-center justify-center font-bold text-xs ring-2 ring-transparent group-hover:ring-brand/20 transition-all shrink-0 text-brand">
                      {ev.user.charAt(0).toUpperCase()}
                    </div>
                    <div className="flex-1 border-b border-border/30 pb-4 last:border-0 min-w-0">
                      <div className="flex justify-between items-start gap-2">
                        <div className="min-w-0">
                          <p className="text-sm font-bold text-foreground truncate">
                            {ev.event}
                          </p>
                          <p className="text-xs text-muted-foreground truncate">
                            {ev.user}
                          </p>
                        </div>
                      </div>
                      <div className="flex items-center gap-1.5 mt-1">
                        <Clock className="h-3 w-3 text-muted-foreground" />
                        <span className="text-[10px] text-muted-foreground font-medium uppercase tracking-tight">
                          {ev.time}
                        </span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <Button className="w-full mt-6 bg-zinc-900 dark:bg-zinc-100 dark:text-zinc-900 hover:bg-zinc-800 rounded-2xl h-12 font-bold shadow-lg shadow-zinc-200 dark:shadow-none transition-all hover:-translate-y-1">
              Ver Logs Completos
              <ChevronRight className="ml-2 h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
