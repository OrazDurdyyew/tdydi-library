<script setup lang="ts">
definePageMeta({ layout: false })

const route = useRoute()
const config = useRuntimeConfig()
const api = useApi()
const { t } = useI18n()

const slug = computed(() => route.params.slug as string)

const { data: book, error } = await useAsyncData(
  () => `read-book-${slug.value}`,
  () => api.get<any>(`/books/slug/${slug.value}`)
)

if (error.value || !book.value || !book.value.file_format) {
  throw createError({ statusCode: 404, statusMessage: 'Book file not found', fatal: true })
}

const isEpub = computed(() => book.value?.file_format === 'epub')
const fileSrc = computed(() => `${config.public.apiBase}/books/${book.value?.id}/file`)
const downloadUrl = computed(() =>
  book.value?.allow_download ? `${config.public.apiBase}/books/${book.value?.id}/download` : undefined
)

useHead({ title: () => `${book.value?.title || ''} — ${t('book.read')}` })
</script>

<template>
  <PdfViewer
    v-if="book"
    :title="book.title"
    :src="fileSrc"
    :storage-key="`book-bookmark-${book.id}`"
    :back-to="`/books/${slug}`"
    :download-url="downloadUrl"
    :is-epub="isEpub"
  />
</template>
