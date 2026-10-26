<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { useVModel } from '@vueuse/core'
import { ChevronDown } from 'lucide-vue-next'
import { cn } from '~/lib/utils'

const props = defineProps<{
  defaultValue?: string | number | null
  modelValue?: string | number | null
  class?: HTMLAttributes['class']
}>()

const emits = defineEmits<{
  'update:modelValue': [value: string | number | null]
}>()

const modelValue = useVModel(props, 'modelValue', emits, {
  passive: true,
  defaultValue: props.defaultValue
})
</script>

<template>
  <div class="relative">
    <select
      v-model="modelValue"
      :class="cn(
        'flex h-10 w-full appearance-none rounded-lg border border-input bg-card px-4 pr-9 py-2 text-sm transition-colors',
        'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 ring-offset-background',
        'disabled:cursor-not-allowed disabled:opacity-50',
        props.class
      )"
    >
      <slot />
    </select>
    <ChevronDown class="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground pointer-events-none" />
  </div>
</template>
