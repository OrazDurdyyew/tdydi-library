<script setup lang="ts">
import { Check, Image as ImageIcon, Plus } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const route = useRoute()
const api = useApi()
const toast = useToast()
const { ask } = useConfirm()
const fileUrl = useFileUrl()
const id = Number(route.params.id)

const { data: periodical, refresh } = await useAsyncData(`periodical-${id}`, () =>
  api.get<any>(`/periodicals/${id}`)
)
const ISSUES_LIMIT = 20
const issuesPage = ref(1)
const { data: issuesResp, refresh: refreshIssues } = await useAsyncData(
  `periodical-issues-${id}`,
  () =>
    api.get<{ data: any[]; total: number }>(
      `/periodicals/${id}/issues?page=${issuesPage.value}&limit=${ISSUES_LIMIT}`
    ),
  { watch: [issuesPage] }
)

const issues = computed(() => issuesResp.value?.data ?? [])
const issuesTotal = computed(() => issuesResp.value?.total ?? 0)

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

useAutoSlug(toRef(form, 'name_tm'), toRef(form, 'slug'))

const saving = ref(false)
const coverInput = ref<HTMLInputElement | null>(null)
const uploading = ref(false)

async function submit() {
  saving.value = true
  try {
    await api.put(`/admin/periodicals/${id}`, form)
    await refresh()
    toast.success(t('admin.common.saved'))
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error'))
  } finally {
    saving.value = false
  }
}

async function uploadCover() {
  const file = coverInput.value?.files?.[0]
  if (!file) return
  uploading.value = true
  try {
    await api.uploadField(`/admin/periodicals/${id}/cover`, 'cover', file)
    await refresh()
    toast.success(t('admin.periodicals.cover_uploaded'))
    if (coverInput.value) coverInput.value.value = ''
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error_upload'))
  } finally {
    uploading.value = false
  }
}

const newIssue = reactive({ issue_number: '', issue_date: '', annotation_ru: '' })
const addingIssue = ref(false)
const issueFileInputs = ref<Record<number, HTMLInputElement | null>>({})
const issueCoverInputs = ref<Record<number, HTMLInputElement | null>>({})

async function createIssue() {
  addingIssue.value = true
  try {
    await api.post(`/admin/periodicals/${id}/issues`, {
      periodical_id: id,
      ...newIssue,
    })
    await refreshIssues()
    toast.success(t('admin.periodicals.issue_added'))
    newIssue.issue_number = ''
    newIssue.issue_date = ''
    newIssue.annotation_ru = ''
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error'))
  } finally {
    addingIssue.value = false
  }
}

async function deleteIssue(issueId: number, num: string) {
  const ok = await ask({
    title: t('admin.periodicals.confirm_delete_issue_title'),
    description: t('admin.periodicals.confirm_delete_issue_description', { num: num || '#' + issueId }),
    confirmText: t('admin.common.delete'),
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/periodicals/issues/${issueId}`)
    await refreshIssues()
    toast.success(t('admin.periodicals.issue_deleted'))
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error'))
  }
}

async function uploadIssueFile(issueId: number) {
  const file = issueFileInputs.value[issueId]?.files?.[0]
  if (!file) return
  try {
    await api.uploadField(`/admin/periodicals/issues/${issueId}/file`, 'file', file)
    await refreshIssues()
    toast.success(t('admin.books.file_uploaded_cover_auto'))
    if (issueFileInputs.value[issueId]) issueFileInputs.value[issueId]!.value = ''
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error_upload'))
  }
}

async function uploadIssueCover(issueId: number) {
  const file = issueCoverInputs.value[issueId]?.files?.[0]
  if (!file) return
  try {
    await api.uploadField(`/admin/periodicals/issues/${issueId}/cover`, 'cover', file)
    await refreshIssues()
    toast.success(t('admin.periodicals.cover_uploaded'))
    if (issueCoverInputs.value[issueId]) issueCoverInputs.value[issueId]!.value = ''
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error_upload'))
  }
}

useHead({ title: () => `${t('admin.periodicals.edit_title')}: ${periodical.value?.name_ru || periodical.value?.name_tm}` })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/periodicals" class="text-muted-foreground hover:text-foreground text-sm">← {{ t('admin.periodicals.title') }}</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">{{ periodical?.name_ru || periodical?.name_tm }}</h1>
    </div>

    <Card class="mb-6 max-w-3xl">
      <CardHeader>
        <CardTitle class="text-lg">{{ t('admin.common.cover') }}</CardTitle>
        <CardDescription v-if="periodical?.cover_path" class="font-mono">{{ periodical.cover_path }}</CardDescription>
        <CardDescription v-else>{{ t('admin.periodicals.cover_not_uploaded') }}</CardDescription>
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
            {{ uploading ? t('admin.common.uploading') : t('admin.common.upload') }}
          </Button>
        </div>
      </CardContent>
    </Card>

    <form @submit.prevent="submit" class="space-y-6 max-w-3xl mb-10">
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">{{ t('admin.sections.main') }}</CardTitle>
        </CardHeader>
        <CardContent class="space-y-5">
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <Label for="type">{{ t('admin.periodicals.type') }}</Label>
              <Select id="type" v-model="form.type">
                <option value="magazine">{{ t('admin.periodicals.type_magazine') }}</option>
                <option value="newspaper">{{ t('admin.periodicals.type_newspaper') }}</option>
              </Select>
            </div>
            <div class="space-y-2">
              <Label for="slug">{{ t('admin.common.slug') }}</Label>
              <Input id="slug" v-model="form.slug" readonly tabindex="-1" class="font-mono text-sm bg-muted text-muted-foreground cursor-not-allowed" />
            </div>
          </div>
          <div class="space-y-2">
            <Label for="period">{{ t('admin.periodicals.periodicity') }}</Label>
            <Input id="period" v-model="form.periodicity" />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-lg">{{ t('admin.sections.names_3lang') }}</CardTitle>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="name-tm">{{ t('admin.common.lang_tm') }}</Label>
            <Input id="name-tm" v-model="form.name_tm" />
          </div>
          <div class="space-y-2">
            <Label for="name-ru">{{ t('admin.common.lang_ru') }}</Label>
            <Input id="name-ru" v-model="form.name_ru" />
          </div>
          <div class="space-y-2">
            <Label for="name-en">{{ t('admin.common.lang_en') }}</Label>
            <Input id="name-en" v-model="form.name_en" />
          </div>
        </CardContent>
      </Card>

      <div class="flex gap-3">
        <Button type="submit" :disabled="saving">
          {{ saving ? t('admin.common.saving') : t('admin.common.save') }}
        </Button>
        <NuxtLink to="/admin/periodicals">
          <Button type="button" variant="outline">{{ t('admin.common.back') }}</Button>
        </NuxtLink>
      </div>
    </form>

    <div class="max-w-3xl max-w-3xl lg:max-w-none">
      <h2 class="text-xl font-bold font-serif mb-4">{{ t('admin.periodicals.issues') }}</h2>

      <Card class="mb-4">
        <CardHeader>
          <CardTitle class="text-base">{{ t('admin.periodicals.add_issue') }}</CardTitle>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-2">
              <Label for="issue-num">{{ t('admin.periodicals.issue_number') }}</Label>
              <Input id="issue-num" v-model="newIssue.issue_number" placeholder="№ 1" />
            </div>
            <div class="space-y-2">
              <Label for="issue-date">{{ t('admin.periodicals.issue_date') }}</Label>
              <Input id="issue-date" v-model="newIssue.issue_date" type="date" />
            </div>
          </div>
          <Button @click="createIssue" :disabled="addingIssue" size="sm">
            <Plus class="w-4 h-4" /> {{ addingIssue ? t('admin.common.saving') : t('admin.common.add') }}
          </Button>
        </CardContent>
      </Card>

      <Card class="p-0 overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>{{ t('admin.common.cover') }}</TableHead>
              <TableHead>{{ t('admin.periodicals.issue_number') }}</TableHead>
              <TableHead>{{ t('admin.periodicals.issue_date') }}</TableHead>
              <TableHead>{{ t('admin.periodicals.file') }}</TableHead>
              <TableHead class="text-right">{{ t('admin.common.actions') }}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="issue in issues" :key="issue.id">
              <TableCell>
                <div class="flex items-center gap-2">
                  <div class="w-10 h-14 shrink-0 rounded bg-muted overflow-hidden border border-border flex items-center justify-center">
                    <img
                      v-if="issue.cover_path"
                      :src="fileUrl(issue.cover_path)"
                      :alt="issue.issue_number"
                      class="w-full h-full object-cover"
                    />
                    <ImageIcon v-else class="w-4 h-4 text-muted-foreground" />
                  </div>
                  <div class="flex flex-col gap-1">
                    <input
                      :ref="(el) => { if (el) issueCoverInputs[issue.id] = el as HTMLInputElement }"
                      type="file"
                      accept="image/*"
                      class="text-xs w-32 file:mr-2 file:rounded file:border-0 file:bg-muted file:px-2 file:py-1 file:text-xs"
                    />
                    <Button @click="uploadIssueCover(issue.id)" variant="outline" size="sm" class="w-fit">
                      {{ issue.cover_path ? t('admin.periodicals.file_replace') : t('admin.common.upload') }}
                    </Button>
                  </div>
                </div>
              </TableCell>
              <TableCell class="font-medium">{{ issue.issue_number || '—' }}</TableCell>
              <TableCell class="text-muted-foreground">{{ issue.issue_date || '—' }}</TableCell>
              <TableCell>
                <div class="flex items-center gap-2">
                  <span
                    v-if="issue.has_file"
                    class="inline-flex items-center gap-1.5 text-xs font-medium text-green-600"
                  >
                    <Check class="w-3.5 h-3.5" /> {{ t('admin.periodicals.file_uploaded') }}
                  </span>
                  <input
                    :ref="(el) => { if (el) issueFileInputs[issue.id] = el as HTMLInputElement }"
                    type="file"
                    accept=".pdf"
                    class="text-xs flex-1 file:mr-2 file:rounded file:border-0 file:bg-muted file:px-2 file:py-1 file:text-xs"
                  />
                  <Button @click="uploadIssueFile(issue.id)" variant="outline" size="sm">
                    {{ issue.has_file ? t('admin.periodicals.file_replace') : t('admin.common.upload') }}
                  </Button>
                </div>
              </TableCell>
              <TableCell class="text-right">
                <Button
                  variant="ghost"
                  size="sm"
                  class="text-destructive hover:text-destructive"
                  @click="deleteIssue(issue.id, issue.issue_number)"
                >
                  {{ t('admin.common.delete') }}
                </Button>
              </TableCell>
            </TableRow>
            <TableEmpty v-if="issues.length === 0" :colspan="5">{{ t('admin.periodicals.no_issues') }}</TableEmpty>
          </TableBody>
        </Table>
      </Card>

      <Pagination
        :page="issuesPage"
        :total="issuesTotal"
        :limit="ISSUES_LIMIT"
        @change="(p) => (issuesPage = p)"
      />
    </div>
  </div>
</template>
