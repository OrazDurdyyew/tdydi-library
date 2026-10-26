<script setup lang="ts">
interface DepartmentFormData {
  slug: string
  name_tm: string
  name_ru: string
  name_en: string
  description_tm: string
  description_ru: string
  description_en: string
}

const props = defineProps<{
  modelValue: DepartmentFormData
  loading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: DepartmentFormData]
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
          <Label for="dept-slug">{{ t('admin.common.slug') }}</Label>
          <Input id="dept-slug" v-model="form.slug" readonly tabindex="-1" class="font-mono bg-muted text-muted-foreground cursor-not-allowed" :placeholder="t('admin.common.slug_auto_hint')" />
          <p class="text-xs text-muted-foreground">{{ t('admin.common.slug_auto_hint') }}</p>
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
            <Label for="dept-name-tm">{{ t('admin.common.lang_tm') }} <span class="text-destructive">*</span></Label>
            <Input id="dept-name-tm" v-model="form.name_tm" required />
          </div>
          <div class="space-y-2">
            <Label for="dept-name-ru">{{ t('admin.common.lang_ru') }}</Label>
            <Input id="dept-name-ru" v-model="form.name_ru" />
          </div>
          <div class="space-y-2">
            <Label for="dept-name-en">{{ t('admin.common.lang_en') }}</Label>
            <Input id="dept-name-en" v-model="form.name_en" />
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
            <Label for="dept-desc-tm">{{ t('admin.common.lang_tm') }}</Label>
            <Textarea id="dept-desc-tm" v-model="form.description_tm" rows="3" />
          </div>
          <div class="space-y-2">
            <Label for="dept-desc-ru">{{ t('admin.common.lang_ru') }}</Label>
            <Textarea id="dept-desc-ru" v-model="form.description_ru" rows="3" />
          </div>
          <div class="space-y-2">
            <Label for="dept-desc-en">{{ t('admin.common.lang_en') }}</Label>
            <Textarea id="dept-desc-en" v-model="form.description_en" rows="3" />
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
      <NuxtLink to="/admin/departments">
        <Button type="button" variant="outline">{{ t('admin.common.cancel') }}</Button>
      </NuxtLink>
    </div>
  </form>
</template>
