<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const route = useRoute()
const api = useApi()
const localePath = useLocalePath()
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
    alert('Сохранено')
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
    alert('Файл загружен. Обложка автоматически сгенерирована (для PDF).')
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
    <h1 class="text-3xl font-bold mb-8">Редактировать книгу</h1>

    <div class="card mb-6 max-w-3xl">
      <h3 class="font-semibold mb-4">Файл книги</h3>
      <div v-if="book?.file_format" class="mb-4 text-sm text-muted-foreground">
        Текущий формат: <span class="font-mono text-foreground">{{ book.file_format }}</span>
      </div>
      <div v-else class="mb-4 text-sm text-muted-foreground">Файл ещё не загружен.</div>

      <div class="flex items-center gap-3 flex-wrap">
        <input
          ref="fileInput"
          type="file"
          accept=".pdf,.epub"
          class="text-sm flex-1 min-w-[200px]"
        />
        <button
          @click="uploadFile"
          :disabled="uploading"
          class="btn-primary btn-sm"
        >
          {{ uploading ? 'Загрузка…' : 'Загрузить' }}
        </button>
      </div>
      <div v-if="uploadError" class="text-sm text-red-600 mt-3">{{ uploadError }}</div>

      <div v-if="book?.cover_path" class="mt-4 text-sm text-muted-foreground">
        Обложка: <span class="font-mono text-foreground">{{ book.cover_path }}</span>
      </div>
    </div>

    <BookForm
      v-model="form"
      :sections="sections"
      :loading="saving"
      :error="error"
      @submit="submit"
    />
  </div>
</template>
