<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const api = useApi()
const localePath = useLocalePath()

const { data: allSectionsResp } = await useAsyncData('all-sections-for-form', () =>
  api.get<{ data: any[] }>('/sections')
)
const rootSections = computed(() => allSectionsResp.value?.data ?? [])

const form = reactive({
  slug: '',
  parent_id: null as number | null,
  sort_order: 0,
  image_path: '',
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
  saving.value = true
  error.value = ''
  try {
    await api.post('/admin/sections', form)
    await navigateTo(localePath('/admin/sections'))
  } catch (e: any) {
    error.value = e.data?.error || 'Ошибка сохранения'
  } finally {
    saving.value = false
  }
}

useHead({ title: 'Создать раздел' })
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">Создать раздел</h1>
    <SectionForm
      v-model="form"
      :root-sections="rootSections"
      :loading="saving"
      :error="error"
      @submit="submit"
    />
  </div>
</template>
