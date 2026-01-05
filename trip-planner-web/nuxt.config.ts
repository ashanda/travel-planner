export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',

  devtools: { enabled: true },

  modules: ['@nuxtjs/tailwindcss', '@pinia/nuxt'],

  runtimeConfig: {
    public: {
      // Docker+Nginx => "/api"
      // Local without Nginx => "http://localhost:8080"
      apiBase: '/api',
    },
  },

  nitro: {
    preset: 'node-server',
  },
})
