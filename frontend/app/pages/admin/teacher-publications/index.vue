<script setup lang="ts">
import { GraduationCap, Users, FileText, Tag, ArrowRight } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const api = useApi()

interface PubStats {
  total_publications: number
  total_teachers: number
  total_departments: number
}

// Счётчики — только на клиенте: токен админа доступен в localStorage, не на SSR.
const { data: stats } = await useAsyncData(
  'admin-teacher-pub-hub-stats',
  () => api.get<PubStats>('/admin/stats/publications').catch(() => null),
  { server: false }
)

const sections = computed(() => [
  {
    to: '/admin/departments',
    icon: GraduationCap,
    title: t('admin.departments.title'),
    desc: t('admin.teacher_publications.departments_desc'),
    count: stats.value?.total_departments
  },
  {
    to: '/admin/teachers',
    icon: Users,
    title: t('admin.teachers.title'),
    desc: t('admin.teacher_publications.teachers_desc'),
    count: stats.value?.total_teachers
  },
  {
    to: '/admin/publications',
    icon: FileText,
    title: t('admin.publications.title'),
    desc: t('admin.teacher_publications.publications_desc'),
    count: stats.value?.total_publications
  },
  {
    to: '/admin/publication-categories',
    icon: Tag,
    title: t('admin.publication_categories.title'),
    desc: t('admin.teacher_publications.categories_desc'),
    count: undefined
  }
])

useHead({ title: () => t('admin.teacher_publications.title') })
</script>

<template>
  <div>
    <div class="mb-8">
      <h1 class="text-3xl font-bold font-serif">{{ t('admin.teacher_publications.title') }}</h1>
      <p class="text-muted-foreground mt-1">{{ t('admin.teacher_publications.subtitle') }}</p>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <NuxtLink
        v-for="s in sections"
        :key="s.to"
        :to="$localePath(s.to)"
        class="group card card-hover p-5 flex items-start gap-4"
      >
        <span class="inline-flex items-center justify-center w-12 h-12 rounded-lg bg-primary/10 text-primary shrink-0">
          <component :is="s.icon" class="w-5 h-5" />
        </span>
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2">
            <h2 class="font-semibold text-base">{{ s.title }}</h2>
            <span
              v-if="s.count !== undefined"
              class="text-xs font-medium px-2 py-0.5 rounded-full bg-muted text-muted-foreground"
            >
              {{ s.count }}
            </span>
          </div>
          <p class="text-sm text-muted-foreground mt-1">{{ s.desc }}</p>
        </div>
        <ArrowRight class="w-4 h-4 text-muted-foreground transition-transform group-hover:translate-x-0.5 shrink-0 mt-1" />
      </NuxtLink>
    </div>
  </div>
</template>
