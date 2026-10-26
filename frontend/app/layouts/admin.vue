<script setup lang="ts">
import {
  BookMarked,
  BookOpen,
  FileText,
  FolderTree,
  GraduationCap,
  HelpCircle,
  LayoutDashboard,
  LayoutTemplate,
  Library,
  LogOut,
  Newspaper,
  Settings,
  Users
} from 'lucide-vue-next'

const { admin, fetchMe, logout } = useAuth()
const localePath = useLocalePath()
const settings = useSiteSettings()
const fileUrl = useFileUrl()
const colorMode = useColorMode()
const loc = useLocalized()
const { t } = useI18n()

if (import.meta.client && !admin.value) {
  await fetchMe()
}

async function doLogout() {
  logout()
  await navigateTo(localePath('/admin/login'))
}

const isSuperadmin = computed(() => admin.value?.role === 'superadmin')
const isPublicationManager = computed(() => admin.value?.role === 'publication_manager')

const siteName = computed(() => loc(settings.value, 'name') || t('site.name'))
const appVersion = useRuntimeConfig().public.appVersion

const PUB_ROLES = ['superadmin', 'publication_manager']

// Все подстраницы раздела «Публикации преподавателей» — для подсветки одного пункта меню.
const PUB_PATHS = [
  '/admin/teacher-publications',
  '/admin/departments',
  '/admin/teachers',
  '/admin/publications',
  '/admin/publication-categories'
]

interface NavItem {
  to: string
  label: string
  icon: any
  exact?: boolean
  roles?: string[]
  match?: string[]
}

const navItems = computed<NavItem[]>(() => {
  const items: NavItem[] = [
    { to: '/admin', label: t('admin.nav.dashboard'), icon: LayoutDashboard, exact: true },
    { to: '/admin/home', label: t('admin.nav.home_page'), icon: LayoutTemplate, roles: ['superadmin'] },
    { to: '/admin/books', label: t('admin.nav.books'), icon: BookMarked },
    { to: '/admin/sections', label: t('admin.nav.sections'), icon: FolderTree, roles: ['superadmin'] },
    { to: '/admin/periodicals', label: t('admin.nav.periodicals'), icon: Newspaper, roles: ['superadmin'] },
    { to: '/admin/collections', label: t('admin.nav.collections'), icon: Library, roles: ['superadmin'] },
    { to: '/admin/teacher-publications', label: t('admin.nav.teacher_publications'), icon: GraduationCap, roles: PUB_ROLES, match: PUB_PATHS },
    { to: '/admin/posts', label: t('admin.nav.posts'), icon: BookOpen, roles: ['superadmin'] },
    { to: '/admin/pages', label: t('admin.nav.pages'), icon: FileText, roles: ['superadmin'] },
    { to: '/admin/faq', label: t('admin.nav.faq'), icon: HelpCircle, roles: ['superadmin'] },
    { to: '/admin/admins', label: t('admin.nav.admins'), icon: Users, roles: ['superadmin'] },
    { to: '/admin/settings', label: t('admin.nav.settings'), icon: Settings, roles: ['superadmin'] },
  ]
  const role = admin.value?.role
  return items.filter((i) => !i.roles || (role !== undefined && i.roles.includes(role)))
})

const route = useRoute()
function isItemActive(item: NavItem): boolean {
  if (!item.match) return false
  const path = route.path
  return item.match.some((m) => {
    const lp = localePath(m)
    return path === lp || path.startsWith(lp + '/')
  })
}
const logoUrl = computed(() => {
  if (!settings.value) return ''
  const path =
    colorMode.value === 'dark' && settings.value.logo_dark_path
      ? settings.value.logo_dark_path
      : settings.value.logo_light_path
  return fileUrl(path)
})
</script>

<template>
  <div class="h-screen flex overflow-hidden">
    <aside class="app-sidebar w-60 flex flex-col shrink-0 relative h-screen border-r border-border bg-card">
      <!-- Логотип -->
      <div class="px-5 py-5">
        <div class="flex items-center gap-3 min-w-0">
          <img
            v-if="logoUrl"
            :src="logoUrl"
            :alt="siteName"
            class="h-10 w-auto object-contain shrink-0"
          />
          <span
            v-else
            class="inline-flex items-center justify-center w-10 h-10 rounded-md bg-primary text-primary-foreground shrink-0"
          >
            <BookMarked class="w-5 h-5" />
          </span>
          <div class="min-w-0">
            <div class="font-semibold text-sm leading-tight text-foreground line-clamp-2">
              {{ siteName }}
            </div>
            <div class="text-[10px] uppercase tracking-wider text-muted-foreground mt-0.5">
              {{ t('admin.label') }}
            </div>
          </div>
        </div>
        <div class="mt-4 h-px bg-border" />
      </div>

      <!-- Навигация -->
      <nav class="flex-1 px-3 pb-4 space-y-0.5 overflow-y-auto">
        <NuxtLink
          v-for="item in navItems"
          :key="$localePath(item.to)"
          :to="$localePath(item.to)"
          :exact-active-class="item.exact ? 'router-link-exact-active' : ''"
          class="book-spine"
          :class="isItemActive(item) ? 'router-link-active' : ''"
        >
          <component :is="item.icon" class="w-4 h-4 shrink-0 opacity-75" />
          <span class="text-[0.9rem]">{{ item.label }}</span>
        </NuxtLink>
      </nav>

      <!-- Тема + язык -->
      <div class="px-4 py-3 border-t border-border flex items-center justify-between gap-2">
        <LangSwitcher />
        <ThemeToggle />
      </div>

      <!-- Пользователь -->
      <div v-if="admin" class="px-4 py-4 border-t border-border">
        <div class="flex items-center gap-2.5">
          <span class="inline-flex items-center justify-center w-8 h-8 rounded-full bg-primary text-primary-foreground font-semibold text-sm shrink-0">
            {{ admin.username.charAt(0).toUpperCase() }}
          </span>
          <div class="min-w-0 flex-1">
            <div class="text-sm text-foreground truncate">{{ admin.username }}</div>
            <div class="text-[10px] uppercase tracking-wider text-muted-foreground">{{ admin.role }}</div>
          </div>
          <button
            @click="doLogout"
            class="p-1.5 rounded-md transition-colors text-muted-foreground hover:text-foreground hover:bg-muted"
            :title="t('admin.logout')"
          >
            <LogOut class="w-4 h-4" />
          </button>
        </div>
        <div class="mt-2 text-[10px] text-muted-foreground/70 text-center tabular-nums">
          v{{ appVersion }}
        </div>
      </div>
    </aside>

    <!-- Контент — скроллится независимо от боковой панели -->
    <main class="admin-main flex-1 min-w-0 overflow-y-auto bg-background">
      <div class="p-8 max-w-7xl">
        <slot />
      </div>
    </main>
  </div>
</template>

<style scoped>
/* Админка — рабочий интерфейс: единый чистый sans-шрифт вместо декоративного Playfair.
   Шрифт управляется из «Настройки → Шрифт» (переменная --font-sans). */
.admin-main :deep(h1),
.admin-main :deep(h2),
.admin-main :deep(h3),
.admin-main :deep(h4),
.admin-main :deep(.font-serif),
.admin-main :deep(.serif) {
  font-family: var(--font-sans);
  letter-spacing: 0;
}
</style>
