<script setup lang="ts">
import { Users } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const api = useApi()
const { t } = useI18n()
const loc = useLocalized()

const LIMIT = 12
const departmentId = ref((route.query.department_id as string) || '')
const page = ref(Number(route.query.page) || 1)

const { data: deptsResp } = await useAsyncData('public-teachers-departments', () =>
  api.get<{ data: any[] }>('/departments?limit=200').catch(() => ({ data: [] }))
)
const departments = computed(() => deptsResp.value?.data ?? [])
const departmentOptions = computed(() =>
  departments.value.map((d) => ({ value: String(d.id), label: loc(d, 'name') }))
)

const teachersQuery = computed(() => {
  const p = new URLSearchParams()
  p.set('page', String(page.value))
  p.set('limit', String(LIMIT))
  if (departmentId.value) p.set('department_id', departmentId.value)
  return p.toString()
})

const { data, pending } = await useAsyncData(
  'public-teachers',
  () =>
    api
      .get<{ data: any[]; meta?: { total: number } }>(`/teachers?${teachersQuery.value}`)
      .catch(() => ({ data: [], meta: { total: 0 } })),
  { watch: [teachersQuery] }
)

const teachers = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.meta?.total ?? 0)

function teacherName(tch: any): string {
  return [tch.last_name, tch.first_name, tch.middle_name].filter(Boolean).join(' ')
}

function syncUrl() {
  const q: Record<string, string> = {}
  if (departmentId.value) q.department_id = departmentId.value
  if (page.value > 1) q.page = String(page.value)
  router.replace({ query: q })
}

function onDepartmentChange() {
  page.value = 1
  syncUrl()
}

function changePage(p: number) {
  page.value = p
  syncUrl()
  if (import.meta.client) window.scrollTo({ top: 0, behavior: 'smooth' })
}

useSeoMeta({
  title: () => t('publications.teachers_title'),
  description: () => t('publications.subtitle')
})
</script>

<template>
  <div>
    <PageHero :title="t('publications.teachers_title')">
      <template #breadcrumbs>
        <Breadcrumbs
          :items="[
            { label: t('nav.home'), to: '/' },
            { label: t('publications.title'), to: '/publications' },
            { label: t('publications.teachers_title') }
          ]"
        />
      </template>
    </PageHero>

    <div class="page-container py-10">
      <div class="grid grid-cols-1 lg:grid-cols-[260px_1fr] gap-8">
        <aside class="card p-5 h-fit lg:sticky lg:top-36">
          <h2 class="font-semibold mb-4">{{ t('publications.department') }}</h2>
          <SelectMenu
            v-model="departmentId"
            :options="departmentOptions"
            :placeholder="t('publications.all_departments')"
            @change="onDepartmentChange"
          />
        </aside>

        <div>
          <div v-if="pending" class="text-muted-foreground py-16 text-center">
            {{ t('common.loading') }}
          </div>

          <div
            v-else-if="teachers.length"
            class="grid grid-cols-1 sm:grid-cols-2 gap-4"
          >
            <NuxtLinkLocale
              v-for="tch in teachers"
              :key="tch.id"
              :to="tch.department ? `/publications/${tch.department.slug}/${tch.slug}` : '/publications'"
              class="card card-hover p-5 flex items-center gap-4"
            >
              <span class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-primary/10 text-primary shrink-0">
                <Users class="w-5 h-5" />
              </span>
              <div class="min-w-0">
                <div class="font-semibold truncate">{{ teacherName(tch) }}</div>
                <div v-if="tch.department" class="text-sm text-muted-foreground truncate">
                  {{ loc(tch.department, 'name') }}
                </div>
              </div>
            </NuxtLinkLocale>
          </div>

          <div v-else class="card p-16 text-center">
            <Users class="w-10 h-10 mx-auto text-muted-foreground/50 mb-3" />
            <p class="text-muted-foreground">{{ t('publications.no_teachers') }}</p>
          </div>

          <Pagination :page="page" :total="total" :limit="LIMIT" @change="changePage" />
        </div>
      </div>
    </div>
  </div>
</template>
