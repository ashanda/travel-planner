<template>
  <button
    class="w-full rounded-2xl bg-slate-900 px-5 py-3 text-sm font-semibold text-white shadow-lg shadow-slate-900/20
           transition hover:-translate-y-0.5 hover:bg-slate-950 active:translate-y-0 disabled:opacity-60"
    :disabled="disabled"
    type="button"
    @click="$emit('click')"
  >
    <slot />
  </button>
</template>

<script setup lang="ts">
defineProps<{ disabled?: boolean }>()
defineEmits<{ (e:'click'): void }>()
</script>
