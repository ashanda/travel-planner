<template>
  <div class="space-y-5">
    <!-- Header -->
    <div class="flex items-center gap-3">
      <div class="h-4 w-4 rounded-full bg-emerald-400 animate-pulse"></div>
      <div class="text-sm font-semibold text-slate-700">
        Generating your itinerary<span class="dots"></span>
     </div>
    </div>

    <!-- Day skeletons -->
    <div
      v-for="i in 3"
      :key="i"
      class="rounded-3xl border bg-white p-5 shadow-sm"
    >
      <div class="flex gap-4">
        <!-- Day number -->
        <div
          class="grid h-10 w-10 place-items-center rounded-2xl bg-slate-200 animate-pulse"
        ></div>

        <div class="flex-1 space-y-3">
          <!-- Title -->
          <div class="h-4 w-40 rounded bg-slate-200 animate-pulse"></div>

          <!-- Activities -->
          <div class="space-y-2">
            <div class="h-3 w-full rounded bg-slate-100 animate-pulse"></div>
            <div class="h-3 w-5/6 rounded bg-slate-100 animate-pulse"></div>
            <div class="h-3 w-4/6 rounded bg-slate-100 animate-pulse"></div>
          </div>

          <!-- Cost -->
          <div class="h-3 w-32 rounded bg-slate-200 animate-pulse"></div>
        </div>
      </div>
    </div>

    <!-- Footer hint -->
    <div class="text-xs text-slate-500 italic">
      AI is analyzing routes, travel time, weather & budget…
    </div>
  </div>
</template>
<style scoped>
.dots::after {
  content: '';
  animation: dots 1.5s steps(4, end) infinite;
}
@keyframes dots {
  0% { content: ''; }
  25% { content: '.'; }
  50% { content: '..'; }
  75% { content: '...'; }
}
</style>

