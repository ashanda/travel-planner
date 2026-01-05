<template>
  <header class="sticky top-0 z-50 border-b bg-white/75 backdrop-blur">
    <div class="mx-auto flex max-w-6xl items-center justify-between px-5 py-4">
      <div class="flex items-center gap-3">
        <div class="grid h-10 w-10 place-items-center rounded-2xl bg-slate-900 text-white shadow-sm">
          <span class="text-sm font-bold">SL</span>
        </div>
        <div class="leading-tight">
          <div class="text-sm font-semibold">Sri Lanka AI Trip Planner</div>
          <div class="text-xs text-slate-500">Generate once • Save JSON • Reuse without extra AI cost</div>
        </div>
      </div>

      <nav class="flex items-center gap-4">
        <NuxtLink to="/" class="text-sm font-semibold text-slate-700 hover:text-slate-900">Planner</NuxtLink>
        <NuxtLink to="/plan" class="text-sm font-semibold text-slate-700 hover:text-slate-900">Plans</NuxtLink>
      </nav>
    </div>
  </header>
</template>
