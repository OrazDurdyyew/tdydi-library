<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const route = useRoute()
const api = useApi()
const toast = useToast()
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
    toast.success(t('admin.common.saved'))
    await navigateTo(localePath('/admin/sections'))
  } catch (e: any) {
    error.value = e.data?.error || t('admin.common.error_save')
  } finally {
    saving.value = false
  }
}

useHead({ title: () => `${t('admin.sections.edit_title')}: ${section.value?.name_ru}` })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/sections" class="text-muted-foreground hover:text-foreground text-sm">← {{ t('admin.sections.title') }}</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">{{ section?.name_ru || section?.name_tm }}</h1>
    </div>
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
