<script setup lang="ts">
import { FileText, SlidersHorizontal, X, Download, GraduationCap } from 'lucide-vue-next'

const route = useRoute()
const config = useRuntimeConfig()
const api = useApi()
const { t } = useI18n()
const loc = useLocalized()
const fileUrl = useFileUrl()
const localePath = useLocalePath()

const teacherSlug = computed(() => route.params.teacherSlug as string)
const kafedraSlug = computed(() => route.params.kafedraSlug as string)

const { data: teacherData, error } = await useAsyncData(
  () => `teacher-${teacherSlug.value}`,
  () => api.get<any>(`/teachers/slug/${teacherSlug.value}`)
)

if (error.value || !teacherData.value) {
  throw createError({ statusCode: 404, statusMessage: 'Teacher not found', fatal: true })
}

const teacher = computed(() => teacherData.value)

// Канонический URL: kafedra-slug в адресе декоративный — если не совпадает, редиректим.
const canonicalKafedra = computed(() => teacher.value?.department?.slug || '')
if (import.meta.server && canonicalKafedra.value && canonicalKafedra.value !== kafedraSlug.value) {
  await navigateTo(
    localePath(`/publications/${canonicalKafedra.value}/${teacherSlug.value}`),
    { redirectCode: 301 }
  )
}

function teacherName(tch: any): string {
  if (!tch) return ''
  return [tch.last_name, tch.first_name, tch.middle_name].filter(Boolean).join(' ')
}

// Категории для фильтра — отдельным запросом, не выводим из загруженных публикаций.
const { data: catsResp } = await useAsyncData('teacher-pub-categories', () =>
  api.get<{ data: any[] }>('/publication-categories?limit=200').catch(() => ({ data: [] }))
)
const categoryOptions = computed(() =>
  (catsResp.value?.data ?? []).map((c) => ({ value: String(c.id), label: loc(c, 'name') }))
)

const LIMIT = 15
const filters = reactive({ category_id: '', year: '' })
const page = ref(1)
watch(() => [filters.category_id, filters.year], () => (page.value = 1))

// Все фильтры применяются на бэке через общий /publications.
const pubsQuery = computed(() => {
  const p = new URLSearchParams()
  p.set('teacher_id', String(teacher.value?.id ?? ''))
  p.set('page', String(page.value))
  p.set('limit', String(LIMIT))
  if (filters.category_id) p.set('category_id', filters.category_id)
  if (filters.year) p.set('year', filters.year)
  return p.toString()
})

const { data: pubsResp, pending } = await useAsyncData(
  () => `teacher-pubs-${teacher.value?.id}-${pubsQuery.value}`,
  () =>
    api
      .get<{ data: any[]; meta?: { total: number } }>(`/publications?${pubsQuery.value}`)
      .catch(() => ({ data: [], meta: { total: 0 } })),
  { watch: [pubsQuery] }
)
const publications = computed(() => pubsResp.value?.data ?? [])
const total = computed(() => pubsResp.value?.meta?.total ?? 0)

function changePage(p: number) {
  page.value = p
  if (import.meta.client) window.scrollTo({ top: 0, behavior: 'smooth' })
}

const hasActiveFilters = computed(() => !!filters.category_id || !!filters.year)
function resetFilters() {
  filters.category_id = ''
  filters.year = ''
}

function downloadUrl(pub: any): string {
  return `${config.public.apiBase}/publications/${pub.id}/download`
}

useSeoMeta({
  title: () => teacherName(teacher.value),
  description: () =>
    teacher.value?.department
      ? `${teacherName(teacher.value)} — ${loc(teacher.value.department, 'name')}`
      : teacherName(teacher.value)
})
</script>

<template>
  <div v-if="teacher">
    <PageHero :title="teacherName(teacher)">
      <template #breadcrumbs>
        <Breadcrumbs
          :items="[
            { label: t('nav.home'), to: '/' },
            { label: t('publications.title'), to: '/publications' },
            { label: t('publications.teachers_title'), to: '/publications/teachers' },
            { label: teacherName(teacher) }
          ]"
        />
      </template>
    </PageHero>

    <div class="page-container py-10">
      <div class="grid grid-cols-1 lg:grid-cols-[260px_1fr] gap-8">
        <!-- Карточка преподавателя + фильтры -->
        <aside class="space-y-6 h-fit lg:sticky lg:top-36">
          <div class="card p-5">
            <div class="flex items-center gap-4">
              <span
                v-if="!teacher.photo_path"
                class="inline-flex items-center justify-center w-14 h-14 rounded-full bg-primary/10 text-primary shrink-0"
              >
                <GraduationCap class="w-6 h-6" />
              </span>
              <img
                v-else
                :src="fileUrl(teacher.photo_path)"
                :alt="teacherName(teacher)"
                class="w-14 h-14 rounded-full object-cover shrink-0"
              />
              <div class="min-w-0">
                <div class="font-semibold leading-tight">{{ teacherName(teacher) }}</div>
                <NuxtLinkLocale
                  v-if="teacher.department"
                  :to="`/publications/teachers?department_id=${teacher.department.id}`"
                  class="text-sm text-muted-foreground hover:text-secondary transition-colors"
                >
                  {{ loc(teacher.department, 'name') }}
                </NuxtLinkLocale>
              </div>
            </div>
            <p v-if="teacher.bio" class="text-sm text-muted-foreground mt-4 leading-relaxed">
              {{ teacher.bio }}
            </p>
          </div>

          <div class="card p-5">
            <div class="flex items-center gap-2 mb-5">
              <SlidersHorizontal class="w-4 h-4 icon-primary" />
              <h2 class="font-semibold">{{ t('publications.filters') }}</h2>
            </div>
            <div class="space-y-4">
              <div>
                <label class="block text-xs font-medium text-muted-foreground mb-1.5">
                  {{ t('publications.category') }}
                </label>
                <SelectMenu
                  v-model="filters.category_id"
                  :options="categoryOptions"
                  :placeholder="t('publications.any')"
                />
              </div>
              <div>
                <label class="block text-xs font-medium text-muted-foreground mb-1.5">
                  {{ t('publications.year') }}
                </label>
                <input
                  v-model="filters.year"
                  type="number"
                  inputmode="numeric"
                  class="input"
                  placeholder="—"
                />
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
          </div>
        </aside>

        <!-- Таблица статей -->
        <div>
          <p class="text-sm text-muted-foreground mb-5">
            {{ t('publications.found', { count: total }) }}
          </p>

          <Card v-if="publications.length" class="p-0 overflow-hidden">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>{{ t('publications.col_title') }}</TableHead>
                  <TableHead>{{ t('publications.col_category') }}</TableHead>
                  <TableHead>{{ t('publications.col_date') }}</TableHead>
                  <TableHead>{{ t('publications.col_pages') }}</TableHead>
                  <TableHead class="text-right">{{ t('publications.col_download') }}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-for="pub in publications" :key="pub.id">
                  <TableCell class="font-medium">{{ loc(pub, 'name') }}</TableCell>
                  <TableCell class="text-muted-foreground">
                    {{ pub.category ? loc(pub.category, 'name') : '—' }}
                  </TableCell>
                  <TableCell class="text-muted-foreground">{{ pub.article_date || '—' }}</TableCell>
                  <TableCell class="text-muted-foreground">{{ pub.page_count || '—' }}</TableCell>
                  <TableCell class="text-right">
                    <a
                      v-if="pub.has_file"
                      :href="downloadUrl(pub)"
                      class="inline-flex items-center gap-1.5 text-sm font-medium text-secondary hover:underline"
                    >
                      <Download class="w-4 h-4" /> {{ t('publications.download') }}
                    </a>
                    <span v-else class="text-xs text-muted-foreground">—</span>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </Card>

          <div v-else-if="pending" class="text-muted-foreground py-16 text-center">
            {{ t('common.loading') }}
          </div>

          <div v-else class="card p-16 text-center">
            <FileText class="w-10 h-10 mx-auto text-muted-foreground/50 mb-3" />
            <p class="text-muted-foreground">{{ t('publications.no_results') }}</p>
            <button v-if="hasActiveFilters" class="btn-outline mt-4 text-sm" @click="resetFilters">
              {{ t('common.reset') }}
            </button>
          </div>

          <Pagination :page="page" :total="total" :limit="LIMIT" @change="changePage" />
        </div>
      </div>
    </div>
  </div>
</template>
