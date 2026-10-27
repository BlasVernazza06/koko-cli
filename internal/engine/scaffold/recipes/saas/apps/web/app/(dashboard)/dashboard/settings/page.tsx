'use client';

import { useState } from 'react';

import { AnimatePresence, motion } from 'framer-motion';
import {
  Bell,
  Building2,
  Camera,
  Globe,
  Lock,
  Mail,
  Palette,
  Save,
  Search,
  Shield,
  Smartphone,
  User,
} from 'lucide-react';

import { Badge } from '@repo/ui/components/ui/badge';
import { Button } from '@repo/ui/components/ui/button';
import { Input } from '@repo/ui/components/ui/input';
import { Label } from '@repo/ui/components/ui/label';
import { Separator } from '@repo/ui/components/ui/separator';

// --- Mock Components for Settings ---

const CustomSwitch = ({
  checked,
  onChange,
}: {
  checked: boolean;
  onChange: (val: boolean) => void;
}) => (
  <div
    onClick={() => onChange(!checked)}
    className={`w-12 h-6 rounded-full transition-colors cursor-pointer relative ${checked ? 'bg-brand' : 'bg-zinc-200 dark:bg-zinc-800'}`}
  >
    <motion.div
      animate={{ x: checked ? 26 : 2 }}
      className="absolute top-1 left-0 w-4 h-4 rounded-full bg-white shadow-sm"
      transition={{ type: 'spring', stiffness: 500, damping: 30 }}
    />
  </div>
);

interface SettingItemProps {
  icon: any;
  title: string;
  description: string;
  children: React.ReactNode;
  delay?: number;
}

const SettingItem = ({
  icon: Icon,
  title,
  description,
  children,
  delay = 0,
}: SettingItemProps) => (
  <motion.div
    initial={{ opacity: 0, y: 10 }}
    animate={{ opacity: 1, y: 0 }}
    transition={{ delay }}
    className="flex items-center justify-between p-6 bg-white dark:bg-card border border-border/40 rounded-3xl hover:border-brand/30 transition-all group"
  >
    <div className="flex items-center gap-4">
      <div className="h-12 w-12 rounded-2xl bg-zinc-50 dark:bg-zinc-900 flex items-center justify-center text-subtitle group-hover:text-brand transition-colors">
        <Icon size={22} />
      </div>
      <div>
        <h4 className="font-bold text-title">{title}</h4>
        <p className="text-xs text-subtitle">{description}</p>
      </div>
    </div>
    <div>{children}</div>
  </motion.div>
);

export default function ConfiguracionPage() {
  const [activeTab, setActiveTab] = useState('perfil');
  const [notis, setNotis] = useState({
    email: true,
    push: false,
    orders: true,
  });

  const menuItems = [
    { id: 'perfil', icon: User, label: 'Perfil' },
    { id: 'negocio', icon: Building2, label: 'Negocio' },
    { id: 'notificaciones', icon: Bell, label: 'Notificaciones' },
    { id: 'seguridad', icon: Lock, label: 'Seguridad' },
    { id: 'apariencia', icon: Palette, label: 'Apariencia' },
  ];

  return (
    <div className="h-full flex flex-col bg-zinc-50/50 dark:bg-zinc-950/20 overflow-y-auto no-scrollbar pb-12">
      {/* Header */}
      <header className="p-8 pb-4 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-6">
        <div>
          <h1 className="text-3xl font-black tracking-tight text-title">
            Configuración
          </h1>
          <p className="text-subtitle font-medium">
            Personaliza tu experiencia y gestiona tu cuenta.
          </p>
        </div>
        <div className="flex items-center gap-3">
          <Button
            variant="outline"
            className="rounded-xl border-border/60 bg-white dark:bg-zinc-900 h-11 px-6 font-bold text-xs uppercase tracking-widest"
          >
            Descartar
          </Button>
          <Button className="bg-brand hover:bg-brand/90 text-white rounded-2xl h-11 px-8 shadow-lg shadow-brand/20 font-bold">
            <Save className="mr-2 h-4 w-4" /> Guardar Cambios
          </Button>
        </div>
      </header>

      <div className="px-8 grid grid-cols-1 lg:grid-cols-4 gap-8">
        {/* Sidebar Menu */}
        <div className="lg:col-span-1 space-y-2">
          {menuItems.map((item) => (
            <button
              key={item.id}
              onClick={() => setActiveTab(item.id)}
              className={`w-full flex items-center gap-3 p-4 rounded-2xl transition-all font-bold text-sm ${
                activeTab === item.id
                  ? 'bg-white dark:bg-card text-brand shadow-sm shadow-brand/5 border border-brand/20'
                  : 'text-subtitle hover:bg-white dark:hover:bg-card/50'
              }`}
            >
              <item.icon
                size={18}
                className={
                  activeTab === item.id ? 'text-brand' : 'text-subtitle'
                }
              />
              {item.label}
              {activeTab === item.id && (
                <motion.div
                  layoutId="activeDot"
                  className="ml-auto w-1.5 h-1.5 rounded-full bg-brand"
                />
              )}
            </button>
          ))}
        </div>

        {/* Main Content Area */}
        <div className="lg:col-span-3 space-y-6">
          <AnimatePresence mode="wait">
            {activeTab === 'perfil' && (
              <motion.div
                key="perfil"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                className="space-y-6"
              >
                {/* Profile Card */}
                <div className="bg-white dark:bg-card border border-border/50 rounded-[32px] p-8 shadow-sm">
                  <div className="flex flex-col md:flex-row items-center gap-8">
                    <div className="relative group">
                      <div className="h-32 w-32 rounded-[40px] bg-linear-to-br from-brand/20 to-brand/5 border-4 border-white dark:border-zinc-900 overflow-hidden shadow-xl">
                        <img
                          src="https://api.dicebear.com/7.x/avataaars/svg?seed=Felix"
                          alt="Avatar"
                          className="w-full h-full object-cover"
                        />
                      </div>
                      <button className="absolute -bottom-2 -right-2 h-10 w-10 bg-brand text-white rounded-2xl flex items-center justify-center shadow-lg hover:scale-110 transition-transform">
                        <Camera size={18} />
                      </button>
                    </div>
                    <div className="flex-1 space-y-4">
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div className="space-y-2">
                          <Label className="text-xs font-black uppercase tracking-widest text-subtitle">
                            Nombre Completo
                          </Label>
                          <Input
                            defaultValue="Felix Anderson"
                            className="rounded-2xl h-12 border-border/60 bg-zinc-50/50 dark:bg-zinc-900/50"
                          />
                        </div>
                        <div className="space-y-2">
                          <Label className="text-xs font-black uppercase tracking-widest text-subtitle">
                            Email Personal
                          </Label>
                          <Input
                            defaultValue="felix@example.com"
                            className="rounded-2xl h-12 border-border/60 bg-zinc-50/50 dark:bg-zinc-900/50"
                          />
                        </div>
                        <div className="space-y-2">
                          <Label className="text-xs font-black uppercase tracking-widest text-subtitle">
                            Cargo / Rol
                          </Label>
                          <Input
                            defaultValue="Lead Designer & Developer"
                            className="rounded-2xl h-12 border-border/60 bg-zinc-50/50 dark:bg-zinc-900/50"
                          />
                        </div>
                        <div className="space-y-2">
                          <Label className="text-xs font-black uppercase tracking-widest text-subtitle">
                            Ubicación
                          </Label>
                          <div className="relative">
                            <Globe
                              className="absolute left-4 top-1/2 -translate-y-1/2 text-subtitle"
                              size={16}
                            />
                            <Input
                              defaultValue="Buenos Aires, Argentina"
                              className="pl-12 rounded-2xl h-12 border-border/60 bg-zinc-50/50 dark:bg-zinc-900/50"
                            />
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                {/* Preferences */}
                <div className="space-y-4">
                  <h3 className="text-lg font-black text-title ml-2">
                    Preferencias de Cuenta
                  </h3>
                  <SettingItem
                    icon={Mail}
                    title="Email Público"
                    description="Mostrar tu email en tu perfil público para que otros te contacten."
                    delay={0.1}
                  >
                    <CustomSwitch checked={true} onChange={() => {}} />
                  </SettingItem>
                  <SettingItem
                    icon={Smartphone}
                    title="Sincronización Móvil"
                    description="Mantener tus datos sincronizados en tiempo real entre dispositivos."
                    delay={0.2}
                  >
                    <CustomSwitch checked={true} onChange={() => {}} />
                  </SettingItem>
                </div>
              </motion.div>
            )}

            {activeTab === 'notificaciones' && (
              <motion.div
                key="notis"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                className="space-y-4"
              >
                <h3 className="text-lg font-black text-title ml-2">
                  Alertas y Notificaciones
                </h3>
                <SettingItem
                  icon={Bell}
                  title="Notificaciones por Email"
                  description="Recibe resúmenes semanales y alertas críticas en tu bandeja."
                  delay={0.1}
                >
                  <CustomSwitch
                    checked={notis.email}
                    onChange={(v) => setNotis({ ...notis, email: v })}
                  />
                </SettingItem>
                <SettingItem
                  icon={Smartphone}
                  title="Notificaciones Push"
                  description="Alertas instantáneas en el navegador y dispositivos móviles."
                  delay={0.2}
                >
                  <CustomSwitch
                    checked={notis.push}
                    onChange={(v) => setNotis({ ...notis, push: v })}
                  />
                </SettingItem>

                <div className="p-8 bg-amber-500/5 border border-amber-500/20 rounded-[32px] mt-8">
                  <div className="flex gap-4">
                    <div className="h-10 w-10 rounded-xl bg-amber-500 flex items-center justify-center text-white shrink-0">
                      <Bell size={20} />
                    </div>
                    <div>
                      <h4 className="font-bold text-amber-900 dark:text-amber-200">
                        ¿Silenciar temporalmente?
                      </h4>
                      <p className="text-sm text-amber-800/70 dark:text-amber-200/60 mt-1">
                        Puedes activar el modo "No molestar" para pausar todas
                        las alertas hasta mañana.
                      </p>
                      <Button className="mt-4 bg-amber-500 hover:bg-amber-600 text-white rounded-xl font-bold text-xs h-10 px-6">
                        Activar Modo Silencio
                      </Button>
                    </div>
                  </div>
                </div>
              </motion.div>
            )}

            {activeTab === 'seguridad' && (
              <motion.div
                key="seguridad"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                className="space-y-6"
              >
                <div className="bg-white dark:bg-card border border-border/50 rounded-[32px] p-8 shadow-sm">
                  <div className="flex items-center gap-4 mb-8">
                    <div className="h-12 w-12 rounded-2xl bg-emerald-500/10 text-emerald-500 flex items-center justify-center">
                      <Shield size={24} />
                    </div>
                    <div>
                      <h3 className="text-xl font-black text-title">
                        Seguridad de la Cuenta
                      </h3>
                      <p className="text-sm text-subtitle">
                        Gestiona tu contraseña y métodos de autenticación.
                      </p>
                    </div>
                  </div>

                  <div className="space-y-4">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="space-y-2">
                        <Label className="text-xs font-black uppercase tracking-widest text-subtitle">
                          Contraseña Actual
                        </Label>
                        <Input
                          type="password"
                          placeholder="••••••••"
                          className="rounded-2xl h-12 border-border/60 bg-zinc-50/50 dark:bg-zinc-900/50"
                        />
                      </div>
                      <div className="space-y-2">
                        <Label className="text-xs font-black uppercase tracking-widest text-subtitle">
                          Nueva Contraseña
                        </Label>
                        <Input
                          type="password"
                          placeholder="••••••••"
                          className="rounded-2xl h-12 border-border/60 bg-zinc-50/50 dark:bg-zinc-900/50"
                        />
                      </div>
                    </div>
                    <Button className="bg-zinc-900 dark:bg-brand text-white rounded-2xl h-11 px-8 font-bold">
                      Actualizar Contraseña
                    </Button>
                  </div>

                  <Separator className="my-8 opacity-50" />

                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <div className="h-10 w-10 rounded-xl bg-blue-500/10 text-blue-500 flex items-center justify-center">
                        <Lock size={20} />
                      </div>
                      <div>
                        <h4 className="font-bold text-title">
                          Autenticación de 2 Factores
                        </h4>
                        <p className="text-xs text-subtitle">
                          Añade una capa extra de seguridad a tu cuenta.
                        </p>
                      </div>
                    </div>
                    <Badge className="bg-emerald-500/10 text-emerald-600 border-none px-4 py-1.5 rounded-full font-bold">
                      Activo
                    </Badge>
                  </div>
                </div>

                <div className="p-8 border-2 border-dashed border-border/50 rounded-[32px] flex flex-col items-center justify-center text-center space-y-4">
                  <div className="h-16 w-16 rounded-full bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center text-subtitle">
                    <Smartphone size={32} />
                  </div>
                  <div>
                    <h4 className="font-bold text-title">
                      Dispositivos Sesión
                    </h4>
                    <p className="text-xs text-subtitle max-w-[300px]">
                      Actualmente tienes 3 sesiones activas en diferentes
                      dispositivos.
                    </p>
                  </div>
                  <Button
                    variant="outline"
                    className="rounded-xl border-border/60 h-10 px-6 font-bold text-xs uppercase tracking-widest"
                  >
                    Cerrar Todas las Sesiones
                  </Button>
                </div>
              </motion.div>
            )}

            {/* Additional tabs can follow the same pattern */}
            {(activeTab === 'negocio' || activeTab === 'apariencia') && (
              <motion.div
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                className="h-[400px] flex flex-col items-center justify-center text-center p-8 bg-white dark:bg-card border border-border/50 rounded-[32px] shadow-sm"
              >
                <div className="h-20 w-20 rounded-3xl bg-brand/10 text-brand flex items-center justify-center mb-6">
                  <Search size={40} />
                </div>
                <h3 className="text-2xl font-black text-title">
                  Sección en Construcción
                </h3>
                <p className="text-subtitle mt-2 max-w-[400px]">
                  Estamos trabajando para traerte más opciones de
                  personalización muy pronto.
                </p>
                <Button
                  onClick={() => setActiveTab('perfil')}
                  className="mt-8 bg-brand hover:bg-brand/90 text-white rounded-2xl h-12 px-8 shadow-lg shadow-brand/20 font-bold"
                >
                  Volver al Perfil
                </Button>
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </div>
    </div>
  );
}
