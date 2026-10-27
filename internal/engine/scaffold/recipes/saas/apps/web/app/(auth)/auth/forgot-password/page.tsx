'use client';

import Link from 'next/link';

import { useState } from 'react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { CheckCircle2, Loader2 } from 'lucide-react';

import { Button } from '@repo/ui/components/ui/button';
import { Input } from '@repo/ui/components/ui/input';
import { Label } from '@repo/ui/components/ui/label';
import { authClient } from '@repo/auth/client';
import {
  RequestPasswordResetFormValues,
  requestPasswordResetSchema,
} from '@repo/validators/auth';

export default function ForgotPasswordPage() {
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [isSubmitted, setIsSubmitted] = useState<boolean>(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors },
    getValues,
  } = useForm<RequestPasswordResetFormValues>({
    resolver: zodResolver(requestPasswordResetSchema),
    defaultValues: {
      email: '',
    },
  });

  const onSubmit = async (values: RequestPasswordResetFormValues) => {
    setIsLoading(true);
    setErrorMsg(null);
    try {
      const { error } = await authClient.requestPasswordReset({
        email: values.email,
        redirectTo: `${window.location.origin}/auth/reset-password`,
      });
      if (error) {
        setErrorMsg(error.message);
      } else {
        setIsSubmitted(true);
      }
    } catch (err) {
      setErrorMsg('Ocurrió un error inesperado');
    } finally {
      setIsLoading(false);
    }
  };

  if (isSubmitted) {
    return (
      <div className="text-center space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
        <div className="w-16 h-16 bg-green-100 dark:bg-green-950/20 text-green-600 rounded-full flex items-center justify-center mx-auto mb-4">
          <CheckCircle2 size={32} />
        </div>
        <h2 className="text-2xl font-bold text-foreground">¡Correo enviado!</h2>
        <p className="text-muted-foreground max-w-xs mx-auto">
          Hemos enviado las instrucciones para restablecer tu contraseña a{' '}
          <span className="font-medium text-foreground">
            {getValues('email')}
          </span>
        </p>
        <Link href="/auth/login" className="w-full block">
          <Button variant="outline" className="mt-4 w-full h-11 rounded-xl">
            Volver al inicio de sesión
          </Button>
        </Link>
      </div>
    );
  }

  return (
    <div className="space-y-6 w-full font-sans">
      <header className="flex flex-col items-center gap-1 text-center mb-6">
        <h1 className="text-3xl font-black tracking-tight text-foreground">¿Olvidaste tu contraseña?</h1>
        <p className="text-muted-foreground text-sm leading-relaxed">
          Ingresa tu correo electrónico y te enviaremos un enlace para restablecerla.
        </p>
      </header>

      {errorMsg && (
        <div className="p-3 bg-red-500/10 border border-red-500/20 text-red-500 text-xs rounded-xl font-medium">
          {errorMsg}
        </div>
      )}

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="email" className="text-xs font-bold text-muted-foreground ml-1 uppercase tracking-widest">
            Email
          </Label>
          <div className="relative group">
            <Input
              id="email"
              type="email"
              {...register('email')}
              placeholder="nombre@empresa.com"
              className={`bg-zinc-50 dark:bg-zinc-900 border-zinc-200 dark:border-zinc-800 h-12 rounded-xl px-4 focus:ring-brand focus:border-brand transition-all text-sm font-medium ${errors.email ? 'border-red-500 focus:border-red-500' : ''}`}
            />
          </div>
          {errors.email && (
            <p className="text-[11px] text-red-500 font-medium ml-1">
              {errors.email.message}
            </p>
          )}
        </div>

        <Button type="submit" disabled={isLoading} className="cursor-pointer h-12 rounded-xl mt-5 w-full bg-[#a855f7] hover:bg-[#9333ea] text-white font-bold">
          {isLoading ? (
            <Loader2 className="h-5 w-5 animate-spin mx-auto" />
          ) : (
            'Enviar Instrucciones'
          )}
        </Button>
      </form>

      <div className="text-center pt-4">
        <Link href="/auth/login" className="text-xs font-bold text-brand hover:underline">
          Volver al Login
        </Link>
      </div>
    </div>
  );
}
