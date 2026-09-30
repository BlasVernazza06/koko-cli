'use client';

import Link from 'next/link';
import { useSearchParams } from 'next/navigation';

import { useState } from 'react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { AlertCircle, CheckCircle2, Eye, EyeOff, Loader2 } from 'lucide-react';

import { Button } from '@repo/ui/components/ui/button';
import { Input } from '@repo/ui/components/ui/input';
import { Label } from '@repo/ui/components/ui/label';
import { authClient } from '@repo/auth/client';
import {
  ResetPasswordFormValues,
  resetPasswordSchema,
} from '@repo/validators/auth';

export default function ResetPasswordPage() {
  const searchParams = useSearchParams();
  const token = searchParams.get('token');

  const [isSubmitted, setIsSubmitted] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [showPassword, setShowPassword] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<ResetPasswordFormValues>({
    resolver: zodResolver(resetPasswordSchema),
    defaultValues: {
      password: '',
      confirmPassword: '',
    },
  });

  const onSubmit = async (values: ResetPasswordFormValues) => {
    if (!token) {
      setErrorMsg('Token de recuperación no encontrado');
      return;
    }

    setErrorMsg(null);
    try {
      const result = await authClient.resetPassword({
        newPassword: values.password,
        token: token,
      });

      if (result.error) {
        setErrorMsg(result.error.message || 'Error al restablecer la contraseña');
      } else {
        setIsSubmitted(true);
      }
    } catch (err) {
      setErrorMsg('Ocurrió un error inesperado');
    }
  };

  // Success state
  if (isSubmitted) {
    return (
      <div className="text-center space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
        <div className="w-16 h-16 bg-green-100 dark:bg-green-950/20 text-green-600 rounded-full flex items-center justify-center mx-auto mb-4">
          <CheckCircle2 size={32} />
        </div>
        <h2 className="text-2xl font-bold text-foreground">
          ¡Contraseña restablecida!
        </h2>
        <p className="text-muted-foreground max-w-xs mx-auto">
          Tu contraseña ha sido actualizada correctamente. Ya puedes iniciar
          sesión con tu nueva contraseña.
        </p>
        <Link href="/auth/login" className="w-full block">
          <Button className="mt-4 w-full h-11 rounded-xl bg-[#a855f7] hover:bg-[#9333ea] text-white font-bold">
            Iniciar sesión
          </Button>
        </Link>
      </div>
    );
  }

  // No token state
  if (!token) {
    return (
      <div className="text-center space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
        <div className="w-16 h-16 bg-red-100 dark:bg-red-950/20 text-red-600 rounded-full flex items-center justify-center mx-auto mb-4">
          <AlertCircle size={32} />
        </div>
        <h2 className="text-2xl font-bold text-foreground">Enlace inválido</h2>
        <p className="text-muted-foreground max-w-xs mx-auto">
          El enlace de recuperación no es válido o ha expirado. Por favor,
          solicita uno nuevo.
        </p>
        <Link href="/auth/forgot-password" className="w-full block">
          <Button variant="outline" className="mt-4 w-full h-11 rounded-xl">
            Solicitar nuevo enlace
          </Button>
        </Link>
      </div>
    );
  }

  return (
    <div className="space-y-6 w-full font-sans">
      <header className="flex flex-col items-center gap-1 text-center mb-6">
        <h1 className="text-3xl font-black tracking-tight text-foreground">Restablecer contraseña</h1>
        <p className="text-muted-foreground text-sm">
          Ingresa tu nueva contraseña para restablecerla.
        </p>
      </header>

      {errorMsg && (
        <div className="p-3 bg-red-500/10 border border-red-500/20 text-red-500 text-xs rounded-xl font-medium">
          {errorMsg}
        </div>
      )}

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="password">Nueva Contraseña</Label>
          <div className="relative group">
            <Input
              id="password"
              type={showPassword ? 'text' : 'password'}
              {...register('password')}
              placeholder="Ingresa una nueva contraseña"
              disabled={isSubmitting}
              className={`bg-zinc-50 dark:bg-zinc-900 border-zinc-200 dark:border-zinc-800 h-12 rounded-xl px-4 focus:ring-brand focus:border-brand transition-all text-sm font-medium ${errors.password ? 'border-red-500 focus:border-red-500' : ''}`}
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-4 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
            >
              {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
            </button>
          </div>
          {errors.password && (
            <p className="text-xs text-red-500 mt-1 font-medium">
              {errors.password.message}
            </p>
          )}
        </div>

        <div className="space-y-2">
          <Label htmlFor="confirmPassword">Confirmar contraseña</Label>
          <Input
            id="confirmPassword"
            type={showPassword ? 'text' : 'password'}
            {...register('confirmPassword')}
            placeholder="Confirma tu nueva contraseña"
            disabled={isSubmitting}
            className={`bg-zinc-50 dark:bg-zinc-900 border-zinc-200 dark:border-zinc-800 h-12 rounded-xl px-4 focus:ring-brand focus:border-brand transition-all text-sm font-medium ${errors.confirmPassword ? 'border-red-500 focus:border-red-500' : ''}`}
          />
          {errors.confirmPassword && (
            <p className="text-xs text-red-500 mt-1 font-medium">
              {errors.confirmPassword.message}
            </p>
          )}
        </div>

        <Button type="submit" disabled={isSubmitting} className="cursor-pointer h-12 rounded-xl mt-5 w-full bg-[#a855f7] hover:bg-[#9333ea] text-white font-bold">
          {isSubmitting ? (
            <Loader2 className="h-5 w-5 animate-spin mx-auto" />
          ) : (
            'Restablecer contraseña'
          )}
        </Button>
      </form>
    </div>
  );
}
