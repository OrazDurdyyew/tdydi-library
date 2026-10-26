<script setup lang="ts">
const route = useRoute()
const api = useApi()
const { t } = useI18n()
const loc = useLocalized()
const fileUrl = useFileUrl()

const slug = computed(() => route.params.slug as string)

const { data: collection, error } = await useAsyncData(
  () => `collection-${slug.value}`,
  () => api.get<any>(`/collections/slug/${slug.value}`)
)

if (error.value || !collection.value) {
  throw createError({ statusCode: 404, statusMessage: 'Collection not found', fatal: true })
}

const name = computed(() => loc(collection.value, 'name'))
const description = computed(() => loc(collection.value, 'description'))
const coverUrl = computed(() => fileUrl(collection.value?.cover_path))
const allBooks = computed(() => collection.value?.books ?? [])

const LIMIT = 18
const page = ref(1)
const books = computed(() => allBooks.value.slice((page.value - 1) * LIMIT, page.value * LIMIT))
function changePage(p: number) {
  page.value = p
  if (import.meta.client) window.scrollTo({ top: 0, behavior: 'smooth' })
}

useSeoMeta({
  title: () => name.value,
  description: () => description.value || name.value
})
</script>

<template>
  <div v-if="collection">
    <PageHero :title="name" :description="description">
      <template #breadcrumbs>
        <Breadcrumbs
          :items="[
            { label: t('nav.home'), to: '/' },
            { label: t('collections.title'), to: '/collections' },
            { label: name }
          ]"
        />
      </template>
    </PageHero>

    <div class="page-container py-10">
      <div
        v-if="coverUrl"
        class="rounded-xl overflow-hidden border border-border max-h-[340px] mb-10"
      >
        <img :src="coverUrl" :alt="name" class="w-full h-full object-cover" />
      </div>

      <h2 class="text-xl md:text-2xl font-bold tracking-tight mb-6">
        {{ t('collections.in_collection') }}
      </h2>

      <div
        v-if="books.length"
        class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-x-5 gap-y-8"
      >
        <BookCard v-for="book in books" :key="book.id" :book="book" />
      </div>

      <div v-else class="card p-12 text-center text-muted-foreground">
        {{ t('collections.empty') }}
      </div>

      <Pagination :page="page" :total="allBooks.length" :limit="LIMIT" @change="changePage" />
    </div>
  </div>
</template>
