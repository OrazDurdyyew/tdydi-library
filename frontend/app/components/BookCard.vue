<script setup lang="ts">
import { BookOpen } from 'lucide-vue-next'

interface Book {
  id: number
  slug: string
  title: string
  authors?: string
  year?: number
  cover_path?: string
}

const props = defineProps<{ book: Book }>()

const fileUrl = useFileUrl()
const coverUrl = computed(() => fileUrl(props.book.cover_path))
</script>

<template>
  <NuxtLinkLocale :to="`/books/${book.slug}`" class="group block">
    <div class="book-card aspect-[3/4]">
      <img
        v-if="coverUrl"
        :src="coverUrl"
        :alt="book.title"
        class="w-full h-full object-cover"
        loading="lazy"
      />
      <div
        v-else
        class="w-full h-full flex flex-col items-center justify-center gap-2 p-4 text-muted-foreground"
      >
        <BookOpen class="w-8 h-8" />
        <span class="text-xs text-center line-clamp-3 leading-snug">{{ book.title }}</span>
      </div>
    </div>

    <div class="mt-3 px-0.5">
      <h3 class="text-sm font-semibold leading-snug line-clamp-2 text-foreground group-hover:text-secondary transition-colors">
        {{ book.title }}
      </h3>
      <p v-if="book.authors" class="text-xs text-muted-foreground mt-1 line-clamp-1">
        {{ book.authors }}
      </p>
      <p v-if="book.year" class="text-xs text-muted-foreground/80 mt-0.5">
        {{ book.year }}
      </p>
    </div>
  </NuxtLinkLocale>
</template>
