<script setup lang="ts">
import { Newspaper, BookText, Calendar, ChevronRight } from 'lucide-vue-next'

const route = useRoute()
const api = useApi()
const { t, locale } = useI18n()
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

const { data: issuesResp } = await useAsyncData(
  () => `periodical-issues-${periodical.value?.id}`,
  () =>
    api
      .get<{ data: any[]; total: number }>(`/periodicals/${periodical.value.id}/issues?limit=100`)
      .catch(() => ({ data: [], total: 0 }))
)

const issues = computed(() => issuesResp.value?.data ?? [])
const name = computed(() => loc(periodical.value, 'name'))
const description = computed(() => loc(periodical.value, 'description'))
const coverUrl = computed(() => fileUrl(periodical.value?.cover_path))
const typeLabel = computed(() =>
  periodical.value?.type === 'magazine' ? t('periodicals.magazine') : t('periodicals.newspaper')
)

const byYear = computed(() => {
  const groups: Record<string, any[]> = {}
  for (const issue of issues.value) {
    const year = new Date(issue.issue_date).getFullYear().toString()
    ;(groups[year] ||= []).push(issue)
  }
  return Object.entries(groups).sort((a, b) => Number(b[0]) - Number(a[0]))
})

const activeYear = ref<string>('')
watchEffect(() => {
  if (!activeYear.value && byYear.value.length) {
    activeYear.value = byYear.value[0]?.[0] ?? ''
  }
})
const activeIssues = computed(
  () => byYear.value.find(([y]) => y === activeYear.value)?.[1] ?? []
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
    <div class="grid grid-cols-1 md:grid-cols-[220px_1fr] gap-8 md:gap-10 mt-2 mb-12">
      <div class="book-card aspect-[3/4] max-w-[220px] w-full mx-auto md:mx-0">
        <img v-if="coverUrl" :src="coverUrl" :alt="name" class="w-full h-full object-cover" />
        <div v-else class="w-full h-full flex items-center justify-center text-muted-foreground">
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

      <div v-if="issues.length" class="grid grid-cols-1 lg:grid-cols-[180px_1fr] gap-8">
        <!-- Архив по годам -->
        <aside class="lg:sticky lg:top-36 h-fit">
          <h3 class="eyebrow mb-3">{{ t('periodicals.archive') }}</h3>
          <ul class="space-y-1">
            <li v-for="[year, list] in byYear" :key="year">
              <button
                class="year-btn"
                :class="{ 'year-btn-active': year === activeYear }"
                @click="activeYear = year"
              >
                <span>{{ year }}</span>
                <span class="text-xs opacity-60">{{ list.length }}</span>
              </button>
            </li>
          </ul>
        </aside>

        <!-- Список выпусков выбранного года -->
        <div class="space-y-2.5">
          <NuxtLinkLocale
            v-for="issue in activeIssues"
            :key="issue.id"
            :to="`/periodicals/issues/${issue.id}`"
            class="group flex items-center gap-4 card card-hover p-4"
          >
            <div class="w-12 h-16 shrink-0 rounded-md overflow-hidden bg-muted border border-border">
              <img
                v-if="issue.cover_path"
                :src="fileUrl(issue.cover_path)"
                :alt="issue.issue_number"
                class="w-full h-full object-cover"
              />
              <div v-else class="w-full h-full flex items-center justify-center text-muted-foreground">
                <Calendar class="w-4 h-4" />
              </div>
            </div>
            <div class="flex-1 min-w-0">
              <div class="font-semibold group-hover:text-secondary transition-colors">
                {{ issue.issue_number || t('periodicals.issue') }}
              </div>
              <div class="text-sm text-muted-foreground">
                {{ formatDate(issue.issue_date, locale) }}
              </div>
              <p
                v-if="loc(issue, 'annotation')"
                class="text-sm text-muted-foreground/80 mt-0.5 line-clamp-1"
              >
                {{ loc(issue, 'annotation') }}
              </p>
            </div>
            <ChevronRight class="w-4 h-4 text-muted-foreground group-hover:text-secondary transition-colors shrink-0" />
          </NuxtLinkLocale>
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
