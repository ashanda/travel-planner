<template>
  <div class="grid gap-6 md:grid-cols-3">
    <div v-for="c in cards" :key="c.title" class="rounded-3xl border bg-white p-6 shadow-sm">
      <div class="flex items-center gap-3">
        <div class="grid h-10 w-10 place-items-center rounded-2xl bg-slate-900 text-white">
          <span class="text-sm font-bold">{{ c.icon }}</span>
        </div>
        <div class="min-w-0">
          <div class="text-sm font-semibold text-slate-900">{{ c.title }}</div>
          <div class="text-xs text-slate-500">{{ c.subtitle }}</div>
        </div>
      </div>

      <div class="mt-4 text-sm text-slate-700">
        {{ c.body }}
      </div>

      <a v-if="c.link" :href="c.link" class="mt-4 inline-flex text-sm font-semibold text-slate-900 underline">
        {{ c.linkText || 'Contact' }}
      </a>
    </div>
  </div>
</template>

<script setup lang="ts">
type Card = { title: string; subtitle: string; body: string; icon: string; link?: string; linkText?: string }

const props = withDefaults(defineProps<{ cards?: Card[] }>(), {
  cards: () => []
})
</script>
