<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const route = useRoute()
const api = useApi()
const toast = useToast()
const id = Number(route.params.id)

const { data: book, refresh: refreshBook } = await useAsyncData(`book-${id}`, () =>
  api.get<any>(`/books/${id}`)
)

const { data: sectionsResp } = await useAsyncData('all-sections-for-book-edit', () =>
  api.get<{ data: any[] }>('/sections')
)
const sections = computed(() => sectionsResp.value?.data ?? [])

const form = reactive({
  slug: book.value?.slug ?? '',
  isbn: book.value?.isbn ?? '',
  year: book.value?.year ?? 0,
  pages: book.value?.pages ?? 0,
  language: book.value?.language ?? 'tm',
  title: book.value?.title ?? '',
  annotation: book.value?.annotation ?? '',
  authors: book.value?.authors ?? '',
  publisher: book.value?.publisher ?? '',
  allow_download: book.value?.allow_download ?? true,
  is_published: book.value?.is_published ?? false,
  section_ids: (book.value?.sections ?? []).map((s: any) => s.id) as number[]
})

const saving = ref(false)
const error = ref('')

async function submit() {
  if (form.section_ids.length === 0) {
    error.value = 'Выберите хотя бы один раздел'
    return
  }
  saving.value = true
  error.value = ''
  try {
    await api.put(`/admin/books/${id}`, form)
    await refreshBook()
    toast.success('Сохранено')
  } catch (e: any) {
    error.value = e.data?.error || 'Ошибка сохранения'
  } finally {
    saving.value = false
  }
}

const uploading = ref(false)
const uploadError = ref('')
const fileInput = ref<HTMLInputElement | null>(null)

async function uploadFile() {
  const file = fileInput.value?.files?.[0]
  if (!file) return
  uploading.value = true
  uploadError.value = ''
  try {
    await api.upload(`/admin/books/${id}/file`, file)
    await refreshBook()
    toast.success('Файл загружен. Обложка сгенерирована автоматически (для PDF).')
    if (fileInput.value) fileInput.value.value = ''
  } catch (e: any) {
    uploadError.value = e.data?.error || 'Ошибка загрузки'
  } finally {
    uploading.value = false
  }
}

useHead({ title: `Изм. книги: ${book.value?.title}` })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/books" class="text-muted-foreground hover:text-foreground text-sm">← Книги</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">Редактировать книгу</h1>
    </div>

    <Card class="mb-6 max-w-3xl">
      <CardHeader>
        <CardTitle class="text-lg">Файл книги</CardTitle>
        <CardDescription v-if="book?.file_format">
          Текущий формат: <span class="font-mono text-foreground">{{ book.file_format }}</span>
        </CardDescription>
        <CardDescription v-else>Файл ещё не загружен.</CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="flex items-center gap-3 flex-wrap">
          <input
            ref="fileInput"
            type="file"
            accept=".pdf,.epub"
            class="text-sm flex-1 min-w-[200px] file:mr-3 file:rounded-md file:border-0 file:bg-muted file:px-3 file:py-1.5 file:text-sm file:font-medium file:text-foreground hover:file:bg-muted/80"
          />
          <Button @click="uploadFile" :disabled="uploading" size="sm">
            {{ uploading ? 'Загрузка…' : 'Загрузить' }}
          </Button>
        </div>
        <p v-if="uploadError" class="text-sm text-destructive">{{ uploadError }}</p>
        <p v-if="book?.cover_path" class="text-sm text-muted-foreground">
          Обложка: <span class="font-mono text-foreground">{{ book.cover_path }}</span>
        </p>
      </CardContent>
    </Card>

    <BookForm
      v-model="form"
      :sections="sections"
      :loading="saving"
      :error="error"
      @submit="submit"
    />
  </div>
</template>
