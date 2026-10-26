<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const route = useRoute()
const api = useApi()
const toast = useToast()
const loc = useLocalized()
const id = Number(route.params.id)

const { data: category, refresh } = await useAsyncData(`publication-category-${id}`, () =>
  api.get<any>(`/publication-categories/${id}`)
)

const form = reactive({
  slug: category.value?.slug ?? '',
  name_tm: category.value?.name_tm ?? '',
  name_ru: category.value?.name_ru ?? '',
  name_en: category.value?.name_en ?? ''
})

const saving = ref(false)
const error = ref('')

async function submit() {
  if (!form.name_tm) {
    error.value = t('admin.common.title_required')
    return
  }
  saving.value = true
  error.value = ''
  try {
    await api.put(`/admin/publication-categories/${id}`, form)
    await refresh()
    toast.success(t('admin.common.saved'))
  } catch (e: any) {
    error.value = e.data?.error || t('admin.common.error_save')
  } finally {
    saving.value = false
  }
}

useHead({ title: () => `${t('admin.publication_categories.edit_title')}: ${loc(category.value, 'name')}` })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/publication-categories" class="text-muted-foreground hover:text-foreground text-sm">← {{ t('admin.publication_categories.title') }}</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">{{ loc(category, 'name') }}</h1>
    </div>
    <PublicationCategoryForm v-model="form" :loading="saving" :error="error" @submit="submit" />
  </div>
</template>
