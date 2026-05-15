<script setup lang="ts">
const route = useRoute()
const api = useApi()
const { t } = useI18n()
const loc = useLocalized()
const fileUrl = useFileUrl()

interface Section {
  id: number
  slug: string
  image_path?: string
  book_count?: number
  parent?: Section | null
  children?: Section[]
  [key: string]: any
}

interface SectionResponse {
  section: Section
  books: any[]
}

const slug = computed(() => route.params.slug as string)

const { data, error } = await useAsyncData(
  () => `section-${slug.value}`,
  () => api.get<SectionResponse>(`/sections/slug/${slug.value}`)
)

if (error.value) {
  throw createError({ statusCode: 404, statusMessage: 'Section not found', fatal: true })
}

const section = computed(() => data.value?.section)
const books = computed(() => data.value?.books ?? [])
const children = computed(() => section.value?.children ?? [])

const sectionName = computed(() => loc(section.value, 'name'))
const sectionDesc = computed(() => loc(section.value, 'description'))
const imageUrl = computed(() => fileUrl(section.value?.image_path))

const breadcrumbs = computed(() => {
  const items: { label: string; to?: string }[] = [{ label: t('nav.home'), to: '/' }]
  if (section.value?.parent) {
    items.push({
      label: loc(section.value.parent, 'name'),
      to: `/sections/${section.value.parent.slug}`
    })
  }
  items.push({ label: sectionName.value })
  return items
})

useSeoMeta({
  title: () => sectionName.value,
  description: () => sectionDesc.value || t('home.subtitle')
})
</script>

<template>
  <div v-if="section">
    <PageHero :title="sectionName" :description="sectionDesc">
      <template #breadcrumbs>
        <Breadcrumbs :items="breadcrumbs" />
      </template>
    </PageHero>

    <div class="page-container py-10 space-y-14">
      <!-- Изображение раздела -->
      <div v-if="imageUrl" class="rounded-xl overflow-hidden border border-border max-h-[380px]">
        <img :src="imageUrl" :alt="sectionName" class="w-full h-full object-cover" />
      </div>

      <!-- Подразделы -->
      <section v-if="children.length">
        <h2 class="text-xl md:text-2xl font-bold tracking-tight mb-6">
          {{ t('section.subsections') }}
        </h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
          <SectionCard v-for="child in children" :key="child.id" :section="child" />
        </div>
      </section>

      <!-- Книги раздела -->
      <section>
        <h2 class="text-xl md:text-2xl font-bold tracking-tight mb-6">
          {{ t('section.books') }}
        </h2>
        <div
          v-if="books.length"
          class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-x-5 gap-y-8"
        >
          <BookCard v-for="book in books" :key="book.id" :book="book" />
        </div>
        <div v-else class="card p-12 text-center text-muted-foreground">
          {{ t('section.no_books') }}
        </div>
      </section>
    </div>
  </div>
</template>
