<script setup lang="ts">
import { Menu, X, Search, BookMarked } from 'lucide-vue-next'

const settings = useSiteSettings()
const { locale, t } = useI18n()
const fileUrl = useFileUrl()
const colorMode = useColorMode()
const localePath = useLocalePath()
const loc = useLocalized()

const siteName = computed(() => loc(settings.value, 'name') || t('site.name'))

const logoUrl = computed(() => {
  if (!settings.value) return ''
  const path =
    colorMode.value === 'dark' && settings.value.logo_dark_path
      ? settings.value.logo_dark_path
      : settings.value.logo_light_path
  return fileUrl(path)
})

const query = ref('')
function submitSearch() {
  const q = query.value.trim()
  if (!q) return
  navigateTo({ path: localePath('/search'), query: { q } })
  mobileOpen.value = false
}

const mobileOpen = ref(false)
const route = useRoute()
watch(
  () => route.fullPath,
  () => {
    mobileOpen.value = false
  }
)
</script>

<template>
  <header class="sticky top-0 z-50 border-b border-border bg-card/85 backdrop-blur-md">
    <!-- Верхняя строка: логотип · поиск · контролы -->
    <div class="page-container">
      <div class="h-16 flex items-center gap-4">
        <!-- Логотип + название -->
        <NuxtLinkLocale to="/" class="flex items-center gap-2.5 shrink-0">
          <img
            v-if="logoUrl"
            :src="logoUrl"
            :alt="siteName"
            class="h-8 w-auto object-contain"
          />
          <span
            v-else
            class="inline-flex items-center justify-center w-9 h-9 rounded-lg bg-primary text-primary-foreground"
          >
            <BookMarked class="w-4.5 h-4.5" />
          </span>
          <span class="font-semibold text-[0.95rem] tracking-tight hidden sm:block">
            {{ siteName }}
          </span>
        </NuxtLinkLocale>

        <!-- Глобальный поиск (десктоп) -->
        <form
          class="hidden md:flex flex-1 max-w-xl mx-auto items-center h-10 px-3.5 rounded-lg border border-border bg-background focus-within:border-ring focus-within:ring-2 focus-within:ring-ring/15 transition-colors"
          @submit.prevent="submitSearch"
        >
          <Search class="w-4 h-4 text-muted-foreground shrink-0" />
          <input
            v-model="query"
            type="search"
            :placeholder="t('search.placeholder')"
            class="flex-1 h-full bg-transparent border-0 outline-none px-3 text-sm placeholder:text-muted-foreground"
          />
        </form>

        <!-- Контролы -->
        <div class="flex items-center gap-1 shrink-0 ml-auto md:ml-0">
          <LangSwitcher />
          <ThemeToggle />
          <button
            type="button"
            class="md:hidden inline-flex items-center justify-center w-9 h-9 rounded-lg text-foreground hover:bg-muted transition-colors"
            :aria-expanded="mobileOpen"
            :aria-label="mobileOpen ? t('common.close') : 'Menu'"
            @click="mobileOpen = !mobileOpen"
          >
            <X v-if="mobileOpen" class="w-5 h-5" />
            <Menu v-else class="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>

    <!-- Нижняя строка: навигация по разделам (десктоп) -->
    <div class="hidden md:block border-t border-border">
      <div class="page-container">
        <SectionsMenu />
      </div>
    </div>

    <!-- Мобильная панель -->
    <Transition name="drawer">
      <div v-if="mobileOpen" class="md:hidden border-t border-border bg-card">
        <div class="page-container py-4 space-y-4">
          <form
            class="flex items-center h-10 px-3.5 rounded-lg border border-border bg-background"
            @submit.prevent="submitSearch"
          >
            <Search class="w-4 h-4 text-muted-foreground shrink-0" />
            <input
              v-model="query"
              type="search"
              :placeholder="t('search.placeholder')"
              class="flex-1 h-full bg-transparent border-0 outline-none px-3 text-sm placeholder:text-muted-foreground"
            />
          </form>
          <SectionsMenu mobile />
        </div>
      </div>
    </Transition>
  </header>
</template>

<style scoped>
.drawer-enter-active,
.drawer-leave-active {
  transition: opacity 180ms ease, transform 180ms ease;
}
.drawer-enter-from,
.drawer-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
