<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const api = useApi()
const localePath = useLocalePath()

const { data: sectionsResp } = await useAsyncData('all-sections-for-book-form', () =>
  api.get<{ data: any[] }>('/sections')
)
const sections = computed(() => sectionsResp.value?.data ?? [])

const form = reactive({
  slug: '',
  isbn: '',
  year: new Date().getFullYear(),
  pages: 0,
  language: 'tm',
  title: '',
  annotation: '',
  authors: '',
  publisher: '',
  allow_download: true,
  is_published: false,
  section_ids: [] as number[]
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
    const created = await api.post<{ id: number }>('/admin/books', form)
    await navigateTo(localePath(`/admin/books/${created.id}`))
  } catch (e: any) {
    error.value = e.data?.error || 'Ошибка сохранения'
  } finally {
    saving.value = false
  }
}

useHead({ title: 'Добавить книгу' })
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">Добавить книгу</h1>
    <BookForm
      v-model="form"
      :sections="sections"
      :loading="saving"
      :error="error"
      @submit="submit"
    />
    <div class="mt-4 text-sm text-muted-foreground">
      После сохранения откроется страница редактирования, где можно загрузить файл PDF/EPUB.
    </div>
  </div>
</template>
