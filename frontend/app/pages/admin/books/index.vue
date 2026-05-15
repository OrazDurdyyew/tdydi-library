<script setup lang="ts">
import { Plus } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

interface Book {
  id: number
  slug: string
  title: string
  authors: string
  year: number
  language: string
  is_published: boolean
  sections?: { id: number; name_ru: string }[]
}

const api = useApi()
const toast = useToast()
const { ask } = useConfirm()
const { data, refresh, pending } = await useAsyncData('admin-books', () =>
  api.get<{ data: Book[]; meta: { total: number } }>('/books?limit=100')
)

const books = computed(() => data.value?.data ?? [])

async function remove(id: number, title: string) {
  const ok = await ask({
    title: 'Удалить книгу?',
    description: `«${title}» будет удалена безвозвратно.`,
    confirmText: 'Удалить',
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/books/${id}`)
    await refresh()
    toast.success('Книга удалена')
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка')
  }
}

useHead({ title: 'Книги' })
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold font-serif">Книги</h1>
      <NuxtLink to="/admin/books/new">
        <Button><Plus class="w-4 h-4" /> Добавить книгу</Button>
      </NuxtLink>
    </div>

    <div v-if="pending" class="text-muted-foreground">Загрузка…</div>

    <Card v-else class="p-0 overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>Название</TableHead>
            <TableHead>Автор</TableHead>
            <TableHead>Год</TableHead>
            <TableHead>Язык</TableHead>
            <TableHead>Статус</TableHead>
            <TableHead class="text-right">Действия</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="b in books" :key="b.id">
            <TableCell class="text-muted-foreground">{{ b.id }}</TableCell>
            <TableCell class="font-medium">{{ b.title }}</TableCell>
            <TableCell class="text-muted-foreground">{{ b.authors || '—' }}</TableCell>
            <TableCell>{{ b.year || '—' }}</TableCell>
            <TableCell class="uppercase text-xs">{{ b.language }}</TableCell>
            <TableCell>
              <Badge :variant="b.is_published ? 'success' : 'warning'">
                {{ b.is_published ? 'Опубликована' : 'Черновик' }}
              </Badge>
            </TableCell>
            <TableCell class="text-right">
              <div class="flex justify-end gap-2">
                <NuxtLink :to="`/admin/books/${b.id}`">
                  <Button variant="ghost" size="sm">Изменить</Button>
                </NuxtLink>
                <Button variant="ghost" size="sm" class="text-destructive hover:text-destructive" @click="remove(b.id, b.title)">
                  Удалить
                </Button>
              </div>
            </TableCell>
          </TableRow>
          <TableEmpty v-if="books.length === 0" :colspan="7">
            Книг пока нет. Добавьте первую.
          </TableEmpty>
        </TableBody>
      </Table>
    </Card>
  </div>
</template>
