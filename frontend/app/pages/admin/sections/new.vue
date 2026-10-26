<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const api = useApi()
const toast = useToast()
const route = useRoute()
const localePath = useLocalePath()

const { data: allSectionsResp } = await useAsyncData('all-sections-for-form', () =>
  api.get<{ data: any[] }>('/sections')
)
const rootSections = computed(() => allSectionsResp.value?.data ?? [])

// ?parent=<id> — предвыбор родителя при создании подраздела из списка разделов.
const presetParent = route.query.parent ? Number(route.query.parent) : null

const form = reactive({
  slug: '',
  parent_id: presetParent,
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
    toast.success(t('admin.sections.created'))
    await navigateTo(localePath('/admin/sections'))
  } catch (e: any) {
    error.value = e.data?.error || t('admin.common.error_save')
  } finally {
    saving.value = false
  }
}

useHead({ title: () => t('admin.sections.add_new') })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/sections" class="text-muted-foreground hover:text-foreground text-sm">← {{ t('admin.sections.title') }}</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">
        {{ presetParent ? t('admin.sections.add_subsection') : t('admin.sections.add_new') }}
      </h1>
    </div>
    <SectionForm
      v-model="form"
      :root-sections="rootSections"
      :loading="saving"
      :error="error"
      @submit="submit"
    />
  </div>
</template>
