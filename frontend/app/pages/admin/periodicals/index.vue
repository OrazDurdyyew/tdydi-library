<script setup lang="ts">
import { Plus } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

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

const { data, refresh, pending } = await useAsyncData('admin-periodicals', () =>
  api.get<{ data: Periodical[] }>('/periodicals?limit=100')
)

const items = computed(() => data.value?.data ?? [])

const typeLabel = (t: string) => (t === 'newspaper' ? 'Газета' : 'Журнал')

async function remove(id: number, name: string) {
  const ok = await ask({
    title: 'Удалить издание?',
    description: `«${name}» и все его выпуски будут удалены.`,
    confirmText: 'Удалить',
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/periodicals/${id}`)
    await refresh()
    toast.success('Издание удалено')
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка удаления')
  }
}

useHead({ title: 'Периодика' })
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold font-serif">Периодические издания</h1>
      <NuxtLink to="/admin/periodicals/new">
        <Button><Plus class="w-4 h-4" /> Добавить издание</Button>
      </NuxtLink>
    </div>

    <div v-if="pending" class="text-muted-foreground">Загрузка…</div>

    <Card v-else class="p-0 overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Название</TableHead>
            <TableHead>Тип</TableHead>
            <TableHead>Периодичность</TableHead>
            <TableHead>Slug</TableHead>
            <TableHead class="text-right">Действия</TableHead>
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
                  <Button variant="ghost" size="sm">Управлять</Button>
                </NuxtLink>
                <Button
                  variant="ghost"
                  size="sm"
                  class="text-destructive hover:text-destructive"
                  @click="remove(p.id, p.name_ru || p.name_tm)"
                >
                  Удалить
                </Button>
              </div>
            </TableCell>
          </TableRow>
          <TableEmpty v-if="items.length === 0" :colspan="5">Изданий пока нет.</TableEmpty>
        </TableBody>
      </Table>
    </Card>
  </div>
</template>
