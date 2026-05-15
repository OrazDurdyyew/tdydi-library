<script setup lang="ts">
import { Search, SlidersHorizontal, X, BookOpen } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const api = useApi()
const { t } = useI18n()
const loc = useLocalized()

const LIMIT = 24

const filters = reactive({
  q: (route.query.q as string) || '',
  section_id: (route.query.section_id as string) || '',
  year: (route.query.year as string) || '',
  language: (route.query.language as string) || '',
  author: (route.query.author as string) || '',
  sort: (route.query.sort as string) || '-created_at',
  page: Number(route.query.page) || 1
})

const queryString = computed(() => {
  const p = new URLSearchParams()
  p.set('limit', String(LIMIT))
  p.set('page', String(filters.page))
  p.set('sort', filters.sort)
  if (filters.q) p.set('q', filters.q)
  if (filters.section_id) p.set('section_id', filters.section_id)
  if (filters.year) p.set('year', filters.year)
  if (filters.language) p.set('language', filters.language)
  if (filters.author) p.set('author', filters.author)
  return p.toString()
})

const { data, pending } = await useAsyncData(
  'catalog-books',
  () =>
    api
      .get<{ data: any[]; meta: { total: number; page: number; limit: number } }>(
        `/books?${queryString.value}`
      )
      .catch(() => ({ data: [], meta: { total: 0, page: 1, limit: LIMIT } })),
  { watch: [queryString] }
)

const { data: sectionsResp } = await useAsyncData('catalog-sections', () =>
  api.get<{ data: any[] }>('/sections').catch(() => ({ data: [] }))
)

const books = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.meta.total ?? 0)
const sections = computed(() => sectionsResp.value?.data ?? [])

const sectionOptions = computed(() => {
  const roots = sections.value.filter((s) => s.parent_id == null)
  const out: { id: number; label: string }[] = []
  for (const r of roots) {
    out.push({ id: r.id, label: loc(r, 'name') })
    for (const c of sections.value.filter((s) => s.parent_id === r.id)) {
      out.push({ id: c.id, label: '— ' + loc(c, 'name') })
    }
  }
  return out
})

function syncUrl() {
  const q: Record<string, string> = {}
  if (filters.q) q.q = filters.q
  if (filters.section_id) q.section_id = filters.section_id
  if (filters.year) q.year = filters.year
  if (filters.language) q.language = filters.language
  if (filters.author) q.author = filters.author
  if (filters.sort !== '-created_at') q.sort = filters.sort
  if (filters.page > 1) q.page = String(filters.page)
  router.replace({ query: q })
}

function applyFilters() {
  filters.page = 1
  syncUrl()
}

function resetFilters() {
  filters.q = ''
  filters.section_id = ''
  filters.year = ''
  filters.language = ''
  filters.author = ''
  filters.sort = '-created_at'
  filters.page = 1
  syncUrl()
}

function changePage(p: number) {
  filters.page = p
  syncUrl()
  if (import.meta.client) window.scrollTo({ top: 0, behavior: 'smooth' })
}

const hasActiveFilters = computed(
  () =>
    !!filters.q ||
    !!filters.section_id ||
    !!filters.year ||
    !!filters.language ||
    !!filters.author
)

useSeoMeta({
  title: () => t('catalog.title'),
  description: () => t('home.subtitle')
})
</script>

<template>
  <div>
    <PageHero :title="t('catalog.title')">
      <template #breadcrumbs>
        <Breadcrumbs :items="[{ label: t('nav.home'), to: '/' }, { label: t('catalog.title') }]" />
      </template>
    </PageHero>

    <div class="page-container py-10">
      <!-- Поиск -->
      <form
        class="flex items-center h-12 px-4 mb-8 rounded-xl border border-border bg-card focus-within:border-ring focus-within:ring-2 focus-within:ring-ring/15 transition-colors max-w-2xl"
        @submit.prevent="applyFilters"
      >
        <Search class="w-5 h-5 text-muted-foreground shrink-0" />
        <input
          v-model="filters.q"
          type="search"
          :placeholder="t('search.placeholder')"
          class="flex-1 h-full bg-transparent border-0 outline-none px-3 text-sm placeholder:text-muted-foreground"
        />
        <button
          type="submit"
          class="h-9 px-4 rounded-lg bg-primary text-primary-foreground text-sm font-medium hover:bg-secondary transition-colors"
        >
          {{ t('common.search') }}
        </button>
      </form>

      <div class="grid grid-cols-1 lg:grid-cols-[260px_1fr] gap-8">
        <!-- Фильтры -->
        <aside class="card p-5 h-fit lg:sticky lg:top-36">
          <div class="flex items-center gap-2 mb-5">
            <SlidersHorizontal class="w-4 h-4 icon-primary" />
            <h2 class="font-semibold">{{ t('catalog.filters') }}</h2>
          </div>

          <div class="space-y-4">
            <div>
              <label class="block text-xs font-medium text-muted-foreground mb-1.5">
                {{ t('catalog.section') }}
              </label>
              <select v-model="filters.section_id" class="input" @change="applyFilters">
                <option value="">{{ t('catalog.any') }}</option>
                <option v-for="o in sectionOptions" :key="o.id" :value="String(o.id)">
                  {{ o.label }}
                </option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-medium text-muted-foreground mb-1.5">
                {{ t('catalog.year') }}
              </label>
              <input
                v-model="filters.year"
                type="number"
                inputmode="numeric"
                class="input"
                placeholder="—"
                @change="applyFilters"
              />
            </div>

            <div>
              <label class="block text-xs font-medium text-muted-foreground mb-1.5">
                {{ t('catalog.author') }}
              </label>
              <input
                v-model="filters.author"
                type="text"
                class="input"
                placeholder="—"
                @keyup.enter="applyFilters"
              />
            </div>

            <div>
              <label class="block text-xs font-medium text-muted-foreground mb-1.5">
                {{ t('catalog.language') }}
              </label>
              <select v-model="filters.language" class="input" @change="applyFilters">
                <option value="">{{ t('catalog.any') }}</option>
                <option value="tm">{{ t('catalog.lang_tm') }}</option>
                <option value="ru">{{ t('catalog.lang_ru') }}</option>
                <option value="en">{{ t('catalog.lang_en') }}</option>
              </select>
            </div>

            <button
              v-if="hasActiveFilters"
              type="button"
              class="btn-ghost w-full text-sm"
              @click="resetFilters"
            >
              <X class="w-4 h-4" /> {{ t('common.reset') }}
            </button>
          </div>
        </aside>

        <!-- Результаты -->
        <div>
          <div class="flex flex-wrap items-center justify-between gap-3 mb-6">
            <p class="text-sm text-muted-foreground">
              {{ t('catalog.found', { count: total }) }}
            </p>
            <select v-model="filters.sort" class="input max-w-[220px]" @change="applyFilters">
              <option value="-created_at">{{ t('catalog.sort_newest') }}</option>
              <option value="created_at">{{ t('catalog.sort_oldest') }}</option>
              <option value="title">{{ t('catalog.sort_title') }}</option>
              <option value="-year">{{ t('catalog.sort_year') }}</option>
              <option value="views">{{ t('catalog.sort_popular') }}</option>
            </select>
          </div>

          <div v-if="pending" class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-x-5 gap-y-8">
            <div v-for="n in 8" :key="n" class="animate-pulse">
              <div class="aspect-[3/4] rounded-lg bg-muted" />
              <div class="h-3.5 bg-muted rounded mt-3 w-3/4" />
              <div class="h-3 bg-muted rounded mt-2 w-1/2" />
            </div>
          </div>

          <div
            v-else-if="books.length"
            class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-x-5 gap-y-8"
          >
            <BookCard v-for="book in books" :key="book.id" :book="book" />
          </div>

          <div v-else class="card p-16 text-center">
            <BookOpen class="w-10 h-10 mx-auto text-muted-foreground/50 mb-3" />
            <p class="text-muted-foreground">{{ t('catalog.no_results') }}</p>
            <button v-if="hasActiveFilters" class="btn-outline mt-4 text-sm" @click="resetFilters">
              {{ t('common.reset') }}
            </button>
          </div>

          <Pagination :page="filters.page" :total="total" :limit="LIMIT" @change="changePage" />
        </div>
      </div>
    </div>
  </div>
</template>
