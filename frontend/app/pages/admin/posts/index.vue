<script setup lang="ts">
import { Plus, Image as ImageIcon } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

const { t } = useI18n()

interface Post {
  id: number
  slug: string
  title_ru: string
  status: string
  created_at: string
  cover_path?: string
}

const api = useApi()
const toast = useToast()
const { ask } = useConfirm()
const fileUrl = useFileUrl()

const LIMIT = 20
const page = ref(1)
const { data, refresh, pending } = await useAsyncData(
  'admin-posts',
  () => api.get<{ data: Post[]; total: number }>(`/posts?page=${page.value}&limit=${LIMIT}`),
  { watch: [page] }
)

const posts = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.total ?? 0)

async function remove(id: number, title: string) {
  const ok = await ask({
    title: t('admin.posts.confirm_delete_title'),
    description: t('admin.posts.confirm_delete_description', { title }),
    confirmText: t('admin.common.delete'),
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/posts/${id}`)
    await refresh()
    toast.success(t('admin.posts.deleted'))
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error_delete'))
  }
}

function formatDate(d: string) {
  return new Date(d).toLocaleDateString('ru-RU')
}

useHead({ title: () => t('admin.posts.title') })
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold font-serif">{{ t('admin.posts.title') }}</h1>
      <NuxtLink to="/admin/posts/new">
        <Button><Plus class="w-4 h-4" /> {{ t('admin.posts.add_new') }}</Button>
      </NuxtLink>
    </div>

    <div v-if="pending" class="text-muted-foreground">{{ t('admin.common.loading') }}</div>

    <Card v-else class="p-0 overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead class="w-16">{{ t('admin.common.image') }}</TableHead>
            <TableHead>{{ t('admin.common.title') }}</TableHead>
            <TableHead>{{ t('admin.common.slug') }}</TableHead>
            <TableHead>{{ t('admin.posts.date') }}</TableHead>
            <TableHead>{{ t('admin.posts.status') }}</TableHead>
            <TableHead class="text-right">{{ t('admin.common.actions') }}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="p in posts" :key="p.id">
            <TableCell>
              <div class="w-12 h-12 rounded-md overflow-hidden bg-muted border border-border flex items-center justify-center">
                <img
                  v-if="p.cover_path"
                  :src="fileUrl(p.cover_path)"
                  :alt="p.title_ru || p.slug"
                  class="w-full h-full object-cover"
                />
                <ImageIcon v-else class="w-4 h-4 text-muted-foreground" />
              </div>
            </TableCell>
            <TableCell class="font-medium">{{ p.title_ru || '—' }}</TableCell>
            <TableCell class="font-mono text-xs text-muted-foreground">{{ p.slug }}</TableCell>
            <TableCell class="text-muted-foreground">{{ formatDate(p.created_at) }}</TableCell>
            <TableCell>
              <Badge :variant="p.status === 'published' ? 'success' : 'warning'">
                {{ p.status === 'published' ? t('admin.common.published') : t('admin.common.draft') }}
              </Badge>
            </TableCell>
            <TableCell class="text-right">
              <div class="flex justify-end gap-2">
                <NuxtLink :to="`/admin/posts/${p.id}`">
                  <Button variant="ghost" size="sm">{{ t('admin.common.edit') }}</Button>
                </NuxtLink>
                <Button
                  variant="ghost"
                  size="sm"
                  class="text-destructive hover:text-destructive"
                  @click="remove(p.id, p.title_ru || p.slug)"
                >
                  {{ t('admin.common.delete') }}
                </Button>
              </div>
            </TableCell>
          </TableRow>
          <TableEmpty v-if="posts.length === 0" :colspan="6">{{ t('admin.common.no_data') }}</TableEmpty>
        </TableBody>
      </Table>
    </Card>

    <Pagination :page="page" :total="total" :limit="LIMIT" @change="(p) => (page = p)" />
  </div>
</template>
