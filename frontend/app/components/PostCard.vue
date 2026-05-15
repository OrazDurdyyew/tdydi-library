<script setup lang="ts">
import { Newspaper } from 'lucide-vue-next'

interface Post {
  id: number
  slug: string
  cover_path?: string
  author?: string
  published_at?: string
  [key: string]: any
}

const props = defineProps<{ post: Post }>()

const fileUrl = useFileUrl()
const loc = useLocalized()
const { locale } = useI18n()

const coverUrl = computed(() => fileUrl(props.post.cover_path))
const title = computed(() => loc(props.post, 'title'))
const date = computed(() => formatDate(props.post.published_at, locale.value))
</script>

<template>
  <NuxtLinkLocale
    :to="`/news/${post.slug}`"
    class="group flex flex-col card card-hover overflow-hidden"
  >
    <div class="aspect-[16/9] bg-muted overflow-hidden">
      <img
        v-if="coverUrl"
        :src="coverUrl"
        :alt="title"
        class="w-full h-full object-cover"
        loading="lazy"
      />
      <div v-else class="w-full h-full flex items-center justify-center text-muted-foreground">
        <Newspaper class="w-8 h-8" />
      </div>
    </div>
    <div class="p-5 flex flex-col flex-1">
      <p v-if="date" class="text-xs text-muted-foreground">{{ date }}</p>
      <h3 class="font-semibold text-base leading-snug mt-1.5 text-foreground group-hover:text-secondary transition-colors line-clamp-2">
        {{ title }}
      </h3>
      <p v-if="post.author" class="text-sm text-muted-foreground mt-auto pt-3">
        {{ post.author }}
      </p>
    </div>
  </NuxtLinkLocale>
</template>
