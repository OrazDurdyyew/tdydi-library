<script setup lang="ts">
import { Plus } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

const route = useRoute()
const api = useApi()
const toast = useToast()
const { ask } = useConfirm()
const id = Number(route.params.id)

const { data: periodical, refresh } = await useAsyncData(`periodical-${id}`, () =>
  api.get<any>(`/periodicals/${id}`)
)
const { data: issuesResp, refresh: refreshIssues } = await useAsyncData(`periodical-issues-${id}`, () =>
  api.get<{ data: any[] }>(`/periodicals/${id}/issues`)
)

const issues = computed(() => issuesResp.value?.data ?? [])

const form = reactive({
  slug: periodical.value?.slug ?? '',
  type: periodical.value?.type ?? 'magazine',
  periodicity: periodical.value?.periodicity ?? '',
  name_tm: periodical.value?.name_tm ?? '',
  name_ru: periodical.value?.name_ru ?? '',
  name_en: periodical.value?.name_en ?? '',
  description_tm: periodical.value?.description_tm ?? '',
  description_ru: periodical.value?.description_ru ?? '',
  description_en: periodical.value?.description_en ?? '',
})

const saving = ref(false)
const coverInput = ref<HTMLInputElement | null>(null)
const uploading = ref(false)

async function submit() {
  saving.value = true
  try {
    await api.put(`/admin/periodicals/${id}`, form)
    await refresh()
    toast.success('Сохранено')
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка')
  } finally {
    saving.value = false
  }
}

async function uploadCover() {
  const file = coverInput.value?.files?.[0]
  if (!file) return
  uploading.value = true
  try {
    await api.upload(`/admin/periodicals/${id}/cover`, file)
    await refresh()
    toast.success('Обложка загружена')
    if (coverInput.value) coverInput.value.value = ''
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка загрузки')
  } finally {
    uploading.value = false
  }
}

const newIssue = reactive({ issue_number: '', issue_date: '', annotation_ru: '' })
const addingIssue = ref(false)
const issueFileInputs = ref<Record<number, HTMLInputElement | null>>({})

async function createIssue() {
  addingIssue.value = true
  try {
    await api.post(`/admin/periodicals/${id}/issues`, {
      periodical_id: id,
      ...newIssue,
    })
    await refreshIssues()
    toast.success('Выпуск добавлен')
    newIssue.issue_number = ''
    newIssue.issue_date = ''
    newIssue.annotation_ru = ''
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка')
  } finally {
    addingIssue.value = false
  }
}

async function deleteIssue(issueId: number, num: string) {
  const ok = await ask({
    title: 'Удалить выпуск?',
    description: `Выпуск ${num || '#' + issueId} будет удалён.`,
    confirmText: 'Удалить',
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/periodicals/issues/${issueId}`)
    await refreshIssues()
    toast.success('Выпуск удалён')
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка')
  }
}

async function uploadIssueFile(issueId: number) {
  const file = issueFileInputs.value[issueId]?.files?.[0]
  if (!file) return
  try {
    await api.upload(`/admin/periodicals/issues/${issueId}/file`, file)
    await refreshIssues()
    toast.success('Файл загружен')
    if (issueFileInputs.value[issueId]) issueFileInputs.value[issueId]!.value = ''
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка загрузки')
  }
}

useHead({ title: `Периодика: ${periodical.value?.name_ru || periodical.value?.name_tm}` })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/periodicals" class="text-muted-foreground hover:text-foreground text-sm">← Периодика</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">{{ periodical?.name_ru || periodical?.name_tm }}</h1>
    </div>

    <Card class="mb-6 max-w-3xl">
      <CardHeader>
        <CardTitle class="text-lg">Обложка</CardTitle>
        <CardDescription v-if="periodical?.cover_path" class="font-mono">{{ periodical.cover_path }}</CardDescription>
        <CardDescription v-else>Обложка не загружена</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="flex items-center gap-3 flex-wrap">
          <input
            ref="coverInput"
            type="file"
            accept="image/*"
            class="text-sm flex-1 min-w-[200px] file:mr-3 file:rounded-md file:border-0 file:bg-muted file:px-3 file:py-1.5 file:text-sm file:font-medium hover:file:bg-muted/80"
          />
          <Button @click="uploadCover" :disabled="uploading" size="sm">
            {{ uploading ? 'Загрузка…' : 'Загрузить' }}
          </Button>
        </div>
      </CardContent>
    </Card>

    <form @submit.prevent="submit" class="space-y-6 max-w-3xl mb-10">
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Основное</CardTitle>
        </CardHeader>
        <CardContent class="space-y-5">
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <Label for="type">Тип</Label>
              <Select id="type" v-model="form.type">
                <option value="magazine">Журнал</option>
                <option value="newspaper">Газета</option>
              </Select>
            </div>
            <div class="space-y-2">
              <Label for="slug">Slug</Label>
              <Input id="slug" v-model="form.slug" class="font-mono text-sm" />
            </div>
          </div>
          <div class="space-y-2">
            <Label for="period">Периодичность</Label>
            <Input id="period" v-model="form.periodicity" />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Названия</CardTitle>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="name-tm">TM</Label>
            <Input id="name-tm" v-model="form.name_tm" />
          </div>
          <div class="space-y-2">
            <Label for="name-ru">RU</Label>
            <Input id="name-ru" v-model="form.name_ru" />
          </div>
          <div class="space-y-2">
            <Label for="name-en">EN</Label>
            <Input id="name-en" v-model="form.name_en" />
          </div>
        </CardContent>
      </Card>

      <div class="flex gap-3">
        <Button type="submit" :disabled="saving">
          {{ saving ? 'Сохранение…' : 'Сохранить' }}
        </Button>
        <NuxtLink to="/admin/periodicals">
          <Button type="button" variant="outline">Назад</Button>
        </NuxtLink>
      </div>
    </form>

    <div class="max-w-3xl">
      <h2 class="text-xl font-bold font-serif mb-4">Выпуски</h2>

      <Card class="mb-4">
        <CardHeader>
          <CardTitle class="text-base">Добавить выпуск</CardTitle>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-2">
              <Label for="issue-num">Номер</Label>
              <Input id="issue-num" v-model="newIssue.issue_number" placeholder="№ 1" />
            </div>
            <div class="space-y-2">
              <Label for="issue-date">Дата</Label>
              <Input id="issue-date" v-model="newIssue.issue_date" type="date" />
            </div>
          </div>
          <Button @click="createIssue" :disabled="addingIssue" size="sm">
            <Plus class="w-4 h-4" /> {{ addingIssue ? 'Добавление…' : 'Добавить' }}
          </Button>
        </CardContent>
      </Card>

      <Card class="p-0 overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Номер</TableHead>
              <TableHead>Дата</TableHead>
              <TableHead>Файл</TableHead>
              <TableHead class="text-right">Действия</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="issue in issues" :key="issue.id">
              <TableCell class="font-medium">{{ issue.issue_number || '—' }}</TableCell>
              <TableCell class="text-muted-foreground">{{ issue.issue_date || '—' }}</TableCell>
              <TableCell>
                <div v-if="issue.file_path" class="text-xs font-mono text-muted-foreground">
                  {{ issue.file_path }}
                </div>
                <div v-else class="flex items-center gap-2">
                  <input
                    :ref="(el) => { if (el) issueFileInputs[issue.id] = el as HTMLInputElement }"
                    type="file"
                    accept=".pdf"
                    class="text-xs flex-1 file:mr-2 file:rounded file:border-0 file:bg-muted file:px-2 file:py-1 file:text-xs"
                  />
                  <Button @click="uploadIssueFile(issue.id)" variant="outline" size="sm">Загрузить</Button>
                </div>
              </TableCell>
              <TableCell class="text-right">
                <Button
                  variant="ghost"
                  size="sm"
                  class="text-destructive hover:text-destructive"
                  @click="deleteIssue(issue.id, issue.issue_number)"
                >
                  Удалить
                </Button>
              </TableCell>
            </TableRow>
            <TableEmpty v-if="issues.length === 0" :colspan="4">Выпусков пока нет</TableEmpty>
          </TableBody>
        </Table>
      </Card>
    </div>
  </div>
</template>
