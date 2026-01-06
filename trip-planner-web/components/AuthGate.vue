<template>
  <div>
    <!-- If logged in -->
    <slot v-if="auth.user" />

    <!-- If not logged in -->
    <div
      v-else
      class="rounded-3xl border border-white/20 bg-white/10 p-6 text-center backdrop-blur"
    >
      <div class="text-sm font-semibold text-white mb-2">
        Login required
      </div>
      <div class="text-xs text-white/70 mb-5">
        Sign in with Google to generate your trip plan.
      </div>

      <ClientOnly>
        <GoogleLoginButton />
    </ClientOnly>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import GoogleLoginButton from '@/components/GoogleLoginButton.vue'

const auth = useAuthStore()
</script>
