<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const route = useRoute()
const api = useApi()
const toast = useToast()
const loc = useLocalized()
const id = Number(route.params.id)

const { data: dept, refresh } = await useAsyncData(`department-${id}`, () =>
  api.get<any>(`/departments/${id}`)
)

const form = reactive({
  slug: dept.value?.slug ?? '',
  name_tm: dept.value?.name_tm ?? '',
  name_ru: dept.value?.name_ru ?? '',
  name_en: dept.value?.name_en ?? '',
  description_tm: dept.value?.description_tm ?? '',
  description_ru: dept.value?.description_ru ?? '',
  description_en: dept.value?.description_en ?? ''
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
    await api.put(`/admin/departments/${id}`, form)
    await refresh()
    toast.success(t('admin.common.saved'))
  } catch (e: any) {
    error.value = e.data?.error || t('admin.common.error_save')
  } finally {
    saving.value = false
  }
}

useHead({ title: () => `${t('admin.departments.edit_title')}: ${loc(dept.value, 'name')}` })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/departments" class="text-muted-foreground hover:text-foreground text-sm">← {{ t('admin.departments.title') }}</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">{{ loc(dept, 'name') }}</h1>
    </div>
    <DepartmentForm v-model="form" :loading="saving" :error="error" @submit="submit" />
  </div>
</template>
