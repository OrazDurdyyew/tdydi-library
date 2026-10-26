<script setup lang="ts">
import { Plus } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

const { t } = useI18n()

interface Page {
  id: number
  slug: string
  title_ru: string
  title_tm: string
  in_menu: boolean
}

const api = useApi()
const toast = useToast()
const { ask } = useConfirm()

// /pages возвращает массив (не {data}).
const { data, refresh, pending } = await useAsyncData('admin-pages', () =>
  api.get<Page[]>('/pages?limit=100').catch(() => [])
)

const pages = computed(() => data.value ?? [])

async function remove(id: number, title: string) {
  const ok = await ask({
    title: t('admin.pages.confirm_delete_title'),
    description: t('admin.pages.confirm_delete_description', { title }),
    confirmText: t('admin.common.delete'),
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/pages/${id}`)
    await refresh()
    toast.success(t('admin.pages.deleted'))
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error_delete'))
  }
}

useHead({ title: () => t('admin.pages.title') })
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold font-serif">{{ t('admin.pages.title') }}</h1>
      <NuxtLink to="/admin/pages/new">
        <Button><Plus class="w-4 h-4" /> {{ t('admin.pages.add_new') }}</Button>
      </NuxtLink>
    </div>

    <div v-if="pending" class="text-muted-foreground">{{ t('admin.common.loading') }}</div>

    <Card v-else class="p-0 overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>{{ t('admin.common.title') }}</TableHead>
            <TableHead>{{ t('admin.common.slug') }}</TableHead>
            <TableHead>{{ t('admin.pages.in_menu') }}</TableHead>
            <TableHead class="text-right">{{ t('admin.common.actions') }}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="p in pages" :key="p.id">
            <TableCell class="font-medium">{{ p.title_ru || p.title_tm }}</TableCell>
            <TableCell class="font-mono text-xs text-muted-foreground">{{ p.slug }}</TableCell>
            <TableCell>
              <Badge :variant="p.in_menu ? 'success' : 'warning'">
                {{ p.in_menu ? t('admin.pages.in_menu') : t('admin.pages.hidden') }}
              </Badge>
            </TableCell>
            <TableCell class="text-right">
              <div class="flex justify-end gap-2">
                <NuxtLink :to="`/admin/pages/${p.id}`">
                  <Button variant="ghost" size="sm">{{ t('admin.common.edit') }}</Button>
                </NuxtLink>
                <Button
                  variant="ghost"
                  size="sm"
                  class="text-destructive hover:text-destructive"
                  @click="remove(p.id, p.title_ru || p.title_tm)"
                >
                  {{ t('admin.common.delete') }}
                </Button>
              </div>
            </TableCell>
          </TableRow>
          <TableEmpty v-if="pages.length === 0" :colspan="4">{{ t('admin.common.no_data') }}</TableEmpty>
        </TableBody>
      </Table>
    </Card>
  </div>
</template>
