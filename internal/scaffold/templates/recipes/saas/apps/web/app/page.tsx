import Link from 'next/link';

export default function Page() {
  return (
    <div className="flex flex-col min-h-screen bg-zinc-950 text-zinc-100 font-sans">
      <header className="px-4 lg:px-6 h-14 flex items-center border-b border-zinc-850">
        <span className="font-bold text-lg">[[.ProjectName]]</span>
        <nav className="ml-auto flex gap-4 sm:gap-6">
          <Link className="text-sm font-medium hover:underline underline-offset-4" href="/auth/login">
            Login
          </Link>
          <Link className="text-sm font-medium hover:underline underline-offset-4" href="/auth/register">
            Register
          </Link>
        </nav>
      </header>
      <main className="flex-1 flex flex-col items-center justify-center text-center p-4">
        <h1 className="text-4xl md:text-6xl font-bold tracking-tighter max-w-3xl mb-4 bg-gradient-to-r from-white to-zinc-400 bg-clip-text text-transparent">
          Welcome to [[.ProjectName]]
        </h1>
        <p className="text-zinc-400 max-w-[600px] md:text-xl mb-8">
          The ultimate boilerplate for your Next.js SaaS project, configured with Drizzle ORM, better-auth, and tailwind.
        </p>
        <div className="flex gap-4 justify-center">
          <Link
            className="inline-flex h-10 items-center justify-center rounded-md bg-zinc-50 px-8 text-sm font-medium text-zinc-900 shadow transition-colors hover:bg-zinc-200"
            href="/auth/login"
          >
            Get Started
          </Link>
        </div>
      </main>
      <footer className="flex flex-col gap-2 sm:flex-row py-6 w-full shrink-0 items-center px-4 md:px-6 border-t border-zinc-900">
        <p className="text-xs text-zinc-500">© 2026 [[.ProjectName]]. All rights reserved.</p>
      </footer>
    </div>
  );
}
