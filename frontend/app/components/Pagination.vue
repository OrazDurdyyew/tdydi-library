<script setup lang="ts">
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'

const props = defineProps<{
  page: number
  total: number
  limit: number
}>()

const emit = defineEmits<{ (e: 'change', page: number): void }>()

const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.limit)))

const pageList = computed(() => {
  const tp = totalPages.value
  const cur = props.page
  const out: (number | '…')[] = []
  for (let i = 1; i <= tp; i++) {
    if (i === 1 || i === tp || (i >= cur - 1 && i <= cur + 1)) {
      out.push(i)
    } else if (out[out.length - 1] !== '…') {
      out.push('…')
    }
  }
  return out
})

function go(p: number) {
  if (p < 1 || p > totalPages.value || p === props.page) return
  emit('change', p)
}
</script>

<template>
  <nav v-if="totalPages > 1" class="flex items-center justify-center gap-1 mt-10">
    <button class="pg-btn" :disabled="page <= 1" :aria-label="'Previous'" @click="go(page - 1)">
      <ChevronLeft class="w-4 h-4" />
    </button>
    <template v-for="(p, i) in pageList" :key="i">
      <span v-if="p === '…'" class="px-2 text-muted-foreground select-none">…</span>
      <button
        v-else
        class="pg-btn"
        :class="{ 'pg-btn-active': p === page }"
        @click="go(p as number)"
      >
        {{ p }}
      </button>
    </template>
    <button
      class="pg-btn"
      :disabled="page >= totalPages"
      :aria-label="'Next'"
      @click="go(page + 1)"
    >
      <ChevronRight class="w-4 h-4" />
    </button>
  </nav>
</template>

<style scoped>
.pg-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 2.5rem;
  height: 2.5rem;
  padding: 0 0.5rem;
  border-radius: 0.5rem;
  border: 1px solid rgb(var(--border));
  background-color: rgb(var(--card));
  font-size: 0.875rem;
  font-weight: 500;
  color: rgb(var(--foreground));
  transition: background-color 140ms ease, border-color 140ms ease, color 140ms ease;
}
.pg-btn:hover:not(:disabled) {
  background-color: rgb(var(--muted));
  border-color: rgb(var(--secondary) / 0.4);
}
.pg-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
.pg-btn-active {
  background-color: rgb(var(--primary));
  color: rgb(var(--primary-foreground));
  border-color: rgb(var(--primary));
}
.pg-btn-active:hover {
  background-color: rgb(var(--primary));
  color: rgb(var(--primary-foreground));
}
</style>
