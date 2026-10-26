<script setup lang="ts">
import { ArrowLeft, Newspaper } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const api = useApi()
const { t } = useI18n()

const LIMIT = 24

const type = ref((route.query.type as string) || '')
const page = ref(Number(route.query.page) || 1)

const queryString = computed(() => {
  const p = new URLSearchParams()
  p.set('limit', String(LIMIT))
  p.set('page', String(page.value))
  if (type.value) p.set('type', type.value)
  return p.toString()
})

const { data, pending } = await useAsyncData(
  'periodicals-titles',
  () =>
    api
      .get<{ data: any[]; total: number }>(`/periodicals?${queryString.value}`)
      .catch(() => ({ data: [], total: 0 })),
  { watch: [queryString] }
)

const items = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.total ?? 0)

const tabs = [
  { value: '', label: 'periodicals.all' },
  { value: 'newspaper', label: 'periodicals.newspapers' },
  { value: 'magazine', label: 'periodicals.magazines' }
]

function syncUrl() {
  const q: Record<string, string> = {}
  if (type.value) q.type = type.value
  if (page.value > 1) q.page = String(page.value)
  router.replace({ query: q })
}

function selectType(v: string) {
  type.value = v
  page.value = 1
  syncUrl()
}

function changePage(p: number) {
  page.value = p
  syncUrl()
  if (import.meta.client) window.scrollTo({ top: 0, behavior: 'smooth' })
}

useSeoMeta({
  title: () => t('periodicals.titles_title'),
  description: () => t('periodicals.titles_title')
})
</script>

<template>
  <div>
    <PageHero :title="t('periodicals.titles_title')">
      <template #breadcrumbs>
        <Breadcrumbs
          :items="[
            { label: t('nav.home'), to: '/' },
            { label: t('periodicals.title'), to: '/periodicals' },
            { label: t('periodicals.titles_title') }
          ]"
        />
      </template>
    </PageHero>

    <div class="page-container py-10">
      <div class="flex flex-wrap items-center justify-between gap-4 mb-8">
        <!-- Фильтр по типу -->
        <div class="inline-flex p-1 rounded-lg border border-border bg-card">
          <button
            v-for="tab in tabs"
            :key="tab.value"
            class="type-tab"
            :class="{ 'type-tab-active': type === tab.value }"
            @click="selectType(tab.value)"
          >
            {{ t(tab.label) }}
          </button>
        </div>

        <NuxtLinkLocale
          to="/periodicals"
          class="inline-flex items-center gap-1.5 h-10 px-4 rounded-lg border border-border bg-card text-sm font-medium hover:border-ring hover:text-primary transition-colors"
        >
          <ArrowLeft class="w-4 h-4" />
          {{ t('periodicals.back_to_issues') }}
        </NuxtLinkLocale>
      </div>

      <div v-if="pending" class="text-muted-foreground py-16 text-center">
        {{ t('common.loading') }}
      </div>

      <div
        v-else-if="items.length"
        class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-x-5 gap-y-8"
      >
        <PeriodicalCard v-for="p in items" :key="p.id" :periodical="p" />
      </div>

      <div v-else class="card p-16 text-center">
        <Newspaper class="w-10 h-10 mx-auto text-muted-foreground/50 mb-3" />
        <p class="text-muted-foreground">{{ t('periodicals.no_periodicals') }}</p>
      </div>

      <Pagination :page="page" :total="total" :limit="LIMIT" @change="changePage" />
    </div>
  </div>
</template>

<style scoped>
.type-tab {
  padding: 0.5rem 1rem;
  border-radius: 0.4rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: rgb(var(--muted-foreground));
  transition: background-color 140ms ease, color 140ms ease;
}
.type-tab:hover {
  color: rgb(var(--foreground));
}
.type-tab-active {
  background-color: rgb(var(--primary));
  color: rgb(var(--primary-foreground));
}
</style>
