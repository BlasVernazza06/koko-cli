import { headers } from 'next/headers';
import Link from 'next/link';
import { redirect } from 'next/navigation';

export default async function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen w-full flex flex-col md:flex-row bg-background overflow-hidden font-sans">
      
      {/* Left Panel: Visual & Branding (Hidden on mobile) */}
      <div className="hidden md:flex md:w-1/2 lg:w-[60%] relative overflow-hidden bg-zinc-950">
        {/* Modern animated-like gradient background */}
        <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top_right,_var(--tw-gradient-stops))] from-zinc-850 via-zinc-900 to-black opacity-80" />
        <div className="absolute top-[-20%] left-[-20%] w-[80%] h-[80%] rounded-full bg-indigo-500/10 blur-[120px]" />
        <div className="absolute bottom-[-20%] right-[-20%] w-[80%] h-[80%] rounded-full bg-violet-600/15 blur-[120px]" />

        <div className="relative z-10 w-full h-full flex flex-col justify-between p-12 lg:p-20 text-white">
          <Link href="/" className="flex items-center gap-3 group">
            <div className="h-10 w-10 rounded-lg bg-zinc-800 border border-zinc-700 flex items-center justify-center text-white shadow-xl group-hover:scale-105 transition-transform duration-300">
              <span className="font-extrabold text-xl">[[.ProjectName]]</span>
            </div>
          </Link>

          <div className="max-w-lg">
            <h2 className="text-4xl lg:text-5xl font-black leading-tight mb-6 animate-in slide-in-from-left duration-700">
              Construye tu SaaS a la <span className="text-indigo-400">velocidad de la luz</span>.
            </h2>
            <p className="text-lg lg:text-xl text-zinc-400 font-medium leading-relaxed mb-8">
              Una arquitectura modular y robusta con base de datos preconfigurada, autenticación y diseño responsivo para iniciar de inmediato.
            </p>
            
            <div className="flex gap-12 pt-8 border-t border-zinc-850">
              <div>
                <p className="text-3xl font-black text-white">100%</p>
                <p className="text-xs uppercase tracking-widest text-zinc-500 font-bold">Producción Ready</p>
              </div>
              <div>
                <p className="text-3xl font-black text-white">Drizzle</p>
                <p className="text-xs uppercase tracking-widest text-zinc-500 font-bold">ORM</p>
              </div>
              <div>
                <p className="text-3xl font-black text-white">BetterAuth</p>
                <p className="text-xs uppercase tracking-widest text-zinc-500 font-bold">Seguridad</p>
              </div>
            </div>
          </div>

          <div className="text-xs font-semibold text-zinc-600 uppercase tracking-widest">
            © 2026 [[.ProjectName]] Inc. All rights reserved.
          </div>
        </div>
      </div>

      {/* Right Panel: Auth Form */}
      <div className="flex-1 flex flex-col items-center justify-center p-6 md:p-12 lg:p-20 bg-zinc-950 relative">
        <div className="w-full max-w-sm relative z-10">
          {/* Mobile Logo */}
          <div className="flex md:hidden justify-center mb-10">
            <Link href="/" className="flex flex-col items-center gap-2">
              <div className="h-12 w-12 rounded-xl bg-zinc-800 border border-zinc-700 flex items-center justify-center text-white">
                <span className="font-extrabold text-xl">[[.ProjectName]]</span>
              </div>
            </Link>
          </div>

          <div className="animate-in fade-in zoom-in-95 duration-500">
            {children}
          </div>
        </div>

        {/* Decorative background for right panel (subtle) */}
        <div className="absolute top-0 right-0 w-64 h-64 bg-indigo-500/5 blur-[80px] rounded-full pointer-events-none" />
        <div className="absolute bottom-0 left-0 w-64 h-64 bg-violet-500/5 blur-[80px] rounded-full pointer-events-none" />
      </div>
    </div>
  );
}
