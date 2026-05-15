<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const api = useApi()
const { t } = useI18n()

const LIMIT = 12
const page = ref(Number(route.query.page) || 1)

const { data, pending } = await useAsyncData(
  'news-list',
  () =>
    api
      .get<{ data: any[]; total: number }>(`/posts?page=${page.value}&limit=${LIMIT}`)
      .catch(() => ({ data: [], total: 0 })),
  { watch: [page] }
)

const posts = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.total ?? 0)

function changePage(p: number) {
  page.value = p
  router.replace({ query: p > 1 ? { page: String(p) } : {} })
  if (import.meta.client) window.scrollTo({ top: 0, behavior: 'smooth' })
}

useSeoMeta({
  title: () => t('news.title'),
  description: () => t('home.latest_news')
})
</script>

<template>
  <div>
    <PageHero :title="t('news.title')">
      <template #breadcrumbs>
        <Breadcrumbs :items="[{ label: t('nav.home'), to: '/' }, { label: t('news.title') }]" />
      </template>
    </PageHero>

    <div class="page-container py-10">
      <div v-if="pending" class="text-muted-foreground py-16 text-center">
        {{ t('common.loading') }}
      </div>

      <div
        v-else-if="posts.length"
        class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5"
      >
        <PostCard v-for="post in posts" :key="post.id" :post="post" />
      </div>

      <div v-else class="card p-16 text-center text-muted-foreground">
        {{ t('news.no_news') }}
      </div>

      <Pagination :page="page" :total="total" :limit="LIMIT" @change="changePage" />
    </div>
  </div>
</template>
