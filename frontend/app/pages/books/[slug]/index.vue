<script setup lang="ts">
import { BookOpen, Download, Eye, ArrowDownToLine, QrCode } from 'lucide-vue-next'
import QRCode from 'qrcode'

const route = useRoute()
const config = useRuntimeConfig()
const api = useApi()
const { t } = useI18n()
const loc = useLocalized()
const fileUrl = useFileUrl()

const slug = computed(() => route.params.slug as string)

const { data: book, error } = await useAsyncData(
  () => `book-${slug.value}`,
  () => api.get<any>(`/books/slug/${slug.value}`)
)

if (error.value || !book.value) {
  throw createError({ statusCode: 404, statusMessage: 'Book not found', fatal: true })
}

const { data: similarResp } = await useAsyncData(
  () => `book-similar-${book.value?.id}`,
  () =>
    book.value
      ? api.get<{ data: any[] }>(`/books/${book.value.id}/similar`).catch(() => ({ data: [] }))
      : Promise.resolve({ data: [] })
)

const similar = computed(() => similarResp.value?.data ?? [])
const coverUrl = computed(() => fileUrl(book.value?.cover_path))
const hasFile = computed(() => !!book.value?.file_format)
const downloadUrl = computed(() => `${config.public.apiBase}/books/${book.value?.id}/download`)

const languageLabel = computed(() => {
  const l = book.value?.language
  if (l === 'ru') return t('catalog.lang_ru')
  if (l === 'en') return t('catalog.lang_en')
  return t('catalog.lang_tm')
})

const breadcrumbs = computed(() => [
  { label: t('nav.home'), to: '/' },
  { label: t('nav.books'), to: '/books' },
  { label: book.value?.title || '' }
])

const details = computed(() => {
  const b = book.value
  if (!b) return []
  return [
    { label: t('book.authors'), value: b.authors },
    { label: t('book.year'), value: b.year },
    { label: t('book.publisher'), value: b.publisher },
    { label: t('book.isbn'), value: b.isbn },
    { label: t('book.pages'), value: b.pages },
    { label: t('book.language'), value: languageLabel.value }
  ].filter((d) => !!d.value)
})

const qrDataUrl = ref('')
onMounted(async () => {
  try {
    qrDataUrl.value = await QRCode.toDataURL(window.location.href, {
      margin: 1,
      width: 160,
      color: { dark: '#1E3A5F', light: '#FFFFFF' }
    })
  } catch {
    /* QR — необязательный элемент */
  }
})

useSeoMeta({
  title: () => book.value?.title || '',
  description: () => book.value?.annotation || '',
  ogImage: () => coverUrl.value || undefined
})
</script>

<template>
  <div v-if="book" class="page-container py-10">
    <Breadcrumbs :items="breadcrumbs" />

    <div class="grid grid-cols-1 lg:grid-cols-[300px_1fr] gap-10 lg:gap-14 mt-2">
      <!-- Левая колонка: обложка + действия + метаданные -->
      <aside class="space-y-6">
        <div class="book-card aspect-[3/4] max-w-[300px] mx-auto lg:mx-0 w-full">
          <img
            v-if="coverUrl"
            :src="coverUrl"
            :alt="book.title"
            class="w-full h-full object-cover"
          />
          <div
            v-else
            class="w-full h-full flex flex-col items-center justify-center gap-3 p-6 text-muted-foreground"
          >
            <BookOpen class="w-12 h-12" />
            <span class="text-sm text-center">{{ book.title }}</span>
          </div>
        </div>

        <div class="flex flex-col gap-2.5">
          <NuxtLinkLocale
            v-if="hasFile"
            :to="`/books/${book.slug}/read`"
            class="btn-primary w-full"
          >
            <BookOpen class="w-4 h-4" /> {{ t('book.read') }}
          </NuxtLinkLocale>
          <p v-else class="text-sm text-muted-foreground text-center py-2">
            {{ t('book.no_file') }}
          </p>

          <a
            v-if="hasFile && book.allow_download"
            :href="downloadUrl"
            class="btn-outline w-full"
          >
            <Download class="w-4 h-4" /> {{ t('book.download') }}
          </a>
          <p v-else-if="hasFile" class="text-xs text-muted-foreground text-center">
            {{ t('book.download_off') }}
          </p>
        </div>

        <!-- Метаданные -->
        <div class="card p-5">
          <h2 class="text-sm font-semibold mb-3">{{ t('book.details') }}</h2>
          <dl class="space-y-2.5 text-sm">
            <div
              v-for="d in details"
              :key="d.label"
              class="flex justify-between gap-4 py-1.5 border-b border-border last:border-0"
            >
              <dt class="text-muted-foreground shrink-0">{{ d.label }}</dt>
              <dd class="text-right font-medium">{{ d.value }}</dd>
            </div>
          </dl>
          <div class="flex items-center gap-5 mt-4 pt-4 border-t border-border text-sm text-muted-foreground">
            <span class="inline-flex items-center gap-1.5">
              <Eye class="w-4 h-4 icon-primary" /> {{ book.views_count }}
            </span>
            <span class="inline-flex items-center gap-1.5">
              <ArrowDownToLine class="w-4 h-4 icon-primary" /> {{ book.downloads_count }}
            </span>
          </div>
        </div>

        <!-- QR -->
        <div v-if="qrDataUrl" class="card p-5 flex items-center gap-4">
          <img :src="qrDataUrl" :alt="t('book.share')" class="w-20 h-20 rounded-lg border border-border" />
          <div>
            <div class="inline-flex items-center gap-1.5 text-sm font-semibold">
              <QrCode class="w-4 h-4 icon-primary" /> {{ t('book.share') }}
            </div>
            <p class="text-xs text-muted-foreground mt-1">{{ t('book.qr_hint') }}</p>
          </div>
        </div>
      </aside>

      <!-- Основная колонка -->
      <div>
        <h1 class="text-2xl md:text-4xl font-bold tracking-tight leading-tight">
          {{ book.title }}
        </h1>
        <p v-if="book.authors" class="text-lg text-muted-foreground mt-2">
          {{ book.authors }}
        </p>

        <!-- Разделы -->
        <div v-if="book.sections?.length" class="flex flex-wrap gap-2 mt-5">
          <NuxtLinkLocale
            v-for="s in book.sections"
            :key="s.id"
            :to="`/sections/${s.slug}`"
            class="inline-flex items-center px-2.5 py-1 rounded-md border border-border bg-muted text-sm hover:border-secondary/40 hover:text-secondary transition-colors"
          >
            {{ loc(s, 'name') }}
          </NuxtLinkLocale>
        </div>

        <!-- Аннотация -->
        <div v-if="book.annotation" class="mt-8">
          <h2 class="text-lg font-semibold mb-3">{{ t('book.annotation') }}</h2>
          <p class="serif text-[1.05rem] leading-[1.75] text-foreground/90 whitespace-pre-line max-w-[68ch]">
            {{ book.annotation }}
          </p>
        </div>

        <!-- Похожие книги -->
        <section v-if="similar.length" class="mt-12 pt-10 border-t border-border">
          <h2 class="text-xl font-bold tracking-tight mb-6">{{ t('book.similar') }}</h2>
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-x-5 gap-y-8">
            <BookCard v-for="b in similar" :key="b.id" :book="b" />
          </div>
        </section>
      </div>
    </div>
  </div>
</template>
