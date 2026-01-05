<template>
  <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
    <div v-for="(img,i) in images" :key="i"
      class="group relative overflow-hidden rounded-3xl border bg-white shadow-sm">
      <img :src="img" class="h-44 w-full object-cover transition duration-500 group-hover:scale-105" />
      <div class="absolute inset-0 bg-gradient-to-t from-black/25 via-transparent to-transparent"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ images: string[] }>()
</script>
