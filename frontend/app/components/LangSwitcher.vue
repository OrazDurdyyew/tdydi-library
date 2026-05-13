<script setup lang="ts">
const switchLocalePath = useSwitchLocalePath()
const { locale } = useI18n()

const locales = ['tm', 'ru', 'en'] as const
</script>

<template>
  <div class="inline-flex items-center gap-1 font-serif text-xs">
    <template v-for="(code, idx) in locales" :key="code">
      <span v-if="idx > 0" class="text-border/80 select-none">/</span>
      <NuxtLink
        :to="switchLocalePath(code)"
        :class="locale === code
          ? 'text-primary font-bold'
          : 'text-muted-foreground hover:text-primary'"
        class="uppercase tracking-wider px-1 transition-colors"
      >
        {{ code }}
      </NuxtLink>
    </template>
  </div>
</template>
