<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const route = useRoute()
const api = useApi()
const localePath = useLocalePath()
const id = Number(route.params.id)

const { data: section } = await useAsyncData(`section-${id}`, () =>
  api.get<any>(`/sections/${id}`)
)

const { data: allSectionsResp } = await useAsyncData('all-sections-for-edit', () =>
  api.get<{ data: any[] }>('/sections')
)
const rootSections = computed(() => allSectionsResp.value?.data ?? [])

const form = reactive({
  slug: section.value?.slug ?? '',
  parent_id: section.value?.parent_id ?? null,
  sort_order: section.value?.sort_order ?? 0,
  image_path: section.value?.image_path ?? '',
  name_tm: section.value?.name_tm ?? '',
  name_ru: section.value?.name_ru ?? '',
  name_en: section.value?.name_en ?? '',
  description_tm: section.value?.description_tm ?? '',
  description_ru: section.value?.description_ru ?? '',
  description_en: section.value?.description_en ?? ''
})

const saving = ref(false)
const error = ref('')

async function submit() {
  saving.value = true
  error.value = ''
  try {
    await api.put(`/admin/sections/${id}`, form)
    await navigateTo(localePath('/admin/sections'))
  } catch (e: any) {
    error.value = e.data?.error || 'Ошибка сохранения'
  } finally {
    saving.value = false
  }
}

useHead({ title: `Изм. раздела: ${section.value?.name_ru}` })
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">Редактировать раздел</h1>
    <SectionForm
      v-model="form"
      :root-sections="rootSections"
      :exclude-id="id"
      :loading="saving"
      :error="error"
      @submit="submit"
    />
  </div>
</template>
