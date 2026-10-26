<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { onClickOutside } from '@vueuse/core'
import { ChevronDown, Check } from 'lucide-vue-next'
import { cn } from '~/lib/utils'

interface Option {
  value: string
  label: string
}

const props = defineProps<{
  modelValue?: string | null
  options: Option[]
  placeholder?: string
  disabled?: boolean
  class?: HTMLAttributes['class']
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
  change: [value: string]
}>()

const root = ref<HTMLElement | null>(null)
const open = ref(false)
const highlighted = ref(-1)

const selected = computed(() => props.options.find((o) => o.value === props.modelValue) ?? null)
const buttonLabel = computed(() => selected.value?.label ?? props.placeholder ?? '')

function toggle() {
  if (props.disabled) return
  open.value = !open.value
  if (open.value) {
    highlighted.value = props.options.findIndex((o) => o.value === props.modelValue)
  }
}

function close() {
  open.value = false
}

function choose(opt: Option) {
  emit('update:modelValue', opt.value)
  emit('change', opt.value)
  close()
}

function onKeydown(e: KeyboardEvent) {
  if (props.disabled) return
  if (!open.value && (e.key === 'ArrowDown' || e.key === 'Enter' || e.key === ' ')) {
    e.preventDefault()
    toggle()
    return
  }
  if (!open.value) return
  if (e.key === 'Escape') {
    e.preventDefault()
    close()
  } else if (e.key === 'ArrowDown') {
    e.preventDefault()
    highlighted.value = Math.min(highlighted.value + 1, props.options.length - 1)
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    highlighted.value = Math.max(highlighted.value - 1, 0)
  } else if (e.key === 'Enter') {
    e.preventDefault()
    const opt = props.options[highlighted.value]
    if (opt) choose(opt)
  }
}

onClickOutside(root, close)
</script>

<template>
  <div ref="root" class="relative">
    <button
      type="button"
      :disabled="disabled"
      :aria-expanded="open"
      aria-haspopup="listbox"
      :class="cn(
        'flex h-10 w-full items-center justify-between gap-2 rounded-lg border border-input bg-card px-4 py-2 text-sm transition-colors',
        'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 ring-offset-background',
        'disabled:cursor-not-allowed disabled:opacity-50 hover:border-ring/60',
        open && 'border-ring ring-2 ring-ring/15',
        props.class
      )"
      @click="toggle"
      @keydown="onKeydown"
    >
      <span :class="selected ? 'text-foreground truncate' : 'text-muted-foreground truncate'">
        {{ buttonLabel }}
      </span>
      <ChevronDown
        class="w-4 h-4 text-muted-foreground shrink-0 transition-transform"
        :class="open && 'rotate-180'"
      />
    </button>

    <Transition
      enter-active-class="transition duration-150 ease-out"
      enter-from-class="opacity-0 -translate-y-1"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-100 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-1"
    >
      <ul
        v-if="open"
        role="listbox"
        class="absolute z-50 mt-1.5 w-full max-h-64 overflow-y-auto rounded-lg border border-border bg-card p-1 shadow-lg shadow-black/5"
      >
        <li v-if="placeholder">
          <button
            type="button"
            role="option"
            :aria-selected="!modelValue"
            class="flex w-full items-center justify-between gap-2 rounded-md px-3 py-2 text-left text-sm text-muted-foreground transition-colors hover:bg-muted"
            @click="choose({ value: '', label: placeholder })"
          >
            {{ placeholder }}
            <Check v-if="!modelValue" class="w-4 h-4 text-primary shrink-0" />
          </button>
        </li>
        <li v-for="(opt, i) in options" :key="opt.value">
          <button
            type="button"
            role="option"
            :aria-selected="opt.value === modelValue"
            :class="cn(
              'flex w-full items-center justify-between gap-2 rounded-md px-3 py-2 text-left text-sm transition-colors',
              i === highlighted ? 'bg-muted' : 'hover:bg-muted',
              opt.value === modelValue ? 'font-medium text-primary' : 'text-foreground'
            )"
            @click="choose(opt)"
            @mouseenter="highlighted = i"
          >
            <span class="truncate">{{ opt.label }}</span>
            <Check v-if="opt.value === modelValue" class="w-4 h-4 text-primary shrink-0" />
          </button>
        </li>
      </ul>
    </Transition>
  </div>
</template>
