<script setup lang="ts">
import { BookText, Newspaper } from 'lucide-vue-next'

const route = useRoute()
const api = useApi()
const { t } = useI18n()
const loc = useLocalized()
const fileUrl = useFileUrl()

const slug = computed(() => route.params.slug as string)

const { data: periodical, error } = await useAsyncData(
  () => `periodical-${slug.value}`,
  () => api.get<any>(`/periodicals/slug/${slug.value}`)
)

if (error.value || !periodical.value) {
  throw createError({ statusCode: 404, statusMessage: 'Periodical not found', fatal: true })
}

// Архив по годам — счётчики и список годов приходят с бэка, выпуски грузятся постранично по выбранному году.
const { data: yearsResp } = await useAsyncData(
  () => `periodical-years-${periodical.value?.id}`,
  () =>
    api
      .get<{ data: { year: number; count: number }[] }>(
        `/periodicals/${periodical.value.id}/issue-years`
      )
      .catch(() => ({ data: [] }))
)
const years = computed(() => yearsResp.value?.data ?? [])

const PER_PAGE = 24
const activeYear = ref<number | null>(years.value[0]?.year ?? null)
const page = ref(1)

const { data: issuesResp, pending: issuesPending } = await useAsyncData(
  () => `periodical-issues-${periodical.value?.id}-${activeYear.value}-${page.value}`,
  () =>
    activeYear.value !== null
      ? api
          .get<{ data: any[]; total: number }>(
            `/periodicals/${periodical.value.id}/issues?year=${activeYear.value}&page=${page.value}&limit=${PER_PAGE}`
          )
          .catch(() => ({ data: [], total: 0 }))
      : Promise.resolve({ data: [] as any[], total: 0 }),
  { watch: [activeYear, page] }
)

const issues = computed(() => issuesResp.value?.data ?? [])
const issuesTotal = computed(() => issuesResp.value?.total ?? 0)

function selectYear(y: number) {
  activeYear.value = y
  page.value = 1
}
function changePage(p: number) {
  page.value = p
  if (import.meta.client) window.scrollTo({ top: 0, behavior: 'smooth' })
}

const name = computed(() => loc(periodical.value, 'name'))
const description = computed(() => loc(periodical.value, 'description'))
const coverUrl = computed(() => fileUrl(periodical.value?.cover_path))
const typeLabel = computed(() =>
  periodical.value?.type === 'magazine' ? t('periodicals.magazine') : t('periodicals.newspaper')
)

useSeoMeta({
  title: () => name.value,
  description: () => description.value || name.value
})
</script>

<template>
  <div v-if="periodical" class="page-container py-10">
    <Breadcrumbs
      :items="[
        { label: t('nav.home'), to: '/' },
        { label: t('periodicals.title'), to: '/periodicals' },
        { label: name }
      ]"
    />

    <!-- Шапка издания -->
    <div class="grid grid-cols-1 md:grid-cols-[550px_1fr] gap-8 md:gap-10 mt-2 mb-12">
      <div class="book-card max-w-[550px] w-full mx-auto md:mx-0 overflow-hidden">
        <img v-if="coverUrl" :src="coverUrl" :alt="name" class="block w-full h-auto" />
        <div v-else class="aspect-[350/145] w-full flex items-center justify-center text-muted-foreground">
          <Newspaper v-if="periodical.type === 'newspaper'" class="w-12 h-12" />
          <BookText v-else class="w-12 h-12" />
        </div>
      </div>
      <div>
        <span class="eyebrow">{{ typeLabel }}</span>
        <h1 class="text-2xl md:text-4xl font-bold tracking-tight mt-1.5 mb-3">{{ name }}</h1>
        <p v-if="periodical.periodicity" class="text-sm text-muted-foreground mb-4">
          {{ t('periodicals.periodicity') }}: {{ periodical.periodicity }}
        </p>
        <p
          v-if="description"
          class="serif text-[1.05rem] leading-[1.75] text-foreground/90 whitespace-pre-line max-w-[68ch]"
        >
          {{ description }}
        </p>
      </div>
    </div>

    <!-- Выпуски -->
    <section>
      <h2 class="text-xl md:text-2xl font-bold tracking-tight mb-6">{{ t('periodicals.issues') }}</h2>

      <div v-if="years.length" class="grid grid-cols-1 lg:grid-cols-[180px_1fr] gap-8">
        <!-- Архив по годам -->
        <aside class="lg:sticky lg:top-36 h-fit">
          <h3 class="eyebrow mb-3">{{ t('periodicals.archive') }}</h3>
          <ul class="space-y-1">
            <li v-for="y in years" :key="y.year">
              <button
                class="year-btn"
                :class="{ 'year-btn-active': y.year === activeYear }"
                @click="selectYear(y.year)"
              >
                <span>{{ y.year }}</span>
                <span class="text-xs opacity-60">{{ y.count }}</span>
              </button>
            </li>
          </ul>
        </aside>

        <!-- Сетка выпусков выбранного года -->
        <div>
          <div
            v-if="issues.length"
            class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-x-5 gap-y-8"
          >
            <IssueCard v-for="issue in issues" :key="issue.id" :issue="issue" />
          </div>
          <div v-else-if="issuesPending" class="text-muted-foreground py-10 text-center">
            {{ t('common.loading') }}
          </div>

          <Pagination :page="page" :total="issuesTotal" :limit="PER_PAGE" @change="changePage" />
        </div>
      </div>

      <div v-else class="card p-12 text-center text-muted-foreground">
        {{ t('periodicals.no_issues') }}
      </div>
    </section>
  </div>
</template>

<style scoped>
.year-btn {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 0.5rem 0.75rem;
  border-radius: 0.5rem;
  font-size: 0.9rem;
  font-weight: 500;
  color: rgb(var(--foreground));
  transition: background-color 140ms ease, color 140ms ease;
}
.year-btn:hover {
  background-color: rgb(var(--muted));
}
.year-btn-active {
  background-color: rgb(var(--primary));
  color: rgb(var(--primary-foreground));
}
</style>
