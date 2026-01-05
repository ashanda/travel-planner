<template>
  <SectionCard>
    <div class="flex items-start justify-between">
      <div>
        <div class="text-sm font-semibold">Weather</div>
        <div class="mt-1 text-xs text-slate-500">{{ destination }}</div>
      </div>

      <span v-if="enabled" class="rounded-full border bg-emerald-50 px-3 py-1 text-xs font-semibold text-emerald-700">
        Live
      </span>
      <span v-else class="rounded-full border bg-slate-50 px-3 py-1 text-xs font-semibold text-slate-600">
        Disabled
      </span>
    </div>

    <div v-if="enabled" class="mt-5 grid grid-cols-2 gap-3 md:grid-cols-4">
      <div class="rounded-2xl border bg-white p-4">
        <div class="text-xs text-slate-500">Temp</div>
        <div class="mt-1 text-lg font-bold">{{ weather.temp_c ?? '—' }}°C</div>
      </div>
      <div class="rounded-2xl border bg-white p-4">
        <div class="text-xs text-slate-500">Feels like</div>
        <div class="mt-1 text-lg font-bold">{{ weather.feels_like_c ?? '—' }}°C</div>
      </div>
      <div class="rounded-2xl border bg-white p-4">
        <div class="text-xs text-slate-500">Humidity</div>
        <div class="mt-1 text-lg font-bold">{{ weather.humidity ?? '—' }}%</div>
      </div>
      <div class="rounded-2xl border bg-white p-4">
        <div class="text-xs text-slate-500">Condition</div>
        <div class="mt-1 text-lg font-bold">{{ weather.condition || '—' }}</div>
      </div>
    </div>

    <div v-else class="mt-4 text-sm text-slate-600">
      Weather info unavailable (missing key or disabled).
    </div>
  </SectionCard>
</template>

<script setup lang="ts">
import SectionCard from '@/components/SectionCard.vue'
const props = defineProps<{ weather?: any; destination?: string }>()
const enabled = computed(() => !!props.weather?.enabled)
</script>
