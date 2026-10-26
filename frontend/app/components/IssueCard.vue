<script setup lang="ts">
import { BookText, Calendar, Newspaper } from 'lucide-vue-next';

defineProps<{ issue: any }>()

const fileUrl = useFileUrl()
const loc = useLocalized()
const { locale, t } = useI18n()
</script>

<template>
  <NuxtLinkLocale
    :to="`/periodicals/issues/${issue.id}`"
    class="group flex flex-col card card-hover overflow-hidden"
  >
    <div class="aspect-[3/4] bg-muted overflow-hidden flex items-center justify-center p-2">
      <img
        v-if="issue.cover_path"
        :src="fileUrl(issue.cover_path)"
        :alt="issue.issue_number"
        class="max-w-full max-h-full object-contain"
        loading="lazy"
      />
      <div v-else class="w-full h-full flex items-center justify-center text-muted-foreground">
        <Newspaper v-if="issue.periodical?.type === 'newspaper'" class="w-8 h-8" />
        <BookText v-else class="w-8 h-8" />
      </div>
    </div>
    <div class="p-3.5 flex flex-col flex-1">
      <div v-if="issue.periodical" class="text-xs font-medium text-secondary truncate">
        {{ loc(issue.periodical, 'name') }}
      </div>
      <div class="font-semibold text-sm leading-snug mt-0.5 group-hover:text-secondary transition-colors line-clamp-2">
        {{ issue.issue_number || t('periodicals.issue') }}
      </div>
      <div class="text-xs text-muted-foreground mt-auto pt-2 inline-flex items-center gap-1">
        <Calendar class="w-3 h-3 shrink-0" /> {{ formatDate(issue.issue_date, locale) }}
      </div>
    </div>
  </NuxtLinkLocale>
</template>
