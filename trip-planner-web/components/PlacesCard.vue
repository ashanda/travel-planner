<template>
  <div class="rounded-2xl border bg-white p-5 shadow-sm">
    <div class="text-lg font-semibold">Top places</div>
    <div class="text-sm text-slate-600 mt-1">Based on Google Places</div>

    <div class="mt-4 space-y-2">
      <div
        v-for="(p, idx) in topPlaces"
        :key="idx"
        class="rounded-xl border p-3"
      >
        <div class="font-semibold text-sm">{{ p.name }}</div>
        <div class="text-xs text-slate-500">{{ p.address }}</div>
        <div class="text-xs text-slate-500">
          ⭐ {{ p.rating ?? '—' }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ places?: any }>()

const topPlaces = computed(() => props.places?.top_places || [])
</script>
