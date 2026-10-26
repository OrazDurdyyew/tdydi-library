<script setup lang="ts">
import { Plus } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

const { t } = useI18n()

interface Collection {
  id: number
  slug: string
  name_ru: string
  name_tm: string
  show_on_home: boolean
}

const api = useApi()
const toast = useToast()
const { ask } = useConfirm()

const { data, refresh, pending } = await useAsyncData('admin-collections', () =>
  api.get<{ data: Collection[] }>('/collections?limit=100')
)

const items = computed(() => data.value?.data ?? [])

async function remove(id: number, name: string) {
  const ok = await ask({
    title: t('admin.collections.confirm_delete_title'),
    description: t('admin.collections.confirm_delete_description', { name }),
    confirmText: t('admin.common.delete'),
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/collections/${id}`)
    await refresh()
    toast.success(t('admin.collections.deleted'))
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error_delete'))
  }
}

useHead({ title: () => t('admin.collections.title') })
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold font-serif">{{ t('admin.collections.title') }}</h1>
      <NuxtLink to="/admin/collections/new">
        <Button><Plus class="w-4 h-4" /> {{ t('admin.collections.add_new') }}</Button>
      </NuxtLink>
    </div>

    <div v-if="pending" class="text-muted-foreground">{{ t('admin.common.loading') }}</div>

    <Card v-else class="p-0 overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>{{ t('admin.common.title') }}</TableHead>
            <TableHead>{{ t('admin.common.slug') }}</TableHead>
            <TableHead>{{ t('admin.collections.show_on_home') }}</TableHead>
            <TableHead class="text-right">{{ t('admin.common.actions') }}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="c in items" :key="c.id">
            <TableCell class="font-medium">{{ c.name_ru || c.name_tm }}</TableCell>
            <TableCell class="font-mono text-xs text-muted-foreground">{{ c.slug }}</TableCell>
            <TableCell>
              <Badge :variant="c.show_on_home ? 'success' : 'warning'">
                {{ c.show_on_home ? t('admin.collections.on_home') : t('admin.collections.hidden') }}
              </Badge>
            </TableCell>
            <TableCell class="text-right">
              <div class="flex justify-end gap-2">
                <NuxtLink :to="`/admin/collections/${c.id}`">
                  <Button variant="ghost" size="sm">{{ t('admin.common.edit') }}</Button>
                </NuxtLink>
                <Button
                  variant="ghost"
                  size="sm"
                  class="text-destructive hover:text-destructive"
                  @click="remove(c.id, c.name_ru || c.name_tm)"
                >
                  {{ t('admin.common.delete') }}
                </Button>
              </div>
            </TableCell>
          </TableRow>
          <TableEmpty v-if="items.length === 0" :colspan="4">{{ t('admin.common.no_data') }}</TableEmpty>
        </TableBody>
      </Table>
    </Card>
  </div>
</template>
