<script setup lang="ts">
import { Plus } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

const route = useRoute()
const api = useApi()
const toast = useToast()
const { ask } = useConfirm()
const id = Number(route.params.id)

const { data: col, refresh } = await useAsyncData(`col-${id}`, () =>
  api.get<any>(`/collections/${id}`)
)

const form = reactive({
  slug: col.value?.slug ?? '',
  name_tm: col.value?.name_tm ?? '',
  name_ru: col.value?.name_ru ?? '',
  name_en: col.value?.name_en ?? '',
  description_tm: col.value?.description_tm ?? '',
  description_ru: col.value?.description_ru ?? '',
  description_en: col.value?.description_en ?? '',
  show_on_home: col.value?.show_on_home ?? false,
  sort_order: col.value?.sort_order ?? 0,
})

const saving = ref(false)
const coverInput = ref<HTMLInputElement | null>(null)
const uploading = ref(false)

const books = computed(() => col.value?.books ?? [])

const addBookId = ref('')
const addingBook = ref(false)

async function submit() {
  saving.value = true
  try {
    await api.put(`/admin/collections/${id}`, form)
    await refresh()
    toast.success('Сохранено')
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка')
  } finally {
    saving.value = false
  }
}

async function uploadCover() {
  const file = coverInput.value?.files?.[0]
  if (!file) return
  uploading.value = true
  try {
    await api.upload(`/admin/collections/${id}/cover`, file)
    await refresh()
    toast.success('Обложка загружена')
    if (coverInput.value) coverInput.value.value = ''
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка загрузки')
  } finally {
    uploading.value = false
  }
}

async function addBook() {
  if (!addBookId.value) return
  addingBook.value = true
  try {
    await api.post(`/admin/collections/${id}/books`, { book_id: Number(addBookId.value) })
    await refresh()
    toast.success('Книга добавлена')
    addBookId.value = ''
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка')
  } finally {
    addingBook.value = false
  }
}

async function removeBook(bookId: number, title: string) {
  const ok = await ask({
    title: 'Убрать книгу из подборки?',
    description: `«${title}» будет удалена из этой подборки.`,
    confirmText: 'Убрать',
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/collections/${id}/books/${bookId}`)
    await refresh()
    toast.success('Книга удалена из подборки')
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка')
  }
}

useHead({ title: `Подборка: ${col.value?.name_ru || col.value?.name_tm}` })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/collections" class="text-muted-foreground hover:text-foreground text-sm">← Подборки</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">{{ col?.name_ru || col?.name_tm }}</h1>
    </div>

    <Card class="mb-6 max-w-3xl">
      <CardHeader>
        <CardTitle class="text-lg">Обложка</CardTitle>
        <CardDescription v-if="col?.cover_path" class="font-mono">{{ col.cover_path }}</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="flex items-center gap-3 flex-wrap">
          <input
            ref="coverInput"
            type="file"
            accept="image/*"
            class="text-sm flex-1 min-w-[200px] file:mr-3 file:rounded-md file:border-0 file:bg-muted file:px-3 file:py-1.5 file:text-sm file:font-medium hover:file:bg-muted/80"
          />
          <Button @click="uploadCover" :disabled="uploading" size="sm">
            {{ uploading ? 'Загрузка…' : 'Загрузить' }}
          </Button>
        </div>
      </CardContent>
    </Card>

    <form @submit.prevent="submit" class="space-y-6 max-w-3xl mb-10">
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Основное</CardTitle>
        </CardHeader>
        <CardContent class="space-y-5">
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <Label for="slug">Slug</Label>
              <Input id="slug" v-model="form.slug" class="font-mono text-sm" />
            </div>
            <div class="space-y-2">
              <Label for="order">Порядок</Label>
              <Input id="order" v-model.number="form.sort_order" type="number" />
            </div>
          </div>
          <label class="flex items-center gap-3 cursor-pointer">
            <Checkbox v-model="form.show_on_home" />
            <span class="text-sm">Показывать на главной</span>
          </label>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Названия</CardTitle>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="name-tm">TM</Label>
            <Input id="name-tm" v-model="form.name_tm" />
          </div>
          <div class="space-y-2">
            <Label for="name-ru">RU</Label>
            <Input id="name-ru" v-model="form.name_ru" />
          </div>
          <div class="space-y-2">
            <Label for="name-en">EN</Label>
            <Input id="name-en" v-model="form.name_en" />
          </div>
        </CardContent>
      </Card>

      <div class="flex gap-3">
        <Button type="submit" :disabled="saving">
          {{ saving ? 'Сохранение…' : 'Сохранить' }}
        </Button>
        <NuxtLink to="/admin/collections">
          <Button type="button" variant="outline">Назад</Button>
        </NuxtLink>
      </div>
    </form>

    <div class="max-w-3xl">
      <h2 class="text-xl font-bold font-serif mb-4">Книги в подборке</h2>

      <Card class="mb-4">
        <CardContent class="pt-6">
          <div class="flex gap-3">
            <Input v-model="addBookId" type="number" placeholder="ID книги" class="flex-1" />
            <Button @click="addBook" :disabled="addingBook">
              <Plus class="w-4 h-4" /> Добавить
            </Button>
          </div>
        </CardContent>
      </Card>

      <Card class="p-0 overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Название</TableHead>
              <TableHead class="text-right">Действия</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="b in books" :key="b.id">
              <TableCell class="text-muted-foreground">{{ b.id }}</TableCell>
              <TableCell class="font-medium">{{ b.title }}</TableCell>
              <TableCell class="text-right">
                <Button
                  variant="ghost"
                  size="sm"
                  class="text-destructive hover:text-destructive"
                  @click="removeBook(b.id, b.title)"
                >
                  Убрать
                </Button>
              </TableCell>
            </TableRow>
            <TableEmpty v-if="books.length === 0" :colspan="3">Книг пока нет</TableEmpty>
          </TableBody>
        </Table>
      </Card>
    </div>
  </div>
</template>
