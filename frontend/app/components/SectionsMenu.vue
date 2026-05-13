<script setup lang="ts">
import { ChevronDown } from 'lucide-vue-next'

interface MenuSection {
  id: number
  slug: string
  name_tm: string
  name_ru: string
  name_en: string
  children?: MenuSection[]
}

interface MenuPage {
  slug: string
  title_tm: string
  title_ru: string
  title_en: string
}

defineProps<{ mobile?: boolean }>()

const api = useApi()
const { locale, t } = useI18n()

const { data: menu } = await useAsyncData('public-menu', () =>
  api.get<{ sections: MenuSection[]; pages: MenuPage[] }>('/menu').catch(() => ({ sections: [], pages: [] }))
)

function sectionName(s: MenuSection): string {
  if (locale.value === 'tm') return s.name_tm
  if (locale.value === 'en') return s.name_en
  return s.name_ru
}

function pageTitle(p: MenuPage): string {
  if (locale.value === 'tm') return p.title_tm
  if (locale.value === 'en') return p.title_en
  return p.title_ru
}
</script>

<template>
  <!-- ===== Desktop: горизонтальный список внутри pill-а ===== -->
  <ul v-if="!mobile" class="hidden md:flex items-center gap-1">
    <li>
      <NuxtLinkLocale to="/" class="nav-link" exact-active-class="nav-link-active">
        {{ t('nav.home') }}
      </NuxtLinkLocale>
    </li>

    <li
      v-for="section in menu?.sections"
      :key="section.id"
      :class="['relative', section.children && section.children.length ? 'group' : '']"
    >
      <NuxtLinkLocale
        :to="`/sections/${section.slug}`"
        class="nav-link"
        active-class="nav-link-active"
      >
        <span class="inline-flex items-center gap-1">
          {{ sectionName(section) }}
          <ChevronDown
            v-if="section.children && section.children.length"
            class="w-3.5 h-3.5 opacity-60 transition-transform group-hover:rotate-180"
          />
        </span>
      </NuxtLinkLocale>

      <!-- Подменю -->
      <div
        v-if="section.children && section.children.length"
        class="nav-dropdown"
      >
        <NuxtLinkLocale
          v-for="child in section.children"
          :key="child.id"
          :to="`/sections/${child.slug}`"
          class="nav-dropdown-item"
        >
          {{ sectionName(child) }}
        </NuxtLinkLocale>
      </div>
    </li>

    <li v-for="page in menu?.pages" :key="page.slug">
      <NuxtLinkLocale
        :to="`/${page.slug}`"
        class="nav-link"
        active-class="nav-link-active"
      >
        {{ pageTitle(page) }}
      </NuxtLinkLocale>
    </li>
  </ul>

  <!-- ===== Mobile: вертикальный список ===== -->
  <ul v-else class="flex flex-col">
    <li>
      <NuxtLinkLocale to="/" class="nav-link nav-link-mobile" exact-active-class="nav-link-active">
        {{ t('nav.home') }}
      </NuxtLinkLocale>
    </li>
    <template v-for="section in menu?.sections" :key="section.id">
      <li>
        <NuxtLinkLocale
          :to="`/sections/${section.slug}`"
          class="nav-link nav-link-mobile"
          active-class="nav-link-active"
        >
          {{ sectionName(section) }}
        </NuxtLinkLocale>
      </li>
      <li v-for="child in section.children" :key="child.id">
        <NuxtLinkLocale
          :to="`/sections/${child.slug}`"
          class="nav-link nav-link-mobile nav-link-sub"
          active-class="nav-link-active"
        >
          — {{ sectionName(child) }}
        </NuxtLinkLocale>
      </li>
    </template>
    <li v-for="page in menu?.pages" :key="page.slug">
      <NuxtLinkLocale
        :to="`/${page.slug}`"
        class="nav-link nav-link-mobile"
        active-class="nav-link-active"
      >
        {{ pageTitle(page) }}
      </NuxtLinkLocale>
    </li>
  </ul>
</template>

<style scoped>
.nav-link {
  position: relative;
  display: inline-flex;
  align-items: center;
  height: 40px;
  padding: 0 1rem;
  border-radius: 999px;
  font-family: var(--font-serif);
  font-size: 0.98rem;
  font-weight: 500;
  color: rgb(var(--ink));
  white-space: nowrap;
  text-decoration: none;
  transition: color 160ms ease;
}

/* Тонкое подчёркивание #16ab65 при наведении и в активном состоянии */
.nav-link::after {
  content: '';
  position: absolute;
  left: 1rem;
  right: 1rem;
  bottom: 8px;
  height: 1.5px;
  background-color: rgb(var(--primary));
  border-radius: 1px;
  transform: scaleX(0);
  transform-origin: center;
  transition: transform 220ms cubic-bezier(0.2, 0.7, 0.2, 1);
}
.nav-link:hover {
  color: rgb(var(--primary));
}
.nav-link:hover::after {
  transform: scaleX(1);
}
.nav-link-active {
  color: rgb(var(--primary));
}
.nav-link-active::after {
  transform: scaleX(1);
}

/* === Мобильный вариант === */
.nav-link-mobile {
  height: 44px;
  width: 100%;
  justify-content: flex-start;
  border-radius: 0.75rem;
  padding: 0 1rem;
}
.nav-link-mobile::after {
  left: 1rem;
  right: auto;
  width: 16px;
  bottom: 10px;
}
.nav-link-sub {
  font-size: 0.9rem;
  color: rgb(var(--muted-foreground));
  padding-left: 1.5rem;
}

/* === Dropdown === */
.nav-dropdown {
  position: absolute;
  top: calc(100% + 0.5rem);
  left: 50%;
  transform: translate(-50%, 4px);
  min-width: 240px;
  padding: 0.5rem;
  background-color: rgba(255, 253, 247, 0.98);
  backdrop-filter: saturate(180%) blur(12px);
  -webkit-backdrop-filter: saturate(180%) blur(12px);
  border: 1px solid rgb(var(--border));
  border-radius: 1rem;
  box-shadow: 0 18px 40px rgba(0, 0, 0, 0.12);
  opacity: 0;
  visibility: hidden;
  pointer-events: none;
  transition: opacity 180ms ease, transform 180ms ease, visibility 0s linear 180ms;
  z-index: 50;
}
:global(.dark) .nav-dropdown {
  background-color: rgba(47, 38, 28, 0.96);
}
.group:hover .nav-dropdown,
.group:focus-within .nav-dropdown {
  opacity: 1;
  visibility: visible;
  pointer-events: auto;
  transform: translate(-50%, 0);
  transition-delay: 0s;
}
.nav-dropdown-item {
  display: block;
  padding: 0.55rem 0.9rem;
  border-radius: 0.6rem;
  font-family: var(--font-serif);
  font-size: 0.92rem;
  color: rgb(var(--ink));
  text-decoration: none;
  transition: background-color 140ms ease, color 140ms ease;
}
.nav-dropdown-item:hover {
  background-color: rgba(22, 171, 101, 0.08);
  color: rgb(var(--primary));
}
</style>
