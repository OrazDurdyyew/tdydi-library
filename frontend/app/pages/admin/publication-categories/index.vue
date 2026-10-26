<script setup lang="ts">
import { Plus, ArrowLeft } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const loc = useLocalized()

interface PublicationCategory {
  id: number
  slug: string
  name_tm: string
  name_ru: string
  name_en: string
}

const api = useApi()
const toast = useToast()
const { ask } = useConfirm()

const { data, refresh, pending } = await useAsyncData('admin-publication-categories', () =>
  api.get<{ data: PublicationCategory[] }>('/publication-categories?limit=200')
)

const items = computed(() => data.value?.data ?? [])

async function remove(id: number, name: string) {
  const ok = await ask({
    title: t('admin.publication_categories.confirm_delete_title'),
    description: t('admin.publication_categories.confirm_delete_description', { name }),
    confirmText: t('admin.common.delete'),
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/publication-categories/${id}`)
    await refresh()
    toast.success(t('admin.publication_categories.deleted'))
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error_delete'))
  }
}

useHead({ title: () => t('admin.publication_categories.title') })
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
      <h1 class="text-3xl font-bold font-serif">{{ t('admin.publication_categories.title') }}</h1>
      <NuxtLink to="/admin/publication-categories/new">
        <Button><Plus class="w-4 h-4" /> {{ t('admin.publication_categories.add_new') }}</Button>
      </NuxtLink>
    </div>

    <div v-if="pending" class="text-muted-foreground">{{ t('admin.common.loading') }}</div>

    <Card v-else class="p-0 overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>{{ t('admin.common.title') }}</TableHead>
            <TableHead>{{ t('admin.common.slug') }}</TableHead>
            <TableHead class="text-right">{{ t('admin.common.actions') }}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="c in items" :key="c.id">
            <TableCell class="font-medium">{{ loc(c, 'name') }}</TableCell>
            <TableCell class="font-mono text-xs text-muted-foreground">{{ c.slug }}</TableCell>
            <TableCell class="text-right">
              <div class="flex justify-end gap-2">
                <NuxtLink :to="`/admin/publication-categories/${c.id}`">
                  <Button variant="ghost" size="sm">{{ t('admin.common.edit') }}</Button>
                </NuxtLink>
                <Button
                  variant="ghost"
                  size="sm"
                  class="text-destructive hover:text-destructive"
                  @click="remove(c.id, loc(c, 'name'))"
                >
                  {{ t('admin.common.delete') }}
                </Button>
              </div>
            </TableCell>
          </TableRow>
          <TableEmpty v-if="items.length === 0" :colspan="3">{{ t('admin.publication_categories.empty_hint') }}</TableEmpty>
        </TableBody>
      </Table>
    </Card>
  </div>
</template>
