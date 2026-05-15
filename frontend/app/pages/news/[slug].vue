<script setup lang="ts">
import { ArrowLeft, Calendar, User } from 'lucide-vue-next'

const route = useRoute()
const api = useApi()
const { t, locale } = useI18n()
const loc = useLocalized()
const fileUrl = useFileUrl()

const slug = computed(() => route.params.slug as string)

const { data: post, error } = await useAsyncData(
  () => `post-${slug.value}`,
  () => api.get<any>(`/posts/slug/${slug.value}`)
)

if (error.value || !post.value) {
  throw createError({ statusCode: 404, statusMessage: 'Post not found', fatal: true })
}

const title = computed(() => loc(post.value, 'title'))
const content = computed(() => loc(post.value, 'content'))
const coverUrl = computed(() => fileUrl(post.value?.cover_path))
const date = computed(() => formatDate(post.value?.published_at, locale.value))

useSeoMeta({
  title: () => title.value,
  description: () => (content.value || '').replace(/<[^>]+>/g, '').slice(0, 160),
  ogImage: () => coverUrl.value || undefined
})
</script>

<template>
  <article v-if="post" class="page-container py-10">
    <div class="max-w-3xl">
      <Breadcrumbs
        :items="[
          { label: t('nav.home'), to: '/' },
          { label: t('news.title'), to: '/news' },
          { label: title }
        ]"
      />

      <h1 class="text-3xl md:text-4xl font-bold tracking-tight leading-tight mt-1 mb-4">
        {{ title }}
      </h1>

      <div class="flex flex-wrap items-center gap-5 text-sm text-muted-foreground mb-8">
        <span v-if="date" class="inline-flex items-center gap-1.5">
          <Calendar class="w-4 h-4 icon-primary" /> {{ date }}
        </span>
        <span v-if="post.author" class="inline-flex items-center gap-1.5">
          <User class="w-4 h-4 icon-primary" /> {{ post.author }}
        </span>
      </div>

      <div v-if="coverUrl" class="rounded-xl overflow-hidden border border-border mb-10">
        <img :src="coverUrl" :alt="title" class="w-full object-cover" />
      </div>

      <div
        v-if="content"
        class="prose-content serif text-[1.08rem] leading-[1.8] text-foreground/90"
        v-html="content"
      />

      <div class="divider-double mt-12 mb-6" />

      <NuxtLinkLocale to="/news" class="btn-ghost">
        <ArrowLeft class="w-4 h-4" /> {{ t('news.back_to_news') }}
      </NuxtLinkLocale>
    </div>
  </article>
</template>

<style scoped>
.prose-content :deep(h2),
.prose-content :deep(h3) {
  font-family: var(--font-serif);
  font-weight: 700;
  margin: 1.75rem 0 0.75rem;
}
.prose-content :deep(h2) {
  font-size: 1.5rem;
}
.prose-content :deep(h3) {
  font-size: 1.25rem;
}
.prose-content :deep(p) {
  margin: 0.85rem 0;
}
.prose-content :deep(ul),
.prose-content :deep(ol) {
  margin: 0.85rem 0;
  padding-left: 1.5rem;
  list-style: revert;
}
.prose-content :deep(a) {
  color: rgb(var(--primary));
  text-decoration: underline;
}
.prose-content :deep(img) {
  border-radius: 0.75rem;
  margin: 1.25rem 0;
}
.prose-content :deep(blockquote) {
  border-left: 3px solid rgb(var(--primary));
  padding-left: 1rem;
  font-style: italic;
  color: rgb(var(--muted-foreground));
  margin: 1.25rem 0;
}
</style>
