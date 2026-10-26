<script setup lang="ts">
interface NamedEntity {
  id: number
  name_tm?: string
  name_ru?: string
  name_en?: string
}

interface Teacher {
  id: number
  first_name: string
  last_name: string
  middle_name?: string
}

interface Issue {
  id: number
  issue_number?: string
  issue_date?: string
}

interface PublicationFormData {
  department_id: number | null
  teacher_id: number | null
  category_id: number | null
  periodical_id: number | null
  periodical_issue_id: number | null
  name_tm: string
  name_ru: string
  name_en: string
  article_date: string
  page_count: number
  is_published: boolean
}

const props = defineProps<{
  modelValue: PublicationFormData
  departments: NamedEntity[]
  categories: NamedEntity[]
  periodicals: NamedEntity[]
  loading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: PublicationFormData]
  submit: []
}>()

const { t } = useI18n()
const loc = useLocalized()
const api = useApi()

const form = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v)
})

const teachers = ref<Teacher[]>([])
const issues = ref<Issue[]>([])

function teacherName(tch: Teacher): string {
  return [tch.last_name, tch.first_name, tch.middle_name].filter(Boolean).join(' ')
}

async function loadTeachers(deptId: number | null) {
  if (!deptId) {
    teachers.value = []
    return
  }
  try {
    const res = await api.get<{ data: Teacher[] }>(`/teachers?department_id=${deptId}`)
    teachers.value = res.data ?? []
  } catch {
    teachers.value = []
  }
}

async function loadIssues(periodicalId: number | null) {
  if (!periodicalId) {
    issues.value = []
    return
  }
  try {
    const res = await api.get<{ data: Issue[] }>(`/periodicals/${periodicalId}/issues?limit=200`)
    issues.value = res.data ?? []
  } catch {
    issues.value = []
  }
}

watch(
  () => props.modelValue.department_id,
  async (val, oldVal) => {
    await loadTeachers(val)
    // Сбрасываем выбранного преподавателя только при реальной смене кафедры,
    // чтобы не затирать значение при инициализации формы редактирования.
    if (oldVal !== undefined && val !== oldVal) {
      form.value.teacher_id = null
    }
  }
)

watch(
  () => props.modelValue.periodical_id,
  async (val, oldVal) => {
    await loadIssues(val)
    if (oldVal !== undefined && val !== oldVal) {
      form.value.periodical_issue_id = null
    }
  }
)

onMounted(async () => {
  await Promise.all([
    loadTeachers(props.modelValue.department_id),
    loadIssues(props.modelValue.periodical_id)
  ])
})

function issueLabel(i: Issue): string {
  return [i.issue_number, i.issue_date].filter(Boolean).join(' · ') || `#${i.id}`
}
</script>

<template>
  <form @submit.prevent="emit('submit')" class="space-y-6 max-w-3xl">
    <Card>
      <CardHeader>
        <CardTitle class="text-lg">{{ t('admin.publications.attribution') }}</CardTitle>
      </CardHeader>
      <CardContent class="space-y-5">
        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-2">
            <Label for="pub-dept">{{ t('admin.publications.department') }} <span class="text-destructive">*</span></Label>
            <Select id="pub-dept" v-model="form.department_id">
              <option :value="null">{{ t('admin.teachers.department_select') }}</option>
              <option v-for="d in departments" :key="d.id" :value="d.id">{{ loc(d, 'name') }}</option>
            </Select>
          </div>
          <div class="space-y-2">
            <Label for="pub-teacher">{{ t('admin.publications.teacher') }} <span class="text-destructive">*</span></Label>
            <Select id="pub-teacher" v-model="form.teacher_id" :disabled="!form.department_id">
              <option :value="null">{{ t('admin.publications.teacher_select') }}</option>
              <option v-for="tch in teachers" :key="tch.id" :value="tch.id">{{ teacherName(tch) }}</option>
            </Select>
          </div>
        </div>

        <div class="space-y-2">
          <Label for="pub-category">{{ t('admin.publications.category') }} <span class="text-destructive">*</span></Label>
          <Select id="pub-category" v-model="form.category_id">
            <option :value="null">{{ t('admin.publications.category_select') }}</option>
            <option v-for="c in categories" :key="c.id" :value="c.id">{{ loc(c, 'name') }}</option>
          </Select>
        </div>
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle class="text-lg">{{ t('admin.publications.source') }}</CardTitle>
        <CardDescription>{{ t('admin.publications.source_hint') }}</CardDescription>
      </CardHeader>
      <CardContent class="space-y-5">
        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-2">
            <Label for="pub-periodical">{{ t('admin.publications.periodical') }}</Label>
            <Select id="pub-periodical" v-model="form.periodical_id">
              <option :value="null">{{ t('admin.publications.periodical_none') }}</option>
              <option v-for="p in periodicals" :key="p.id" :value="p.id">{{ loc(p, 'name') }}</option>
            </Select>
          </div>
          <div class="space-y-2">
            <Label for="pub-issue">{{ t('admin.publications.issue') }}</Label>
            <Select id="pub-issue" v-model="form.periodical_issue_id" :disabled="!form.periodical_id">
              <option :value="null">{{ t('admin.publications.issue_none') }}</option>
              <option v-for="i in issues" :key="i.id" :value="i.id">{{ issueLabel(i) }}</option>
            </Select>
          </div>
        </div>
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle class="text-lg">{{ t('admin.publications.article') }}</CardTitle>
      </CardHeader>
      <CardContent class="space-y-5">
        <div class="space-y-2">
          <Label for="pub-name-tm">{{ t('admin.common.title') }} ({{ t('admin.common.lang_tm') }}) <span class="text-destructive">*</span></Label>
          <Input id="pub-name-tm" v-model="form.name_tm" required />
        </div>
        <div class="space-y-2">
          <Label for="pub-name-ru">{{ t('admin.common.title') }} ({{ t('admin.common.lang_ru') }})</Label>
          <Input id="pub-name-ru" v-model="form.name_ru" />
        </div>
        <div class="space-y-2">
          <Label for="pub-name-en">{{ t('admin.common.title') }} ({{ t('admin.common.lang_en') }})</Label>
          <Input id="pub-name-en" v-model="form.name_en" />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-2">
            <Label for="pub-date">{{ t('admin.publications.article_date') }}</Label>
            <Input id="pub-date" v-model="form.article_date" type="date" />
          </div>
          <div class="space-y-2">
            <Label for="pub-pages">{{ t('admin.publications.page_count') }}</Label>
            <Input id="pub-pages" v-model.number="form.page_count" type="number" min="0" />
          </div>
        </div>
      </CardContent>
    </Card>

    <Card>
      <CardContent class="pt-6">
        <label class="flex items-center gap-3 cursor-pointer">
          <Checkbox v-model="form.is_published" />
          <span class="text-sm">{{ t('admin.publications.published') }}</span>
        </label>
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
      <NuxtLink to="/admin/publications">
        <Button type="button" variant="outline">{{ t('admin.common.cancel') }}</Button>
      </NuxtLink>
    </div>
  </form>
</template>
