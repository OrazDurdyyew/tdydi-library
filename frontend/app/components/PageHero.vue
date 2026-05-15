<script setup lang="ts">
defineProps<{
  eyebrow?: string
  title: string
  description?: string
}>()
</script>

<template>
  <header class="border-b border-border bg-card">
    <div class="page-container py-10 md:py-12">
      <slot name="breadcrumbs" />
      <span v-if="eyebrow" class="eyebrow">{{ eyebrow }}</span>
      <h1 class="text-3xl md:text-4xl font-bold tracking-tight mt-1.5">
        {{ title }}
      </h1>
      <p
        v-if="description"
        class="text-base md:text-lg text-muted-foreground leading-relaxed max-w-2xl mt-3"
      >
        {{ description }}
      </p>
      <slot />
    </div>
  </header>
</template>
