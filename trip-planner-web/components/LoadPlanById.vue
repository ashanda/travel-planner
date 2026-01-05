<template>
  <div class="rounded-2xl border bg-white p-5 shadow-sm space-y-3">
    <div class="flex items-center justify-between">
      <div class="text-lg font-semibold">Open saved plan</div>
      <span class="text-xs text-slate-500">Paste plan ID</span>
    </div>

    <input
      v-model="id"
      class="w-full rounded-xl border p-3 font-mono text-sm"
      placeholder="Plan ID"
    />

    <button
      class="w-full rounded-xl border py-3 font-semibold"
      type="button"
      @click="$emit('load', id)"
    >
      Load plan
    </button>
  </div>
</template>

<script setup lang="ts">
defineEmits<{ (e: 'load', id: string): void }>()
const id = ref('')
</script>
