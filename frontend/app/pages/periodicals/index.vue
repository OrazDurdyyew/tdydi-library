<script setup lang="ts">
import { ArrowRight, Newspaper } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const api = useApi()
const { t } = useI18n()
const localePath = useLocalePath()

const LIMIT = 24

function today() {
  return new Date().toISOString().slice(0, 10)
}

const type = ref((route.query.type as string) || '')
const before = ref((route.query.before as string) || today())
const page = ref(Number(route.query.page) || 1)

const queryString = computed(() => {
  const p = new URLSearchParams()
  p.set('limit', String(LIMIT))
  p.set('page', String(page.value))
  if (type.value) p.set('type', type.value)
  if (before.value) p.set('before', before.value)
  return p.toString()
})

const { data, pending } = await useAsyncData(
  'periodicals-issues-feed',
  () =>
    api
      .get<{ data: any[]; total: number }>(`/periodicals/issues?${queryString.value}`)
      .catch(() => ({ data: [], total: 0 })),
  { watch: [queryString] }
)

const issues = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.total ?? 0)

const tabs = [
  { value: '', label: 'periodicals.all' },
  { value: 'newspaper', label: 'periodicals.newspapers' },
  { value: 'magazine', label: 'periodicals.magazines' }
]

function syncUrl() {
  const q: Record<string, string> = {}
  if (type.value) q.type = type.value
  if (before.value && before.value !== today()) q.before = before.value
  if (page.value > 1) q.page = String(page.value)
  router.replace({ query: q })
}

function selectType(v: string) {
  type.value = v
  page.value = 1
  syncUrl()
}

function onDateChange() {
  page.value = 1
  syncUrl()
}

function changePage(p: number) {
  page.value = p
  syncUrl()
  if (import.meta.client) window.scrollTo({ top: 0, behavior: 'smooth' })
}

useSeoMeta({
  title: () => t('periodicals.latest_issues'),
  description: () => t('periodicals.title')
})
</script>

<template>
  <div>
    <PageHero :title="t('periodicals.title')">
      <template #breadcrumbs>
        <Breadcrumbs
          :items="[{ label: t('nav.home'), to: '/' }, { label: t('periodicals.title') }]"
        />
      </template>
    </PageHero>

    <div class="page-container py-10">
      <!-- Управление: тип + дата + ссылка на список изданий -->
      <div class="flex flex-wrap items-end justify-between gap-4 mb-8">
        <div class="flex flex-wrap items-end gap-4">
          <!-- Тип -->
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

          <!-- Дата -->
          <div>
            <label class="block text-xs font-medium text-muted-foreground mb-1.5">
              {{ t('periodicals.issues_up_to') }}
            </label>
            <input
              v-model="before"
              type="date"
              class="input max-w-[180px]"
              @change="onDateChange"
            />
          </div>
        </div>

        <NuxtLinkLocale
          to="/periodicals/titles"
          class="inline-flex items-center gap-1.5 h-10 px-4 rounded-lg border border-border bg-card text-sm font-medium hover:border-ring hover:text-primary transition-colors"
        >
          {{ t('periodicals.all_titles') }}
          <ArrowRight class="w-4 h-4" />
        </NuxtLinkLocale>
      </div>

      <h2 class="text-lg font-bold tracking-tight mb-5">{{ t('periodicals.latest_issues') }}</h2>

      <div v-if="pending" class="text-muted-foreground py-16 text-center">
        {{ t('common.loading') }}
      </div>

      <div
        v-else-if="issues.length"
        class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-x-5 gap-y-8"
      >
        <IssueCard v-for="issue in issues" :key="issue.id" :issue="issue" />
      </div>

      <div v-else class="card p-16 text-center">
        <Newspaper class="w-10 h-10 mx-auto text-muted-foreground/50 mb-3" />
        <p class="text-muted-foreground">{{ t('periodicals.no_issues_for_date') }}</p>
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
