import { useState } from "react";
import { Sparkles } from "lucide-react";

export default function App() {
  const [count, setCount] = useState(0);

  return (
    <div className="flex min-h-screen flex-col items-center justify-center p-8 text-center">
      <div className="max-w-2xl space-y-6">
        <div className="inline-flex items-center gap-2 rounded-full border border-cyan-500/30 bg-cyan-500/10 px-4 py-1 text-sm text-cyan-400">
          <Sparkles className="h-4 w-4" />
          React + Vite SPA
        </div>
        <h1 className="text-4xl sm:text-6xl font-bold tracking-tight bg-gradient-to-r from-white via-cyan-200 to-cyan-400 bg-clip-text text-transparent">
          [[.ProjectName]]
        </h1>
        <p className="text-zinc-400 text-lg">
          Fast Single Page Application with React, Vite & Tailwind CSS.
        </p>
        <div className="pt-4">
          <button
            onClick={() => setCount((c) => c + 1)}
            className="rounded-lg bg-cyan-600 px-6 py-2.5 font-medium text-white shadow-lg transition hover:bg-cyan-500 active:scale-95"
          >
            Count is {count}
          </button>
        </div>
      </div>
    </div>
  );
}
