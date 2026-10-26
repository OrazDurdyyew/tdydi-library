<script setup lang="ts">
import { Newspaper, BookText } from 'lucide-vue-next'

interface Periodical {
  id: number
  slug: string
  type: 'newspaper' | 'magazine'
  cover_path?: string
  periodicity?: string
  [key: string]: any
}

const props = defineProps<{ periodical: Periodical }>()

const fileUrl = useFileUrl()
const loc = useLocalized()
const { t } = useI18n()

const coverUrl = computed(() => fileUrl(props.periodical.cover_path))
const name = computed(() => loc(props.periodical, 'name'))
const typeLabel = computed(() =>
  props.periodical.type === 'magazine' ? t('periodicals.magazine') : t('periodicals.newspaper')
)
</script>

<template>
  <NuxtLinkLocale :to="`/periodicals/${periodical.slug}`" class="group block">
    <div class="book-card aspect-[350/145] bg-white overflow-hidden flex items-center justify-center px-2">
      <img
        v-if="coverUrl"
        :src="coverUrl"
        :alt="name"
        class="max-w-full max-h-full object-contain"
        loading="lazy"
      />
      <div v-else class="flex items-center justify-center text-muted-foreground">
        <Newspaper v-if="periodical.type === 'newspaper'" class="w-9 h-9" />
        <BookText v-else class="w-9 h-9" />
      </div>
    </div>
    <div class="mt-3 px-0.5">
      <span class="eyebrow text-[0.65rem]">{{ typeLabel }}</span>
      <h3 class="text-sm font-semibold leading-snug mt-0.5 text-foreground group-hover:text-secondary transition-colors line-clamp-2">
        {{ name }}
      </h3>
      <p v-if="periodical.periodicity" class="text-xs text-muted-foreground mt-1">
        {{ periodical.periodicity }}
      </p>
    </div>
  </NuxtLinkLocale>
</template>
