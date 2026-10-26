<script setup lang="ts">
import { ChevronRight } from 'lucide-vue-next'

defineProps<{
  items: { label: string; to?: string }[]
}>()
</script>

<template>
  <nav class="flex flex-wrap items-center gap-1.5 text-sm text-muted-foreground mb-4" aria-label="breadcrumb">
    <template v-for="(item, i) in items" :key="i">
      <NuxtLinkLocale
        v-if="item.to"
        :to="item.to"
        class="hover:text-secondary transition-colors"
      >
        {{ item.label }}
      </NuxtLinkLocale>
      <span v-else class="text-foreground font-medium">{{ item.label }}</span>
      <ChevronRight
        v-if="i < items.length - 1"
        class="w-3.5 h-3.5 opacity-50 shrink-0"
        aria-hidden="true"
      />
    </template>
  </nav>
</template>
