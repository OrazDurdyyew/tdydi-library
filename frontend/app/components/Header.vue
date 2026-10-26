<script setup lang="ts">
import { Menu, X, BookMarked } from 'lucide-vue-next'

const settings = useSiteSettings()
const { t } = useI18n()
const fileUrl = useFileUrl()
const colorMode = useColorMode()
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
        <NuxtLinkLocale to="/" class="flex items-center gap-3 shrink-0 min-w-0 max-w-[50%] md:max-w-[28%]">
          <img
            v-if="logoUrl"
            :src="logoUrl"
            :alt="siteName"
            class="h-11 w-auto object-contain shrink-0"
          />
          <span
            v-else
            class="inline-flex items-center justify-center w-11 h-11 rounded-lg bg-primary text-primary-foreground shrink-0"
          >
            <BookMarked class="w-5 h-5" />
          </span>
          <span class="font-semibold text-[0.95rem] leading-tight tracking-tight hidden sm:block line-clamp-2">
            {{ siteName }}
          </span>
        </NuxtLinkLocale>

        <!-- Глобальный поиск (десктоп) с live-подсказками -->
        <div class="hidden md:block flex-1 max-w-xl mx-auto">
          <SearchBox variant="desktop" />
        </div>

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
          <SearchBox variant="mobile" @navigated="mobileOpen = false" />
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
