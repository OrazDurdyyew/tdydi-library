<script setup lang="ts">
import { Plus } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

const { t } = useI18n()

interface Periodical {
  id: number
  slug: string
  type: 'newspaper' | 'magazine'
  name_ru: string
  name_tm: string
  periodicity: string
}

const api = useApi()
const toast = useToast()
const { ask } = useConfirm()

const LIMIT = 20
const page = ref(1)
const { data, refresh, pending } = await useAsyncData(
  'admin-periodicals',
  () =>
    api.get<{ data: Periodical[]; total: number }>(
      `/periodicals?page=${page.value}&limit=${LIMIT}`
    ),
  { watch: [page] }
)

const items = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.total ?? 0)

const typeLabel = (type: string) => (type === 'newspaper' ? t('admin.periodicals.type_newspaper') : t('admin.periodicals.type_magazine'))

async function remove(id: number, name: string) {
  const ok = await ask({
    title: t('admin.periodicals.confirm_delete_title'),
    description: t('admin.periodicals.confirm_delete_description', { name }),
    confirmText: t('admin.common.delete'),
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/periodicals/${id}`)
    await refresh()
    toast.success(t('admin.periodicals.deleted'))
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error_delete'))
  }
}

useHead({ title: () => t('admin.periodicals.title') })
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold font-serif">{{ t('admin.periodicals.title') }}</h1>
      <NuxtLink to="/admin/periodicals/new">
        <Button><Plus class="w-4 h-4" /> {{ t('admin.periodicals.add_new') }}</Button>
      </NuxtLink>
    </div>

    <div v-if="pending" class="text-muted-foreground">{{ t('admin.common.loading') }}</div>

    <Card v-else class="p-0 overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>{{ t('admin.common.title') }}</TableHead>
            <TableHead>{{ t('admin.periodicals.type') }}</TableHead>
            <TableHead>{{ t('admin.periodicals.periodicity') }}</TableHead>
            <TableHead>{{ t('admin.common.slug') }}</TableHead>
            <TableHead class="text-right">{{ t('admin.common.actions') }}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="p in items" :key="p.id">
            <TableCell class="font-medium">{{ p.name_ru || p.name_tm }}</TableCell>
            <TableCell>
              <Badge :variant="p.type === 'newspaper' ? 'info' : 'secondary'">
                {{ typeLabel(p.type) }}
              </Badge>
            </TableCell>
            <TableCell class="text-muted-foreground">{{ p.periodicity || '—' }}</TableCell>
            <TableCell class="font-mono text-xs text-muted-foreground">{{ p.slug }}</TableCell>
            <TableCell class="text-right">
              <div class="flex justify-end gap-2">
                <NuxtLink :to="`/admin/periodicals/${p.id}`">
                  <Button variant="ghost" size="sm">{{ t('admin.common.edit') }}</Button>
                </NuxtLink>
                <Button
                  variant="ghost"
                  size="sm"
                  class="text-destructive hover:text-destructive"
                  @click="remove(p.id, p.name_ru || p.name_tm)"
                >
                  {{ t('admin.common.delete') }}
                </Button>
              </div>
            </TableCell>
          </TableRow>
          <TableEmpty v-if="items.length === 0" :colspan="5">{{ t('admin.common.no_data') }}</TableEmpty>
        </TableBody>
      </Table>
    </Card>

    <Pagination :page="page" :total="total" :limit="LIMIT" @change="(p) => (page = p)" />
  </div>
</template>
