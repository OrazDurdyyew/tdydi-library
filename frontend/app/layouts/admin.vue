<script setup lang="ts">
import { LayoutDashboard, BookMarked, FolderTree, LogOut } from 'lucide-vue-next'

const { admin, fetchMe, logout } = useAuth()
const localePath = useLocalePath()

if (import.meta.client && !admin.value) {
  await fetchMe()
}

async function doLogout() {
  logout()
  await navigateTo(localePath('/admin/login'))
}

const navItems = [
  { to: '/admin', label: 'Дашборд', icon: LayoutDashboard, exact: true },
  { to: '/admin/books', label: 'Книги', icon: BookMarked },
  { to: '/admin/sections', label: 'Разделы', icon: FolderTree }
]
</script>

<template>
  <div class="min-h-screen flex">
    <aside class="app-sidebar w-64 flex flex-col shrink-0 relative">
      <!-- Тиснёная "табличка" -->
      <div class="px-5 py-6 relative">
        <div class="flex items-center gap-3">
          <span class="inline-flex items-center justify-center w-11 h-11 rounded-sm bg-primary text-primary-foreground border border-black/20 shadow-paper">
            <BookMarked class="w-5 h-5" />
          </span>
          <div>
            <div class="font-serif font-bold text-base sidebar-title leading-tight">
              Bibliotheca
            </div>
            <div class="small-caps text-[10px] sidebar-muted mt-0.5">
              Administratio
            </div>
          </div>
        </div>
        <!-- Двойная разделительная линия в духе титульного листа -->
        <div class="mt-5 h-1 border-t border-b" style="border-color: rgba(252, 250, 245, 0.18);" />
      </div>

      <!-- Полка с корешками -->
      <nav class="flex-1 px-3 space-y-1 pb-4">
        <NuxtLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          :exact-active-class="item.exact ? 'router-link-exact-active' : ''"
          class="book-spine"
        >
          <component :is="item.icon" class="w-4 h-4 shrink-0 opacity-80" />
          <span>{{ item.label }}</span>
        </NuxtLink>
      </nav>

      <!-- Шильдик пользователя -->
      <div v-if="admin" class="px-4 py-4 mt-auto" style="border-top: 1px solid rgba(252, 250, 245, 0.15);">
        <div class="flex items-center gap-3">
          <span class="inline-flex items-center justify-center w-10 h-10 rounded-full bg-primary/85 text-primary-foreground font-serif font-bold text-base border border-black/20 shadow-paper">
            {{ admin.username.charAt(0).toUpperCase() }}
          </span>
          <div class="min-w-0 flex-1">
            <div class="font-serif text-sm sidebar-title truncate">{{ admin.username }}</div>
            <div class="text-[10px] small-caps sidebar-muted">{{ admin.role }}</div>
          </div>
          <button
            @click="doLogout"
            class="p-2 rounded-sm transition-colors"
            style="color: rgba(252, 250, 245, 0.7);"
            title="Выйти"
            onmouseover="this.style.color='rgb(178,53,49)'; this.style.background='rgba(255,255,255,0.06)'"
            onmouseout="this.style.color='rgba(252,250,245,0.7)'; this.style.background='transparent'"
          >
            <LogOut class="w-4 h-4" />
          </button>
        </div>
      </div>
    </aside>

    <main class="flex-1 min-w-0 overflow-x-auto">
      <div class="p-8 max-w-7xl">
        <slot />
      </div>
    </main>
  </div>
</template>
