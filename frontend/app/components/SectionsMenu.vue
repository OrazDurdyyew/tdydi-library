<script setup lang="ts">
import { ChevronDown } from 'lucide-vue-next'
import type { MenuSection, MenuPage } from '~/composables/useMenu'

defineProps<{ mobile?: boolean }>()

const { t } = useI18n()
const loc = useLocalized()
const { data: menu } = await useMenu()

function sectionName(s: MenuSection) {
  return loc(s, 'name')
}
function pageTitle(p: MenuPage) {
  return loc(p, 'title')
}
</script>

<template>
  <!-- ===== Desktop ===== -->
  <ul v-if="!mobile" class="hidden md:flex items-center gap-0.5 h-12">
    <li>
      <NuxtLinkLocale to="/" class="nav-link" exact-active-class="nav-link-active">
        {{ t('nav.home') }}
      </NuxtLinkLocale>
    </li>
    <li>
      <NuxtLinkLocale to="/books" class="nav-link" active-class="nav-link-active">
        {{ t('nav.books') }}
      </NuxtLinkLocale>
    </li>

    <li
      v-for="section in menu?.sections"
      :key="section.id"
      :class="['relative', section.children && section.children.length ? 'group' : '']"
    >
      <NuxtLinkLocale
        :to="`/sections/${section.slug}`"
        class="nav-link inline-flex items-center gap-1"
        active-class="nav-link-active"
      >
        {{ sectionName(section) }}
        <ChevronDown
          v-if="section.children && section.children.length"
          class="w-3.5 h-3.5 opacity-50 transition-transform group-hover:rotate-180"
        />
      </NuxtLinkLocale>

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

    <li>
      <NuxtLinkLocale to="/periodicals" class="nav-link" active-class="nav-link-active">
        {{ t('nav.periodicals') }}
      </NuxtLinkLocale>
    </li>
    <li>
      <NuxtLinkLocale to="/news" class="nav-link" active-class="nav-link-active">
        {{ t('nav.news') }}
      </NuxtLinkLocale>
    </li>

    <li v-for="page in menu?.pages" :key="page.slug">
      <NuxtLinkLocale :to="`/${page.slug}`" class="nav-link" active-class="nav-link-active">
        {{ pageTitle(page) }}
      </NuxtLinkLocale>
    </li>
  </ul>

  <!-- ===== Mobile ===== -->
  <ul v-else class="flex flex-col gap-0.5">
    <li>
      <NuxtLinkLocale to="/" class="nav-link-mobile" exact-active-class="nav-link-active">
        {{ t('nav.home') }}
      </NuxtLinkLocale>
    </li>
    <li>
      <NuxtLinkLocale to="/books" class="nav-link-mobile" active-class="nav-link-active">
        {{ t('nav.books') }}
      </NuxtLinkLocale>
    </li>
    <template v-for="section in menu?.sections" :key="section.id">
      <li>
        <NuxtLinkLocale
          :to="`/sections/${section.slug}`"
          class="nav-link-mobile"
          active-class="nav-link-active"
        >
          {{ sectionName(section) }}
        </NuxtLinkLocale>
      </li>
      <li v-for="child in section.children" :key="child.id">
        <NuxtLinkLocale
          :to="`/sections/${child.slug}`"
          class="nav-link-mobile nav-link-sub"
          active-class="nav-link-active"
        >
          {{ sectionName(child) }}
        </NuxtLinkLocale>
      </li>
    </template>
    <li>
      <NuxtLinkLocale to="/periodicals" class="nav-link-mobile" active-class="nav-link-active">
        {{ t('nav.periodicals') }}
      </NuxtLinkLocale>
    </li>
    <li>
      <NuxtLinkLocale to="/news" class="nav-link-mobile" active-class="nav-link-active">
        {{ t('nav.news') }}
      </NuxtLinkLocale>
    </li>
    <li>
      <NuxtLinkLocale to="/collections" class="nav-link-mobile" active-class="nav-link-active">
        {{ t('nav.collections') }}
      </NuxtLinkLocale>
    </li>
    <li>
      <NuxtLinkLocale to="/faq" class="nav-link-mobile" active-class="nav-link-active">
        {{ t('nav.faq') }}
      </NuxtLinkLocale>
    </li>
    <li v-for="page in menu?.pages" :key="page.slug">
      <NuxtLinkLocale :to="`/${page.slug}`" class="nav-link-mobile" active-class="nav-link-active">
        {{ pageTitle(page) }}
      </NuxtLinkLocale>
    </li>
  </ul>
</template>

<style scoped>
.nav-link {
  display: inline-flex;
  align-items: center;
  height: 3rem;
  padding: 0 0.875rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: rgb(var(--muted-foreground));
  white-space: nowrap;
  border-bottom: 2px solid transparent;
  transition: color 150ms ease, border-color 150ms ease;
}
.nav-link:hover {
  color: rgb(var(--foreground));
}
.nav-link-active {
  color: rgb(var(--primary));
  border-bottom-color: rgb(var(--accent));
}

.nav-link-mobile {
  display: block;
  padding: 0.625rem 0.75rem;
  border-radius: 0.5rem;
  font-size: 0.95rem;
  font-weight: 500;
  color: rgb(var(--foreground));
  transition: background-color 140ms ease;
}
.nav-link-mobile:hover {
  background-color: rgb(var(--muted));
}
.nav-link-sub {
  padding-left: 1.75rem;
  font-size: 0.9rem;
  font-weight: 400;
  color: rgb(var(--muted-foreground));
}
.nav-link-mobile.nav-link-active {
  color: rgb(var(--primary));
  background-color: rgb(var(--muted));
  border-bottom: 0;
}

/* Dropdown */
.nav-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  min-width: 14rem;
  padding: 0.375rem;
  background-color: rgb(var(--card));
  border: 1px solid rgb(var(--border));
  border-radius: 0.625rem;
  box-shadow: 0 8px 24px rgb(15 23 42 / 0.1);
  opacity: 0;
  visibility: hidden;
  transform: translateY(4px);
  transition: opacity 160ms ease, transform 160ms ease, visibility 0s linear 160ms;
  z-index: 50;
}
.group:hover .nav-dropdown,
.group:focus-within .nav-dropdown {
  opacity: 1;
  visibility: visible;
  transform: translateY(0);
  transition-delay: 0s;
}
.nav-dropdown-item {
  display: block;
  padding: 0.5rem 0.75rem;
  border-radius: 0.4rem;
  font-size: 0.875rem;
  color: rgb(var(--foreground));
  transition: background-color 130ms ease, color 130ms ease;
}
.nav-dropdown-item:hover {
  background-color: rgb(var(--muted));
  color: rgb(var(--primary));
}
</style>
