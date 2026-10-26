<script setup lang="ts">
import { Plus, ArrowLeft } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const loc = useLocalized()

interface Teacher {
  id: number
  first_name: string
  last_name: string
  middle_name: string
}

interface Publication {
  id: number
  slug: string
  name_tm: string
  name_ru: string
  name_en: string
  article_date: string
  is_published: boolean
  has_file: boolean
  teacher?: Teacher
  category?: { id: number; name_tm: string; name_ru: string; name_en: string }
}

const api = useApi()
const toast = useToast()
const { ask } = useConfirm()

const LIMIT = 20
const page = ref(1)

const { data, refresh, pending } = await useAsyncData(
  'admin-publications',
  () =>
    api.get<{ data: Publication[]; meta: { total: number } }>(
      `/publications?page=${page.value}&limit=${LIMIT}`
    ),
  { watch: [page] }
)

const items = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.meta?.total ?? 0)

function teacherName(tch?: Teacher): string {
  if (!tch) return '—'
  return [tch.last_name, tch.first_name, tch.middle_name].filter(Boolean).join(' ')
}

async function remove(id: number, name: string) {
  const ok = await ask({
    title: t('admin.publications.confirm_delete_title'),
    description: t('admin.publications.confirm_delete_description', { name }),
    confirmText: t('admin.common.delete'),
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/publications/${id}`)
    await refresh()
    toast.success(t('admin.publications.deleted'))
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error_delete'))
  }
}

useHead({ title: () => t('admin.publications.title') })
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
      <h1 class="text-3xl font-bold font-serif">{{ t('admin.publications.title') }}</h1>
      <NuxtLink to="/admin/publications/new">
        <Button><Plus class="w-4 h-4" /> {{ t('admin.publications.add_new') }}</Button>
      </NuxtLink>
    </div>

    <div v-if="pending" class="text-muted-foreground">{{ t('admin.common.loading') }}</div>

    <Card v-else class="p-0 overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>{{ t('admin.common.title') }}</TableHead>
            <TableHead>{{ t('admin.publications.teacher') }}</TableHead>
            <TableHead>{{ t('admin.publications.category') }}</TableHead>
            <TableHead>{{ t('admin.publications.article_date') }}</TableHead>
            <TableHead>{{ t('admin.publications.file') }}</TableHead>
            <TableHead>{{ t('admin.common.status') }}</TableHead>
            <TableHead class="text-right">{{ t('admin.common.actions') }}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="p in items" :key="p.id">
            <TableCell class="font-medium">{{ loc(p, 'name') }}</TableCell>
            <TableCell class="text-muted-foreground">{{ teacherName(p.teacher) }}</TableCell>
            <TableCell class="text-muted-foreground">{{ p.category ? loc(p.category, 'name') : '—' }}</TableCell>
            <TableCell class="text-muted-foreground">{{ p.article_date || '—' }}</TableCell>
            <TableCell>
              <Badge :variant="p.has_file ? 'success' : 'secondary'">
                {{ p.has_file ? t('admin.publications.file_uploaded') : t('admin.publications.no_file') }}
              </Badge>
            </TableCell>
            <TableCell>
              <Badge :variant="p.is_published ? 'success' : 'warning'">
                {{ p.is_published ? t('admin.common.published') : t('admin.common.draft') }}
              </Badge>
            </TableCell>
            <TableCell class="text-right">
              <div class="flex justify-end gap-2">
                <NuxtLink :to="`/admin/publications/${p.id}`">
                  <Button variant="ghost" size="sm">{{ t('admin.common.edit') }}</Button>
                </NuxtLink>
                <Button
                  variant="ghost"
                  size="sm"
                  class="text-destructive hover:text-destructive"
                  @click="remove(p.id, loc(p, 'name'))"
                >
                  {{ t('admin.common.delete') }}
                </Button>
              </div>
            </TableCell>
          </TableRow>
          <TableEmpty v-if="items.length === 0" :colspan="7">{{ t('admin.publications.empty_hint') }}</TableEmpty>
        </TableBody>
      </Table>
    </Card>

    <Pagination :page="page" :total="total" :limit="LIMIT" @change="(p) => (page = p)" />
  </div>
</template>
