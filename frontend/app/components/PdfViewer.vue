<script setup lang="ts">
import {
  ArrowLeft,
  ChevronLeft,
  ChevronRight,
  Download,
  ZoomIn,
  ZoomOut,
  Loader2
} from 'lucide-vue-next'

const props = defineProps<{
  title: string
  src: string
  storageKey: string
  backTo: string
  downloadUrl?: string
  isEpub?: boolean
}>()

const { t } = useI18n()

const canvas = ref<HTMLCanvasElement | null>(null)
const pdfDoc = shallowRef<any>(null)
const currentPage = ref(1)
const totalPages = ref(0)
const scale = ref(1.3)
const loading = ref(true)
const pageInput = ref('1')
let renderTask: any = null

async function loadPdf() {
  try {
    const pdfjsLib: any = await import('pdfjs-dist')
    const worker = await import('pdfjs-dist/build/pdf.worker.min.mjs?url')
    pdfjsLib.GlobalWorkerOptions.workerSrc = worker.default

    const doc = await pdfjsLib.getDocument({ url: props.src, withCredentials: false }).promise
    pdfDoc.value = doc
    totalPages.value = doc.numPages

    const saved = Number(localStorage.getItem(props.storageKey))
    if (saved >= 1 && saved <= doc.numPages) currentPage.value = saved

    loading.value = false
    await nextTick()
    await renderPage()
  } catch (e) {
    loading.value = false
    console.error('PDF load error', e)
  }
}

async function renderPage() {
  if (!pdfDoc.value || !canvas.value) return
  const page = await pdfDoc.value.getPage(currentPage.value)
  const viewport = page.getViewport({ scale: scale.value })
  const ctx = canvas.value.getContext('2d')
  if (!ctx) return
  canvas.value.width = viewport.width
  canvas.value.height = viewport.height
  if (renderTask) {
    try {
      renderTask.cancel()
    } catch {
      /* предыдущий рендер уже завершён */
    }
  }
  renderTask = page.render({ canvasContext: ctx, viewport })
  try {
    await renderTask.promise
  } catch {
    return // рендер отменён более новым переходом
  }
  pageInput.value = String(currentPage.value)
  localStorage.setItem(props.storageKey, String(currentPage.value))
}

function go(delta: number) {
  const next = currentPage.value + delta
  if (next >= 1 && next <= totalPages.value) currentPage.value = next
}

function jumpTo() {
  const n = Number(pageInput.value)
  if (n >= 1 && n <= totalPages.value) currentPage.value = n
  else pageInput.value = String(currentPage.value)
}

function zoom(delta: number) {
  scale.value = Math.min(3, Math.max(0.5, +(scale.value + delta).toFixed(2)))
}

watch([currentPage, scale], renderPage)

function onKey(e: KeyboardEvent) {
  if (e.key === 'ArrowLeft') go(-1)
  if (e.key === 'ArrowRight') go(1)
}

onMounted(() => {
  if (!props.isEpub) loadPdf()
  else loading.value = false
  window.addEventListener('keydown', onKey)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKey)
  try {
    renderTask?.cancel?.()
  } catch {
    /* noop */
  }
  pdfDoc.value?.destroy?.()
})
</script>

<template>
  <div class="reader-root">
    <header class="reader-bar">
      <NuxtLinkLocale :to="backTo" class="reader-btn" :aria-label="t('reader.close')">
        <ArrowLeft class="w-4 h-4" />
        <span class="hidden sm:inline ml-1.5">{{ t('common.back') }}</span>
      </NuxtLinkLocale>

      <div class="flex-1 min-w-0 px-3 text-center">
        <p class="font-serif font-semibold truncate text-sm md:text-base">{{ title }}</p>
      </div>

      <div v-if="!isEpub && !loading" class="flex items-center gap-1.5">
        <button class="reader-btn" :disabled="currentPage <= 1" @click="go(-1)">
          <ChevronLeft class="w-4 h-4" />
        </button>
        <div class="reader-page-box">
          <input
            v-model="pageInput"
            class="reader-page-input"
            inputmode="numeric"
            @keyup.enter="jumpTo"
            @blur="jumpTo"
          />
          <span class="text-muted-foreground text-sm">/ {{ totalPages }}</span>
        </div>
        <button class="reader-btn" :disabled="currentPage >= totalPages" @click="go(1)">
          <ChevronRight class="w-4 h-4" />
        </button>

        <span class="w-px h-6 bg-border mx-1" aria-hidden="true" />

        <button class="reader-btn" @click="zoom(-0.2)"><ZoomOut class="w-4 h-4" /></button>
        <button class="reader-btn" @click="zoom(0.2)"><ZoomIn class="w-4 h-4" /></button>
      </div>

      <a
        v-if="downloadUrl"
        :href="downloadUrl"
        class="reader-btn ml-1.5"
        :aria-label="t('reader.download')"
      >
        <Download class="w-4 h-4" />
      </a>
    </header>

    <main class="reader-canvas-wrap">
      <div v-if="loading" class="flex flex-col items-center gap-3 text-muted-foreground py-32">
        <Loader2 class="w-8 h-8 animate-spin" />
        <p class="serif italic">{{ t('reader.loading') }}</p>
      </div>

      <iframe
        v-else-if="isEpub"
        :src="src"
        class="w-full h-full border-0 bg-white rounded"
        :title="title"
      />

      <canvas v-else ref="canvas" class="reader-canvas" />
    </main>
  </div>
</template>

<style scoped>
.reader-root {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: rgb(var(--paper-deep));
}
.reader-bar {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  height: 56px;
  padding: 0 0.75rem;
  background-color: rgba(255, 253, 247, 0.96);
  backdrop-filter: saturate(180%) blur(12px);
  -webkit-backdrop-filter: saturate(180%) blur(12px);
  border-bottom: 1px solid rgb(var(--border));
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  z-index: 10;
  flex-shrink: 0;
}
:global(.dark) .reader-bar {
  background-color: rgba(47, 38, 28, 0.96);
}
.reader-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 38px;
  height: 38px;
  padding: 0 0.5rem;
  border-radius: 999px;
  background-color: rgb(var(--card));
  border: 1px solid rgb(var(--border));
  color: rgb(var(--ink));
  font-family: var(--font-serif);
  font-size: 0.9rem;
  transition: background-color 140ms ease, color 140ms ease, border-color 140ms ease;
}
.reader-btn:hover:not(:disabled) {
  border-color: rgba(22, 171, 101, 0.5);
  color: rgb(var(--primary));
}
.reader-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
.reader-page-box {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0 0.5rem;
}
.reader-page-input {
  width: 44px;
  height: 34px;
  text-align: center;
  border-radius: 0.5rem;
  background-color: rgb(var(--card));
  border: 1px solid rgb(var(--border));
  font-family: var(--font-serif);
  outline: none;
}
.reader-page-input:focus {
  border-color: rgb(var(--primary));
}
.reader-canvas-wrap {
  flex: 1;
  overflow: auto;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding: 2rem 1rem;
}
.reader-canvas {
  max-width: 100%;
  height: auto;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.25);
  border-radius: 2px;
  background: #fff;
}
</style>
