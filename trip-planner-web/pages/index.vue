<template>
  <!-- HERO -->
  <section class="relative overflow-hidden">
    <!-- Background image -->
    <img src="/hero.png" alt="hero" class="absolute inset-0 h-full w-full object-cover" />
    <div class="absolute inset-0 bg-slate-950/55"></div>

    <div class="relative mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
      <div class="py-10 sm:py-14 lg:py-16">
        <div class="inline-flex items-center gap-2 rounded-full border border-white/15 bg-white/10 px-3 py-1 text-xs font-semibold text-white/90 backdrop-blur">
          <span class="h-2 w-2 rounded-full bg-emerald-400"></span>
          Sri Lanka • AI Trip Planner
        </div>

        <div class="mt-6 grid gap-8 lg:grid-cols-2 lg:items-start">
          <!-- LEFT: copy -->
          <div class="pt-2">
            <h1 class="text-3xl font-extrabold tracking-tight text-white sm:text-4xl lg:text-5xl">
              Plan Your Sri Lanka Trip
              <span class="block text-white/90">in Minutes</span>
            </h1>

            <p class="mt-4 max-w-xl text-sm leading-6 text-white/80 sm:text-base">
              Day-by-day routes, budget estimates, places, tips, and weather-aware plans.
              Generate once → saved JSON → reuse without extra AI cost.
            </p>

            <div class="mt-6">
              <FeaturePills :items="['Day-by-day plan','Budget estimate','Local tips','Weather-aware']" />
            </div>

            <div class="mt-8 flex items-center gap-3 text-white/80">
              <div class="flex -space-x-2">
                <div v-for="i in 3" :key="i" class="h-9 w-9 rounded-full border border-white/20 bg-white/10 backdrop-blur"></div>
              </div>
              <div class="text-xs sm:text-sm">
                <span class="font-semibold text-white">Fast & practical</span> itineraries for families, couples, and groups.
              </div>
            </div>

            <div class="mt-8 flex items-center gap-3">
              <NuxtLink to="/plan"
                class="inline-flex items-center justify-center rounded-2xl bg-white px-4 py-2 text-sm font-semibold text-slate-900 shadow hover:-translate-y-0.5 transition">
                Open Plan Page
              </NuxtLink>

              <a href="#results"
                class="inline-flex items-center justify-center rounded-2xl border border-white/20 bg-white/10 px-4 py-2 text-sm font-semibold text-white/90 backdrop-blur hover:border-white/30 transition">
                View Results
              </a>
            </div>
          </div>

          <!-- RIGHT: your existing TripForm inside glass card -->
          <div class="lg:justify-self-end w-full">
            <GlassCard class="w-full max-w-md">
              <div class="text-sm font-semibold text-white">Create your trip</div>
              <div class="mt-1 text-xs text-white/70">Tell us destination & preferences</div>

              <div class="mt-5">
                <!-- ✅ Same logic as before -->
                <TripForm @submit="onSubmit" @regenerate="onRegenerate" />
              </div>
            </GlassCard>
          </div>
        </div>
      </div>
    </div>

    <div class="pointer-events-none absolute bottom-0 left-0 right-0 h-24 bg-gradient-to-b from-transparent to-white"></div>
  </section>

  <!-- RESULTS -->
  <section id="results" class="bg-white">
    <SectionShell>
      <div class="grid gap-6 lg:grid-cols-12">
        <!-- LEFT: itinerary -->
        <div class="lg:col-span-7 space-y-6">
          <WeatherCard
            v-if="store.plan"
            :weather="store.plan?.itinerary?.weather"
            :destination="store.plan?.request?.destination"
          />

          <SectionCard title="Itinerary" subtitle="Day-by-day plan">
            <ItineraryView v-if="store.plan?.itinerary" :itinerary="store.plan.itinerary" />
            <div v-else class="text-sm text-slate-600">
              Create a plan to see itinerary here.
            </div>
          </SectionCard>
        </div>

        <!-- RIGHT: summary + load by id -->
        <div class="lg:col-span-5 space-y-6">
          <SectionCard>
            <div class="flex items-start justify-between gap-4">
              <div>
                <div class="text-sm font-semibold">Plan summary</div>
                <div class="mt-1 text-xs text-slate-500">Saved plan info</div>
              </div>
              <NuxtLink to="/plan" class="text-xs font-semibold underline text-slate-700 hover:text-slate-900">
                Open plans
              </NuxtLink>
            </div>

            <div class="mt-4">
              <div v-if="store.loading" class="rounded-2xl border bg-white p-4 text-sm">
                Planning…
              </div>

              <div v-else-if="store.error" class="rounded-2xl border bg-white p-4 text-sm text-red-600">
                {{ store.error }}
              </div>

              <PlanSummaryCard v-else-if="store.plan" :plan="store.plan" />

              <div v-else class="rounded-2xl border bg-white p-4 text-sm text-slate-600">
                Create a plan to see summary here.
              </div>
            </div>
          </SectionCard>

          <LoadPlanById @load="loadById" />
        </div>
      </div>
    </SectionShell>
  </section>

  <!-- GALLERY -->
  <section class="bg-slate-50">
    <SectionShell>
      <div class="text-center">
        <div class="inline-flex items-center gap-2 rounded-full bg-white px-3 py-1 text-xs font-semibold text-slate-700 border">
          Our Gallery
        </div>
        <h2 class="mt-4 text-2xl font-extrabold tracking-tight text-slate-900">
          Beautiful moments across Sri Lanka
        </h2>
        <p class="mx-auto mt-2 max-w-2xl text-sm text-slate-600">
          Replace these images with your own. This is a premium landing layout.
        </p>
      </div>

      <div class="mt-8">
        <GalleryGrid :images="gallery" />
      </div>
    </SectionShell>
  </section>

  <!-- CONTACT -->
  <section class="bg-white">
    <SectionShell>
      <div class="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
        <div>
          <div class="inline-flex items-center gap-2 rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold text-slate-700">
            Get in touch
          </div>
          <h2 class="mt-4 text-3xl font-extrabold tracking-tight text-slate-900">
            We’d love to hear from you
          </h2>
          <p class="mt-2 max-w-xl text-sm text-slate-600">
            Support, bookings, and questions — contact us anytime.
          </p>
        </div>

        <div class="flex items-center gap-3 text-slate-500">
          <a href="#" class="rounded-2xl border bg-white px-3 py-2 text-xs font-semibold hover:text-slate-900">f</a>
          <a href="#" class="rounded-2xl border bg-white px-3 py-2 text-xs font-semibold hover:text-slate-900">ig</a>
          <a href="#" class="rounded-2xl border bg-white px-3 py-2 text-xs font-semibold hover:text-slate-900">in</a>
        </div>
      </div>

      <div class="mt-8">
        <ContactCards />
      </div>

      <div class="mt-10 border-t pt-6 text-center text-xs text-slate-500">
        © {{ new Date().getFullYear() }} Travel Planner — All rights reserved.
      </div>
    </SectionShell>
  </section>
</template>

<script setup lang="ts">
import { usePlanStore } from '@/stores/plan'
import type { TripRequest } from '@/types/trip'

import TripForm from '@/components/TripForm.vue'
import PlanSummaryCard from '@/components/PlanSummaryCard.vue'
import ItineraryView from '@/components/ItineraryView.vue'
import LoadPlanById from '@/components/LoadPlanById.vue'
import WeatherCard from '@/components/WeatherCard.vue'

import FeaturePills from '@/components/FeaturePills.vue'
import GlassCard from '@/components/GlassCard.vue'
import GalleryGrid from '@/components/GalleryGrid.vue'
import ContactCards from '@/components/ContactCards.vue'
import SectionShell from '@/components/SectionShell.vue'
import SectionCard from '@/components/SectionCard.vue'

const store = usePlanStore()

const onSubmit = async (payload: TripRequest) => {
  await store.generate(payload)
}
const onRegenerate = async (payload: TripRequest) => {
  await store.regenerate(payload)
}
const loadById = async (id: string) => {
  if (!id) return
  await store.loadPlan(id)
}

const gallery = [
  'https://images.unsplash.com/photo-1544986581-efac024faf62?auto=format&fit=crop&w=1200&q=70',
  'https://images.unsplash.com/photo-1526772662000-3f88f10405ff?auto=format&fit=crop&w=1200&q=70',
  'https://images.unsplash.com/photo-1526778548025-fa2f459cd5c1?auto=format&fit=crop&w=1200&q=70',
  'https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=1200&q=70',
  'https://images.unsplash.com/photo-1489515217757-5fd1be406fef?auto=format&fit=crop&w=1200&q=70',
  'https://images.unsplash.com/photo-1526481280695-3c687fd643ed?auto=format&fit=crop&w=1200&q=70',
]
</script>
