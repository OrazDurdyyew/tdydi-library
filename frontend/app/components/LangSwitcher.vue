<script setup lang="ts">
const switchLocalePath = useSwitchLocalePath()
const { locale } = useI18n()

const locales = ['tm', 'ru', 'en'] as const
</script>

<template>
  <div class="inline-flex items-center text-xs font-medium">
    <template v-for="(code, idx) in locales" :key="code">
      <span v-if="idx > 0" class="text-border select-none">/</span>
      <NuxtLink
        :to="switchLocalePath(code)"
        :class="locale === code
          ? 'text-primary font-semibold'
          : 'text-muted-foreground hover:text-foreground'"
        class="uppercase tracking-wide px-1.5 py-1 transition-colors"
      >
        {{ code }}
      </NuxtLink>
    </template>
  </div>
</template>
