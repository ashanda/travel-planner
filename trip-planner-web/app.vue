<template>
  <div class="min-h-screen bg-slate-50 text-slate-900">
    <!-- background blobs -->
    <div class="pointer-events-none fixed inset-0 -z-10 overflow-hidden">
      <div class="absolute -top-32 left-1/2 h-96 w-[52rem] -translate-x-1/2 rounded-full bg-gradient-to-r from-indigo-200 via-sky-200 to-emerald-200 blur-3xl opacity-70"></div>
      <div class="absolute bottom-0 right-0 h-72 w-72 rounded-full bg-gradient-to-tr from-amber-200 to-rose-200 blur-3xl opacity-60"></div>
      <div class="absolute bottom-10 left-10 h-64 w-64 rounded-full bg-gradient-to-tr from-purple-200 to-cyan-200 blur-3xl opacity-40"></div>
    </div>

    <NuxtPage />
  </div>
</template>
