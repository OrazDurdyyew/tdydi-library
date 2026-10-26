<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const route = useRoute()
const api = useApi()
const toast = useToast()
const loc = useLocalized()
const id = Number(route.params.id)

const { data: publication, refresh } = await useAsyncData(`publication-${id}`, () =>
  api.get<any>(`/publications/${id}`)
)

const { data: deptsResp } = await useAsyncData('publications-edit-departments', () =>
  api.get<{ data: any[] }>('/departments?limit=200')
)
const { data: catsResp } = await useAsyncData('publications-edit-categories', () =>
  api.get<{ data: any[] }>('/publication-categories?limit=200')
)
const { data: periodicalsResp } = await useAsyncData('publications-edit-periodicals', () =>
  api.get<{ data: any[] }>('/periodicals?limit=200')
)

const departments = computed(() => deptsResp.value?.data ?? [])
const categories = computed(() => catsResp.value?.data ?? [])
const periodicals = computed(() => periodicalsResp.value?.data ?? [])

const form = reactive({
  department_id: publication.value?.teacher?.department_id ?? null,
  teacher_id: publication.value?.teacher_id ?? null,
  category_id: publication.value?.category_id ?? null,
  periodical_id: publication.value?.periodical_issue?.periodical_id ?? null,
  periodical_issue_id: publication.value?.periodical_issue_id ?? null,
  name_tm: publication.value?.name_tm ?? '',
  name_ru: publication.value?.name_ru ?? '',
  name_en: publication.value?.name_en ?? '',
  article_date: publication.value?.article_date ?? '',
  page_count: publication.value?.page_count ?? 0,
  is_published: publication.value?.is_published ?? true
})

const saving = ref(false)
const error = ref('')

async function submit() {
  if (!form.teacher_id) {
    error.value = t('admin.publications.teacher_required')
    return
  }
  if (!form.category_id) {
    error.value = t('admin.publications.category_required')
    return
  }
  if (!form.name_tm) {
    error.value = t('admin.common.title_required')
    return
  }
  saving.value = true
  error.value = ''
  try {
    await api.put(`/admin/publications/${id}`, form)
    await refresh()
    toast.success(t('admin.common.saved'))
  } catch (e: any) {
    error.value = e.data?.error || t('admin.common.error_save')
  } finally {
    saving.value = false
  }
}

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
      `/admin/publications/${id}/file`,
      'file',
      file,
      (p) => (uploadProgress.value = p)
    )
    await refresh()
    toast.success(t('admin.publications.file_uploaded_cover_auto'))
    if (fileInput.value) fileInput.value.value = ''
  } catch (e: any) {
    uploadError.value = e.data?.error || t('admin.common.error_upload')
  } finally {
    uploading.value = false
    uploadProgress.value = 0
  }
}

useHead({ title: () => `${t('admin.publications.edit_title')}: ${loc(publication.value, 'name')}` })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/publications" class="text-muted-foreground hover:text-foreground text-sm">← {{ t('admin.publications.title') }}</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">{{ t('admin.publications.edit_title') }}</h1>
    </div>

    <Card class="mb-6 max-w-3xl">
      <CardHeader>
        <CardTitle class="text-lg">{{ t('admin.publications.file_card') }}</CardTitle>
        <CardDescription v-if="publication?.has_file">{{ t('admin.publications.file_uploaded') }}</CardDescription>
        <CardDescription v-else>{{ t('admin.publications.no_file') }}</CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="flex items-center gap-3 flex-wrap">
          <input
            ref="fileInput"
            type="file"
            accept=".pdf"
            :disabled="uploading"
            class="text-sm flex-1 min-w-[200px] file:mr-3 file:rounded-md file:border-0 file:bg-muted file:px-3 file:py-1.5 file:text-sm file:font-medium file:text-foreground hover:file:bg-muted/80"
          />
          <Button @click="uploadFile" :disabled="uploading" size="sm">
            {{ uploading ? `${t('admin.common.uploading')} ${uploadProgress}%` : t('admin.common.upload') }}
          </Button>
        </div>
        <div v-if="uploading" class="h-2 w-full bg-muted rounded overflow-hidden">
          <div class="h-full bg-primary transition-all duration-150" :style="{ width: `${uploadProgress}%` }" />
        </div>
        <p v-if="uploadError" class="text-sm text-destructive">{{ uploadError }}</p>
      </CardContent>
    </Card>

    <PublicationForm
      v-model="form"
      :departments="departments"
      :categories="categories"
      :periodicals="periodicals"
      :loading="saving"
      :error="error"
      @submit="submit"
    />
  </div>
</template>
