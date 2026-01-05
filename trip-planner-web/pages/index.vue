<template>
  <AppHeader />

  <main class="mx-auto max-w-6xl px-5 py-8">
    <!-- hero -->
    <section class="grid gap-6 md:grid-cols-2 md:items-center">
      <div>
        <div class="inline-flex items-center gap-2 rounded-full border bg-white px-3 py-1 text-xs font-semibold text-slate-700">
          <span class="h-2 w-2 rounded-full bg-emerald-500"></span>
          Modern planner • Sri Lanka • Low cost AI
        </div>

        <h1 class="mt-4 text-3xl font-bold tracking-tight md:text-4xl">
          Build a beautiful itinerary in seconds.
        </h1>

        <p class="mt-3 text-base leading-relaxed text-slate-600">
          Create a plan once, save it, and reload anytime without extra AI cost. Perfect for families, couples, or solo travel.
        </p>

        <div class="mt-5 flex flex-wrap gap-2">
          <Chip>Day-by-day plan</Chip>
          <Chip>Budget estimate</Chip>
          <Chip>Local tips</Chip>
          <Chip>Weather-aware</Chip>
        </div>
      </div>

      <TripForm @submit="onSubmit" @regenerate="onRegenerate" />
    </section>

    <!-- content -->
    <section class="mt-8 grid gap-6 lg:grid-cols-[1.1fr_.9fr]">
      <!-- left -->
      <div class="space-y-6">
        <SkeletonPlan v-if="store.loading" />

        <SectionCard v-else-if="store.error">
          <div class="text-sm font-semibold text-red-600">Request failed</div>
          <div class="mt-2 text-sm text-slate-600">{{ store.error }}</div>
        </SectionCard>

        <template v-else-if="store.plan">
          <WeatherCard
            v-if="store.plan.weather"
            :weather="store.plan.weather"
            :destination="store.plan.request.destination"
          />
          <ItineraryView v-if="store.plan.itinerary" :itinerary="store.plan.itinerary" />
        </template>

        <SectionCard v-else>
          <div class="text-sm text-slate-600">Create a plan to see results here.</div>
        </SectionCard>
      </div>

      <!-- right: sticky summary -->
      <aside class="h-fit space-y-6 lg:sticky lg:top-24">
        <SectionCard>
          <div class="flex items-start justify-between">
            <div>
              <div class="text-sm font-semibold">Plan summary</div>
              <div class="mt-1 text-xs text-slate-500">Saved plan info</div>
            </div>
            <NuxtLink to="/plan" class="text-xs font-semibold underline text-slate-700 hover:text-slate-900">
              Open plans
            </NuxtLink>
          </div>

          <div v-if="store.plan" class="mt-4">
            <PlanSummaryCard :plan="store.plan" />
          </div>
          <div v-else class="mt-4 text-sm text-slate-600">
            No plan loaded.
          </div>
        </SectionCard>

        <SectionCard>
          <div class="text-sm font-semibold">Load plan by ID</div>
          <div class="mt-1 text-xs text-slate-500">Share plan ID with others</div>
          <div class="mt-4">
            <LoadPlanById @load="loadById" />
          </div>
        </SectionCard>
      </aside>
    </section>
  </main>
</template>

<script setup lang="ts">
import AppHeader from '@/components/AppHeader.vue'
import Chip from '@/components/Chip.vue'
import TripForm from '@/components/TripForm.vue'
import ItineraryView from '@/components/ItineraryView.vue'
import WeatherCard from '@/components/WeatherCard.vue'
import SkeletonPlan from '@/components/SkeletonPlan.vue'
import SectionCard from '@/components/SectionCard.vue'
import PlanSummaryCard from '@/components/PlanSummaryCard.vue'
import LoadPlanById from '@/components/LoadPlanById.vue'
import { usePlanStore } from '@/stores/plan'
import type { TripRequest } from '@/types/trip'

const store = usePlanStore()

const onSubmit = async (payload: TripRequest) => { await store.generate(payload) }
const onRegenerate = async (payload: TripRequest) => { await store.regenerate(payload) }
const loadById = async (id: string) => { if (!id) return; await store.loadPlan(id) }
</script>
