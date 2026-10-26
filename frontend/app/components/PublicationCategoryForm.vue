<script setup lang="ts">
interface PublicationCategoryFormData {
  slug: string
  name_tm: string
  name_ru: string
  name_en: string
}

const props = defineProps<{
  modelValue: PublicationCategoryFormData
  loading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: PublicationCategoryFormData]
  submit: []
}>()

const { t } = useI18n()
const form = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v)
})

useAutoSlug(toRef(props.modelValue, 'name_tm'), toRef(props.modelValue, 'slug'))
</script>

<template>
  <form @submit.prevent="emit('submit')" class="space-y-6 max-w-3xl">
    <Card>
      <CardHeader>
        <CardTitle class="text-lg">{{ t('admin.sections.main') }}</CardTitle>
      </CardHeader>
      <CardContent class="space-y-5">
        <div class="space-y-2">
          <Label for="cat-slug">{{ t('admin.common.slug') }}</Label>
          <Input id="cat-slug" v-model="form.slug" readonly tabindex="-1" class="font-mono bg-muted text-muted-foreground cursor-not-allowed" :placeholder="t('admin.common.slug_auto_hint')" />
          <p class="text-xs text-muted-foreground">{{ t('admin.common.slug_auto_hint') }}</p>
        </div>

        <div class="grid grid-cols-3 gap-4">
          <div class="space-y-2">
            <Label for="cat-name-tm">{{ t('admin.common.lang_tm') }} <span class="text-destructive">*</span></Label>
            <Input id="cat-name-tm" v-model="form.name_tm" required />
          </div>
          <div class="space-y-2">
            <Label for="cat-name-ru">{{ t('admin.common.lang_ru') }}</Label>
            <Input id="cat-name-ru" v-model="form.name_ru" />
          </div>
          <div class="space-y-2">
            <Label for="cat-name-en">{{ t('admin.common.lang_en') }}</Label>
            <Input id="cat-name-en" v-model="form.name_en" />
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
      <NuxtLink to="/admin/publication-categories">
        <Button type="button" variant="outline">{{ t('admin.common.cancel') }}</Button>
      </NuxtLink>
    </div>
  </form>
</template>
