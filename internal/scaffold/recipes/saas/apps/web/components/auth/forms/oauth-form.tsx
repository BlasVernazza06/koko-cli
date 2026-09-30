'use client';

import { useState } from 'react';

import { GitHubLight, Google } from '@ridemountainpig/svgl-react';
import { Loader2 } from 'lucide-react';

import { signIn } from '@repo/auth/client';
import { Button } from '@repo/ui/components/ui/button';

interface OAuthButtonsProps {
  disabled?: boolean;
}

export default function OAuthButtons({ disabled }: OAuthButtonsProps) {
  const [isLoading, setIsLoading] = useState<'google' | 'github' | null>(null);

  const handleSignIn = async (provider: 'google' | 'github') => {
    setIsLoading(provider);
    const webUrl = process.env.NEXT_PUBLIC_APP_URL || 'http://localhost:3000';
    await signIn.social({
      provider,
      callbackURL: `${webUrl}/onboarding`,
    });
    setIsLoading(null);
  };

  return (
    <div className="grid grid-cols-2 gap-3 w-full font-sans">
      <Button
        variant="outline"
        type="button"
        disabled={disabled}
        className="bg-white border-[#E2E8F0] hover:bg-[#FAFBFC] h-12 rounded-xl transition-all hover:border-primary/20 flex items-center justify-center gap-2 shadow-sm"
        onClick={() => handleSignIn('google')}
      >
        {isLoading === 'google' ? (
          <Loader2 className="size-5 animate-spin" />
        ) : (
          <Google className="size-5" />
        )}
        <span className="text-xs font-semibold text-[#1A1C1E]">Google</span>
      </Button>
      <Button
        variant="outline"
        type="button"
        disabled={disabled}
        className="bg-white border-[#E2E8F0] hover:bg-[#FAFBFC] h-12 rounded-xl transition-all hover:border-primary/20 flex items-center justify-center gap-2 shadow-sm"
        onClick={() => handleSignIn('github')}
      >
        {isLoading === 'github' ? (
          <Loader2 className="size-5 animate-spin" />
        ) : (
          <GitHubLight className="size-5" />
        )}
        <span className="text-xs font-semibold text-[#1A1C1E]">GitHub</span>
      </Button>
    </div>
  );
}
