<template>
  <button
    type="button"
    @click="signIn"
    :disabled="auth.loading || !ready"
    class="w-full rounded-2xl bg-white text-slate-900 font-semibold py-3 shadow hover:-translate-y-0.5 transition flex items-center justify-center gap-3 disabled:opacity-60 disabled:cursor-not-allowed"
  >
    <img src="/google.svg" class="h-5 w-5" alt="google" />
    {{ ready ? 'Continue with Google' : 'Loading Google…' }}
  </button>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const cfg = useRuntimeConfig()
const ready = ref(false)

declare global {
  interface Window {
    google?: any
  }
}

onMounted(() => {
  // ✅ only in browser
  if (!import.meta.client) return

  const check = () => {
    if (window.google?.accounts?.id) {
      ready.value = true
      return
    }
    setTimeout(check, 150)
  }
  check()
})

const signIn = () => {
  if (!import.meta.client) return
  if (!window.google?.accounts?.id) return

  window.google.accounts.id.initialize({
    client_id: cfg.public.googleClientId,
    callback: async (response: any) => {
      await auth.loginWithGoogle(response.credential)
    }
  })

  window.google.accounts.id.prompt()
}
</script>
