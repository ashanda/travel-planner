import { defineStore } from 'pinia'

export type User = {
  id: string
  email: string
  name: string
  picture?: string
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    loading: false
  }),

  actions: {
    async fetchMe() {
      try {
        const cfg = useRuntimeConfig()
        const data = await $fetch<User>(`${cfg.public.apiBase}/v1/auth/me`, {
          credentials: 'include'
        })
        this.user = data
      } catch {
        this.user = null
      }
    },

    async loginWithGoogle(idToken: string) {
      const cfg = useRuntimeConfig()
      this.loading = true
      try {
        const res = await $fetch<{ user: User }>(`${cfg.public.apiBase}/v1/auth/google`, {
          method: 'POST',
          body: { id_token: idToken },
          credentials: 'include'
        })
        this.user = res.user
      } finally {
        this.loading = false
      }
    },

    async logout() {
      const cfg = useRuntimeConfig()
      await $fetch(`${cfg.public.apiBase}/v1/auth/logout`, {
        method: 'POST',
        credentials: 'include'
      })
      this.user = null
    }
  }
})
