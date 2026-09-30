'use client';

import { useRouter } from 'next/navigation';

import { useEffect, useState } from 'react';

import { Loader2, LogOut, Boxes } from 'lucide-react';
import { AnimatePresence, motion } from 'motion/react';

import { authClient, useSession } from '@repo/auth/client';
import { Button } from '@repo/ui/components/ui/button';

import CreateTenantForm from '@/components/onboarding/create-tenant-form';
import DecisionStep from '@/components/onboarding/decision-step';

export default function OnboardingPage() {
  const router = useRouter();
  const { data: session, isPending: isSessionLoading } = useSession();
  const [isCheckingTenants, setIsCheckingTenants] = useState(true);
  const [isCreating, setIsCreating] = useState(false);
  const [step, setStep] = useState<'decision' | 'create'>('decision');

  useEffect(() => {
    async function checkTenants() {
      if (isSessionLoading || !session) return;

      try {
        const response = await fetch('/api/tenants/me');

        const tenants = await response.json();
        if (tenants && tenants.length > 0) {
          router.push('/dashboard');
        } else {
          setIsCheckingTenants(false);
        }
      } catch (error) {
        console.error('Error checking tenants:', error);
        setIsCheckingTenants(false);
      }
    }

    checkTenants();
  }, [session, isSessionLoading, router]);

  const handleLogout = async () => {
    await authClient.signOut({
      fetchOptions: {
        onSuccess: () => {
          router.push('/auth/login');
        },
      },
    });
  };

  if (isSessionLoading || isCheckingTenants) {
    return (
      <div className="h-screen w-full flex flex-col items-center justify-center gap-4 bg-zinc-50 dark:bg-zinc-950">
        <Loader2 className="size-8 animate-spin text-brand" />
        <p className="text-sm font-medium text-muted-foreground animate-pulse">
          Preparando tu espacio...
        </p>
      </div>
    );
  }

  return (
    <div className="min-h-screen w-full bg-zinc-50 dark:bg-zinc-950 flex flex-col items-center justify-center p-6 bg-[radial-gradient(#e5e7eb_1px,transparent_1px)] dark:bg-[radial-gradient(#18181b_1px,transparent_1px)] [background-size:16px_16px]">
      <div className="absolute top-8 right-8">
        <Button
          variant="ghost"
          size="sm"
          onClick={handleLogout}
          className="text-muted-foreground hover:text-foreground font-bold"
        >
          <LogOut className="size-4 mr-2" />
          Cerrar Sesión
        </Button>
      </div>

      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="max-w-md w-full"
      >
        <div className="text-center mb-10">
          <div className="h-16 w-16 bg-brand/10 text-brand rounded-2xl flex items-center justify-center mx-auto mb-6">
            <Boxes size={32} />
          </div>
          <h1 className="text-3xl font-black tracking-tight text-foreground mb-2">
            ¡Hola, {session?.user.name.split(' ')[0]}!
          </h1>
          <p className="text-muted-foreground font-medium text-balance">
            Parece que aún no eres parte de ninguna organización. ¿Cómo quieres
            empezar hoy?
          </p>
        </div>

        <AnimatePresence mode="wait">
          {step === 'decision' ? (
            <DecisionStep session={session} setStep={setStep} />
          ) : (
            <CreateTenantForm setStep={setStep} />
          )}
        </AnimatePresence>
      </motion.div>
    </div>
  );
}
