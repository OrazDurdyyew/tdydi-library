<script setup lang="ts">
import { Plus } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

const { t } = useI18n()

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
const LIMIT = 20
const page = ref(1)
const { data, refresh, pending } = await useAsyncData(
  'admin-books',
  () => api.get<{ data: Book[]; meta: { total: number } }>(`/books?page=${page.value}&limit=${LIMIT}`),
  { watch: [page] }
)

const books = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.meta?.total ?? 0)

async function remove(id: number, title: string) {
  const ok = await ask({
    title: t('admin.books.confirm_delete_title'),
    description: t('admin.books.confirm_delete_description', { title }),
    confirmText: t('admin.common.delete'),
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/books/${id}`)
    await refresh()
    toast.success(t('admin.books.deleted'))
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error'))
  }
}

useHead({ title: () => t('admin.books.title') })
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold font-serif">{{ t('admin.books.title') }}</h1>
      <NuxtLink to="/admin/books/new">
        <Button><Plus class="w-4 h-4" /> {{ t('admin.books.add_new') }}</Button>
      </NuxtLink>
    </div>

    <div v-if="pending" class="text-muted-foreground">{{ t('admin.common.loading') }}</div>

    <Card v-else class="p-0 overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>{{ t('admin.common.title') }}</TableHead>
            <TableHead>{{ t('admin.books.authors') }}</TableHead>
            <TableHead>{{ t('admin.books.year') }}</TableHead>
            <TableHead>{{ t('admin.books.language') }}</TableHead>
            <TableHead>{{ t('admin.common.status') }}</TableHead>
            <TableHead class="text-right">{{ t('admin.common.actions') }}</TableHead>
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
                {{ b.is_published ? t('admin.common.published') : t('admin.common.draft') }}
              </Badge>
            </TableCell>
            <TableCell class="text-right">
              <div class="flex justify-end gap-2">
                <NuxtLink :to="`/admin/books/${b.id}`">
                  <Button variant="ghost" size="sm">{{ t('admin.common.edit') }}</Button>
                </NuxtLink>
                <Button variant="ghost" size="sm" class="text-destructive hover:text-destructive" @click="remove(b.id, b.title)">
                  {{ t('admin.common.delete') }}
                </Button>
              </div>
            </TableCell>
          </TableRow>
          <TableEmpty v-if="books.length === 0" :colspan="7">
            {{ t('admin.books.empty_hint') }}
          </TableEmpty>
        </TableBody>
      </Table>
    </Card>

    <Pagination :page="page" :total="total" :limit="LIMIT" @change="(p) => (page = p)" />
  </div>
</template>
