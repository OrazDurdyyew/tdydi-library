<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { computed, watch, onUnmounted } from 'vue'
import { cn } from '~/lib/utils'

const props = defineProps<{
  open: boolean
  class?: HTMLAttributes['class']
}>()

const emits = defineEmits<{
  'update:open': [value: boolean]
  close: []
}>()

const isOpen = computed(() => props.open)

function close() {
  emits('update:open', false)
  emits('close')
}

function onBackdropClick(e: MouseEvent) {
  if (e.target === e.currentTarget) close()
}

function onKey(e: KeyboardEvent) {
  if (e.key === 'Escape' && isOpen.value) close()
}

watch(isOpen, (v) => {
  if (import.meta.client) {
    document.body.style.overflow = v ? 'hidden' : ''
    if (v) document.addEventListener('keydown', onKey)
    else document.removeEventListener('keydown', onKey)
  }
})

onUnmounted(() => {
  if (import.meta.client) {
    document.body.style.overflow = ''
    document.removeEventListener('keydown', onKey)
  }
})
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-150 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-100 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="isOpen"
        class="fixed inset-0 z-50 bg-black/50 backdrop-blur-sm flex items-center justify-center p-4"
        @click="onBackdropClick"
      >
        <Transition
          enter-active-class="transition duration-150 ease-out"
          enter-from-class="opacity-0 scale-95"
          enter-to-class="opacity-100 scale-100"
          leave-active-class="transition duration-100 ease-in"
          leave-from-class="opacity-100 scale-100"
          leave-to-class="opacity-0 scale-95"
        >
          <div
            v-if="isOpen"
            :class="cn(
              'relative w-full max-w-lg rounded-xl bg-card text-card-foreground shadow-2xl border border-border p-6',
              props.class
            )"
          >
            <slot :close="close" />
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>
