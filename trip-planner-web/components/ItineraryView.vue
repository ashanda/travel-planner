<template>
  <SectionCard>
    <div class="flex items-start justify-between">
      <div>
        <div class="text-sm font-semibold">Itinerary</div>
        <div class="mt-1 text-xs text-slate-500">Day-by-day plan</div>
      </div>
    </div>

    <div class="mt-6 space-y-6">
      <div
        v-for="(day, idx) in days"
        :key="idx"
        class="rounded-3xl border p-6"
      >
        <!-- Day header -->
        <div class="flex items-center gap-4">
          <div
            class="grid h-11 w-11 place-items-center rounded-2xl bg-slate-900 text-white text-sm font-bold"
          >
            {{ Number(idx) + 1 }}
          </div>

          <div>
            <div class="text-sm font-semibold">
              Day {{ Number(idx) + 1 }} – {{ day.theme || 'Planned activities' }}
            </div>
            <div class="text-xs text-slate-500">
              {{ day.base_city }} • {{ day.date }}
            </div>
          </div>
        </div>

        <!-- Activities -->
        <div class="mt-5 space-y-3">
          <div
            v-for="(item, iidx) in day.items || []"
            :key="iidx"
            class="rounded-2xl bg-slate-50 p-4"
          >
            <div class="flex items-start justify-between gap-3">
              <div>
                <div class="text-sm font-semibold text-slate-900">
                  {{ item.title }}
                </div>
                <div class="mt-1 text-xs text-slate-600">
                  📍 {{ item.location }}
                </div>
              </div>

              <div class="shrink-0 text-right text-xs text-slate-500">
                <div>{{ item.time_block }}</div>
                <div>
                  {{ item.travel_mode }}
                  <span v-if="item.travel_mins">({{ item.travel_mins }} min)</span>
                </div>
              </div>
            </div>

            <div v-if="item.description" class="mt-2 text-sm text-slate-700">
              {{ item.description }}
            </div>
          </div>
        </div>

        <!-- Meals -->
        <div v-if="day.meals?.length" class="mt-5">
          <div class="text-xs font-semibold text-slate-600 mb-2">Meals</div>
          <div class="grid gap-2 md:grid-cols-3">
            <div
              v-for="(m, midx) in day.meals"
              :key="midx"
              class="rounded-xl border bg-white p-3 text-sm"
            >
              <div class="font-semibold capitalize">{{ m.meal_type }}</div>
              <div class="text-xs text-slate-600">{{ m.suggestion }}</div>
            </div>
          </div>
        </div>

        <!-- Cost -->
        <div v-if="day.cost_range" class="mt-5 text-xs text-slate-600">
          Estimated daily cost:
          <span class="font-semibold text-slate-900">
            {{ day.cost_range.low }} – {{ day.cost_range.high }}
            {{ day.cost_range.currency }}
          </span>
        </div>
      </div>
    </div>
  </SectionCard>
</template>

<script setup lang="ts">
import SectionCard from '@/components/SectionCard.vue'

const props = defineProps<{ itinerary: any }>()

const days = computed<any[]>(() => {
  return props.itinerary?.days || []
})
</script>
