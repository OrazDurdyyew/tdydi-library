<script setup lang="ts">
const api = useApi()
const { t } = useI18n()

const LIMIT = 12
const page = ref(1)
const { data, pending } = await useAsyncData(
  'collections-list',
  () =>
    api
      .get<{ data: any[]; total: number }>(`/collections?page=${page.value}&limit=${LIMIT}`)
      .catch(() => ({ data: [], total: 0 })),
  { watch: [page] }
)

const collections = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.total ?? 0)
function changePage(p: number) {
  page.value = p
  if (import.meta.client) window.scrollTo({ top: 0, behavior: 'smooth' })
}

useSeoMeta({
  title: () => t('collections.title'),
  description: () => t('collections.title')
})
</script>

<template>
  <div>
    <PageHero :title="t('collections.title')">
      <template #breadcrumbs>
        <Breadcrumbs
          :items="[{ label: t('nav.home'), to: '/' }, { label: t('collections.title') }]"
        />
      </template>
    </PageHero>

    <div class="page-container py-10">
      <div v-if="pending" class="text-muted-foreground py-16 text-center">
        {{ t('common.loading') }}
      </div>

      <div
        v-else-if="collections.length"
        class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5"
      >
        <CollectionCard v-for="col in collections" :key="col.id" :collection="col" />
      </div>

      <div v-else class="card p-16 text-center text-muted-foreground">
        {{ t('collections.no_collections') }}
      </div>

      <Pagination :page="page" :total="total" :limit="LIMIT" @change="changePage" />
    </div>
  </div>
</template>
