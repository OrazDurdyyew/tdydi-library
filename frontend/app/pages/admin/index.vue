<script setup lang="ts">
import { BookMarked, Eye, Download, TrendingUp, Plus, FileText, Users, GraduationCap, Landmark } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const { admin } = useAuth()
const api = useApi()

const isSuperadmin = computed(() => admin.value?.role === 'superadmin')
const canSeePublications = computed(
  () => admin.value?.role === 'superadmin' || admin.value?.role === 'publication_manager'
)

interface Stats {
  total_books: number
  institute_books: number
  total_views: number
  total_downloads: number
  top_books: { id: number; title: string; authors: string; views_count: number; downloads_count: number }[]
  daily_activity: { date: string; views: number }[]
  top_sections: { section_id: number; section_name: string; book_count: number; views_total: number }[]
  period_days: number
}

interface PubStats {
  total_publications: number
  total_teachers: number
  total_departments: number
  top_teachers: { id: number; name: string; publication_count: number }[]
  per_department: { department: string; count: number }[]
}

const { data: stats, pending } = await useAsyncData<Stats | null>('admin-stats', () =>
  isSuperadmin.value ? api.get<Stats>('/admin/stats') : Promise.resolve(null)
)

const { data: pubStats } = await useAsyncData<PubStats | null>('admin-pub-stats', () =>
  canSeePublications.value
    ? api.get<PubStats>('/admin/stats/publications').catch(() => null)
    : Promise.resolve(null)
)

const CHART_COLORS = [
  '#2563eb', '#16a34a', '#d97706', '#dc2626', '#7c3aed',
  '#0891b2', '#db2777', '#65a30d', '#ea580c', '#4f46e5'
]

const topTeachersChart = computed(() => {
  const list = pubStats.value?.top_teachers ?? []
  return {
    labels: list.map((x) => x.name),
    datasets: [
      {
        label: t('admin.dashboard.publications_count_label'),
        data: list.map((x) => x.publication_count),
        backgroundColor: '#2563eb',
        borderRadius: 4
      }
    ]
  }
})

const perDepartmentBarChart = computed(() => {
  const list = pubStats.value?.per_department ?? []
  return {
    labels: list.map((x) => x.department),
    datasets: [
      {
        label: t('admin.dashboard.publications_count_label'),
        data: list.map((x) => x.count),
        backgroundColor: '#16a34a',
        borderRadius: 4
      }
    ]
  }
})

const perDepartmentPieChart = computed(() => {
  const list = pubStats.value?.per_department ?? []
  return {
    labels: list.map((x) => x.department),
    datasets: [
      {
        data: list.map((x) => x.count),
        backgroundColor: list.map((_, i) => CHART_COLORS[i % CHART_COLORS.length])
      }
    ]
  }
})

const hasPubChartData = computed(
  () =>
    (pubStats.value?.top_teachers?.length ?? 0) > 0 ||
    (pubStats.value?.per_department?.length ?? 0) > 0
)

useHead({ title: () => t('admin.dashboard.title') })
</script>

<template>
  <div>
    <div class="mb-8">
      <h1 class="text-3xl font-bold font-serif">{{ t('admin.dashboard.title') }}</h1>
      <p v-if="admin" class="text-muted-foreground mt-1">
        {{ t('admin.dashboard.welcome') }} <strong class="text-foreground">{{ admin.username }}</strong>
        <span class="text-sm opacity-70"> · {{ admin.role }}</span>
      </p>
    </div>

    <!-- Метрики -->
    <div v-if="!pending && stats" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5 mb-8">
      <Card>
        <CardContent class="flex items-center gap-4 pt-6">
          <span class="p-3 rounded-md bg-primary/10">
            <BookMarked class="w-5 h-5 text-primary" />
          </span>
          <div>
            <div class="text-2xl font-bold font-serif">{{ stats.total_books }}</div>
            <div class="text-sm text-muted-foreground">{{ t('admin.dashboard.books_total') }}</div>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent class="flex items-center gap-4 pt-6">
          <span class="p-3 rounded-md bg-primary/10">
            <Landmark class="w-5 h-5 text-primary" />
          </span>
          <div>
            <div class="text-2xl font-bold font-serif">{{ stats.institute_books }}</div>
            <div class="text-sm text-muted-foreground">{{ t('admin.dashboard.institute_books') }}</div>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent class="flex items-center gap-4 pt-6">
          <span class="p-3 rounded-md bg-primary/10">
            <Eye class="w-5 h-5 text-primary" />
          </span>
          <div>
            <div class="text-2xl font-bold font-serif">{{ stats.total_views.toLocaleString() }}</div>
            <div class="text-sm text-muted-foreground">{{ t('admin.dashboard.views_total') }}</div>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent class="flex items-center gap-4 pt-6">
          <span class="p-3 rounded-md bg-primary/10">
            <Download class="w-5 h-5 text-primary" />
          </span>
          <div>
            <div class="text-2xl font-bold font-serif">{{ stats.total_downloads.toLocaleString() }}</div>
            <div class="text-sm text-muted-foreground">{{ t('admin.dashboard.downloads_total') }}</div>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Скелетон загрузки -->
    <div v-if="pending" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5 mb-8">
      <Card v-for="i in 4" :key="i" class="h-20 animate-pulse bg-muted" />
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Топ книг -->
      <Card v-if="stats" class="p-0 overflow-hidden">
        <CardHeader class="border-b border-border">
          <div class="flex items-center gap-2">
            <TrendingUp class="w-4 h-4 text-primary" />
            <CardTitle class="text-base">{{ t('admin.dashboard.top_books') }}</CardTitle>
            <span class="text-xs text-muted-foreground ml-auto">{{ t('admin.dashboard.period_days_label', { days: stats.period_days }) }}</span>
          </div>
        </CardHeader>
        <CardContent class="p-0">
          <div v-if="stats.top_books.length === 0" class="px-6 py-10 text-center text-muted-foreground text-sm">
            {{ t('admin.dashboard.no_data') }}
          </div>
          <ul v-else>
            <li
              v-for="(book, idx) in stats.top_books"
              :key="book.id"
              class="flex items-center gap-3 px-6 py-3 border-b border-border last:border-0 hover:bg-muted/30 transition-colors"
            >
              <span class="w-5 text-center text-sm font-mono text-muted-foreground shrink-0">{{ idx + 1 }}</span>
              <div class="min-w-0 flex-1">
                <NuxtLink :to="`/admin/books/${book.id}`" class="font-medium text-sm hover:text-primary transition-colors truncate block">
                  {{ book.title }}
                </NuxtLink>
                <div class="text-xs text-muted-foreground truncate">{{ book.authors || '—' }}</div>
              </div>
              <div class="text-right shrink-0">
                <div class="text-sm font-medium">{{ book.views_count }}</div>
                <div class="text-xs text-muted-foreground">{{ t('admin.dashboard.views_short') }}</div>
              </div>
            </li>
          </ul>
        </CardContent>
      </Card>

      <!-- Топ разделов -->
      <Card v-if="stats && stats.top_sections.length > 0" class="p-0 overflow-hidden">
        <CardHeader class="border-b border-border">
          <CardTitle class="text-base">{{ t('admin.dashboard.active_sections') }}</CardTitle>
        </CardHeader>
        <CardContent class="p-0">
          <ul>
            <li
              v-for="s in stats.top_sections"
              :key="s.section_id"
              class="flex items-center gap-3 px-6 py-3 border-b border-border last:border-0"
            >
              <div class="min-w-0 flex-1">
                <div class="text-sm font-medium truncate">{{ s.section_name }}</div>
                <div class="text-xs text-muted-foreground">{{ t('admin.dashboard.books_count', { count: s.book_count }) }}</div>
              </div>
              <div class="text-right shrink-0 text-sm">{{ s.views_total.toLocaleString() }} {{ t('admin.dashboard.views_short') }}</div>
            </li>
          </ul>
        </CardContent>
      </Card>

      <!-- Быстрый доступ -->
      <Card v-if="isSuperadmin && !pending && (!stats || stats.top_books.length === 0)">
        <CardHeader>
          <CardTitle class="text-base">{{ t('admin.dashboard.quick_access') }}</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="grid grid-cols-2 gap-3">
            <NuxtLink to="/admin/books/new">
              <Button class="w-full"><Plus class="w-4 h-4" /> {{ t('admin.books.add_new') }}</Button>
            </NuxtLink>
            <NuxtLink to="/admin/sections">
              <Button variant="outline" class="w-full">{{ t('admin.nav.sections') }}</Button>
            </NuxtLink>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- ===== Публикации преподавателей ===== -->
    <section v-if="canSeePublications && pubStats" class="mt-10">
      <h2 class="text-xl font-bold font-serif mb-5">{{ t('admin.dashboard.publications_section') }}</h2>

      <div class="grid grid-cols-1 sm:grid-cols-3 gap-5 mb-6">
        <Card>
          <CardContent class="flex items-center gap-4 pt-6">
            <span class="p-3 rounded-md bg-primary/10">
              <FileText class="w-5 h-5 text-primary" />
            </span>
            <div>
              <div class="text-2xl font-bold font-serif">{{ pubStats.total_publications.toLocaleString() }}</div>
              <div class="text-sm text-muted-foreground">{{ t('admin.dashboard.publications_total') }}</div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent class="flex items-center gap-4 pt-6">
            <span class="p-3 rounded-md bg-primary/10">
              <Users class="w-5 h-5 text-primary" />
            </span>
            <div>
              <div class="text-2xl font-bold font-serif">{{ pubStats.total_teachers.toLocaleString() }}</div>
              <div class="text-sm text-muted-foreground">{{ t('admin.dashboard.teachers_total') }}</div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent class="flex items-center gap-4 pt-6">
            <span class="p-3 rounded-md bg-primary/10">
              <GraduationCap class="w-5 h-5 text-primary" />
            </span>
            <div>
              <div class="text-2xl font-bold font-serif">{{ pubStats.total_departments.toLocaleString() }}</div>
              <div class="text-sm text-muted-foreground">{{ t('admin.dashboard.departments_total') }}</div>
            </div>
          </CardContent>
        </Card>
      </div>

      <div v-if="hasPubChartData" class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card v-if="pubStats.top_teachers.length">
          <CardHeader class="border-b border-border">
            <CardTitle class="text-base">{{ t('admin.dashboard.top_teachers') }}</CardTitle>
          </CardHeader>
          <CardContent class="pt-6">
            <ClientOnly>
              <BarChart :data="topTeachersChart" />
              <template #fallback>
                <div class="h-64 animate-pulse bg-muted rounded" />
              </template>
            </ClientOnly>
          </CardContent>
        </Card>

        <Card v-if="pubStats.per_department.length">
          <CardHeader class="border-b border-border">
            <CardTitle class="text-base">{{ t('admin.dashboard.per_department') }}</CardTitle>
          </CardHeader>
          <CardContent class="pt-6">
            <ClientOnly>
              <BarChart :data="perDepartmentBarChart" />
              <template #fallback>
                <div class="h-64 animate-pulse bg-muted rounded" />
              </template>
            </ClientOnly>
          </CardContent>
        </Card>

        <Card v-if="pubStats.per_department.length">
          <CardHeader class="border-b border-border">
            <CardTitle class="text-base">{{ t('admin.dashboard.per_department_share') }}</CardTitle>
          </CardHeader>
          <CardContent class="pt-6">
            <ClientOnly>
              <PieChart :data="perDepartmentPieChart" />
              <template #fallback>
                <div class="h-64 animate-pulse bg-muted rounded" />
              </template>
            </ClientOnly>
          </CardContent>
        </Card>
      </div>

      <Card v-else>
        <CardContent class="py-10 text-center text-muted-foreground text-sm">
          {{ t('admin.dashboard.no_data') }}
        </CardContent>
      </Card>
    </section>
  </div>
</template>
