<script setup lang="ts">
import {
  BookMarked,
  BookOpen,
  FileText,
  FolderTree,
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

if (import.meta.client && !admin.value) {
  await fetchMe()
}

async function doLogout() {
  logout()
  await navigateTo(localePath('/admin/login'))
}

const isSuperadmin = computed(() => admin.value?.role === 'superadmin')

const navItems = computed(() => {
  const items = [
    { to: '/admin', label: 'Дашборд', icon: LayoutDashboard, exact: true, super: false },
    { to: '/admin/home', label: 'Главная', icon: LayoutTemplate, super: true },
    { to: '/admin/books', label: 'Книги', icon: BookMarked, super: false },
    { to: '/admin/sections', label: 'Разделы', icon: FolderTree, super: true },
    { to: '/admin/periodicals', label: 'Периодика', icon: Newspaper, super: true },
    { to: '/admin/collections', label: 'Подборки', icon: Library, super: true },
    { to: '/admin/posts', label: 'Новости', icon: BookOpen, super: true },
    { to: '/admin/pages', label: 'Страницы', icon: FileText, super: true },
    { to: '/admin/faq', label: 'FAQ', icon: HelpCircle, super: true },
    { to: '/admin/admins', label: 'Администраторы', icon: Users, super: true },
    { to: '/admin/settings', label: 'Настройки', icon: Settings, super: true },
  ]
  return items.filter((i) => !i.super || isSuperadmin.value)
})
</script>

<template>
  <div class="h-screen flex overflow-hidden">
    <aside class="app-sidebar w-60 flex flex-col shrink-0 relative h-screen">
      <!-- Логотип -->
      <div class="px-5 py-5 relative">
        <div class="flex items-center gap-3">
          <span class="inline-flex items-center justify-center w-10 h-10 rounded-sm bg-primary text-primary-foreground border border-black/20 shadow-paper">
            <BookMarked class="w-4.5 h-4.5" />
          </span>
          <div>
            <div class="font-serif font-bold text-sm sidebar-title leading-tight">Bibliotheca</div>
            <div class="small-caps text-[10px] sidebar-muted mt-0.5">Administratio</div>
          </div>
        </div>
        <div class="mt-4 h-px" style="background: rgba(252,250,245,0.15);" />
      </div>

      <!-- Навигация -->
      <nav class="flex-1 px-3 pb-4 space-y-0.5 overflow-y-auto">
        <NuxtLink
          v-for="item in navItems"
          :key="$localePath(item.to)"
          :to="$localePath(item.to)"
          :exact-active-class="item.exact ? 'router-link-exact-active' : ''"
          class="book-spine"
        >
          <component :is="item.icon" class="w-4 h-4 shrink-0 opacity-75" />
          <span class="text-[0.9rem]">{{ item.label }}</span>
        </NuxtLink>
      </nav>

      <!-- Пользователь -->
      <div v-if="admin" class="px-4 py-4" style="border-top: 1px solid rgba(252,250,245,0.15);">
        <div class="flex items-center gap-2.5">
          <span class="inline-flex items-center justify-center w-8 h-8 rounded-full bg-primary/80 text-primary-foreground font-serif font-bold text-sm border border-black/20 shrink-0">
            {{ admin.username.charAt(0).toUpperCase() }}
          </span>
          <div class="min-w-0 flex-1">
            <div class="font-serif text-sm sidebar-title truncate">{{ admin.username }}</div>
            <div class="text-[10px] small-caps sidebar-muted">{{ admin.role }}</div>
          </div>
          <button
            @click="doLogout"
            class="p-1.5 rounded-sm transition-colors hover:bg-white/10"
            style="color: rgba(252,250,245,0.65);"
            title="Выйти"
          >
            <LogOut class="w-4 h-4" />
          </button>
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
