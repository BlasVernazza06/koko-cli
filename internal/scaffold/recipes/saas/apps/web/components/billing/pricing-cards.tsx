'use client';

import React, { useState } from 'react';
import { Check, Sparkles, Zap, Shield } from 'lucide-react';

interface PricingTier {
  id: string;
  name: string;
  priceMonthly: number;
  priceYearly: number;
  description: string;
  features: string[];
  popular?: boolean;
  cta: string;
}

const tiers: PricingTier[] = [
  {
    id: 'starter',
    name: 'Starter',
    priceMonthly: 19,
    priceYearly: 15,
    description: 'Perfect for indie hackers and small side projects.',
    features: [
      'Up to 5 team members',
      '10,000 API requests / month',
      'Community support',
      'Basic analytics dashboard',
    ],
    cta: 'Start Free Trial',
  },
  {
    id: 'pro',
    name: 'Pro',
    priceMonthly: 49,
    priceYearly: 39,
    description: 'For growing businesses requiring scalability and advanced analytics.',
    popular: true,
    features: [
      'Unlimited team members',
      '500,000 API requests / month',
      'Priority 24/7 support',
      'Advanced analytics & exports',
      'Custom domains & webhooks',
      'Better-Auth multi-tenant SSO',
    ],
    cta: 'Get Started with Pro',
  },
  {
    id: 'enterprise',
    name: 'Enterprise',
    priceMonthly: 199,
    priceYearly: 159,
    description: 'Dedicated infrastructure, custom SLAs and enterprise compliance.',
    features: [
      'Dedicated cluster & database',
      'Unlimited requests',
      '99.99% Uptime SLA',
      'Dedicated Slack channel',
      'Custom invoice billing',
    ],
    cta: 'Contact Sales',
  },
];

export function PricingCards() {
  const [billingCycle, setBillingCycle] = useState<'monthly' | 'yearly'>('monthly');
  const [loadingTier, setLoadingTier] = useState<string | null>(null);

  const handleCheckout = async (tierId: string) => {
    setLoadingTier(tierId);
    try {
      // Example call to /api/checkout/session with Stripe
      const res = await fetch('/api/checkout', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ tierId, billingCycle }),
      });
      const data = await res.json();
      if (data.url) {
        window.location.href = data.url;
      }
    } catch (err) {
      console.error('Checkout error:', err);
    } finally {
      setLoadingTier(null);
    }
  };

  return (
    <div className="py-12 px-4 max-w-7xl mx-auto">
      <div className="text-center space-y-4 mb-10">
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-violet-500/10 text-violet-400 text-sm font-medium">
          <Sparkles className="w-4 h-4" /> Flexible Pricing for High Growth
        </div>
        <h2 className="text-3xl sm:text-4xl font-extrabold text-white tracking-tight">
          Simple, Transparent Plans
        </h2>
        <p className="text-zinc-400 max-w-xl mx-auto text-base">
          Choose the plan that fits your business stage. Upgrade or downgrade anytime with instant Stripe billing.
        </p>

        {/* Toggle Billing Cycle */}
        <div className="flex items-center justify-center gap-3 pt-4">
          <span className={`text-sm ${billingCycle === 'monthly' ? 'text-white font-medium' : 'text-zinc-400'}`}>
            Monthly
          </span>
          <button
            type="button"
            onClick={() => setBillingCycle(billingCycle === 'monthly' ? 'yearly' : 'monthly')}
            className="relative w-14 h-7 rounded-full bg-zinc-800 p-1 transition-colors hover:bg-zinc-700"
          >
            <div
              className={`w-5 h-5 rounded-full bg-violet-500 transition-transform ${
                billingCycle === 'yearly' ? 'translate-x-7' : 'translate-x-0'
              }`}
            />
          </button>
          <span className={`text-sm ${billingCycle === 'yearly' ? 'text-white font-medium' : 'text-zinc-400'}`}>
            Yearly <span className="text-xs text-emerald-400 font-semibold">(Save 20%)</span>
          </span>
        </div>
      </div>

      {/* Cards Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-8 items-stretch">
        {tiers.map((tier) => {
          const price = billingCycle === 'monthly' ? tier.priceMonthly : tier.priceYearly;
          return (
            <div
              key={tier.id}
              className={`relative rounded-2xl p-8 flex flex-col justify-between transition-all duration-300 ${
                tier.popular
                  ? 'bg-zinc-900 border-2 border-violet-500 shadow-xl shadow-violet-500/10'
                  : 'bg-zinc-900/60 border border-zinc-800 hover:border-zinc-700'
              }`}
            >
              {tier.popular && (
                <div className="absolute -top-3.5 left-1/2 -translate-x-1/2 bg-violet-600 text-white text-xs font-semibold px-3 py-1 rounded-full uppercase tracking-wider">
                  Most Popular
                </div>
              )}

              <div>
                <h3 className="text-xl font-bold text-white mb-2">{tier.name}</h3>
                <p className="text-sm text-zinc-400 min-h-[40px]">{tier.description}</p>

                <div className="my-6 flex items-baseline gap-1">
                  <span className="text-4xl font-extrabold text-white">${price}</span>
                  <span className="text-zinc-400 text-sm">/ month</span>
                </div>

                <ul className="space-y-3 mb-8">
                  {tier.features.map((feat, i) => (
                    <li key={i} className="flex items-center gap-3 text-sm text-zinc-300">
                      <Check className="w-4 h-4 text-emerald-400 shrink-0" />
                      <span>{feat}</span>
                    </li>
                  ))}
                </ul>
              </div>

              <button
                onClick={() => handleCheckout(tier.id)}
                disabled={loadingTier === tier.id}
                className={`w-full py-3 px-4 rounded-xl font-semibold text-sm transition-all duration-200 flex items-center justify-center gap-2 ${
                  tier.popular
                    ? 'bg-violet-600 hover:bg-violet-500 text-white shadow-lg shadow-violet-600/30'
                    : 'bg-zinc-800 hover:bg-zinc-700 text-zinc-100'
                }`}
              >
                {loadingTier === tier.id ? (
                  <span>Processing...</span>
                ) : (
                  <>
                    <Zap className="w-4 h-4" />
                    <span>{tier.cta}</span>
                  </>
                )}
              </button>
            </div>
          );
        })}
      </div>
    </div>
  );
}
