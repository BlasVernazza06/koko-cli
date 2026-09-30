<script>
  let count = 0;
</script>

<div class="flex min-h-screen flex-col items-center justify-center p-8 text-center">
  <div class="max-w-2xl space-y-6">
    <div class="inline-flex items-center gap-2 rounded-full border border-orange-500/30 bg-orange-500/10 px-4 py-1 text-sm text-orange-400">
      <span class="h-2 w-2 rounded-full bg-orange-400 animate-pulse"></span>
      SvelteKit Web Application
    </div>
    <h1 class="text-4xl sm:text-6xl font-bold tracking-tight bg-gradient-to-r from-white via-orange-200 to-orange-400 bg-clip-text text-transparent">
      [[.ProjectName]]
    </h1>
    <p class="text-zinc-400 text-lg">
      Cybernetically enhanced web apps powered by Svelte & Tailwind CSS.
    </p>
    <div class="pt-4">
      <button
        on:click={() => count++}
        class="rounded-lg bg-orange-600 px-6 py-2.5 font-medium text-white shadow-lg transition hover:bg-orange-500 active:scale-95"
      >
        Clicks: {count}
      </button>
    </div>
  </div>
</div>
