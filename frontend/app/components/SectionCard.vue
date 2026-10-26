<script setup lang="ts">
import { FolderOpen, ArrowRight } from 'lucide-vue-next'

interface Section {
  id: number
  slug: string
  image_path?: string
  book_count?: number
  [key: string]: any
}

const props = defineProps<{ section: Section }>()

const fileUrl = useFileUrl()
const loc = useLocalized()
const { t } = useI18n()

const imageUrl = computed(() => fileUrl(props.section.image_path))
const name = computed(() => loc(props.section, 'name'))
const description = computed(() => loc(props.section, 'description'))
</script>

<template>
  <NuxtLinkLocale
    :to="`/sections/${section.slug}`"
    class="group flex flex-col card card-hover overflow-hidden"
  >
    <div class="aspect-[16/9] bg-muted overflow-hidden">
      <img
        v-if="imageUrl"
        :src="imageUrl"
        :alt="name"
        class="w-full h-full object-cover"
        loading="lazy"
      />
      <div v-else class="w-full h-full flex items-center justify-center text-muted-foreground">
        <FolderOpen class="w-8 h-8" />
      </div>
    </div>
    <div class="p-5 flex flex-col flex-1">
      <h3 class="font-semibold text-base leading-snug text-foreground group-hover:text-secondary transition-colors">
        {{ name }}
      </h3>
      <p v-if="description" class="text-sm text-muted-foreground mt-1.5 line-clamp-2 flex-1">
        {{ description }}
      </p>
      <div class="flex items-center justify-between mt-4 pt-3 border-t border-border">
        <span class="text-xs text-muted-foreground">
          {{ t('catalog.found', { count: section.book_count ?? 0 }) }}
        </span>
        <ArrowRight class="w-4 h-4 text-secondary opacity-0 group-hover:opacity-100 transition-opacity" />
      </div>
    </div>
  </NuxtLinkLocale>
</template>
