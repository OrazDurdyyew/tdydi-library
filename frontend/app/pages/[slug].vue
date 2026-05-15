<script setup lang="ts">
const route = useRoute()
const api = useApi()
const { t } = useI18n()
const loc = useLocalized()

const slug = computed(() => route.params.slug as string)

const { data: page, error } = await useAsyncData(
  () => `page-${slug.value}`,
  () => api.get<any>(`/pages/slug/${slug.value}`)
)

if (error.value || !page.value) {
  throw createError({ statusCode: 404, statusMessage: 'Page not found', fatal: true })
}

const title = computed(() => loc(page.value, 'title'))
const content = computed(() => loc(page.value, 'content'))

useSeoMeta({
  title: () => title.value,
  description: () => (content.value || '').replace(/<[^>]+>/g, '').slice(0, 160)
})
</script>

<template>
  <article v-if="page" class="page-container py-10">
    <div class="max-w-3xl">
      <Breadcrumbs :items="[{ label: t('nav.home'), to: '/' }, { label: title }]" />

      <h1 class="text-3xl md:text-4xl font-bold tracking-tight leading-tight mb-8">{{ title }}</h1>

      <div
        v-if="content"
        class="prose-content serif text-[1.08rem] leading-[1.8] text-foreground/90"
        v-html="content"
      />
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
