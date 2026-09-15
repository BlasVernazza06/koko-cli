import Stripe from 'stripe';

export const stripe = process.env.STRIPE_API_KEY
  ? new Stripe(process.env.STRIPE_API_KEY, {
      apiVersion: '2024-11-20.accredited',
      typescript: true,
    })
  : null;
