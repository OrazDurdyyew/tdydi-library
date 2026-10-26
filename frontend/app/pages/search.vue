<script setup lang="ts">
import { FileText, HelpCircle, Search } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const api = useApi()
const { t } = useI18n()
const loc = useLocalized()

const query = ref((route.query.q as string) || '')
const submitted = computed(() => (route.query.q as string) || '')

const { data, pending } = await useAsyncData(
  'search-results',
  () => {
    const q = submitted.value.trim()
    if (q.length < 2) {
      return Promise.resolve({ books: [], periodicals: [], posts: [], faq: [], pages: [] })
    }
    return api
      .get<any>(`/search?q=${encodeURIComponent(q)}`)
      .catch(() => ({ books: [], periodicals: [], posts: [], faq: [], pages: [] }))
  },
  { watch: [submitted] }
)

const results = computed(
  () => data.value ?? { books: [], periodicals: [], posts: [], faq: [], pages: [] }
)
const totalCount = computed(
  () =>
    results.value.books.length +
    results.value.periodicals.length +
    results.value.posts.length +
    results.value.faq.length +
    results.value.pages.length
)

function submitSearch() {
  const q = query.value.trim()
  router.replace({ query: q ? { q } : {} })
}

useSeoMeta({
  title: () => (submitted.value ? `${t('search.title')}: ${submitted.value}` : t('search.title'))
})
</script>

<template>
  <div>
    <PageHero :title="t('search.title')">
      <template #breadcrumbs>
        <Breadcrumbs :items="[{ label: t('nav.home'), to: '/' }, { label: t('search.title') }]" />
      </template>
    </PageHero>

    <div class="page-container py-10">
      <form
        class="flex items-center h-12 px-4 mb-10 rounded-xl border border-border bg-card focus-within:border-ring focus-within:ring-2 focus-within:ring-ring/15 transition-colors max-w-2xl"
        @submit.prevent="submitSearch"
      >
        <Search class="w-5 h-5 text-muted-foreground shrink-0" />
        <input
          v-model="query"
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

      <div v-if="!submitted" class="card p-16 text-center text-muted-foreground">
        {{ t('search.type_query') }}
      </div>

      <div v-else-if="pending" class="text-muted-foreground py-12 text-center">
        {{ t('common.loading') }}
      </div>

      <div v-else-if="totalCount === 0" class="card p-16 text-center text-muted-foreground">
        {{ t('search.nothing') }}
      </div>

      <div v-else class="space-y-12">
        <p class="text-sm text-muted-foreground">
          {{ t('search.results_for', { query: submitted }) }}
        </p>

        <section v-if="results.books.length">
          <h2 class="text-lg font-bold tracking-tight mb-5">{{ t('search.books') }}</h2>
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-x-5 gap-y-8">
            <BookCard v-for="b in results.books" :key="b.id" :book="b" />
          </div>
        </section>

        <section v-if="results.periodicals.length">
          <h2 class="text-lg font-bold tracking-tight mb-5">{{ t('search.periodicals') }}</h2>
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-x-5 gap-y-8">
            <PeriodicalCard v-for="p in results.periodicals" :key="p.id" :periodical="p" />
          </div>
        </section>

        <section v-if="results.posts.length">
          <h2 class="text-lg font-bold tracking-tight mb-5">{{ t('search.news') }}</h2>
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
            <PostCard v-for="p in results.posts" :key="p.id" :post="p" />
          </div>
        </section>

        <section v-if="results.faq.length">
          <h2 class="text-lg font-bold tracking-tight mb-5">{{ t('search.faq') }}</h2>
          <div class="space-y-2">
            <NuxtLinkLocale
              v-for="item in results.faq"
              :key="item.id"
              to="/faq"
              class="flex items-start gap-3 card card-hover p-4"
            >
              <HelpCircle class="w-5 h-5 icon-primary shrink-0 mt-0.5" />
              <div>
                <div class="font-medium">{{ loc(item, 'question') }}</div>
                <div
                  class="text-sm text-muted-foreground line-clamp-2 mt-1"
                  v-html="loc(item, 'answer')"
                />
              </div>
            </NuxtLinkLocale>
          </div>
        </section>

        <section v-if="results.pages.length">
          <h2 class="text-lg font-bold tracking-tight mb-5">{{ t('search.pages') }}</h2>
          <div class="space-y-2">
            <NuxtLinkLocale
              v-for="pg in results.pages"
              :key="pg.id"
              :to="`/${pg.slug}`"
              class="flex items-center gap-3 card card-hover p-4"
            >
              <FileText class="w-5 h-5 icon-primary shrink-0" />
              <span class="font-medium">{{ loc(pg, 'title') }}</span>
            </NuxtLinkLocale>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>
