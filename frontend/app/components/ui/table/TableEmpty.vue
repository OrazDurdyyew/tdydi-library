<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '~/lib/utils'

const props = defineProps<{
  colspan: number
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <tr>
    <td
      :colspan="colspan"
      :class="cn('px-6 py-16 text-center text-sm text-muted-foreground', props.class)"
    >
      <slot />
    </td>
  </tr>
</template>
