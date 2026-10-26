<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const { t } = useI18n()
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
  is_institute_book: false,
  section_ids: [] as number[]
})

const saving = ref(false)
const error = ref('')

async function submit() {
  if (form.section_ids.length === 0) {
    error.value = t('admin.books.select_at_least_one_section')
    return
  }
  saving.value = true
  error.value = ''
  try {
    const created = await api.post<{ id: number }>('/admin/books', form)
    await navigateTo(localePath(`/admin/books/${created.id}`))
  } catch (e: any) {
    error.value = e.data?.error || t('admin.common.error_save')
  } finally {
    saving.value = false
  }
}

useHead({ title: () => t('admin.books.add_new') })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/books" class="text-muted-foreground hover:text-foreground text-sm">← {{ t('admin.books.title') }}</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">{{ t('admin.books.add_new') }}</h1>
    </div>
    <BookForm
      v-model="form"
      :sections="sections"
      :loading="saving"
      :error="error"
      @submit="submit"
    />
    <p class="mt-4 text-sm text-muted-foreground">
      {{ t('admin.books.after_save_hint') }}
    </p>
  </div>
</template>
