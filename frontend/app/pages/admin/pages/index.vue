<script setup lang="ts">
import { Plus } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

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
    title: 'Удалить страницу?',
    description: `«${title}» будет удалена.`,
    confirmText: 'Удалить',
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/pages/${id}`)
    await refresh()
    toast.success('Страница удалена')
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка удаления')
  }
}

useHead({ title: 'Статичные страницы' })
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold font-serif">Статичные страницы</h1>
      <NuxtLink to="/admin/pages/new">
        <Button><Plus class="w-4 h-4" /> Создать страницу</Button>
      </NuxtLink>
    </div>

    <div v-if="pending" class="text-muted-foreground">Загрузка…</div>

    <Card v-else class="p-0 overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Заголовок</TableHead>
            <TableHead>Slug</TableHead>
            <TableHead>В меню</TableHead>
            <TableHead class="text-right">Действия</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="p in pages" :key="p.id">
            <TableCell class="font-medium">{{ p.title_ru || p.title_tm }}</TableCell>
            <TableCell class="font-mono text-xs text-muted-foreground">{{ p.slug }}</TableCell>
            <TableCell>
              <Badge :variant="p.in_menu ? 'success' : 'warning'">
                {{ p.in_menu ? 'В меню' : 'Скрыта' }}
              </Badge>
            </TableCell>
            <TableCell class="text-right">
              <div class="flex justify-end gap-2">
                <NuxtLink :to="`/admin/pages/${p.id}`">
                  <Button variant="ghost" size="sm">Изменить</Button>
                </NuxtLink>
                <Button
                  variant="ghost"
                  size="sm"
                  class="text-destructive hover:text-destructive"
                  @click="remove(p.id, p.title_ru || p.title_tm)"
                >
                  Удалить
                </Button>
              </div>
            </TableCell>
          </TableRow>
          <TableEmpty v-if="pages.length === 0" :colspan="4">Страниц пока нет.</TableEmpty>
        </TableBody>
      </Table>
    </Card>
  </div>
</template>
