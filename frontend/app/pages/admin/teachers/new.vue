<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const api = useApi()
const toast = useToast()
const localePath = useLocalePath()

const { data: deptsResp } = await useAsyncData('teachers-new-departments', () =>
  api.get<{ data: any[] }>('/departments?limit=200')
)
const departments = computed(() => deptsResp.value?.data ?? [])

const form = reactive({
  department_id: null as number | null,
  slug: '',
  first_name: '',
  last_name: '',
  middle_name: '',
  bio: ''
})

const saving = ref(false)
const error = ref('')

async function submit() {
  if (!form.department_id) {
    error.value = t('admin.teachers.department_required')
    return
  }
  if (!form.first_name || !form.last_name) {
    error.value = t('admin.teachers.name_required')
    return
  }
  saving.value = true
  error.value = ''
  try {
    const created = await api.post<{ id: number }>('/admin/teachers', form)
    toast.success(t('admin.teachers.created'))
    await navigateTo(localePath(`/admin/teachers/${created.id}`))
  } catch (e: any) {
    error.value = e.data?.error || t('admin.common.error_save')
  } finally {
    saving.value = false
  }
}

useHead({ title: () => t('admin.teachers.add_new') })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/teachers" class="text-muted-foreground hover:text-foreground text-sm">← {{ t('admin.teachers.title') }}</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">{{ t('admin.teachers.add_new') }}</h1>
    </div>
    <TeacherForm v-model="form" :departments="departments" :loading="saving" :error="error" @submit="submit" />
    <p class="mt-4 text-sm text-muted-foreground">{{ t('admin.teachers.after_save_hint') }}</p>
  </div>
</template>
