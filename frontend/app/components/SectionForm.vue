<script setup lang="ts">
interface SectionFormData {
  slug: string
  parent_id: number | null
  sort_order: number
  image_path: string
  name_tm: string
  name_ru: string
  name_en: string
  description_tm: string
  description_ru: string
  description_en: string
}

interface RootSection {
  id: number
  name_ru: string
  parent_id: number | null
}

const props = defineProps<{
  modelValue: SectionFormData
  rootSections: RootSection[]
  excludeId?: number
  loading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: SectionFormData]
  submit: []
}>()

const { t } = useI18n()
const form = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v)
})

useAutoSlug(toRef(props.modelValue, 'name_tm'), toRef(props.modelValue, 'slug'))

const availableParents = computed(() =>
  props.rootSections.filter((s) => s.parent_id == null && s.id !== props.excludeId)
)
</script>

<template>
  <form @submit.prevent="emit('submit')" class="space-y-6 max-w-3xl">
    <Card>
      <CardHeader>
        <CardTitle class="text-lg">{{ t('admin.sections.main') }}</CardTitle>
      </CardHeader>
      <CardContent class="space-y-5">
        <div class="space-y-2">
          <Label for="sec-slug">{{ t('admin.common.slug') }}</Label>
          <Input id="sec-slug" v-model="form.slug" readonly tabindex="-1" class="font-mono bg-muted text-muted-foreground cursor-not-allowed" placeholder="economics" />
          <p class="text-xs text-muted-foreground">
            {{ t('admin.sections.slug_hint') }}
          </p>
        </div>

        <div class="space-y-2">
          <Label for="sec-parent">{{ t('admin.sections.parent') }}</Label>
          <Select id="sec-parent" v-model="form.parent_id">
            <option :value="null">{{ t('admin.sections.parent_none') }}</option>
            <option v-for="s in availableParents" :key="s.id" :value="s.id">
              {{ s.name_ru }}
            </option>
          </Select>
          <p class="text-xs text-muted-foreground">
            {{ t('admin.sections.parent_hint') }}
          </p>
        </div>

        <div class="space-y-2">
          <Label for="sec-order">{{ t('admin.sections.sort_order') }}</Label>
          <Input id="sec-order" v-model.number="form.sort_order" type="number" class="w-32" />
        </div>
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle class="text-lg">{{ t('admin.sections.names_3lang') }}</CardTitle>
      </CardHeader>
      <CardContent class="space-y-5">
        <div class="grid grid-cols-3 gap-4">
          <div class="space-y-2">
            <Label for="sec-name-tm">{{ t('admin.common.lang_tm') }} *</Label>
            <Input id="sec-name-tm" v-model="form.name_tm" required />
          </div>
          <div class="space-y-2">
            <Label for="sec-name-ru">{{ t('admin.common.lang_ru') }} *</Label>
            <Input id="sec-name-ru" v-model="form.name_ru" required />
          </div>
          <div class="space-y-2">
            <Label for="sec-name-en">{{ t('admin.common.lang_en') }} *</Label>
            <Input id="sec-name-en" v-model="form.name_en" required />
          </div>
        </div>
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle class="text-lg">{{ t('admin.sections.descriptions_3lang') }}</CardTitle>
      </CardHeader>
      <CardContent>
        <div class="grid grid-cols-3 gap-4">
          <div class="space-y-2">
            <Label for="sec-desc-tm">{{ t('admin.common.lang_tm') }}</Label>
            <Textarea id="sec-desc-tm" v-model="form.description_tm" rows="3" />
          </div>
          <div class="space-y-2">
            <Label for="sec-desc-ru">{{ t('admin.common.lang_ru') }}</Label>
            <Textarea id="sec-desc-ru" v-model="form.description_ru" rows="3" />
          </div>
          <div class="space-y-2">
            <Label for="sec-desc-en">{{ t('admin.common.lang_en') }}</Label>
            <Textarea id="sec-desc-en" v-model="form.description_en" rows="3" />
          </div>
        </div>
      </CardContent>
    </Card>

    <div
      v-if="error"
      class="text-sm text-destructive p-3 bg-destructive/10 border border-destructive/20 rounded-lg"
    >
      {{ error }}
    </div>

    <div class="flex gap-3">
      <Button type="submit" :disabled="loading">
        {{ loading ? t('admin.common.saving') : t('admin.common.save') }}
      </Button>
      <NuxtLink to="/admin/sections">
        <Button type="button" variant="outline">{{ t('admin.common.cancel') }}</Button>
      </NuxtLink>
    </div>
  </form>
</template>
