<script setup lang="ts">
import { Plus } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

interface Post {
  id: number
  slug: string
  title_ru: string
  status: string
  created_at: string
}

const api = useApi()
const toast = useToast()
const { ask } = useConfirm()

const { data, refresh, pending } = await useAsyncData('admin-posts', () =>
  api.get<{ data: Post[]; meta: { total: number } }>('/posts?limit=100')
)

const posts = computed(() => data.value?.data ?? [])

async function remove(id: number, title: string) {
  const ok = await ask({
    title: 'Удалить новость?',
    description: `«${title}» будет удалена.`,
    confirmText: 'Удалить',
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/posts/${id}`)
    await refresh()
    toast.success('Новость удалена')
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка удаления')
  }
}

function formatDate(d: string) {
  return new Date(d).toLocaleDateString('ru-RU')
}

useHead({ title: 'Новости' })
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold font-serif">Новости</h1>
      <NuxtLink to="/admin/posts/new">
        <Button><Plus class="w-4 h-4" /> Добавить новость</Button>
      </NuxtLink>
    </div>

    <div v-if="pending" class="text-muted-foreground">Загрузка…</div>

    <Card v-else class="p-0 overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Заголовок</TableHead>
            <TableHead>Slug</TableHead>
            <TableHead>Дата</TableHead>
            <TableHead>Статус</TableHead>
            <TableHead class="text-right">Действия</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="p in posts" :key="p.id">
            <TableCell class="font-medium">{{ p.title_ru || '—' }}</TableCell>
            <TableCell class="font-mono text-xs text-muted-foreground">{{ p.slug }}</TableCell>
            <TableCell class="text-muted-foreground">{{ formatDate(p.created_at) }}</TableCell>
            <TableCell>
              <Badge :variant="p.status === 'published' ? 'success' : 'warning'">
                {{ p.status === 'published' ? 'Опубликована' : 'Черновик' }}
              </Badge>
            </TableCell>
            <TableCell class="text-right">
              <div class="flex justify-end gap-2">
                <NuxtLink :to="`/admin/posts/${p.id}`">
                  <Button variant="ghost" size="sm">Изменить</Button>
                </NuxtLink>
                <Button
                  variant="ghost"
                  size="sm"
                  class="text-destructive hover:text-destructive"
                  @click="remove(p.id, p.title_ru || p.slug)"
                >
                  Удалить
                </Button>
              </div>
            </TableCell>
          </TableRow>
          <TableEmpty v-if="posts.length === 0" :colspan="5">Новостей пока нет.</TableEmpty>
        </TableBody>
      </Table>
    </Card>
  </div>
</template>
