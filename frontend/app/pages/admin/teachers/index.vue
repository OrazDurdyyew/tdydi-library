<script setup lang="ts">
import { Plus, ArrowLeft } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const loc = useLocalized()
const route = useRoute()
const router = useRouter()

interface Department {
  id: number
  name_tm: string
  name_ru: string
  name_en: string
}

interface Teacher {
  id: number
  slug: string
  first_name: string
  last_name: string
  middle_name: string
  department_id: number
  department?: Department
}

const api = useApi()
const toast = useToast()
const { ask } = useConfirm()

const LIMIT = 20
const departmentFilter = ref((route.query.department_id as string) || '')
const page = ref(Number(route.query.page) || 1)

const { data: deptsResp } = await useAsyncData('teachers-departments', () =>
  api.get<{ data: Department[] }>('/departments?limit=200')
)
const departments = computed(() => deptsResp.value?.data ?? [])

const teachersQuery = computed(() => {
  const p = new URLSearchParams()
  p.set('page', String(page.value))
  p.set('limit', String(LIMIT))
  if (departmentFilter.value) p.set('department_id', departmentFilter.value)
  return p.toString()
})

const { data, refresh, pending } = await useAsyncData(
  'admin-teachers',
  () => api.get<{ data: Teacher[]; meta: { total: number } }>(`/teachers?${teachersQuery.value}`),
  { watch: [teachersQuery] }
)

const items = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.meta?.total ?? 0)

function teacherName(tch: Teacher): string {
  return [tch.last_name, tch.first_name, tch.middle_name].filter(Boolean).join(' ')
}

function onFilterChange() {
  page.value = 1
  syncUrl()
}

function changePage(p: number) {
  page.value = p
  syncUrl()
}

function syncUrl() {
  const q: Record<string, string> = {}
  if (departmentFilter.value) q.department_id = departmentFilter.value
  if (page.value > 1) q.page = String(page.value)
  router.replace({ query: q })
}

async function remove(id: number, name: string) {
  const ok = await ask({
    title: t('admin.teachers.confirm_delete_title'),
    description: t('admin.teachers.confirm_delete_description', { name }),
    confirmText: t('admin.common.delete'),
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/teachers/${id}`)
    await refresh()
    toast.success(t('admin.teachers.deleted'))
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error_delete'))
  }
}

useHead({ title: () => t('admin.teachers.title') })
</script>

<template>
  <div>
    <NuxtLink
      :to="$localePath('/admin/teacher-publications')"
      class="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground mb-4"
    >
      <ArrowLeft class="w-4 h-4" /> {{ t('admin.teacher_publications.back') }}
    </NuxtLink>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold font-serif">{{ t('admin.teachers.title') }}</h1>
      <NuxtLink to="/admin/teachers/new">
        <Button><Plus class="w-4 h-4" /> {{ t('admin.teachers.add_new') }}</Button>
      </NuxtLink>
    </div>

    <div class="mb-6 max-w-xs">
      <Label for="filter-dept" class="mb-1.5 block">{{ t('admin.teachers.filter_department') }}</Label>
      <Select id="filter-dept" v-model="departmentFilter" @change="onFilterChange">
        <option value="">{{ t('admin.teachers.all_departments') }}</option>
        <option v-for="d in departments" :key="d.id" :value="String(d.id)">{{ loc(d, 'name') }}</option>
      </Select>
    </div>

    <div v-if="pending" class="text-muted-foreground">{{ t('admin.common.loading') }}</div>

    <Card v-else class="p-0 overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>{{ t('admin.teachers.full_name') }}</TableHead>
            <TableHead>{{ t('admin.teachers.department') }}</TableHead>
            <TableHead>{{ t('admin.common.slug') }}</TableHead>
            <TableHead class="text-right">{{ t('admin.common.actions') }}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="tch in items" :key="tch.id">
            <TableCell class="font-medium">{{ teacherName(tch) }}</TableCell>
            <TableCell class="text-muted-foreground">{{ tch.department ? loc(tch.department, 'name') : '—' }}</TableCell>
            <TableCell class="font-mono text-xs text-muted-foreground">{{ tch.slug }}</TableCell>
            <TableCell class="text-right">
              <div class="flex justify-end gap-2">
                <NuxtLink :to="`/admin/teachers/${tch.id}`">
                  <Button variant="ghost" size="sm">{{ t('admin.common.edit') }}</Button>
                </NuxtLink>
                <Button
                  variant="ghost"
                  size="sm"
                  class="text-destructive hover:text-destructive"
                  @click="remove(tch.id, teacherName(tch))"
                >
                  {{ t('admin.common.delete') }}
                </Button>
              </div>
            </TableCell>
          </TableRow>
          <TableEmpty v-if="items.length === 0" :colspan="4">{{ t('admin.teachers.empty_hint') }}</TableEmpty>
        </TableBody>
      </Table>
    </Card>

    <Pagination :page="page" :total="total" :limit="LIMIT" @change="changePage" />
  </div>
</template>
