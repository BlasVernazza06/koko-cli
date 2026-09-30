import Stripe from 'stripe';

if (!process.env.STRIPE_API_KEY) {
  throw new Error('STRIPE_API_KEY environment variable is not defined.');
}

export const stripe = new Stripe(process.env.STRIPE_API_KEY, {
  apiVersion: '2024-11-20.accredited', // or a standard/latest stable API version
  typescript: true,
});
