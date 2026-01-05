<template>
  <AppHeader />

  <main class="mx-auto max-w-6xl px-5 py-8">
    <section class="grid gap-6 lg:grid-cols-[1.1fr_.9fr]">
      <!-- Left: Plan viewer -->
      <div class="space-y-6">
        <SectionCard>
          <div class="flex items-start justify-between gap-4">
            <div>
              <div class="text-sm font-semibold">Plans</div>
              <div class="mt-1 text-xs text-slate-500">Search and open a saved plan</div>
            </div>

            <button
              type="button"
              class="rounded-2xl border bg-white px-4 py-2 text-xs font-semibold text-slate-800 hover:bg-slate-50"
              @click="refresh"
            >
              Refresh
            </button>
          </div>

          <div class="mt-4 grid gap-3 md:grid-cols-[1fr_auto]">
            <input
              v-model="planId"
              class="w-full rounded-2xl border px-4 py-3 text-sm outline-none focus:border-slate-400"
              placeholder="Paste Plan ID (UUID) ..."
            />
            <button
              type="button"
              class="rounded-2xl bg-slate-900 px-5 py-3 text-sm font-semibold text-white hover:bg-slate-800"
              @click="openById"
            >
              Open
            </button>
          </div>

          <div v-if="error" class="mt-3 text-sm text-red-600">
            {{ error }}
          </div>
        </SectionCard>

        <div v-if="loading">
          <SkeletonPlan />
        </div>

        <template v-else-if="activePlan">
          <PlanSummaryCard :plan="activePlan" />
          <WeatherCard
            v-if="activePlan.weather"
            :weather="activePlan.weather"
            :destination="activePlan.request.destination"
          />
          <ItineraryView v-if="activePlan.itinerary" :itinerary="activePlan.itinerary" />
        </template>

        <SectionCard v-else>
          <div class="text-sm text-slate-600">Select a plan from the right panel.</div>
        </SectionCard>
      </div>

      <!-- Right: Recent list -->
      <aside class="h-fit space-y-6 lg:sticky lg:top-24">
        <SectionCard>
          <div class="flex items-start justify-between">
            <div>
              <div class="text-sm font-semibold">Recent plans</div>
              <div class="mt-1 text-xs text-slate-500">Click to open</div>
            </div>
            <span class="rounded-full border bg-slate-50 px-3 py-1 text-xs font-semibold text-slate-700">
              {{ plans.length }}
            </span>
          </div>

          <div class="mt-4 space-y-3">
            <button
              v-for="p in plans"
              :key="p.id"
              type="button"
              class="w-full rounded-2xl border bg-white p-4 text-left hover:bg-slate-50"
              @click="selectPlan(p)"
            >
              <div class="flex items-start justify-between gap-3">
                <div class="min-w-0">
                  <div class="truncate text-sm font-semibold text-slate-900">
                    {{ p.request?.destination || 'Trip' }}
                  </div>
                  <div class="mt-1 text-xs text-slate-500">
                    {{ p.request?.days }} days • {{ p.request?.budget || 'mid' }} • {{ p.request?.pace || 'balanced' }}
                  </div>
                </div>

                <div class="flex shrink-0 flex-col items-end gap-2">
                  <span class="text-[11px] text-slate-500">
                    {{ shortDate(p.updated_at) }}
                  </span>
                  <div class="flex gap-2">
                    <button
                      type="button"
                      class="rounded-xl border bg-white px-2 py-1 text-[11px] font-semibold text-slate-800 hover:bg-slate-50"
                      @click.stop="copy(p.id)"
                    >
                      Copy
                    </button>
                    <a
                      class="rounded-xl bg-slate-900 px-2 py-1 text-[11px] font-semibold text-white hover:bg-slate-800"
                      :href="`/plan?id=${p.id}`"
                      target="_blank"
                      rel="noreferrer"
                      @click.stop
                    >
                      New tab
                    </a>
                  </div>
                </div>
              </div>

              <div class="mt-2 text-[11px] text-slate-500 truncate">
                {{ p.id }}
              </div>
            </button>

            <div v-if="!plans.length" class="text-sm text-slate-600">
              No saved plans yet.
            </div>
          </div>

          <div v-if="copied" class="mt-3 text-xs font-semibold text-emerald-700">
            ✅ Copied!
          </div>
        </SectionCard>
      </aside>
    </section>
  </main>
</template>

<script setup lang="ts">
import AppHeader from '@/components/AppHeader.vue'
import SectionCard from '@/components/SectionCard.vue'
import SkeletonPlan from '@/components/SkeletonPlan.vue'
import PlanSummaryCard from '@/components/PlanSummaryCard.vue'
import WeatherCard from '@/components/WeatherCard.vue'
import ItineraryView from '@/components/ItineraryView.vue'
import { usePlanStore } from '@/stores/plan'

const store = usePlanStore()

const loading = ref(false)
const error = ref('')
const plans = ref<any[]>([])
const activePlan = ref<any | null>(null)
const planId = ref('')

const copied = ref(false)
const copy = async (id: string) => {
  try {
    await navigator.clipboard.writeText(id)
    copied.value = true
    setTimeout(() => (copied.value = false), 1200)
  } catch {
    alert('Copy failed')
  }
}

const shortDate = (ts: number) => {
  if (!ts) return '—'
  const d = new Date(ts * 1000)
  return d.toLocaleDateString()
}

const refresh = async () => {
  error.value = ''
  loading.value = true
  try {
    // use store method or direct fetch — easiest:
    const base = useRuntimeConfig().public.apiBase
    plans.value = await $fetch(`${base}/v1/trip/plans`)
    // newest first
    plans.value.sort((a: any, b: any) => (b.updated_at || 0) - (a.updated_at || 0))
  } catch (e: any) {
    error.value = e?.data?.error || e?.message || 'Failed to load plans'
  } finally {
    loading.value = false
  }
}

const selectPlan = (p: any) => {
  activePlan.value = p
  planId.value = p.id
}

const openById = async () => {
  if (!planId.value) return
  error.value = ''
  loading.value = true
  try {
    await store.loadPlan(planId.value)
    activePlan.value = store.plan
  } catch (e: any) {
    error.value = e?.data?.error || e?.message || 'Plan not found'
  } finally {
    loading.value = false
  }
}

// read ?id=... param
const route = useRoute()
onMounted(async () => {
  await refresh()
  const id = route.query.id as string | undefined
  if (id) {
    planId.value = id
    await openById()
  } else if (plans.value.length) {
    selectPlan(plans.value[0])
  }
})
</script>
