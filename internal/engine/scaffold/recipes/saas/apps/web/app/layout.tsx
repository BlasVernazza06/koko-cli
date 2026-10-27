import { Poppins } from 'next/font/google';

import type { Metadata } from 'next';

import './globals.css';

const poppins = Poppins({
  subsets: ['latin'],
  weight: ['400', '500', '700'],
  variable: '--font-poppins',
});

export const metadata: Metadata = {
  title: {
    default: '[[.ProjectName]] - SaaS App',
    template: '%s | [[.ProjectName]]',
  },
  description: 'Write your application description here...',
  keywords: ['saas', 'app', 'nextjs', 'drizzle'],
  authors: [{ name: '[[.ProjectName]] Team' }],
  creator: '[[.ProjectName]] Team',
  openGraph: {
    type: 'website',
    locale: 'es_ES',
    url: 'https://example.com',
    siteName: '[[.ProjectName]]',
    title: '[[.ProjectName]] - SaaS App',
    description: 'Write your application description here...',
    images: [
      {
        url: '/og-image.png',
        width: 1200,
        height: 630,
        alt: '[[.ProjectName]] Preview',
      },
    ],
  },
  twitter: {
    card: 'summary_large_image',
    title: '[[.ProjectName]] - SaaS App',
    description: 'Write your application description here...',
    images: ['/og-image.png'],
    creator: '@[[.ProjectName]]',
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="es" suppressHydrationWarning>
      <body className={`${poppins.variable} font-sans bg-background text-foreground antialiased`}>
        {children}
      </body>
    </html>
  );
}
