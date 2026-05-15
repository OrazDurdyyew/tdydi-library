<script setup lang="ts">
import { Plus } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

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
    title: 'Удалить подборку?',
    description: `«${name}» будет удалена.`,
    confirmText: 'Удалить',
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/collections/${id}`)
    await refresh()
    toast.success('Подборка удалена')
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка удаления')
  }
}

useHead({ title: 'Подборки' })
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold font-serif">Подборки книг</h1>
      <NuxtLink to="/admin/collections/new">
        <Button><Plus class="w-4 h-4" /> Создать подборку</Button>
      </NuxtLink>
    </div>

    <div v-if="pending" class="text-muted-foreground">Загрузка…</div>

    <Card v-else class="p-0 overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Название</TableHead>
            <TableHead>Slug</TableHead>
            <TableHead>На главной</TableHead>
            <TableHead class="text-right">Действия</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="c in items" :key="c.id">
            <TableCell class="font-medium">{{ c.name_ru || c.name_tm }}</TableCell>
            <TableCell class="font-mono text-xs text-muted-foreground">{{ c.slug }}</TableCell>
            <TableCell>
              <Badge :variant="c.show_on_home ? 'success' : 'warning'">
                {{ c.show_on_home ? 'На главной' : 'Скрыта' }}
              </Badge>
            </TableCell>
            <TableCell class="text-right">
              <div class="flex justify-end gap-2">
                <NuxtLink :to="`/admin/collections/${c.id}`">
                  <Button variant="ghost" size="sm">Управлять</Button>
                </NuxtLink>
                <Button
                  variant="ghost"
                  size="sm"
                  class="text-destructive hover:text-destructive"
                  @click="remove(c.id, c.name_ru || c.name_tm)"
                >
                  Удалить
                </Button>
              </div>
            </TableCell>
          </TableRow>
          <TableEmpty v-if="items.length === 0" :colspan="4">Подборок пока нет.</TableEmpty>
        </TableBody>
      </Table>
    </Card>
  </div>
</template>
