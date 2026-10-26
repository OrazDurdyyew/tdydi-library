<script setup lang="ts">
interface Department {
  id: number
  name_tm: string
  name_ru: string
  name_en: string
}

interface TeacherFormData {
  department_id: number | null
  slug: string
  first_name: string
  last_name: string
  middle_name: string
  bio: string
}

const props = defineProps<{
  modelValue: TeacherFormData
  departments: Department[]
  loading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: TeacherFormData]
  submit: []
}>()

const { t } = useI18n()
const loc = useLocalized()
const form = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v)
})

// Slug подсказывается из имени+фамилии; финальный slug формирует бэк (dept.slug + имя).
const nameSource = computed(() => `${props.modelValue.first_name} ${props.modelValue.last_name}`.trim())
useAutoSlug(nameSource, toRef(props.modelValue, 'slug'))
</script>

<template>
  <form @submit.prevent="emit('submit')" class="space-y-6 max-w-3xl">
    <Card>
      <CardHeader>
        <CardTitle class="text-lg">{{ t('admin.sections.main') }}</CardTitle>
      </CardHeader>
      <CardContent class="space-y-5">
        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-2">
            <Label for="teacher-dept">{{ t('admin.teachers.department') }} <span class="text-destructive">*</span></Label>
            <Select id="teacher-dept" v-model="form.department_id">
              <option :value="null">{{ t('admin.teachers.department_select') }}</option>
              <option v-for="d in departments" :key="d.id" :value="d.id">
                {{ loc(d, 'name') }}
              </option>
            </Select>
          </div>
          <div class="space-y-2">
            <Label for="teacher-slug">{{ t('admin.common.slug') }}</Label>
            <Input id="teacher-slug" v-model="form.slug" readonly tabindex="-1" class="font-mono bg-muted text-muted-foreground cursor-not-allowed" :placeholder="t('admin.common.slug_auto_hint')" />
          </div>
        </div>

        <div class="grid grid-cols-3 gap-4">
          <div class="space-y-2">
            <Label for="teacher-last">{{ t('admin.teachers.last_name') }} <span class="text-destructive">*</span></Label>
            <Input id="teacher-last" v-model="form.last_name" required />
          </div>
          <div class="space-y-2">
            <Label for="teacher-first">{{ t('admin.teachers.first_name') }} <span class="text-destructive">*</span></Label>
            <Input id="teacher-first" v-model="form.first_name" required />
          </div>
          <div class="space-y-2">
            <Label for="teacher-middle">{{ t('admin.teachers.middle_name') }}</Label>
            <Input id="teacher-middle" v-model="form.middle_name" />
          </div>
        </div>

        <div class="space-y-2">
          <Label for="teacher-bio">{{ t('admin.teachers.bio') }}</Label>
          <Textarea id="teacher-bio" v-model="form.bio" rows="4" />
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
      <NuxtLink to="/admin/teachers">
        <Button type="button" variant="outline">{{ t('admin.common.cancel') }}</Button>
      </NuxtLink>
    </div>
  </form>
</template>
