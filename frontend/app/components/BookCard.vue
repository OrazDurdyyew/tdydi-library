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

const props = defineProps<{
  book: Book
  /** Показывать ленту-закладку (для новинок, рекомендованных и т.п.) */
  ribbon?: boolean
}>()

const fileUrl = useFileUrl()
const coverUrl = computed(() => fileUrl(props.book.cover_path))
</script>

<template>
  <NuxtLinkLocale
    :to="`/books/${book.slug}`"
    class="group block"
  >
    <!-- Обложка-том: 3D-объём, лента-закладка опционально -->
    <div
      :class="[
        'book-card relative aspect-[2/3] overflow-hidden mb-4',
        ribbon ? 'ribbon ribbon-flat' : ''
      ]"
    >
      <img
        v-if="coverUrl"
        :src="coverUrl"
        :alt="book.title"
        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-[1.02]"
        loading="lazy"
      />
      <div
        v-else
        class="w-full h-full flex flex-col items-center justify-center bg-gradient-to-b from-muted to-muted/60 text-primary/60"
      >
        <BookOpen class="w-12 h-12 mb-2" />
        <span class="serif italic text-xs text-muted-foreground px-3 text-center line-clamp-2">
          {{ book.title }}
        </span>
      </div>

      <!-- "Тиснение" корешка справа: имитация торца страниц -->
      <span
        aria-hidden="true"
        class="absolute top-1 bottom-1 right-0 w-[3px] bg-gradient-to-l from-black/15 via-black/5 to-transparent pointer-events-none"
      />
    </div>

    <!-- Подпись каталога -->
    <div class="px-1 text-center">
      <h3 class="font-serif font-semibold text-sm leading-snug line-clamp-2 text-foreground group-hover:text-primary transition-colors">
        {{ book.title }}
      </h3>
      <p v-if="book.authors" class="text-xs text-muted-foreground italic mt-1 line-clamp-1">
        {{ book.authors }}
      </p>
      <p v-if="book.year" class="text-[10px] uppercase tracking-widest text-muted-foreground/70 mt-1.5">
        anno {{ book.year }}
      </p>
    </div>
  </NuxtLinkLocale>
</template>
