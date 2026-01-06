<template>
  <button
    @click="signIn"
    class="w-full rounded-2xl bg-white text-slate-900 font-semibold py-3 shadow hover:-translate-y-0.5 transition flex items-center justify-center gap-3"
  >
    <img src="/google.svg" class="h-5 w-5" alt="google" />
    Continue with Google
  </button>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const cfg = useRuntimeConfig()

const signIn = () => {
  // @ts-ignore
  google.accounts.id.initialize({
    client_id: cfg.public.googleClientId,
    callback: async (response: any) => {
      await auth.loginWithGoogle(response.credential)
    }
  })

  // @ts-ignore
  google.accounts.id.prompt()
}
</script>
