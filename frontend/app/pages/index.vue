<script setup lang="ts">
import {
  Search,
  ArrowRight,
  BookMarked,
  Newspaper,
  Library,
  Megaphone,
  FileText,
  GraduationCap,
  Building2
} from 'lucide-vue-next'

const api = useApi()
const { t } = useI18n()
const loc = useLocalized()
const fileUrl = useFileUrl()
const localePath = useLocalePath()

interface PublicationsStats {
  total_publications: number
  total_teachers: number
  total_departments: number
}

interface HomeData {
  slides: any[]
  new_books: any[]
  section_blocks: { section: any; books: any[] }[]
  collections: { collection: any }[]
  latest_posts: any[]
  promos: any[]
  publications_stats?: PublicationsStats
}

const { data: home } = await useAsyncData('home-data', () =>
  api
    .get<HomeData>('/home')
    .catch((): HomeData => ({
      slides: [],
      new_books: [],
      section_blocks: [],
      collections: [],
      latest_posts: [],
      promos: []
    }))
)

const slides = computed(() => home.value?.slides ?? [])
const hasSlides = computed(() => slides.value.some((s) => s.image_path))
const newBooks = computed(() => home.value?.new_books ?? [])
const sectionBlocks = computed(() => home.value?.section_blocks ?? [])
const collections = computed(() =>
  (home.value?.collections ?? []).map((hc) => hc.collection).filter(Boolean)
)
const latestPosts = computed(() => home.value?.latest_posts ?? [])
const promos = computed(() => home.value?.promos ?? [])
const pubStats = computed(() => home.value?.publications_stats ?? null)
const hasPubStats = computed(
  () => !!pubStats.value && pubStats.value.total_publications > 0
)

const query = ref('')
function submitSearch() {
  if (!query.value.trim()) return
  navigateTo({ path: localePath('/search'), query: { q: query.value.trim() } })
}

const categories = [
  { to: '/books', icon: BookMarked, label: 'nav.books', desc: 'home.feat_books' },
  { to: '/periodicals', icon: Newspaper, label: 'nav.periodicals', desc: 'home.feat_periodicals' },
  { to: '/publications', icon: FileText, label: 'nav.publications', desc: 'home.feat_publications' },
  { to: '/collections', icon: Library, label: 'nav.collections', desc: 'home.feat_collections' },
  { to: '/news', icon: Megaphone, label: 'nav.news', desc: 'home.feat_news' }
]

useSeoMeta({
  title: () => t('home.title'),
  description: () => t('home.subtitle')
})
</script>

<template>
  <div>
    <!-- ===== Hero — на md+ карусель как фон-слой, на мобиле как блок-карточка сверху ===== -->
    <section
      class="relative isolate md:border-b md:border-border"
      :class="hasSlides ? 'md:text-white' : 'bg-card'"
    >
      <HeroCarousel v-if="hasSlides" :slides="slides" />

      <div class="page-container py-8 md:py-28">
        <div class="max-w-2xl">
          <span class="eyebrow" :class="hasSlides ? 'md:!text-white/70' : ''">
            {{ t('home.intro_eyebrow') }}
          </span>
          <h1
            class="text-3xl md:text-5xl font-bold tracking-tight mt-2 mb-4"
            :class="hasSlides ? 'md:text-white' : ''"
          >
            {{ t('home.intro_title') }}
          </h1>
          <p
            class="text-base md:text-lg leading-relaxed mb-7"
            :class="hasSlides ? 'text-muted-foreground md:text-white/85' : 'text-muted-foreground'"
          >
            {{ t('home.intro_text') }}
          </p>

          <form
            class="flex items-center h-12 px-4 rounded-xl border border-border bg-card focus-within:border-ring focus-within:ring-2 focus-within:ring-ring/15 transition-colors max-w-xl"
            @submit.prevent="submitSearch"
          >
            <Search class="w-5 h-5 text-muted-foreground shrink-0" />
            <input
              v-model="query"
              type="search"
              :placeholder="t('search.placeholder')"
              class="flex-1 h-full bg-transparent border-0 outline-none px-3 text-sm md:text-base text-foreground placeholder:text-muted-foreground"
            />
            <button
              type="submit"
              class="h-9 px-4 rounded-lg bg-primary text-primary-foreground text-sm font-medium hover:bg-secondary transition-colors"
            >
              {{ t('common.search') }}
            </button>
          </form>
        </div>
      </div>
    </section>

    <div class="page-container py-12 md:py-16 space-y-16">
      <!-- ===== Категории ===== -->
      <section>
        <div class="mb-6">
          <h2 class="text-xl md:text-2xl font-bold tracking-tight">{{ t('home.explore_title') }}</h2>
          <p class="text-sm text-muted-foreground mt-1">{{ t('home.explore_text') }}</p>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <NuxtLinkLocale
            v-for="cat in categories"
            :key="cat.to"
            :to="cat.to"
            class="group card card-hover p-5 flex flex-col"
          >
            <component :is="cat.icon" class="w-6 h-6 text-secondary mb-3" />
            <h3 class="font-semibold text-base">{{ t(cat.label) }}</h3>
            <p class="text-sm text-muted-foreground mt-1 flex-1">{{ t(cat.desc) }}</p>
            <span class="inline-flex items-center gap-1 text-sm font-medium text-secondary mt-3">
              {{ t('common.see_all') }}
              <ArrowRight class="w-4 h-4 transition-transform group-hover:translate-x-0.5" />
            </span>
          </NuxtLinkLocale>
        </div>
      </section>

      <!-- ===== Статистика публикаций преподавателей ===== -->
      <section v-if="hasPubStats">
        <NuxtLinkLocale
          to="/publications"
          class="group block rounded-xl border border-border bg-card p-6 md:p-8 hover:border-ring transition-colors"
        >
          <div class="flex items-end justify-between mb-6">
            <h2 class="text-xl md:text-2xl font-bold tracking-tight">
              {{ t('home.publications_stats_title') }}
            </h2>
            <span class="inline-flex items-center gap-1 text-sm font-medium text-secondary">
              {{ t('common.see_all') }}
              <ArrowRight class="w-4 h-4 transition-transform group-hover:translate-x-0.5" />
            </span>
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <div class="flex items-center gap-4">
              <span class="inline-flex items-center justify-center w-12 h-12 rounded-lg bg-primary/10 text-primary shrink-0">
                <FileText class="w-5 h-5" />
              </span>
              <div>
                <div class="text-2xl font-bold">{{ pubStats?.total_publications }}</div>
                <div class="text-sm text-muted-foreground">{{ t('home.stat_publications') }}</div>
              </div>
            </div>
            <div class="flex items-center gap-4">
              <span class="inline-flex items-center justify-center w-12 h-12 rounded-lg bg-primary/10 text-primary shrink-0">
                <GraduationCap class="w-5 h-5" />
              </span>
              <div>
                <div class="text-2xl font-bold">{{ pubStats?.total_teachers }}</div>
                <div class="text-sm text-muted-foreground">{{ t('home.stat_teachers') }}</div>
              </div>
            </div>
            <div class="flex items-center gap-4">
              <span class="inline-flex items-center justify-center w-12 h-12 rounded-lg bg-primary/10 text-primary shrink-0">
                <Building2 class="w-5 h-5" />
              </span>
              <div>
                <div class="text-2xl font-bold">{{ pubStats?.total_departments }}</div>
                <div class="text-sm text-muted-foreground">{{ t('home.stat_departments') }}</div>
              </div>
            </div>
          </div>
        </NuxtLinkLocale>
      </section>

      <!-- ===== Промо-блоки ===== -->
      <section v-if="promos.length">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-5">
          <component
            :is="promo.link_url ? 'a' : 'div'"
            v-for="promo in promos"
            :key="promo.id"
            :href="promo.link_url || undefined"
            class="group relative block rounded-xl overflow-hidden border border-border min-h-[200px] bg-muted"
          >
            <img
              v-if="promo.image_path"
              :src="fileUrl(promo.image_path)"
              :alt="loc(promo, 'title')"
              class="absolute inset-0 w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
            />
            <div
              class="relative z-10 flex flex-col justify-end min-h-[200px] p-6"
              style="background: linear-gradient(to top, rgba(15,23,42,0.78), transparent 65%)"
            >
              <h3 class="text-lg md:text-xl font-bold text-white">{{ loc(promo, 'title') }}</h3>
              <p v-if="loc(promo, 'subtitle')" class="text-white/85 text-sm mt-1.5">
                {{ loc(promo, 'subtitle') }}
              </p>
            </div>
          </component>
        </div>
      </section>

      <!-- ===== Новинки (карусель) ===== -->
      <section>
        <div class="flex items-end justify-between mb-5">
          <h2 class="text-xl md:text-2xl font-bold tracking-tight">{{ t('home.new_books') }}</h2>
          <NuxtLinkLocale
            to="/books"
            class="inline-flex items-center gap-1 text-sm font-medium text-secondary hover:underline"
          >
            {{ t('common.see_all') }}
            <ArrowRight class="w-4 h-4" />
          </NuxtLinkLocale>
        </div>
        <BookCarousel v-if="newBooks.length" :books="newBooks" />
        <div v-else class="card p-12 text-center text-muted-foreground">
          {{ t('home.no_books_yet') }}
        </div>
      </section>

      <!-- ===== Блоки разделов: баннер + карусель ===== -->
      <section v-for="block in sectionBlocks" :key="block.section.id">
        <NuxtLinkLocale
          v-if="block.section.image_path"
          :to="`/sections/${block.section.slug}`"
          class="group relative block rounded-xl overflow-hidden border border-border mb-5 min-h-[160px] md:min-h-[200px] bg-muted"
        >
          <img
            :src="fileUrl(block.section.image_path)"
            :alt="loc(block.section, 'name')"
            class="absolute inset-0 w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
          />
          <div
            class="relative z-10 flex flex-col justify-end min-h-[160px] md:min-h-[200px] p-6"
            style="background: linear-gradient(to top, rgba(15,23,42,0.7), transparent 60%)"
          >
            <h2 class="text-xl md:text-2xl font-bold text-white">
              {{ loc(block.section, 'name') }}
            </h2>
            <span class="inline-flex items-center gap-1 text-sm font-medium text-white/90 mt-1">
              {{ t('common.see_all') }}
              <ArrowRight class="w-4 h-4" />
            </span>
          </div>
        </NuxtLinkLocale>
        <div v-else class="flex items-end justify-between mb-5">
          <h2 class="text-xl md:text-2xl font-bold tracking-tight">
            {{ loc(block.section, 'name') }}
          </h2>
          <NuxtLinkLocale
            :to="`/sections/${block.section.slug}`"
            class="inline-flex items-center gap-1 text-sm font-medium text-secondary hover:underline"
          >
            {{ t('common.see_all') }}
            <ArrowRight class="w-4 h-4" />
          </NuxtLinkLocale>
        </div>

        <BookCarousel v-if="block.books.length" :books="block.books" />
        <p v-else class="text-sm text-muted-foreground">{{ t('section.no_books') }}</p>
      </section>

      <!-- ===== Подборки ===== -->
      <section v-if="collections.length">
        <div class="flex items-end justify-between mb-6">
          <h2 class="text-xl md:text-2xl font-bold tracking-tight">{{ t('home.collections') }}</h2>
          <NuxtLinkLocale
            to="/collections"
            class="inline-flex items-center gap-1 text-sm font-medium text-secondary hover:underline"
          >
            {{ t('common.see_all') }}
            <ArrowRight class="w-4 h-4" />
          </NuxtLinkLocale>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
          <CollectionCard v-for="col in collections" :key="col.id" :collection="col" />
        </div>
      </section>

      <!-- ===== Новости ===== -->
      <section v-if="latestPosts.length">
        <div class="flex items-end justify-between mb-6">
          <h2 class="text-xl md:text-2xl font-bold tracking-tight">{{ t('home.latest_news') }}</h2>
          <NuxtLinkLocale
            to="/news"
            class="inline-flex items-center gap-1 text-sm font-medium text-secondary hover:underline"
          >
            {{ t('common.see_all') }}
            <ArrowRight class="w-4 h-4" />
          </NuxtLinkLocale>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
          <PostCard v-for="post in latestPosts" :key="post.id" :post="post" />
        </div>
      </section>
    </div>
  </div>
</template>
