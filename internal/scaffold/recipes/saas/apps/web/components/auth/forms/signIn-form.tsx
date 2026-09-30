'use client';

import Link from 'next/link';

import { useState } from 'react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { Eye, EyeOff, Loader, Lock, Mail } from 'lucide-react';
import { motion } from 'motion/react';

import { authClient } from '@repo/auth/client';
import { Button } from '@repo/ui/components/ui/button';
import { Input } from '@repo/ui/components/ui/input';
import { Label } from '@repo/ui/components/ui/label';
import { type LoginFormValues, loginSchema } from '@repo/validators/auth';

import OAuthButtons from '@/components/auth/forms/oauth-form';

export default function SignInForm() {
  const [isLoadingForm, setIsLoadingForm] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
  });

  const onSubmit = async (formData: LoginFormValues) => {
    setIsLoadingForm(true);
    try {
      await authClient.signIn.email({
        email: formData.email,
        password: formData.password,
        callbackURL: '/onboarding',
      });
    } catch (error) {
      console.error(error);
    } finally {
      setIsLoadingForm(false);
    }
  };

  return (
    <div className="flex flex-col gap-8 w-full">
      <div className="flex flex-col gap-2">
        <h1 className="text-3xl font-black tracking-tight text-foreground">
          Bienvenido de nuevo
        </h1>
        <p className="text-muted-foreground font-medium text-balance leading-relaxed">
          Ingresa tus credenciales para acceder a tu panel de control y gestionar tu tienda.
        </p>
      </div>

      <motion.form
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.2 }}
        className="space-y-5"
        onSubmit={handleSubmit(onSubmit)}
      >
        <div className="space-y-2">
          <Label
            htmlFor="email"
            className="text-xs font-bold text-muted-foreground ml-1 uppercase tracking-widest"
          >
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

        <div className="space-y-2">
          <div className="flex items-center justify-between ml-1">
            <Label
              htmlFor="password"
              className="text-xs font-bold text-muted-foreground uppercase tracking-widest"
            >
              Contraseña
            </Label>
            <Link
              href="/auth/forgot-password"
              className="text-xs font-bold text-brand hover:underline underline-offset-4"
            >
              ¿Olvidaste tu contraseña?
            </Link>
          </div>
          <div className="relative group">
            <Input
              id="password"
              type={showPassword ? 'text' : 'password'}
              {...register('password')}
              placeholder="••••••••"
              className={`bg-zinc-50 dark:bg-zinc-900 border-zinc-200 dark:border-zinc-800 h-12 rounded-xl px-4 pr-12 focus:ring-brand focus:border-brand transition-all text-sm font-medium ${errors.password ? 'border-red-500 focus:border-red-500' : ''}`}
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-4 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors z-10"
            >
              {showPassword ? (
                <EyeOff className="w-5 h-5" />
              ) : (
                <Eye className="w-5 h-5" />
              )}
            </button>
          </div>
          {errors.password && (
            <p className="text-[11px] text-red-500 font-medium ml-1">
              {errors.password.message}
            </p>
          )}
        </div>

        <Button
          type="submit"
          disabled={isLoadingForm}
          className="w-full bg-[#a855f7] hover:bg-[#9333ea] text-white font-black h-12 rounded-xl text-base shadow-lg shadow-purple-500/20 transition-all active:scale-[0.98] cursor-pointer"
        >
          {isLoadingForm ? (
            <Loader className="size-5 animate-spin mx-auto" />
          ) : (
            'Entrar al Sistema'
          )}
        </Button>

        <div className="relative py-4">
          <div className="absolute inset-0 flex items-center">
            <span className="w-full border-t border-zinc-200 dark:border-zinc-800" />
          </div>
          <div className="relative flex justify-center text-[10px] uppercase">
            <span className="bg-white dark:bg-zinc-950 px-2 text-muted-foreground font-bold tracking-widest">
              O continuar con
            </span>
          </div>
        </div>

        <OAuthButtons />
      </motion.form>
    </div>
  );
}
