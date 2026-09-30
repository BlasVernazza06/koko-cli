'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import { Loader, Mail } from 'lucide-react';
import { motion } from 'motion/react';
import { useState } from 'react';
import { useForm } from 'react-hook-form';

import { authClient } from '@repo/auth/src/client';
import { Button } from '@repo/ui/components/ui/button';
import { Input } from '@repo/ui/components/ui/input';
import { Label } from '@repo/ui/components/ui/label';
import {
  type RequestPasswordResetFormValues,
  requestPasswordResetSchema,
} from '@repo/validators/auth';

export default function SignInForm() {
  const [isLoadingForm, setIsLoadingForm] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RequestPasswordResetFormValues>({
    resolver: zodResolver(requestPasswordResetSchema),
  });

  const onSubmit = async (formData: RequestPasswordResetFormValues) => {
    setIsLoadingForm(true);
    try {
      await authClient.requestPasswordReset({
        email: formData.email,
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
      transition={{ delay: 0.2 }}
      className="space-y-5"
      onSubmit={handleSubmit(onSubmit)}
    >
      <div className="space-y-2">
        <Label
          htmlFor="email"
          className="text-xs font-bold text-[#4A4C4E] ml-1 uppercase tracking-wider"
        >
          Email
        </Label>
        <div className="relative group">
          <div className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 flex items-center justify-center">
            <Mail className="w-4 h-4 text-muted-foreground transition-colors group-focus-within:text-primary" />
          </div>
          <Input
            id="email"
            type="email"
            {...register('email')}
            placeholder="ej. blas@memo.ai"
            className={`bg-[#FAFBFC] border-[#E2E8F0] h-12 rounded-xl pl-11 focus:ring-primary/10 focus:border-primary transition-all text-sm ${errors.email ? 'border-red-500 focus:border-red-500 focus:ring-red-500/10' : ''}`}
          />
        </div>
        {errors.email && (
          <p className="text-[10px] text-red-500 font-bold ml-1 uppercase">
            {errors.email.message}
          </p>
        )}
      </div>

      <Button
        type="submit"
        className="w-full bg-primary hover:bg-primary/90 text-white py-6 rounded-xl font-bold text-md shadow-lg shadow-primary/25 transition-all active:scale-[0.98]"
      >
        {isLoadingForm ? <Loader className="size-4" /> : 'Enviar'}
      </Button>
    </motion.form>
  );
}
