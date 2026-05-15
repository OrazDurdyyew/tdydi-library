<script setup lang="ts">
import { Plus } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

interface FAQItem {
  id: number
  question_ru: string
  question_tm: string
  answer_ru: string
  answer_tm: string
  sort_order: number
}

interface FAQCategory {
  id: number
  name_ru: string
  name_tm: string
  sort_order: number
  items?: FAQItem[]
}

const api = useApi()
const toast = useToast()
const { ask } = useConfirm()

const { data, refresh } = await useAsyncData('admin-faq', () =>
  api.get<FAQCategory[]>('/faq')
)

const categories = computed(() => data.value ?? [])

const newCat = reactive({ name_tm: '', name_ru: '', name_en: '', sort_order: 0 })
const addingCat = ref(false)

async function createCategory() {
  if (!newCat.name_tm) {
    toast.error('Название (TM) обязательно')
    return
  }
  addingCat.value = true
  try {
    await api.post('/admin/faq/categories', newCat)
    await refresh()
    toast.success('Категория добавлена')
    newCat.name_tm = ''
    newCat.name_ru = ''
    newCat.name_en = ''
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка')
  } finally {
    addingCat.value = false
  }
}

async function deleteCategory(id: number, name: string) {
  const ok = await ask({
    title: 'Удалить категорию?',
    description: `Категория «${name}» и все вопросы в ней будут удалены.`,
    confirmText: 'Удалить',
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/faq/categories/${id}`)
    await refresh()
    toast.success('Категория удалена')
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка')
  }
}

const newItems = reactive<Record<number, { question_tm: string; question_ru: string; answer_tm: string; answer_ru: string }>>({})

function getNewItem(catId: number) {
  if (!newItems[catId]) {
    newItems[catId] = { question_tm: '', question_ru: '', answer_tm: '', answer_ru: '' }
  }
  return newItems[catId]
}

const addingItem = ref<number | null>(null)

async function createItem(catId: number) {
  const item = newItems[catId]
  if (!item?.question_tm) {
    toast.error('Вопрос (TM) обязателен')
    return
  }
  addingItem.value = catId
  try {
    await api.post('/admin/faq/items', { category_id: catId, ...item })
    await refresh()
    toast.success('Вопрос добавлен')
    newItems[catId] = { question_tm: '', question_ru: '', answer_tm: '', answer_ru: '' }
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка')
  } finally {
    addingItem.value = null
  }
}

async function deleteItem(id: number, q: string) {
  const ok = await ask({
    title: 'Удалить вопрос?',
    description: `«${q}» будет удалён.`,
    confirmText: 'Удалить',
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/faq/items/${id}`)
    await refresh()
    toast.success('Вопрос удалён')
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка')
  }
}

useHead({ title: 'FAQ' })
</script>

<template>
  <div>
    <h1 class="text-3xl font-bold font-serif mb-8">FAQ</h1>

    <Card class="mb-8 max-w-2xl">
      <CardHeader>
        <CardTitle class="text-lg">Добавить категорию</CardTitle>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-2">
            <Label for="cat-tm">Название TM <span class="text-destructive">*</span></Label>
            <Input id="cat-tm" v-model="newCat.name_tm" />
          </div>
          <div class="space-y-2">
            <Label for="cat-ru">Название RU</Label>
            <Input id="cat-ru" v-model="newCat.name_ru" />
          </div>
        </div>
        <Button @click="createCategory" :disabled="addingCat" size="sm">
          <Plus class="w-4 h-4" /> {{ addingCat ? 'Добавление…' : 'Добавить категорию' }}
        </Button>
      </CardContent>
    </Card>

    <div class="space-y-6 max-w-3xl">
      <Card v-for="cat in categories" :key="cat.id" class="p-0 overflow-hidden">
        <div class="px-6 py-4 bg-muted/30 flex items-center justify-between border-b border-border">
          <h2 class="font-semibold font-serif">{{ cat.name_ru || cat.name_tm }}</h2>
          <Button
            variant="ghost"
            size="sm"
            class="text-destructive hover:text-destructive"
            @click="deleteCategory(cat.id, cat.name_ru || cat.name_tm)"
          >
            Удалить категорию
          </Button>
        </div>

        <div v-if="cat.items && cat.items.length > 0">
          <div
            v-for="item in cat.items"
            :key="item.id"
            class="px-6 py-4 border-b border-border hover:bg-muted/20"
          >
            <div class="flex justify-between items-start gap-4">
              <div class="flex-1 min-w-0">
                <div class="font-medium text-sm">{{ item.question_ru || item.question_tm }}</div>
                <div class="text-sm text-muted-foreground mt-1 line-clamp-2">
                  {{ item.answer_ru || item.answer_tm }}
                </div>
              </div>
              <Button
                variant="ghost"
                size="sm"
                class="text-destructive hover:text-destructive shrink-0"
                @click="deleteItem(item.id, item.question_ru || item.question_tm)"
              >
                Удалить
              </Button>
            </div>
          </div>
        </div>
        <div v-else class="px-6 py-4 text-sm text-muted-foreground">Вопросов пока нет</div>

        <div class="px-6 py-4 bg-muted/10 border-t border-border space-y-3">
          <h4 class="text-sm font-medium">Добавить вопрос</h4>
          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-2">
              <Label :for="`q-tm-${cat.id}`" class="text-xs">Вопрос TM <span class="text-destructive">*</span></Label>
              <Input :id="`q-tm-${cat.id}`" v-model="getNewItem(cat.id).question_tm" class="text-sm" />
            </div>
            <div class="space-y-2">
              <Label :for="`q-ru-${cat.id}`" class="text-xs">Вопрос RU</Label>
              <Input :id="`q-ru-${cat.id}`" v-model="getNewItem(cat.id).question_ru" class="text-sm" />
            </div>
            <div class="space-y-2">
              <Label :for="`a-tm-${cat.id}`" class="text-xs">Ответ TM</Label>
              <Textarea :id="`a-tm-${cat.id}`" v-model="getNewItem(cat.id).answer_tm" rows="2" class="text-sm" />
            </div>
            <div class="space-y-2">
              <Label :for="`a-ru-${cat.id}`" class="text-xs">Ответ RU</Label>
              <Textarea :id="`a-ru-${cat.id}`" v-model="getNewItem(cat.id).answer_ru" rows="2" class="text-sm" />
            </div>
          </div>
          <Button
            @click="createItem(cat.id)"
            :disabled="addingItem === cat.id"
            variant="outline"
            size="sm"
          >
            <Plus class="w-4 h-4" /> {{ addingItem === cat.id ? 'Добавление…' : 'Добавить вопрос' }}
          </Button>
        </div>
      </Card>

      <div v-if="categories.length === 0" class="text-center py-16 text-muted-foreground">
        Категорий пока нет. Создайте первую выше.
      </div>
    </div>
  </div>
</template>
