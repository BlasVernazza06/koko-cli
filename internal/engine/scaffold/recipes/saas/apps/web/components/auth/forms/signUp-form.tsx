'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';

import { zodResolver } from '@hookform/resolvers/zod';
import { Eye, EyeOff, Lock, Mail, User } from 'lucide-react';
import { motion } from 'motion/react';

import { authClient } from '@repo/auth/client';
import { Button } from '@repo/ui/components/ui/button';
import { Input } from '@repo/ui/components/ui/input';
import { Label } from '@repo/ui/components/ui/label';
import { type RegisterFormValues, registerSchema } from '@repo/validators/auth';

import OAuthButtons from '@/components/auth/forms/oauth-form';

export default function SignUpForm() {
  const [isLoadingForm, setIsLoadingForm] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
  });

  const onSubmit = async (formData: RegisterFormValues) => {
    setIsLoadingForm(true);

    try {
      await authClient.signUp.email({
        email: formData.email,
        password: formData.password,
        name: formData.name,
        callbackURL: '/onboarding',
      });
    } catch (error) {
      console.error(error);
    } finally {
      setIsLoadingForm(false);
    }
  };
  return (
    <motion.form
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.1 }}
      className="space-y-4"
      onSubmit={handleSubmit(onSubmit)}
    >
      <div className="space-y-2">
        <Label
          htmlFor="name"
          className="text-xs font-semibold text-muted-foreground ml-1 uppercase tracking-wider"
        >
          Nombre Completo
        </Label>
        <div className="relative group">
          <div className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 flex items-center justify-center pointer-events-none z-10">
            <User className="w-4 h-4 text-muted-foreground transition-colors group-focus-within:text-brand" />
          </div>
          <Input
            id="name"
            type="text"
            {...register('name')}
            placeholder="Tu nombre"
            className={`bg-zinc-50/50 dark:bg-zinc-900/50 border-zinc-200/80 dark:border-zinc-800/80 h-12 rounded-2xl pl-11 focus:ring-brand/10 focus:border-brand transition-all text-sm font-medium ${errors.name ? 'border-red-500 focus:border-red-500 focus:ring-red-500/10' : ''}`}
          />
        </div>
        {errors.name && (
          <p className="text-[11px] text-red-500 font-medium ml-1">
            {errors.name.message}
          </p>
        )}
      </div>

      <div className="space-y-2">
        <Label
          htmlFor="email"
          className="text-xs font-semibold text-muted-foreground ml-1 uppercase tracking-wider"
        >
          Email
        </Label>
        <div className="relative group">
          <div className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 flex items-center justify-center pointer-events-none z-10">
            <Mail className="w-4 h-4 text-muted-foreground transition-colors group-focus-within:text-brand" />
          </div>
          <Input
            id="email"
            type="email"
            {...register('email')}
            placeholder="tu@email.com"
            className={`bg-zinc-50/50 dark:bg-zinc-900/50 border-zinc-200/80 dark:border-zinc-800/80 h-12 rounded-2xl pl-11 focus:ring-brand/10 focus:border-brand transition-all text-sm font-medium ${errors.email ? 'border-red-500 focus:border-red-500 focus:ring-red-500/10' : ''}`}
          />
        </div>
        {errors.email && (
          <p className="text-[11px] text-red-500 font-medium ml-1">
            {errors.email.message}
          </p>
        )}
      </div>

      <div className="space-y-2">
        <Label
          htmlFor="password"
          className="text-xs font-semibold text-muted-foreground ml-1 uppercase tracking-wider"
        >
          Contraseña
        </Label>
        <div className="relative group">
          <div className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 flex items-center justify-center pointer-events-none z-10">
            <Lock className="w-4 h-4 text-muted-foreground transition-colors group-focus-within:text-brand" />
          </div>
          <Input
            id="password"
            type={showPassword ? 'text' : 'password'}
            {...register('password')}
            placeholder="Mínimo 8 caracteres"
            className={`bg-zinc-50/50 dark:bg-zinc-900/50 border-zinc-200/80 dark:border-zinc-800/80 h-12 rounded-2xl pl-11 pr-10 focus:ring-brand/10 focus:border-brand transition-all text-sm font-medium ${errors.password ? 'border-red-500 focus:border-red-500 focus:ring-red-500/10' : ''}`}
          />
          <button
            type="button"
            onClick={() => setShowPassword(!showPassword)}
            className="absolute right-4 top-1/2 -translate-y-1/2 w-5 h-5 flex items-center justify-center text-muted-foreground hover:text-foreground transition-colors z-10"
          >
            {showPassword ? (
              <Eye className="w-4 h-4" />
            ) : (
              <EyeOff className="w-4 h-4" />
            )}
          </button>
        </div>
        {errors.password && (
          <p className="text-[11px] text-red-500 font-medium ml-1">
            {errors.password.message}
          </p>
        )}
      </div>

      <div className="space-y-2">
        <Label
          htmlFor="confirmPassword"
          className="text-xs font-semibold text-muted-foreground ml-1 uppercase tracking-wider"
        >
          Confirmar Contraseña
        </Label>
        <div className="relative group">
          <div className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 flex items-center justify-center pointer-events-none z-10">
            <Lock className="w-4 h-4 text-muted-foreground transition-colors group-focus-within:text-brand" />
          </div>
          <Input
            id="confirmPassword"
            type={showPassword ? 'text' : 'password'}
            {...register('confirmPassword')}
            placeholder="Repite tu contraseña"
            className={`bg-zinc-50/50 dark:bg-zinc-900/50 border-zinc-200/80 dark:border-zinc-800/80 h-12 rounded-2xl pl-11 pr-10 focus:ring-brand/10 focus:border-brand transition-all text-sm font-medium ${errors.confirmPassword ? 'border-red-500 focus:border-red-500 focus:ring-red-500/10' : ''}`}
          />
        </div>
        {errors.confirmPassword && (
          <p className="text-[11px] text-red-500 font-medium ml-1">
            {errors.confirmPassword.message}
          </p>
        )}
      </div>

      <Button
        type="submit"
        disabled={isLoadingForm}
        className="w-full bg-brand hover:bg-brand/90 text-white font-bold h-12 rounded-2xl text-base shadow-xl shadow-brand/20 transition-all hover:scale-[1.01] active:scale-[0.98] cursor-pointer mt-4"
      >
        {isLoadingForm ? 'Creando Cuenta...' : 'Crear Cuenta'}
      </Button>

      <div className="flex items-center gap-4 py-2">
        <div className="flex-1 border-t border-zinc-200/80 dark:border-zinc-800/80" />
        <span className="text-[10px] font-bold uppercase tracking-widest text-muted-foreground whitespace-nowrap">
          O regístrate con
        </span>
        <div className="flex-1 border-t border-zinc-200/80 dark:border-zinc-800/80" />
      </div>

      <OAuthButtons />
    </motion.form>
  );
}
