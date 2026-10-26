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
  name_en: '',
  description_tm: '',
  description_ru: '',
  description_en: ''
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
    const created = await api.post<{ id: number }>('/admin/departments', form)
    toast.success(t('admin.departments.created'))
    await navigateTo(localePath(`/admin/departments/${created.id}`))
  } catch (e: any) {
    error.value = e.data?.error || t('admin.common.error_save')
  } finally {
    saving.value = false
  }
}

useHead({ title: () => t('admin.departments.add_new') })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/departments" class="text-muted-foreground hover:text-foreground text-sm">← {{ t('admin.departments.title') }}</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">{{ t('admin.departments.add_new') }}</h1>
    </div>
    <DepartmentForm v-model="form" :loading="saving" :error="error" @submit="submit" />
  </div>
</template>
