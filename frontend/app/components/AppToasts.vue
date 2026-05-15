<script setup lang="ts">
import { CheckCircle, XCircle, Info, X } from 'lucide-vue-next'

const { toasts, remove } = useToast()

const icons = { success: CheckCircle, error: XCircle, info: Info }
const colors = {
  success: 'border-primary bg-primary/10 text-primary',
  error: 'border-destructive bg-destructive/10 text-destructive',
  info: 'border-border bg-card text-foreground'
}
</script>

<template>
  <Teleport to="body">
    <div class="fixed bottom-6 right-6 z-[9999] flex flex-col gap-2 w-80 max-w-[calc(100vw-2rem)]">
      <TransitionGroup name="toast">
        <div
          v-for="t in toasts"
          :key="t.id"
          class="flex items-start gap-3 px-4 py-3 rounded-md border shadow-lg text-sm"
          :class="colors[t.type]"
        >
          <component :is="icons[t.type]" class="w-4 h-4 mt-0.5 shrink-0" />
          <span class="flex-1">{{ t.message }}</span>
          <button @click="remove(t.id)" class="opacity-60 hover:opacity-100 transition-opacity mt-0.5">
            <X class="w-3.5 h-3.5" />
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
.toast-enter-active { transition: all 0.2s ease; }
.toast-leave-active { transition: all 0.25s ease; }
.toast-enter-from { opacity: 0; transform: translateY(8px); }
.toast-leave-to { opacity: 0; transform: translateX(100%); }
</style>
