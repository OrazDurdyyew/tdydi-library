<script setup lang="ts">
import { Search, ArrowRight } from 'lucide-vue-next'

const api = useApi()
const { t } = useI18n()
const localePath = useLocalePath()

const { data: booksResp } = await useAsyncData('home-new-books', () =>
  api.get<{ data: any[] }>('/books?sort=-created_at&limit=12').catch(() => ({ data: [] }))
)

const newBooks = computed(() => booksResp.value?.data ?? [])

const query = ref('')
function submitSearch() {
  if (!query.value.trim()) return
  navigateTo({
    path: localePath('/books'),
    query: { q: query.value.trim() }
  })
}

useSeoMeta({
  title: () => t('home.title'),
  description: () => t('home.subtitle')
})
</script>

<template>
  <div>
    <!-- ===== Hero ===== -->
    <section class="pt-24 md:pt-32 pb-12 md:pb-16">
      <div class="page-container">
        <div class="max-w-3xl mx-auto text-center">
          <h1 class="font-serif font-bold tracking-tight text-5xl md:text-7xl leading-[1.02] mb-7">
            {{ $t('home.title') }}
          </h1>
          <p class="text-base md:text-lg text-foreground/75 leading-relaxed max-w-xl mx-auto mb-10">
            {{ $t('home.subtitle') }}
          </p>

          <!-- Pill-CTA -->
          <div class="flex flex-wrap justify-center gap-3 mb-10">
            <NuxtLinkLocale to="/books" class="cta-pill cta-pill-solid">
              {{ $t('nav.books') }}
            </NuxtLinkLocale>
            <NuxtLinkLocale to="/periodicals" class="cta-pill">
              {{ $t('nav.periodicals') }}
            </NuxtLinkLocale>
          </div>

          <!-- Поиск -->
          <form
            @submit.prevent="submitSearch"
            class="search-bar relative max-w-xl mx-auto h-14 flex items-center pl-5 pr-2"
          >
            <Search class="w-5 h-5 text-muted-foreground shrink-0 mr-3" />
            <input
              v-model="query"
              type="search"
              :placeholder="$t('common.search') + '…'"
              class="flex-1 h-full bg-transparent border-0 outline-none font-serif italic text-base placeholder:text-muted-foreground/70"
            />
            <button
              type="submit"
              class="h-10 px-5 rounded-full bg-primary text-primary-foreground text-sm font-semibold hover:opacity-90 active:scale-[0.98] transition"
            >
              {{ $t('common.search') }}
            </button>
          </form>
        </div>
      </div>
    </section>

    <!-- ===== Большой презентационный блок (как карточка-постер) ===== -->
    <section class="pb-20">
      <div class="page-container">
        <article class="hero-poster">
          <div class="hero-poster-art" aria-hidden="true">
            <div class="hero-poster-art-inner">
              <div class="poster-shape shape-1" />
              <div class="poster-shape shape-2" />
              <div class="poster-shape shape-3" />
            </div>
          </div>
          <div class="hero-poster-body">
            <span class="eyebrow">{{ $t('home.new_books') }}</span>
            <h2 class="font-serif font-bold text-3xl md:text-5xl leading-tight mt-3 mb-4">
              {{ $t('home.new_books') }}
            </h2>
            <p class="text-foreground/75 leading-relaxed mb-7 max-w-md">
              {{ $t('home.subtitle') }}
            </p>
            <NuxtLinkLocale to="/books" class="cta-pill cta-pill-solid inline-flex items-center gap-2">
              {{ $t('common.see_all') }}
              <ArrowRight class="w-4 h-4" />
            </NuxtLinkLocale>
          </div>
        </article>
      </div>
    </section>

    <!-- ===== Полка с новинками ===== -->
    <section class="pb-24">
      <div class="page-container">
        <div class="flex items-end justify-between mb-10">
          <div>
            <span class="eyebrow">caput · index</span>
            <h2 class="font-serif font-bold text-3xl md:text-4xl mt-2">
              {{ $t('home.new_books') }}
            </h2>
          </div>
          <NuxtLinkLocale
            to="/books"
            class="hidden sm:inline-flex items-center gap-1.5 text-sm font-semibold text-primary hover:underline"
          >
            {{ $t('common.see_all') }}
            <ArrowRight class="w-4 h-4" />
          </NuxtLinkLocale>
        </div>

        <div v-if="newBooks.length === 0" class="paper-card text-center py-20 text-muted-foreground serif italic">
          {{ $t('home.no_books_yet') }}
        </div>

        <div
          v-else
          class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-x-6 gap-y-12"
        >
          <BookCard
            v-for="book in newBooks"
            :key="book.id"
            :book="book"
          />
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
/* === Pill-кнопки === */
.cta-pill {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 48px;
  padding: 0 1.5rem;
  border-radius: 999px;
  background-color: rgb(var(--card));
  border: 1px solid rgb(var(--border));
  font-family: var(--font-sans);
  font-weight: 500;
  font-size: 0.95rem;
  color: rgb(var(--ink));
  text-decoration: none;
  transition: background-color 160ms ease, color 160ms ease, transform 160ms ease, box-shadow 160ms ease;
  box-shadow: 0 1px 0 rgba(0, 0, 0, 0.04);
}
.cta-pill:hover {
  background-color: rgba(22, 171, 101, 0.08);
  color: rgb(var(--primary));
  border-color: rgba(22, 171, 101, 0.45);
}
.cta-pill-solid {
  background-color: rgb(var(--primary));
  color: rgb(var(--primary-foreground));
  border-color: transparent;
  box-shadow: 0 6px 18px rgba(22, 171, 101, 0.28);
}
.cta-pill-solid:hover {
  background-color: rgb(var(--primary));
  color: rgb(var(--primary-foreground));
  transform: translateY(-1px);
  box-shadow: 0 10px 22px rgba(22, 171, 101, 0.35);
}

/* === Search bar === */
.search-bar {
  border-radius: 999px;
  background-color: rgb(var(--card));
  border: 1px solid rgb(var(--border));
  box-shadow: 0 1px 0 rgba(255, 255, 255, 0.5), 0 8px 18px rgba(0, 0, 0, 0.05);
  transition: border-color 160ms ease, box-shadow 160ms ease;
}
.search-bar:focus-within {
  border-color: rgba(22, 171, 101, 0.5);
  box-shadow: 0 0 0 4px rgba(22, 171, 101, 0.12), 0 8px 18px rgba(0, 0, 0, 0.06);
}

/* === Hero poster (большая презентационная карточка) === */
.hero-poster {
  display: grid;
  grid-template-columns: 1fr;
  gap: 0;
  background-color: rgb(var(--card));
  border-radius: 28px;
  overflow: hidden;
  box-shadow: 0 24px 56px rgba(0, 0, 0, 0.08), inset 0 0 1px rgba(0, 0, 0, 0.04);
}
@media (min-width: 768px) {
  .hero-poster {
    grid-template-columns: 1.1fr 1fr;
  }
}
.hero-poster-art {
  position: relative;
  min-height: 320px;
  background:
    radial-gradient(circle at 30% 30%, rgba(22, 171, 101, 0.15), transparent 60%),
    linear-gradient(135deg, rgb(var(--muted)) 0%, rgb(var(--card)) 100%);
}
.hero-poster-art-inner {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
}
.poster-shape {
  position: absolute;
  border-radius: 12px;
  background-color: rgb(var(--card));
  box-shadow: 4px 6px 18px rgba(0, 0, 0, 0.12), inset 0 0 1px rgba(0, 0, 0, 0.06);
}
.shape-1 {
  width: 140px;
  height: 200px;
  transform: translate(-40px, -10px) rotate(-8deg);
  background: linear-gradient(180deg, rgb(var(--primary)) 0%, rgb(var(--primary)) 35%, rgb(var(--card)) 35%);
}
.shape-2 {
  width: 140px;
  height: 200px;
  transform: translate(40px, 20px) rotate(6deg);
  background: linear-gradient(180deg, rgb(var(--secondary)) 0%, rgb(var(--secondary)) 35%, rgb(var(--card)) 35%);
}
.shape-3 {
  width: 140px;
  height: 200px;
  transform: translate(110px, -20px) rotate(12deg);
  background: linear-gradient(180deg, rgb(var(--ink)) 0%, rgb(var(--ink)) 35%, rgb(var(--card)) 35%);
}
.hero-poster-body {
  padding: 2rem 2rem 2.25rem;
}
@media (min-width: 768px) {
  .hero-poster-body { padding: 3rem 3rem 3.25rem; }
}
</style>
