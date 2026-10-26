<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const { t } = useI18n()
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
  is_institute_book: book.value?.is_institute_book ?? false,
  section_ids: (book.value?.sections ?? []).map((s: any) => s.id) as number[]
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
    await api.put(`/admin/books/${id}`, form)
    await refreshBook()
    toast.success(t('admin.common.saved'))
  } catch (e: any) {
    error.value = e.data?.error || t('admin.common.error_save')
  } finally {
    saving.value = false
  }
}

const fileUrl = useFileUrl()

const uploading = ref(false)
const uploadProgress = ref(0)
const uploadError = ref('')
const fileInput = ref<HTMLInputElement | null>(null)

async function uploadFile() {
  const file = fileInput.value?.files?.[0]
  if (!file) return
  uploading.value = true
  uploadError.value = ''
  uploadProgress.value = 0
  try {
    await api.uploadWithProgress(
      `/admin/books/${id}/file`,
      'file',
      file,
      (p) => (uploadProgress.value = p)
    )
    const hadCustomCover = !!book.value?.cover_path
    await refreshBook()
    toast.success(
      hadCustomCover
        ? t('admin.books.file_uploaded')
        : t('admin.books.file_uploaded_cover_auto')
    )
    if (fileInput.value) fileInput.value.value = ''
  } catch (e: any) {
    uploadError.value = e.data?.error || t('admin.common.error_upload')
  } finally {
    uploading.value = false
    uploadProgress.value = 0
  }
}

const coverInput = ref<HTMLInputElement | null>(null)
const coverUploading = ref(false)
const coverProgress = ref(0)
const coverError = ref('')

async function uploadCover() {
  const file = coverInput.value?.files?.[0]
  if (!file) return
  coverUploading.value = true
  coverError.value = ''
  coverProgress.value = 0
  try {
    await api.uploadWithProgress(
      `/admin/books/${id}/cover`,
      'cover',
      file,
      (p) => (coverProgress.value = p)
    )
    await refreshBook()
    toast.success(t('admin.books.cover_updated'))
    if (coverInput.value) coverInput.value.value = ''
  } catch (e: any) {
    coverError.value = e.data?.error || t('admin.common.error_upload')
  } finally {
    coverUploading.value = false
    coverProgress.value = 0
  }
}

const coverPreview = computed(() => fileUrl(book.value?.cover_path))

useHead({ title: () => `${t('admin.books.edit_title')}: ${book.value?.title}` })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/books" class="text-muted-foreground hover:text-foreground text-sm">← {{ t('admin.books.title') }}</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">{{ t('admin.books.edit_title') }}</h1>
    </div>

    <Card class="mb-6 max-w-3xl">
      <CardHeader>
        <CardTitle class="text-lg">{{ t('admin.books.file_card') }}</CardTitle>
        <CardDescription v-if="book?.file_format">
          {{ t('admin.books.file_format') }}: <span class="font-mono text-foreground">{{ book.file_format }}</span>
        </CardDescription>
        <CardDescription v-else>{{ t('admin.books.file_not_uploaded') }}</CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="flex items-center gap-3 flex-wrap">
          <input
            ref="fileInput"
            type="file"
            accept=".pdf,.epub"
            :disabled="uploading"
            class="text-sm flex-1 min-w-[200px] file:mr-3 file:rounded-md file:border-0 file:bg-muted file:px-3 file:py-1.5 file:text-sm file:font-medium file:text-foreground hover:file:bg-muted/80"
          />
          <Button @click="uploadFile" :disabled="uploading" size="sm">
            {{ uploading ? `${t('admin.common.uploading')} ${uploadProgress}%` : t('admin.common.upload') }}
          </Button>
        </div>
        <div v-if="uploading" class="h-2 w-full bg-muted rounded overflow-hidden">
          <div
            class="h-full bg-primary transition-all duration-150"
            :style="{ width: `${uploadProgress}%` }"
          />
        </div>
        <p v-if="uploadError" class="text-sm text-destructive">{{ uploadError }}</p>
      </CardContent>
    </Card>

    <Card class="mb-6 max-w-3xl">
      <CardHeader>
        <CardTitle class="text-lg">{{ t('admin.books.cover_card') }}</CardTitle>
        <CardDescription>
          {{ t('admin.books.cover_hint') }}
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="flex items-start gap-5 flex-wrap">
          <div class="w-28 aspect-[3/4] rounded-md overflow-hidden bg-muted shrink-0 flex items-center justify-center text-muted-foreground text-xs">
            <img
              v-if="coverPreview"
              :src="coverPreview"
              :alt="book?.title || ''"
              class="w-full h-full object-cover"
            />
            <span v-else>{{ t('admin.common.no') }}</span>
          </div>
          <div class="flex-1 min-w-[200px] space-y-3">
            <input
              ref="coverInput"
              type="file"
              accept="image/*"
              :disabled="coverUploading"
              class="text-sm w-full file:mr-3 file:rounded-md file:border-0 file:bg-muted file:px-3 file:py-1.5 file:text-sm file:font-medium file:text-foreground hover:file:bg-muted/80"
            />
            <Button @click="uploadCover" :disabled="coverUploading" size="sm">
              {{ coverUploading ? `${t('admin.common.uploading')} ${coverProgress}%` : t('admin.books.cover_upload') }}
            </Button>
            <div v-if="coverUploading" class="h-2 w-full bg-muted rounded overflow-hidden">
              <div
                class="h-full bg-primary transition-all duration-150"
                :style="{ width: `${coverProgress}%` }"
              />
            </div>
            <p v-if="coverError" class="text-sm text-destructive">{{ coverError }}</p>
          </div>
        </div>
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
