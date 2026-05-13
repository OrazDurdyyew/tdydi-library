<script setup lang="ts">
import { Sun, Moon, Monitor } from 'lucide-vue-next'

const colorMode = useColorMode()

function cycle() {
  if (colorMode.preference === 'light') colorMode.preference = 'dark'
  else if (colorMode.preference === 'dark') colorMode.preference = 'system'
  else colorMode.preference = 'light'
}

const title = computed(() => {
  if (colorMode.preference === 'light') return 'Светлая тема'
  if (colorMode.preference === 'dark') return 'Тёмная тема'
  return 'Системная'
})
</script>

<template>
  <button
    @click="cycle"
    :title="title"
    class="w-7 h-7 inline-flex items-center justify-center rounded-full text-muted-foreground hover:text-primary transition-colors"
  >
    <Sun v-if="colorMode.preference === 'light'" class="w-4 h-4" />
    <Moon v-else-if="colorMode.preference === 'dark'" class="w-4 h-4" />
    <Monitor v-else class="w-4 h-4" />
  </button>
</template>
