<template>
  <!-- HERO -->
  <section class="relative overflow-hidden">
    <!-- Background -->
    <img src="/hero.png" alt="hero" class="absolute inset-0 h-full w-full object-cover" />
    <div class="absolute inset-0 bg-slate-950/55"></div>

    <div class="relative mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
      <div class="py-10 sm:py-14 lg:py-16">
        <!-- Top label -->
        <div class="inline-flex items-center gap-2 rounded-full border border-white/15 bg-white/10 px-3 py-1 text-xs font-semibold text-white/90 backdrop-blur">
          <span class="h-2 w-2 rounded-full bg-emerald-400"></span>
          Sri Lanka • AI Trip Planner
        </div>

        <div class="mt-6 grid gap-8 lg:grid-cols-2 lg:items-center">
          <!-- Left content -->
          <div>
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
          </div>

          <!-- Right form -->
          <div class="lg:justify-self-end">
            <GlassCard class="w-full max-w-md">
              <div class="text-sm font-semibold text-white">Create your trip</div>
              <div class="mt-1 text-xs text-white/70">Tell us your destination & preferences</div>

              <form class="mt-5 space-y-3" @submit.prevent="submit">
                <div>
                  <label class="text-xs font-semibold text-white/80">Destination</label>
                  <input v-model.trim="form.destination"
                         class="mt-1 w-full rounded-2xl border border-white/15 bg-white/10 px-4 py-3 text-sm text-white placeholder-white/40 outline-none backdrop-blur
                                focus:border-white/30"
                         placeholder="Kandy / Ella / Galle / Colombo" />
                </div>

                <div class="grid gap-3 sm:grid-cols-2">
                  <div>
                    <label class="text-xs font-semibold text-white/80">Days</label>
                    <input v-model.number="form.days"
                           type="number" min="1" max="30"
                           class="mt-1 w-full rounded-2xl border border-white/15 bg-white/10 px-4 py-3 text-sm text-white outline-none backdrop-blur focus:border-white/30" />
                  </div>
                  <div>
                    <label class="text-xs font-semibold text-white/80">Start date</label>
                    <input v-model="form.start_date"
                           type="date"
                           class="mt-1 w-full rounded-2xl border border-white/15 bg-white/10 px-4 py-3 text-sm text-white outline-none backdrop-blur focus:border-white/30" />
                  </div>
                </div>

                <div class="grid gap-3 sm:grid-cols-2">
                  <div>
                    <label class="text-xs font-semibold text-white/80">Budget</label>
                    <select v-model="form.budget"
                            class="mt-1 w-full rounded-2xl border border-white/15 bg-white/10 px-4 py-3 text-sm text-white outline-none backdrop-blur focus:border-white/30">
                      <option class="text-slate-900" value="low">Low</option>
                      <option class="text-slate-900" value="mid">Mid</option>
                      <option class="text-slate-900" value="high">High</option>
                    </select>
                  </div>
                  <div>
                    <label class="text-xs font-semibold text-white/80">Pace</label>
                    <select v-model="form.pace"
                            class="mt-1 w-full rounded-2xl border border-white/15 bg-white/10 px-4 py-3 text-sm text-white outline-none backdrop-blur focus:border-white/30">
                      <option class="text-slate-900" value="chill">Chill</option>
                      <option class="text-slate-900" value="balanced">Balanced</option>
                      <option class="text-slate-900" value="fast">Fast</option>
                    </select>
                  </div>
                </div>

                <div>
                  <label class="text-xs font-semibold text-white/80">Interests</label>
                  <div class="mt-2 flex flex-wrap gap-2">
                    <button v-for="t in interestOptions" :key="t" type="button"
                            @click="toggleInterest(t)"
                            class="rounded-full border px-3 py-1 text-xs font-semibold backdrop-blur transition"
                            :class="form.interests.includes(t)
                              ? 'border-white/30 bg-white text-slate-900'
                              : 'border-white/15 bg-white/10 text-white/85 hover:border-white/25'">
                      {{ t }}
                    </button>
                  </div>
                </div>

                <div>
                  <label class="text-xs font-semibold text-white/80">Notes (optional)</label>
                  <textarea v-model.trim="form.notes" rows="3"
                            class="mt-1 w-full rounded-2xl border border-white/15 bg-white/10 px-4 py-3 text-sm text-white placeholder-white/40 outline-none backdrop-blur focus:border-white/30"
                            placeholder="Kids, elderly, no long hikes, must-see places..."></textarea>
                </div>

                <PrimaryButton :disabled="loading || !form.destination || !form.days" @click="submit">
                  {{ loading ? 'Planning…' : 'Continue' }}
                </PrimaryButton>

                <div v-if="error" class="rounded-2xl border border-red-400/30 bg-red-500/10 p-3 text-xs text-red-100">
                  {{ error }}
                </div>
              </form>
            </GlassCard>
          </div>
        </div>
      </div>
    </div>

    <!-- subtle bottom fade -->
    <div class="pointer-events-none absolute bottom-0 left-0 right-0 h-24 bg-gradient-to-b from-transparent to-white"></div>
  </section>

  <!-- GALLERY -->
  <SectionShell>
    <div class="text-center">
      <div class="inline-flex items-center gap-2 rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold text-slate-700">
        Our Gallery
      </div>
      <h2 class="mt-4 text-2xl font-extrabold tracking-tight text-slate-900">
        Beautiful moments across Sri Lanka
      </h2>
      <p class="mx-auto mt-2 max-w-2xl text-sm text-slate-600">
        Inspiration photos — replace these with your own travel images.
      </p>
    </div>

    <div class="mt-8">
      <GalleryGrid :images="gallery" />
    </div>
  </SectionShell>

  <!-- CONTACT -->
  <section class="bg-slate-50">
    <SectionShell>
      <div class="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
        <div>
          <div class="inline-flex items-center gap-2 rounded-full bg-white px-3 py-1 text-xs font-semibold text-slate-700 border">
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
        <ContactCards :cards="contactCards" />
      </div>

      <div class="mt-10 border-t pt-6 text-center text-xs text-slate-500">
        © {{ new Date().getFullYear() }} Travel Planner — All rights reserved.
      </div>
    </SectionShell>
  </section>
</template>

<script setup lang="ts">
import SectionShell from '@/components/SectionShell.vue'
import FeaturePills from '@/components/FeaturePills.vue'
import GlassCard from '@/components/GlassCard.vue'
import PrimaryButton from '@/components/PrimaryButton.vue'
import GalleryGrid from '@/components/GalleryGrid.vue'
import ContactCards from '@/components/ContactCards.vue'

type TripRequest = {
  destination: string
  start_date?: string
  days: number
  budget: 'low' | 'mid' | 'high'
  pace: 'chill' | 'balanced' | 'fast'
  interests: string[]
  notes?: string
}

const loading = ref(false)
const error = ref<string>('')

const interestOptions = ['Nature', 'Food', 'History', 'Luxury', 'Kids', 'Beaches', 'Wildlife', 'Culture']

const form = reactive<TripRequest>({
  destination: '',
  start_date: '',
  days: 3,
  budget: 'mid',
  pace: 'balanced',
  interests: ['Nature', 'Food'],
  notes: ''
})

const toggleInterest = (t: string) => {
  const idx = form.interests.indexOf(t)
  if (idx >= 0) form.interests.splice(idx, 1)
  else form.interests.push(t)
}

// ✅ Hook to your planner API later
const submit = async () => {
  error.value = ''
  loading.value = true
  try {
    // example call (uncomment when your API is ready)
    // const apiBase = useRuntimeConfig().public.apiBase
    // await $fetch(`${apiBase}/v1/trip/plan`, { method: 'POST', body: form })

    // Demo behavior:
    await new Promise(r => setTimeout(r, 700))
  } catch (e: any) {
    error.value = e?.data?.details || e?.message || 'Request failed'
  } finally {
    loading.value = false
  }
}

const gallery = [
  'https://images.unsplash.com/photo-1544986581-efac024faf62?auto=format&fit=crop&w=1200&q=70',
  'https://images.unsplash.com/photo-1526772662000-3f88f10405ff?auto=format&fit=crop&w=1200&q=70',
  'https://images.unsplash.com/photo-1526778548025-fa2f459cd5c1?auto=format&fit=crop&w=1200&q=70',
  'https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=1200&q=70',
  'https://images.unsplash.com/photo-1489515217757-5fd1be406fef?auto=format&fit=crop&w=1200&q=70',
  'https://images.unsplash.com/photo-1526481280695-3c687fd643ed?auto=format&fit=crop&w=1200&q=70',
]

const contactCards = [
  {
    icon: '✉',
    title: 'Email Support',
    subtitle: 'Fast replies',
    body: 'Send your questions and we’ll respond quickly with trip suggestions and help.',
    link: 'mailto:support@yourdomain.com',
    linkText: 'support@yourdomain.com'
  },
  {
    icon: '💬',
    title: 'WhatsApp',
    subtitle: 'Chat with us',
    body: 'Message us to customize itineraries for families, groups, or special needs.',
    link: 'https://wa.me/94770000000',
    linkText: 'Open WhatsApp'
  },
  {
    icon: '📞',
    title: 'Call',
    subtitle: 'Available hours',
    body: 'Call us for booking help and urgent travel questions.',
    link: 'tel:+94770000000',
    linkText: '+94 77 000 0000'
  }
]
</script>
