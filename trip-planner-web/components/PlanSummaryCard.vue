<template>
  <div class="rounded-3xl border bg-white p-6 shadow-sm">
    <div class="flex items-start justify-between gap-4">
      <div class="min-w-0">
        <div class="text-xs font-semibold text-slate-500">Summary</div>
        <div class="mt-1 text-lg font-bold text-slate-900">
          {{ plan.request?.destination || 'Trip plan' }}
        </div>

        <div class="mt-2 flex flex-wrap gap-2">
          <span class="rounded-full border bg-slate-50 px-3 py-1 text-xs font-semibold text-slate-700">
            {{ plan.request?.days }} days
          </span>
          <span class="rounded-full border bg-slate-50 px-3 py-1 text-xs font-semibold text-slate-700">
            {{ plan.request?.budget || 'mid' }}
          </span>
          <span class="rounded-full border bg-slate-50 px-3 py-1 text-xs font-semibold text-slate-700">
            {{ plan.request?.pace || 'balanced' }}
          </span>
          <span
            v-if="plan.request?.start_date"
            class="rounded-full border bg-slate-50 px-3 py-1 text-xs font-semibold text-slate-700"
          >
            Start: {{ plan.request.start_date }}
          </span>
        </div>

        <div v-if="plan.request?.interests?.length" class="mt-3 text-sm text-slate-600">
          <span class="font-semibold text-slate-900">Interests:</span>
          {{ plan.request.interests.join(', ') }}
        </div>
      </div>

      <!-- Actions -->
      <div class="flex flex-col gap-2">
        <button
          type="button"
          class="rounded-2xl border bg-white px-4 py-2 text-xs font-semibold text-slate-800 hover:bg-slate-50"
          @click="copyId"
        >
          Copy Plan ID
        </button>

        <a
          class="rounded-2xl bg-slate-900 px-4 py-2 text-center text-xs font-semibold text-white hover:bg-slate-800"
          :href="`/plan?id=${plan.id}`"
          target="_blank"
          rel="noreferrer"
        >
          Open in new tab
        </a>
      </div>
    </div>

    <div class="mt-5 grid gap-3 md:grid-cols-3">
      <div class="rounded-2xl border bg-slate-50 p-4">
        <div class="text-xs text-slate-500">Plan ID</div>
        <div class="mt-1 truncate text-sm font-semibold text-slate-900">{{ plan.id }}</div>
      </div>

      <div class="rounded-2xl border bg-slate-50 p-4">
        <div class="text-xs text-slate-500">Created</div>
        <div class="mt-1 text-sm font-semibold text-slate-900">
          {{ formatTs(plan.created_at) }}
        </div>
      </div>

      <div class="rounded-2xl border bg-slate-50 p-4">
        <div class="text-xs text-slate-500">Updated</div>
        <div class="mt-1 text-sm font-semibold text-slate-900">
          {{ formatTs(plan.updated_at) }}
        </div>
      </div>
    </div>

    <div v-if="copied" class="mt-4 text-xs font-semibold text-emerald-700">
      ✅ Copied!
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ plan: any }>()
const copied = ref(false)

const copyId = async () => {
  try {
    await navigator.clipboard.writeText(props.plan?.id || '')
    copied.value = true
    setTimeout(() => (copied.value = false), 1200)
  } catch {
    // fallback
    copied.value = false
    alert('Copy failed. Please copy manually.')
  }
}

const formatTs = (ts: number) => {
  if (!ts) return '—'
  const d = new Date(ts * 1000)
  return d.toLocaleString()
}
</script>
