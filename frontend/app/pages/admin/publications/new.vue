<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const api = useApi()
const toast = useToast()
const localePath = useLocalePath()

const { data: deptsResp } = await useAsyncData('publications-new-departments', () =>
  api.get<{ data: any[] }>('/departments?limit=200')
)
const { data: catsResp } = await useAsyncData('publications-new-categories', () =>
  api.get<{ data: any[] }>('/publication-categories?limit=200')
)
const { data: periodicalsResp } = await useAsyncData('publications-new-periodicals', () =>
  api.get<{ data: any[] }>('/periodicals?limit=200')
)

const departments = computed(() => deptsResp.value?.data ?? [])
const categories = computed(() => catsResp.value?.data ?? [])
const periodicals = computed(() => periodicalsResp.value?.data ?? [])

const form = reactive({
  department_id: null as number | null,
  teacher_id: null as number | null,
  category_id: null as number | null,
  periodical_id: null as number | null,
  periodical_issue_id: null as number | null,
  name_tm: '',
  name_ru: '',
  name_en: '',
  article_date: '',
  page_count: 0,
  is_published: true
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
    const created = await api.post<{ id: number }>('/admin/publications', form)
    toast.success(t('admin.publications.created'))
    await navigateTo(localePath(`/admin/publications/${created.id}`))
  } catch (e: any) {
    error.value = e.data?.error || t('admin.common.error_save')
  } finally {
    saving.value = false
  }
}

useHead({ title: () => t('admin.publications.add_new') })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/publications" class="text-muted-foreground hover:text-foreground text-sm">← {{ t('admin.publications.title') }}</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">{{ t('admin.publications.add_new') }}</h1>
    </div>
    <PublicationForm
      v-model="form"
      :departments="departments"
      :categories="categories"
      :periodicals="periodicals"
      :loading="saving"
      :error="error"
      @submit="submit"
    />
    <p class="mt-4 text-sm text-muted-foreground">{{ t('admin.publications.after_save_hint') }}</p>
  </div>
</template>
