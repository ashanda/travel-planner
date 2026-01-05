<template>
  <SectionCard>
    <div class="flex items-start justify-between">
      <div>
        <div class="text-sm font-semibold">Create a trip</div>
        <div class="mt-1 text-xs text-slate-500">Fast form • Premium output • Saved automatically</div>
      </div>

      <div class="flex items-center gap-2">
        <span class="rounded-full border bg-slate-50 px-3 py-1 text-xs font-semibold text-slate-700">
          Step {{ step }}/3
        </span>
      </div>
    </div>

    <!-- Stepper -->
    <div class="mt-4 grid grid-cols-3 gap-2">
      <div :class="barClass(1)"></div>
      <div :class="barClass(2)"></div>
      <div :class="barClass(3)"></div>
    </div>

    <!-- STEP 1 -->
    <div v-if="step === 1" class="mt-6 space-y-4">
      <div>
        <label class="text-xs font-semibold text-slate-700">Destination</label>
        <input
          v-model="form.destination"
          type="text"
          placeholder="e.g., Kandy / Ella / Galle / Nuwara Eliya"
          class="mt-2 w-full rounded-2xl border px-4 py-3 text-sm outline-none focus:border-slate-400"
        />
        <div class="mt-3 flex flex-wrap gap-2">
          <button v-for="c in suggestedCities" :key="c" type="button" class="rounded-full border bg-white px-3 py-1 text-xs font-semibold text-slate-700 hover:bg-slate-50"
            @click="form.destination = c">
            {{ c }}
          </button>
        </div>
      </div>

      <div class="flex gap-3">
        <button type="button" class="w-full rounded-2xl bg-slate-900 px-4 py-3 text-sm font-semibold text-white hover:bg-slate-800"
          :disabled="!form.destination.trim()"
          @click="step = 2">
          Continue
        </button>
      </div>
    </div>

    <!-- STEP 2 -->
    <div v-else-if="step === 2" class="mt-6 space-y-5">
      <div class="grid gap-4 md:grid-cols-2">
        <div>
          <label class="text-xs font-semibold text-slate-700">Days</label>
          <input v-model.number="form.days" min="1" max="30" type="number"
            class="mt-2 w-full rounded-2xl border px-4 py-3 text-sm outline-none focus:border-slate-400" />
        </div>

        <div>
          <label class="text-xs font-semibold text-slate-700">Start date (optional)</label>
          <input v-model="form.start_date" type="date"
            class="mt-2 w-full rounded-2xl border px-4 py-3 text-sm outline-none focus:border-slate-400" />
        </div>
      </div>

      <div>
        <label class="text-xs font-semibold text-slate-700">Budget</label>
        <div class="mt-2 grid grid-cols-3 gap-2">
          <button type="button" :class="pill(form.budget==='low')" @click="form.budget='low'">Low</button>
          <button type="button" :class="pill(form.budget==='mid')" @click="form.budget='mid'">Mid</button>
          <button type="button" :class="pill(form.budget==='high')" @click="form.budget='high'">High</button>
        </div>
      </div>

      <div>
        <label class="text-xs font-semibold text-slate-700">Pace</label>
        <div class="mt-2 grid grid-cols-3 gap-2">
          <button type="button" :class="pill(form.pace==='chill')" @click="form.pace='chill'">Chill</button>
          <button type="button" :class="pill(form.pace==='balanced')" @click="form.pace='balanced'">Balanced</button>
          <button type="button" :class="pill(form.pace==='fast')" @click="form.pace='fast'">Fast</button>
        </div>
      </div>

      <div class="flex gap-3">
        <button type="button" class="w-full rounded-2xl border bg-white px-4 py-3 text-sm font-semibold text-slate-800 hover:bg-slate-50"
          @click="step = 1">
          Back
        </button>
        <button type="button" class="w-full rounded-2xl bg-slate-900 px-4 py-3 text-sm font-semibold text-white hover:bg-slate-800"
          @click="step = 3">
          Continue
        </button>
      </div>
    </div>

    <!-- STEP 3 -->
    <div v-else class="mt-6 space-y-5">
      <div>
        <label class="text-xs font-semibold text-slate-700">Interests</label>
        <div class="mt-3 flex flex-wrap gap-2">
          <button
            v-for="tag in allInterests"
            :key="tag"
            type="button"
            :class="chipActive(tag)"
            @click="toggleInterest(tag)"
          >
            {{ tag }}
          </button>
        </div>
      </div>

      <div>
        <label class="text-xs font-semibold text-slate-700">Notes (optional)</label>
        <textarea v-model="form.notes" rows="3"
          placeholder="e.g., family friendly, avoid long hikes, prefer local food..."
          class="mt-2 w-full rounded-2xl border px-4 py-3 text-sm outline-none focus:border-slate-400"></textarea>
      </div>

      <div class="grid grid-cols-2 gap-3">
        <button type="button" class="rounded-2xl border bg-white px-4 py-3 text-sm font-semibold text-slate-800 hover:bg-slate-50"
          @click="step = 2">
          Back
        </button>

        <button type="button"
          class="rounded-2xl bg-gradient-to-r from-slate-900 to-slate-700 px-4 py-3 text-sm font-semibold text-white hover:opacity-95"
          @click="submit">
          Generate plan
        </button>
      </div>

      <button
        type="button"
        class="w-full rounded-2xl border bg-white px-4 py-3 text-sm font-semibold text-slate-800 hover:bg-slate-50"
        @click="regenerate"
      >
        Regenerate (force new AI)
      </button>

      <div class="text-xs text-slate-500">
        Tip: Generate once and reuse. Regenerate only if you edit.
      </div>
    </div>
  </SectionCard>
</template>

<script setup lang="ts">
import SectionCard from '@/components/SectionCard.vue'

const emit = defineEmits<{
  (e: 'submit', payload: any): void
  (e: 'regenerate', payload: any): void
}>()

const step = ref(1)

const suggestedCities = ['Kandy', 'Ella', 'Galle', 'Nuwara Eliya', 'Sigiriya', 'Mirissa', 'Jaffna', 'Trincomalee']

const allInterests = [
  'nature', 'food', 'history', 'culture', 'wildlife', 'beach', 'luxury', 'kids', 'shopping', 'photography'
]

const form = reactive({
  destination: '',
  start_date: '',
  days: 3,
  budget: 'mid',
  pace: 'balanced',
  interests: [] as string[],
  notes: ''
})

const pill = (active: boolean) =>
  active
    ? 'rounded-2xl border bg-slate-900 px-4 py-3 text-sm font-semibold text-white'
    : 'rounded-2xl border bg-white px-4 py-3 text-sm font-semibold text-slate-800 hover:bg-slate-50'

const chipActive = (tag: string) => {
  const active = form.interests.includes(tag)
  return active
    ? 'rounded-full border bg-slate-900 px-3 py-1 text-xs font-semibold text-white'
    : 'rounded-full border bg-white px-3 py-1 text-xs font-semibold text-slate-700 hover:bg-slate-50'
}

const toggleInterest = (tag: string) => {
  const i = form.interests.indexOf(tag)
  if (i >= 0) form.interests.splice(i, 1)
  else form.interests.push(tag)
}

const barClass = (n: number) =>
  (step.value >= n ? 'h-2 rounded-full bg-slate-900' : 'h-2 rounded-full bg-slate-200')

const payload = () => ({
  destination: form.destination.trim(),
  start_date: form.start_date || undefined,
  days: Number(form.days || 1),
  budget: form.budget,
  pace: form.pace,
  interests: form.interests,
  notes: form.notes.trim()
})

const submit = () => emit('submit', payload())
const regenerate = () => emit('regenerate', payload())
</script>
