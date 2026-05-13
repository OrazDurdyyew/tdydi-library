<script setup lang="ts">
import { Menu, X } from 'lucide-vue-next'

const settings = useSiteSettings()
const { locale, t } = useI18n()
const fileUrl = useFileUrl()
const colorMode = useColorMode()

const siteName = computed(() => {
  if (!settings.value) return ''
  if (locale.value === 'tm') return settings.value.name_tm
  if (locale.value === 'en') return settings.value.name_en
  return settings.value.name_ru
})

const shortBrand = computed(() => (siteName.value || t('site.name') || 'Bibliotheca').toUpperCase())

const logoUrl = computed(() => {
  if (!settings.value) return ''
  const path = colorMode.value === 'dark' && settings.value.logo_dark_path
    ? settings.value.logo_dark_path
    : settings.value.logo_light_path
  return fileUrl(path)
})

const mobileOpen = ref(false)
const route = useRoute()
watch(() => route.fullPath, () => { mobileOpen.value = false })
</script>

<template>
  <!-- Угловой бренд-плакат сверху-справа (вертикальная пластина) -->
  <NuxtLinkLocale
    to="/"
    class="brand-plate"
    :aria-label="shortBrand"
  >
    <img
      v-if="logoUrl"
      :src="logoUrl"
      :alt="shortBrand"
      class="brand-plate-logo"
    />
    <span v-else class="brand-plate-text">{{ shortBrand }}</span>
  </NuxtLinkLocale>

  <!-- Плавающий floating-control сверху-слева: язык + тема -->
  <div class="floating-controls">
    <LangSwitcher />
    <span class="sep" aria-hidden="true">·</span>
    <ThemeToggle />
  </div>

  <!-- Центрированное pill-меню (sticky + blur) -->
  <header class="pill-nav-wrap">
    <nav class="pill-nav" :aria-label="t('nav.home')">
      <SectionsMenu />
      <button
        type="button"
        class="pill-burger md:hidden"
        :aria-expanded="mobileOpen"
        :aria-label="mobileOpen ? 'Закрыть меню' : 'Открыть меню'"
        @click="mobileOpen = !mobileOpen"
      >
        <X v-if="mobileOpen" class="w-5 h-5" />
        <Menu v-else class="w-5 h-5" />
      </button>
    </nav>

    <!-- Мобильная панель: раскрывается из-под pill-а -->
    <Transition name="mobile-drop">
      <div v-if="mobileOpen" class="pill-mobile">
        <SectionsMenu mobile />
      </div>
    </Transition>
  </header>

  <div class="app-header-spacer" aria-hidden="true" />
</template>

<style scoped>
/* ===== Бренд-плакат в углу ===== */
.brand-plate {
  position: fixed;
  top: 0;
  right: 1rem;
  z-index: 60;
  display: inline-flex;
  align-items: flex-start;
  justify-content: center;
  width: 70px;
  min-height: 130px;
  padding: 16px 8px 18px;
  background-color: rgb(var(--primary));
  color: rgb(var(--primary-foreground));
  box-shadow:
    0 6px 18px rgba(0, 0, 0, 0.18),
    inset 0 -2px 0 rgba(0, 0, 0, 0.15);
  text-decoration: none;
}
.brand-plate-text {
  font-family: var(--font-serif);
  font-weight: 700;
  font-size: 0.95rem;
  letter-spacing: 0.08em;
  line-height: 1.05;
  writing-mode: vertical-rl;
  transform: rotate(180deg);
  text-transform: uppercase;
  white-space: nowrap;
  text-shadow: 0 1px 0 rgba(0, 0, 0, 0.18);
}
.brand-plate-logo {
  max-width: 48px;
  max-height: 110px;
  object-fit: contain;
}

/* ===== Floating-controls сверху-слева ===== */
.floating-controls {
  position: fixed;
  top: 1.5rem;
  left: 1.5rem;
  z-index: 55;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.4rem 0.85rem;
  background-color: rgba(255, 253, 247, 0.85);
  backdrop-filter: saturate(160%) blur(10px);
  -webkit-backdrop-filter: saturate(160%) blur(10px);
  border-radius: 999px;
  border: 1px solid rgb(var(--border));
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.06);
}
:global(.dark) .floating-controls {
  background-color: rgba(47, 38, 28, 0.85);
}
.floating-controls .sep {
  color: rgb(var(--border));
  user-select: none;
}

/* ===== Центрированное pill-меню ===== */
.pill-nav-wrap {
  position: fixed;
  top: 1.25rem;
  left: 50%;
  transform: translateX(-50%);
  z-index: 50;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  width: max-content;
  max-width: calc(100vw - 2rem);
}

.pill-nav {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  height: 56px;
  padding: 0 0.5rem 0 1.25rem;
  border-radius: 999px;
  background-color: rgba(255, 253, 247, 0.92);
  backdrop-filter: saturate(180%) blur(14px);
  -webkit-backdrop-filter: saturate(180%) blur(14px);
  border: 1px solid rgb(var(--border));
  box-shadow:
    0 12px 32px rgba(0, 0, 0, 0.08),
    0 2px 6px rgba(0, 0, 0, 0.05),
    inset 0 1px 0 rgba(255, 255, 255, 0.5);
}
:global(.dark) .pill-nav {
  background-color: rgba(47, 38, 28, 0.85);
  box-shadow:
    0 12px 32px rgba(0, 0, 0, 0.45),
    0 2px 6px rgba(0, 0, 0, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.05);
}

.pill-burger {
  width: 40px;
  height: 40px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 999px;
  background: rgb(var(--muted));
  color: rgb(var(--ink));
  border: 1px solid rgb(var(--border));
  transition: background 150ms ease, color 150ms ease;
}
.pill-burger:hover {
  background: rgb(var(--primary));
  color: rgb(var(--primary-foreground));
  border-color: transparent;
}

/* ===== Мобильное раскрытие ===== */
.pill-mobile {
  width: min(92vw, 360px);
  padding: 0.5rem;
  background-color: rgba(255, 253, 247, 0.96);
  backdrop-filter: saturate(180%) blur(14px);
  -webkit-backdrop-filter: saturate(180%) blur(14px);
  border-radius: 1.25rem;
  border: 1px solid rgb(var(--border));
  box-shadow: 0 16px 36px rgba(0, 0, 0, 0.12);
}
:global(.dark) .pill-mobile {
  background-color: rgba(47, 38, 28, 0.95);
}

.mobile-drop-enter-active,
.mobile-drop-leave-active {
  transition: opacity 180ms ease, transform 180ms ease;
}
.mobile-drop-enter-from,
.mobile-drop-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}

@media (max-width: 640px) {
  .brand-plate {
    width: 52px;
    min-height: 100px;
    padding: 12px 6px 14px;
  }
  .brand-plate-text { font-size: 0.8rem; }
  .floating-controls {
    top: 1rem;
    left: 1rem;
    padding: 0.3rem 0.65rem;
  }
  .pill-nav {
    height: 50px;
    padding: 0 0.4rem 0 0.6rem;
  }
}
</style>
