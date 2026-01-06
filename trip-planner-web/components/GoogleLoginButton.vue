<template>
  <button
    type="button"
    @click="signIn"
    :disabled="auth.loading || !ready"
    class="w-full rounded-2xl bg-white text-slate-900 font-semibold py-3 shadow hover:-translate-y-0.5 transition flex items-center justify-center gap-3 disabled:opacity-60 disabled:cursor-not-allowed pl-4 pr-4"
  >
    <img src="/google.svg" class="h-5 w-5" alt="google" />
    {{ ready ? 'Continue with Google' : 'Loading Google…' }}
  </button>

  <p v-if="error" class="mt-3 text-sm text-red-600">
    {{ error }}
  </p>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const cfg = useRuntimeConfig()

const ready = ref(false)
const error = ref('')

declare global {
  interface Window {
    google?: any
  }
}

const initGSI = () => {
  const gsi = window.google?.accounts?.id
  const clientId = cfg.public.googleClientId

  if (!gsi) return false
  if (!clientId) {
    error.value = 'Missing Google Client ID (NUXT_PUBLIC_GOOGLE_CLIENT_ID).'
    console.error('Missing cfg.public.googleClientId')
    return false
  }

  // ✅ initialize once
  gsi.initialize({
    client_id: clientId,
    callback: async (response: any) => {
      error.value = ''
      try {
        const idToken = response?.credential
        if (!idToken) {
          error.value = 'Google did not return a credential token.'
          console.error('No credential', response)
          return
        }

        await auth.loginWithGoogle(idToken)
        await auth.fetchMe()
      } catch (e: any) {
        // Show useful error from API if available
        error.value =
          e?.data?.details ||
          e?.data?.error ||
          e?.message ||
          'Google login failed'
        console.error('Login error', e)
      }
    }
  })

  ready.value = true
  return true
}

onMounted(() => {
  if (!import.meta.client) return

  const wait = () => {
    if (initGSI()) return
    setTimeout(wait, 150)
  }
  wait()
})

const signIn = () => {
  if (!import.meta.client) return
  const gsi = window.google?.accounts?.id
  if (!gsi || !ready.value) return
  gsi.prompt()
}
</script>
