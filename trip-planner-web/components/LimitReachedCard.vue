<template>
  <div class="rounded-3xl border border-amber-200 bg-amber-50 p-6 shadow-sm">
    <div class="flex items-start gap-4">
      <div class="flex h-12 w-12 items-center justify-center rounded-2xl bg-amber-100">
        🔒
      </div>

      <div class="flex-1">
        <h3 class="text-lg font-bold text-amber-900">Free Limit Reached</h3>
        <p class="mt-1 text-sm text-amber-800">
          You’ve reached your free trip generation limit.
          Upgrade to continue creating unlimited travel plans.
        </p>

        <div class="mt-4 flex flex-wrap gap-3">
          <button
            class="rounded-2xl bg-slate-900 px-5 py-2 text-sm font-semibold text-white hover:bg-slate-800 transition"
            @click="$emit('upgrade')"
          >
            Upgrade Now
          </button>

          <button
            class="rounded-2xl border border-amber-300 bg-white px-5 py-2 text-sm font-semibold text-amber-900 hover:bg-amber-100 transition"
            @click="$emit('close')"
          >
            Maybe Later
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
