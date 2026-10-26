<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const api = useApi()
const toast = useToast()
const localePath = useLocalePath()

const form = reactive({
  slug: '',
  name_tm: '',
  name_ru: '',
  name_en: ''
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
    const created = await api.post<{ id: number }>('/admin/publication-categories', form)
    toast.success(t('admin.publication_categories.created'))
    await navigateTo(localePath(`/admin/publication-categories/${created.id}`))
  } catch (e: any) {
    error.value = e.data?.error || t('admin.common.error_save')
  } finally {
    saving.value = false
  }
}

useHead({ title: () => t('admin.publication_categories.add_new') })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/publication-categories" class="text-muted-foreground hover:text-foreground text-sm">← {{ t('admin.publication_categories.title') }}</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">{{ t('admin.publication_categories.add_new') }}</h1>
    </div>
    <PublicationCategoryForm v-model="form" :loading="saving" :error="error" @submit="submit" />
  </div>
</template>
