<script setup lang="ts">
import { Plus, CornerDownRight, FolderTree } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

interface Section {
  id: number
  slug: string
  parent_id: number | null
  sort_order: number
  name_tm: string
  name_ru: string
  name_en: string
  book_count: number
}

const api = useApi()
const toast = useToast()
const { ask } = useConfirm()
const { data, refresh, pending } = await useAsyncData('admin-sections', () =>
  api.get<{ data: Section[] }>('/sections')
)

const sections = computed(() => data.value?.data ?? [])

// Дерево: корневые разделы по порядку, под каждым — его подразделы.
const tree = computed(() => {
  const roots = sections.value
    .filter((s) => s.parent_id == null)
    .sort((a, b) => a.sort_order - b.sort_order || a.id - b.id)
  return roots.map((root) => ({
    ...root,
    children: sections.value
      .filter((s) => s.parent_id === root.id)
      .sort((a, b) => a.sort_order - b.sort_order || a.id - b.id)
  }))
})

async function remove(id: number, name: string) {
  const ok = await ask({
    title: 'Удалить раздел?',
    description: `«${name}» будет удалён. Если в нём есть книги или подразделы — операция не выполнится.`,
    confirmText: 'Удалить',
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/sections/${id}`)
    await refresh()
    toast.success('Раздел удалён')
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка удаления')
  }
}

useHead({ title: 'Разделы' })
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-2">
      <h1 class="text-3xl font-bold font-serif">Разделы</h1>
      <NuxtLink to="/admin/sections/new">
        <Button><Plus class="w-4 h-4" /> Создать раздел</Button>
      </NuxtLink>
    </div>
    <p class="text-sm text-muted-foreground mb-8">
      Два уровня: раздел верхнего уровня → подразделы. Кнопка «Подраздел» добавляет вложенный раздел.
    </p>

    <div v-if="pending" class="text-muted-foreground">Загрузка…</div>

    <div v-else-if="tree.length === 0" class="paper-card text-center py-16 text-muted-foreground serif italic">
      Пока ни одного раздела. Создайте первый.
    </div>

    <div v-else class="space-y-4">
      <Card v-for="root in tree" :key="root.id" class="p-0 overflow-hidden">
        <!-- Корневой раздел -->
        <div class="flex items-center gap-3 px-5 py-4 bg-muted/40">
          <FolderTree class="w-5 h-5 icon-primary shrink-0" />
          <div class="flex-1 min-w-0">
            <div class="font-serif font-semibold text-lg">{{ root.name_ru }}</div>
            <div class="text-xs text-muted-foreground font-mono">
              /sections/{{ root.slug }} · книг: {{ root.book_count }} · порядок: {{ root.sort_order }}
            </div>
          </div>
          <div class="flex items-center gap-2 shrink-0">
            <NuxtLink :to="`/admin/sections/new?parent=${root.id}`">
              <Button variant="outline" size="sm">
                <Plus class="w-3.5 h-3.5" /> Подраздел
              </Button>
            </NuxtLink>
            <NuxtLink :to="`/admin/sections/${root.id}`">
              <Button variant="ghost" size="sm">Изменить</Button>
            </NuxtLink>
            <Button
              variant="ghost"
              size="sm"
              class="text-destructive hover:text-destructive"
              @click="remove(root.id, root.name_ru)"
            >
              Удалить
            </Button>
          </div>
        </div>

        <!-- Подразделы -->
        <div v-if="root.children.length" class="divide-y divide-border">
          <div
            v-for="child in root.children"
            :key="child.id"
            class="flex items-center gap-3 pl-12 pr-5 py-3"
          >
            <CornerDownRight class="w-4 h-4 text-muted-foreground shrink-0" />
            <div class="flex-1 min-w-0">
              <div class="font-medium">{{ child.name_ru }}</div>
              <div class="text-xs text-muted-foreground font-mono">
                /sections/{{ child.slug }} · книг: {{ child.book_count }} · порядок: {{ child.sort_order }}
              </div>
            </div>
            <div class="flex items-center gap-2 shrink-0">
              <NuxtLink :to="`/admin/sections/${child.id}`">
                <Button variant="ghost" size="sm">Изменить</Button>
              </NuxtLink>
              <Button
                variant="ghost"
                size="sm"
                class="text-destructive hover:text-destructive"
                @click="remove(child.id, child.name_ru)"
              >
                Удалить
              </Button>
            </div>
          </div>
        </div>
        <div v-else class="pl-12 pr-5 py-3 text-sm text-muted-foreground italic">
          Подразделов нет.
        </div>
      </Card>
    </div>
  </div>
</template>
