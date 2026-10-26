<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const route = useRoute()
const api = useApi()
const toast = useToast()
const fileUrl = useFileUrl()
const id = Number(route.params.id)

const { data: teacher, refresh } = await useAsyncData(`teacher-${id}`, () =>
  api.get<any>(`/teachers/${id}`)
)

const { data: deptsResp } = await useAsyncData('teachers-edit-departments', () =>
  api.get<{ data: any[] }>('/departments?limit=200')
)
const departments = computed(() => deptsResp.value?.data ?? [])

const form = reactive({
  department_id: teacher.value?.department_id ?? null,
  slug: teacher.value?.slug ?? '',
  first_name: teacher.value?.first_name ?? '',
  last_name: teacher.value?.last_name ?? '',
  middle_name: teacher.value?.middle_name ?? '',
  bio: teacher.value?.bio ?? ''
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
    await api.put(`/admin/teachers/${id}`, form)
    await refresh()
    toast.success(t('admin.common.saved'))
  } catch (e: any) {
    error.value = e.data?.error || t('admin.common.error_save')
  } finally {
    saving.value = false
  }
}

const photoInput = ref<HTMLInputElement | null>(null)
const uploading = ref(false)

async function uploadPhoto() {
  const file = photoInput.value?.files?.[0]
  if (!file) return
  uploading.value = true
  try {
    await api.uploadField(`/admin/teachers/${id}/photo`, 'photo', file)
    await refresh()
    toast.success(t('admin.teachers.photo_uploaded'))
    if (photoInput.value) photoInput.value.value = ''
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error_upload'))
  } finally {
    uploading.value = false
  }
}

const photoPreview = computed(() => fileUrl(teacher.value?.photo_path))

useHead({ title: () => `${t('admin.teachers.edit_title')}: ${teacher.value?.last_name}` })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/teachers" class="text-muted-foreground hover:text-foreground text-sm">← {{ t('admin.teachers.title') }}</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">{{ t('admin.teachers.edit_title') }}</h1>
    </div>

    <Card class="mb-6 max-w-3xl">
      <CardHeader>
        <CardTitle class="text-lg">{{ t('admin.teachers.photo') }}</CardTitle>
        <CardDescription v-if="!teacher?.photo_path">{{ t('admin.teachers.photo_not_uploaded') }}</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="flex items-start gap-5 flex-wrap">
          <div class="w-24 h-24 rounded-full overflow-hidden bg-muted shrink-0 flex items-center justify-center text-muted-foreground text-xs">
            <img
              v-if="photoPreview"
              :src="photoPreview"
              :alt="teacher?.last_name || ''"
              class="w-full h-full object-cover"
            />
            <span v-else>{{ t('admin.common.no') }}</span>
          </div>
          <div class="flex-1 min-w-[200px] space-y-3">
            <input
              ref="photoInput"
              type="file"
              accept="image/*"
              :disabled="uploading"
              class="text-sm w-full file:mr-3 file:rounded-md file:border-0 file:bg-muted file:px-3 file:py-1.5 file:text-sm file:font-medium file:text-foreground hover:file:bg-muted/80"
            />
            <Button @click="uploadPhoto" :disabled="uploading" size="sm">
              {{ uploading ? t('admin.common.uploading') : t('admin.common.upload') }}
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>

    <TeacherForm v-model="form" :departments="departments" :loading="saving" :error="error" @submit="submit" />
  </div>
</template>
