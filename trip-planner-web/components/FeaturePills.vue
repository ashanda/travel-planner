<template>
  <div class="flex flex-wrap gap-2">
    <span v-for="(t,i) in items" :key="i"
      class="rounded-full border border-white/15 bg-white/10 px-3 py-1 text-xs font-semibold text-white/90 backdrop-blur">
      {{ t }}
    </span>
  </div>
</template>

<script setup lang="ts">
defineProps<{ items: string[] }>()
</script>
