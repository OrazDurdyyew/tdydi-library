<script setup lang="ts">
import { Check, ChevronDown } from 'lucide-vue-next'

interface Option {
  value: string
  label: string
}

const props = withDefaults(
  defineProps<{
    modelValue: string
    options: Option[]
    placeholder?: string
  }>(),
  { placeholder: '' }
)

const emit = defineEmits<{
  'update:modelValue': [value: string]
  change: [value: string]
}>()

const open = ref(false)
const activeIndex = ref(-1)
const rootEl = ref<HTMLElement | null>(null)

const selected = computed(() => props.options.find((o) => o.value === props.modelValue))
const displayLabel = computed(() => selected.value?.label ?? props.placeholder)

function toggle() {
  open.value = !open.value
  if (open.value) {
    activeIndex.value = props.options.findIndex((o) => o.value === props.modelValue)
  }
}

function select(opt: Option) {
  emit('update:modelValue', opt.value)
  emit('change', opt.value)
  open.value = false
}

function onKeydown(e: KeyboardEvent) {
  if (!open.value) {
    if (e.key === 'Enter' || e.key === ' ' || e.key === 'ArrowDown') {
      e.preventDefault()
      toggle()
    }
    return
  }
  switch (e.key) {
    case 'ArrowDown':
      e.preventDefault()
      activeIndex.value = (activeIndex.value + 1) % props.options.length
      break
    case 'ArrowUp':
      e.preventDefault()
      activeIndex.value =
        activeIndex.value <= 0 ? props.options.length - 1 : activeIndex.value - 1
      break
    case 'Enter':
    case ' ': {
      e.preventDefault()
      const opt = props.options[activeIndex.value]
      if (opt) select(opt)
      break
    }
    case 'Escape':
      open.value = false
      break
    case 'Home':
      e.preventDefault()
      activeIndex.value = 0
      break
    case 'End':
      e.preventDefault()
      activeIndex.value = props.options.length - 1
      break
  }
}

onClickOutside(rootEl, () => {
  open.value = false
})
</script>

<template>
  <div ref="rootEl" class="relative">
    <button
      type="button"
      class="flex h-10 w-full items-center justify-between gap-2 rounded-lg border border-input bg-card px-3.5 py-2 text-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 ring-offset-background"
      :aria-expanded="open"
      @click="toggle"
      @keydown="onKeydown"
    >
      <span :class="['truncate text-left', selected ? '' : 'text-muted-foreground']">
        {{ displayLabel }}
      </span>
      <ChevronDown
        class="w-4 h-4 text-muted-foreground shrink-0 transition-transform"
        :class="{ 'rotate-180': open }"
      />
    </button>

    <div
      v-if="open"
      class="absolute left-0 right-0 top-full mt-1.5 z-50 rounded-lg border border-border bg-card shadow-lg overflow-hidden"
    >
      <ul class="max-h-72 overflow-y-auto py-1" role="listbox">
        <li v-for="(opt, i) in options" :key="opt.value">
          <button
            type="button"
            role="option"
            :aria-selected="opt.value === modelValue"
            :class="[
              'w-full flex items-center gap-2 px-3 py-2 text-left text-sm transition-colors',
              i === activeIndex ? 'bg-muted' : 'hover:bg-muted',
              opt.value === modelValue ? 'text-primary font-medium' : 'text-foreground'
            ]"
            @click="select(opt)"
            @mouseenter="activeIndex = i"
          >
            <Check v-if="opt.value === modelValue" class="w-4 h-4 shrink-0" />
            <span v-else class="w-4 h-4 shrink-0" />
            <span class="truncate">{{ opt.label }}</span>
          </button>
        </li>
      </ul>
    </div>
  </div>
</template>
