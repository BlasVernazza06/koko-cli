'use client';

import { useState } from 'react';

import {
  BookOpen,
  ChevronDown,
  ChevronUp,
  CreditCard,
  LifeBuoy,
  MessageCircleQuestion,
  Search,
  Settings,
  ShieldCheck,
  Zap,
} from 'lucide-react';

import { Badge } from '@repo/ui/components/ui/badge';

export default function HelpPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [openFaq, setOpenFaq] = useState<number | null>(null);

  const categories = [
    {
      title: 'Primeros Pasos',
      description:
        'Aprende lo básico para configurar tu cuenta y empezar.',
      icon: Zap,
      color: 'text-amber-500',
      bg: 'bg-amber-500/10',
    },
    {
      title: 'Facturación',
      description:
        'Gestiona tus suscripciones, métodos de pago y descarga tus facturas.',
      icon: CreditCard,
      color: 'text-emerald-500',
      bg: 'bg-emerald-500/10',
    },
    {
      title: 'Configuración',
      description:
        'Personaliza tu perfil, notificaciones y preferencias de la plataforma.',
      icon: Settings,
      color: 'text-blue-500',
      bg: 'bg-blue-500/10',
    },
    {
      title: 'Seguridad',
      description:
        'Protege tu cuenta con autenticación de dos factores y gestión de accesos.',
      icon: ShieldCheck,
      color: 'text-purple-500',
      bg: 'bg-purple-500/10',
    },
  ];

  const faqs = [
    {
      question: '¿Cómo canjeo un código de descuento?',
      answer:
        "Puedes canjear códigos de descuento en la sección de 'Configuración > Plan', ingresando el código antes de confirmar tu suscripción o actualización de plan.",
    },
    {
      question: '¿Puedo exportar mis datos de ventas?',
      answer:
        "Sí, en la sección de 'Analíticas' encontrarás un botón de exportar que te permite descargar tus reportes en formato CSV o PDF.",
    },
    {
      question: '¿Cómo agrego nuevos miembros a mi equipo?',
      answer:
        "Actualmente, la gestión de equipos está disponible en los planes PRO y Enterprise. Puedes hacerlo desde 'Configuración > Equipo'.",
    },
    {
      question: '¿Qué métodos de pago aceptan?',
      answer:
        'Aceptamos todas las tarjetas de crédito y débito principales (Visa, Mastercard, Amex) a través de Stripe, además de transferencias bancarias según tu región.',
    },
  ];

  const toggleFaq = (index: number) => {
    setOpenFaq(openFaq === index ? null : index);
  };

  return (
    <div className="h-full overflow-y-auto no-scrollbar p-6 lg:p-8 space-y-12 bg-zinc-50/50 dark:bg-zinc-950/20">
      {/* Hero Header Section */}
      <div className="relative overflow-hidden rounded-[2.5rem] bg-zinc-900 dark:bg-zinc-900 p-8 md:p-16 text-center shadow-2xl">
        <div className="absolute top-0 left-0 w-full h-full opacity-10 pointer-events-none">
          <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-brand rounded-full blur-[100px]"></div>
          <div className="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-blue-500 rounded-full blur-[100px]"></div>
        </div>

        <div className="relative z-10 max-w-2xl mx-auto space-y-6">
          <Badge
            variant="outline"
            className="bg-brand/10 text-brand border-brand/20 px-3 py-1 font-bold tracking-wider uppercase text-[10px]"
          >
            Centro de Ayuda de [[.ProjectName]]
          </Badge>
          <h1 className="text-4xl md:text-5xl font-black text-white tracking-tight">
            ¿Cómo podemos ayudarte?
          </h1>
          <p className="text-zinc-400 text-lg font-medium">
            Busca en nuestra documentación o explora las categorías para
            resolver tus dudas rápidamente.
          </p>

          <div className="relative group max-w-xl mx-auto mt-8">
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-zinc-500 group-focus-within:text-brand transition-colors" />
            <input
              type="text"
              placeholder="Ej: '¿Cómo conectar mi base de datos?'"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-12 pr-4 py-4 bg-zinc-800/50 border border-zinc-700 rounded-2xl text-white placeholder:text-zinc-500 focus:outline-none focus:ring-2 focus:ring-brand/40 transition-all text-lg shadow-inner"
            />
          </div>
        </div>
      </div>

      {/* Knowledge Categories */}
      <div className="space-y-6">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-brand/10 rounded-xl">
            <BookOpen className="h-5 w-5 text-brand" />
          </div>
          <h2 className="text-2xl font-black text-title italic">
            Explorar Categorías
          </h2>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {categories.map((cat, i) => (
            <div
              key={i}
              className="group bg-white dark:bg-card border border-border/50 rounded-3xl p-6 shadow-sm hover:shadow-xl hover:-translate-y-1 transition-all cursor-pointer"
            >
              <div
                className={`p-3 rounded-2xl ${cat.bg} ${cat.color} w-fit mb-4 group-hover:scale-110 transition-transform`}
              >
                <cat.icon className="h-6 w-6" />
              </div>
              <h3 className="text-lg font-bold mb-2 text-foreground">
                {cat.title}
              </h3>
              <p className="text-sm text-subtitle font-medium leading-relaxed">
                {cat.description}
              </p>
            </div>
          ))}
        </div>
      </div>

      {/* FAQ Section */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-12 pt-4">
        <div className="lg:col-span-4 space-y-4">
          <div className="sticky top-8">
            <div className="p-2 bg-emerald-500/10 rounded-xl w-fit mb-4">
              <MessageCircleQuestion className="h-5 w-5 text-emerald-600" />
            </div>
            <h2 className="text-3xl font-black text-title leading-tight italic">
              Preguntas <br /> Frecuentes
            </h2>
            <p className="text-subtitle font-medium mt-4">
              Respuestas rápidas a las consultas más comunes de nuestros
              usuarios.
            </p>
            <div className="mt-8 p-6 bg-brand/5 border border-brand/10 rounded-3xl space-y-3">
              <LifeBuoy className="h-8 w-8 text-brand opacity-40" />
              <p className="text-xs font-bold text-brand uppercase tracking-widest leading-normal">
                [[.ProjectName]] es intuitivo y fácil de usar. Cada sección cuenta con
                guías interactivas.
              </p>
            </div>
          </div>
        </div>

        <div className="lg:col-span-8 space-y-4">
          {faqs.map((faq, i) => (
            <div
              key={i}
              className={`bg-white dark:bg-card border border-border/50 rounded-2xl overflow-hidden transition-all ${openFaq === i ? 'ring-2 ring-brand/20 shadow-lg' : 'hover:border-brand/30'}`}
            >
              <button
                onClick={() => toggleFaq(i)}
                className="w-full flex items-center justify-between p-6 text-left focus:outline-none"
              >
                <span className="text-lg font-bold text-foreground pr-4">
                  {faq.question}
                </span>
                <div
                  className={`p-2 rounded-xl transition-colors ${openFaq === i ? 'bg-brand/10 text-brand' : 'bg-muted/50 text-muted-foreground'}`}
                >
                  {openFaq === i ? (
                    <ChevronUp className="h-4 w-4" />
                  ) : (
                    <ChevronDown className="h-4 w-4" />
                  )}
                </div>
              </button>

              <div
                className={`overflow-hidden transition-all duration-300 ${openFaq === i ? 'max-h-40' : 'max-h-0'}`}
              >
                <div className="p-6 pt-0 text-subtitle font-medium">
                  {faq.answer}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Footer Note - Simplified */}
      <div className="pb-12 text-center">
        <div className="h-px bg-linear-to-r from-transparent via-border to-transparent mb-8"></div>
        <p className="text-sm text-muted-foreground font-medium">
          ¿No encuentras lo que buscas? Revisa nuestra{' '}
          <span className="text-brand font-bold cursor-pointer hover:underline">
            Documentación Técnica
          </span>{' '}
          completa.
        </p>
      </div>
    </div>
  );
}
