<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { computed } from 'vue'
import { Check } from 'lucide-vue-next'
import { cn } from '~/lib/utils'

const props = defineProps<{
  modelValue?: boolean | (string | number)[]
  value?: string | number
  id?: string
  disabled?: boolean
  class?: HTMLAttributes['class']
}>()

const emits = defineEmits<{
  'update:modelValue': [value: boolean | (string | number)[]]
}>()

const isChecked = computed(() => {
  if (Array.isArray(props.modelValue)) {
    return props.value !== undefined && props.modelValue.includes(props.value)
  }
  return !!props.modelValue
})

function toggle() {
  if (props.disabled) return
  if (Array.isArray(props.modelValue) && props.value !== undefined) {
    const arr = [...props.modelValue]
    const idx = arr.indexOf(props.value)
    if (idx >= 0) arr.splice(idx, 1)
    else arr.push(props.value)
    emits('update:modelValue', arr)
  } else {
    emits('update:modelValue', !props.modelValue)
  }
}
</script>

<template>
  <button
    :id="id"
    type="button"
    role="checkbox"
    :aria-checked="isChecked"
    :disabled="disabled"
    @click="toggle"
    :class="cn(
      'peer inline-flex h-4 w-4 shrink-0 items-center justify-center rounded-sm border border-primary ring-offset-background transition-colors',
      'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
      'disabled:cursor-not-allowed disabled:opacity-50',
      isChecked ? 'bg-primary text-primary-foreground' : 'bg-card',
      props.class
    )"
  >
    <Check v-if="isChecked" class="h-3 w-3" :stroke-width="3" />
  </button>
</template>
