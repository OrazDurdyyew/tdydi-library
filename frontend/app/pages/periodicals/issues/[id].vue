<script setup lang="ts">
definePageMeta({ layout: false })

const route = useRoute()
const config = useRuntimeConfig()
const api = useApi()
const { t } = useI18n()

const id = computed(() => route.params.id as string)

const { data: issue, error } = await useAsyncData(
  () => `read-issue-${id.value}`,
  () => api.get<any>(`/periodicals/issues/${id.value}`)
)

if (error.value || !issue.value) {
  throw createError({ statusCode: 404, statusMessage: 'Issue not found', fatal: true })
}

const { data: periodical } = await useAsyncData(
  () => `read-issue-periodical-${issue.value?.periodical_id}`,
  () =>
    issue.value
      ? api.get<any>(`/periodicals/${issue.value.periodical_id}`).catch(() => null)
      : Promise.resolve(null)
)

const fileSrc = computed(() => `${config.public.apiBase}/periodicals/issues/${id.value}/file`)
const downloadUrl = computed(
  () => `${config.public.apiBase}/periodicals/issues/${id.value}/download`
)
const backTo = computed(() =>
  periodical.value ? `/periodicals/${periodical.value.slug}` : '/periodicals'
)
const title = computed(
  () => issue.value?.issue_number || t('periodicals.issue')
)

useHead({ title: () => `${title.value} — ${t('periodicals.read')}` })
</script>

<template>
  <PdfViewer
    v-if="issue"
    :title="title"
    :src="fileSrc"
    :storage-key="`issue-bookmark-${issue.id}`"
    :back-to="backTo"
    :download-url="downloadUrl"
  />
</template>
