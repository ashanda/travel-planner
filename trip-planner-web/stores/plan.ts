import { defineStore } from 'pinia'
import type { TripPlan, TripRequest } from '@/types/trip'

export const usePlanStore = defineStore('plan', {
  state: () => ({
    loading: false,
    error: '' as string,
    plan: null as TripPlan | null,
  }),

  actions: {
    async generate(payload: TripRequest) {
      const { post } = useApi()
      this.loading = true
      this.error = ''
      try {
        this.plan = await post<TripPlan>('/v1/trip/plan', payload)
      } catch (e: any) {
        this.error = e?.data?.error || e?.message || 'Failed'
      } finally {
        this.loading = false
      }
    },

    async regenerate(payload: TripRequest) {
      const { post } = useApi()
      this.loading = true
      this.error = ''
      try {
        this.plan = await post<TripPlan>('/v1/trip/plan/regenerate', payload)
      } catch (e: any) {
        this.error = e?.data?.error || e?.message || 'Failed'
      } finally {
        this.loading = false
      }
    },

    async loadPlan(id: string) {
      const { get } = useApi()
      this.loading = true
      this.error = ''
      try {
        this.plan = await get<TripPlan>(`/v1/trip/plan/${id}`)
      } catch (e: any) {
        this.error = e?.data?.error || e?.message || 'Not found'
      } finally {
        this.loading = false
      }
    },
  },
})
