<script setup lang="ts">
import { FileText, SlidersHorizontal, X, Download } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const config = useRuntimeConfig()
const api = useApi()
const { t } = useI18n()
const loc = useLocalized()

const LIMIT = 24

const filters = reactive({
  department_id: (route.query.department_id as string) || '',
  teacher_id: (route.query.teacher_id as string) || '',
  category_id: (route.query.category_id as string) || '',
  year: (route.query.year as string) || '',
  page: Number(route.query.page) || 1
})

const queryString = computed(() => {
  const p = new URLSearchParams()
  p.set('limit', String(LIMIT))
  p.set('page', String(filters.page))
  if (filters.department_id) p.set('department_id', filters.department_id)
  if (filters.teacher_id) p.set('teacher_id', filters.teacher_id)
  if (filters.category_id) p.set('category_id', filters.category_id)
  if (filters.year) p.set('year', filters.year)
  return p.toString()
})

const { data, pending } = await useAsyncData(
  'public-publications',
  () =>
    api
      .get<{ data: any[]; meta: { total: number; page: number; limit: number } }>(
        `/publications?${queryString.value}`
      )
      .catch(() => ({ data: [], meta: { total: 0, page: 1, limit: LIMIT } })),
  { watch: [queryString] }
)

const { data: deptsResp } = await useAsyncData('public-pub-departments', () =>
  api.get<{ data: any[] }>('/departments?limit=200').catch(() => ({ data: [] }))
)
const { data: catsResp } = await useAsyncData('public-pub-categories', () =>
  api.get<{ data: any[] }>('/publication-categories?limit=200').catch(() => ({ data: [] }))
)

const publications = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.meta?.total ?? 0)
const departments = computed(() => deptsResp.value?.data ?? [])
const categories = computed(() => catsResp.value?.data ?? [])

// Каскад: преподаватели подгружаются под выбранную кафедру.
const teachers = ref<any[]>([])
async function loadTeachers(deptId: string) {
  if (!deptId) {
    teachers.value = []
    return
  }
  try {
    const res = await api.get<{ data: any[] }>(`/teachers?department_id=${deptId}`)
    teachers.value = res.data ?? []
  } catch {
    teachers.value = []
  }
}
watch(() => filters.department_id, (val) => {
  filters.teacher_id = ''
  loadTeachers(val)
})
onMounted(() => loadTeachers(filters.department_id))

function teacherName(tch: any): string {
  return [tch.last_name, tch.first_name, tch.middle_name].filter(Boolean).join(' ')
}

const departmentOptions = computed(() =>
  departments.value.map((d) => ({ value: String(d.id), label: loc(d, 'name') }))
)
const teacherOptions = computed(() =>
  teachers.value.map((tch) => ({ value: String(tch.id), label: teacherName(tch) }))
)
const categoryOptions = computed(() =>
  categories.value.map((c) => ({ value: String(c.id), label: loc(c, 'name') }))
)

function syncUrl() {
  const q: Record<string, string> = {}
  if (filters.department_id) q.department_id = filters.department_id
  if (filters.teacher_id) q.teacher_id = filters.teacher_id
  if (filters.category_id) q.category_id = filters.category_id
  if (filters.year) q.year = filters.year
  if (filters.page > 1) q.page = String(filters.page)
  router.replace({ query: q })
}

function applyFilters() {
  filters.page = 1
  syncUrl()
}

function resetFilters() {
  filters.department_id = ''
  filters.teacher_id = ''
  filters.category_id = ''
  filters.year = ''
  filters.page = 1
  syncUrl()
}

function changePage(p: number) {
  filters.page = p
  syncUrl()
  if (import.meta.client) window.scrollTo({ top: 0, behavior: 'smooth' })
}

const hasActiveFilters = computed(
  () => !!filters.department_id || !!filters.teacher_id || !!filters.category_id || !!filters.year
)

function downloadUrl(pub: any): string {
  return `${config.public.apiBase}/publications/${pub.id}/download`
}

useSeoMeta({
  title: () => t('publications.title'),
  description: () => t('publications.subtitle')
})
</script>

<template>
  <div>
    <PageHero :title="t('publications.title')" :description="t('publications.subtitle')">
      <template #breadcrumbs>
        <Breadcrumbs :items="[{ label: t('nav.home'), to: '/' }, { label: t('publications.title') }]" />
      </template>
    </PageHero>

    <div class="page-container py-10">
      <div class="grid grid-cols-1 lg:grid-cols-[260px_1fr] gap-8">
        <!-- Фильтры -->
        <aside class="card p-5 h-fit lg:sticky lg:top-36">
          <div class="flex items-center gap-2 mb-5">
            <SlidersHorizontal class="w-4 h-4 icon-primary" />
            <h2 class="font-semibold">{{ t('publications.filters') }}</h2>
          </div>

          <div class="space-y-4">
            <div>
              <label class="block text-xs font-medium text-muted-foreground mb-1.5">
                {{ t('publications.department') }}
              </label>
              <SelectMenu
                v-model="filters.department_id"
                :options="departmentOptions"
                :placeholder="t('publications.any')"
                @change="applyFilters"
              />
            </div>

            <div>
              <label class="block text-xs font-medium text-muted-foreground mb-1.5">
                {{ t('publications.teacher') }}
              </label>
              <SelectMenu
                v-model="filters.teacher_id"
                :options="teacherOptions"
                :placeholder="t('publications.any')"
                :disabled="!filters.department_id"
                @change="applyFilters"
              />
            </div>

            <div>
              <label class="block text-xs font-medium text-muted-foreground mb-1.5">
                {{ t('publications.category') }}
              </label>
              <SelectMenu
                v-model="filters.category_id"
                :options="categoryOptions"
                :placeholder="t('publications.any')"
                @change="applyFilters"
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
                @change="applyFilters"
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
        </aside>

        <!-- Результаты -->
        <div>
          <p class="text-sm text-muted-foreground mb-5">
            {{ t('publications.found', { count: total }) }}
          </p>

          <div v-if="pending" class="text-muted-foreground py-16 text-center">
            {{ t('common.loading') }}
          </div>

          <Card v-else-if="publications.length" class="p-0 overflow-hidden">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>{{ t('publications.col_title') }}</TableHead>
                  <TableHead>{{ t('publications.col_teacher') }}</TableHead>
                  <TableHead>{{ t('publications.col_category') }}</TableHead>
                  <TableHead>{{ t('publications.col_date') }}</TableHead>
                  <TableHead>{{ t('publications.col_periodical') }}</TableHead>
                  <TableHead>{{ t('publications.col_pages') }}</TableHead>
                  <TableHead class="text-right">{{ t('publications.col_download') }}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-for="pub in publications" :key="pub.id">
                  <TableCell class="font-medium">{{ loc(pub, 'name') }}</TableCell>
                  <TableCell class="text-muted-foreground">
                    <NuxtLinkLocale
                      v-if="pub.teacher && pub.teacher.slug && pub.teacher.department"
                      :to="`/publications/${pub.teacher.department.slug}/${pub.teacher.slug}`"
                      class="hover:text-secondary transition-colors"
                    >
                      {{ teacherName(pub.teacher) }}
                    </NuxtLinkLocale>
                    <span v-else>{{ pub.teacher ? teacherName(pub.teacher) : '—' }}</span>
                  </TableCell>
                  <TableCell class="text-muted-foreground">{{ pub.category ? loc(pub.category, 'name') : '—' }}</TableCell>
                  <TableCell class="text-muted-foreground">{{ pub.article_date || '—' }}</TableCell>
                  <TableCell class="text-muted-foreground">
                    {{ pub.periodical_issue && pub.periodical_issue.periodical ? loc(pub.periodical_issue.periodical, 'name') : '—' }}
                  </TableCell>
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

          <div v-else class="card p-16 text-center">
            <FileText class="w-10 h-10 mx-auto text-muted-foreground/50 mb-3" />
            <p class="text-muted-foreground">{{ t('publications.no_results') }}</p>
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
