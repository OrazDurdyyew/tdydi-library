<script setup lang="ts">
import { Search, BookOpen, ArrowRight } from 'lucide-vue-next'

interface Book {
  id: number
  slug: string
  title: string
  authors?: string
  year?: number
  cover_path?: string
}

const props = defineProps<{
  variant?: 'desktop' | 'mobile'
}>()

const emit = defineEmits<{ navigated: [] }>()

const { t } = useI18n()
const api = useApi()
const localePath = useLocalePath()
const fileUrl = useFileUrl()

const query = ref('')
const suggestions = ref<Book[]>([])
const open = ref(false)
const loading = ref(false)
const activeIndex = ref(-1)
const rootEl = ref<HTMLElement | null>(null)

let debounceTimer: ReturnType<typeof setTimeout> | null = null
let reqSeq = 0

async function fetchSuggestions(q: string) {
  const mySeq = ++reqSeq
  loading.value = true
  try {
    const res = await api.get<{ books: Book[] }>(
      `/search?q=${encodeURIComponent(q)}&books_only=1`
    )
    if (mySeq !== reqSeq) return
    suggestions.value = res.books || []
  } catch {
    if (mySeq !== reqSeq) return
    suggestions.value = []
  } finally {
    if (mySeq === reqSeq) loading.value = false
  }
}

watch(query, (val) => {
  activeIndex.value = -1
  const q = val.trim()
  if (debounceTimer) clearTimeout(debounceTimer)
  if (q.length < 2) {
    suggestions.value = []
    loading.value = false
    return
  }
  debounceTimer = setTimeout(() => fetchSuggestions(q), 250)
})

function submitSearch() {
  const q = query.value.trim()
  if (!q) return
  open.value = false
  navigateTo({ path: localePath('/search'), query: { q } })
  emit('navigated')
}

function goToBook(book: Book) {
  open.value = false
  query.value = ''
  suggestions.value = []
  navigateTo(localePath(`/books/${book.slug}`))
  emit('navigated')
}

function onKeyDown(e: KeyboardEvent) {
  if (!open.value || suggestions.value.length === 0) return
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    activeIndex.value = (activeIndex.value + 1) % suggestions.value.length
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    activeIndex.value =
      activeIndex.value <= 0 ? suggestions.value.length - 1 : activeIndex.value - 1
  } else if (e.key === 'Enter') {
    if (activeIndex.value >= 0) {
      e.preventDefault()
      const book = suggestions.value[activeIndex.value]
      if (book) goToBook(book)
    }
  } else if (e.key === 'Escape') {
    open.value = false
  }
}

onClickOutside(rootEl, () => {
  open.value = false
})

const showDropdown = computed(
  () => open.value && query.value.trim().length >= 2
)
const showEmpty = computed(
  () => showDropdown.value && !loading.value && suggestions.value.length === 0
)
</script>

<template>
  <div ref="rootEl" class="relative w-full">
    <form
      :class="[
        'flex items-center px-3.5 rounded-lg border border-border bg-background focus-within:border-ring focus-within:ring-2 focus-within:ring-ring/15 transition-colors',
        variant === 'mobile' ? 'h-10' : 'h-10'
      ]"
      @submit.prevent="submitSearch"
    >
      <Search class="w-4 h-4 text-muted-foreground shrink-0" />
      <input
        v-model="query"
        type="search"
        autocomplete="off"
        :placeholder="t('search.placeholder')"
        class="flex-1 h-full bg-transparent border-0 outline-none px-3 text-sm placeholder:text-muted-foreground"
        @focus="open = true"
        @keydown="onKeyDown"
      />
    </form>

    <div
      v-if="showDropdown"
      class="absolute left-0 right-0 top-full mt-2 z-50 rounded-xl border border-border bg-card shadow-lg overflow-hidden"
    >
      <div v-if="loading && suggestions.length === 0" class="p-4 text-sm text-muted-foreground">
        {{ t('common.loading') }}
      </div>

      <div v-else-if="showEmpty" class="p-4 text-sm text-muted-foreground">
        {{ t('search.no_suggestions') }}
      </div>

      <ul v-else class="py-1 max-h-[70vh] overflow-y-auto">
        <li v-for="(book, i) in suggestions" :key="book.id">
          <button
            type="button"
            :class="[
              'w-full flex items-center gap-3 px-3 py-2 text-left transition-colors',
              activeIndex === i ? 'bg-muted' : 'hover:bg-muted'
            ]"
            @click="goToBook(book)"
            @mouseenter="activeIndex = i"
          >
            <div class="w-10 h-14 shrink-0 rounded bg-muted overflow-hidden flex items-center justify-center">
              <img
                v-if="book.cover_path"
                :src="fileUrl(book.cover_path)"
                :alt="book.title"
                class="w-full h-full object-cover"
                loading="lazy"
              />
              <BookOpen v-else class="w-4 h-4 text-muted-foreground" />
            </div>
            <div class="min-w-0 flex-1">
              <div class="text-sm font-medium truncate">{{ book.title }}</div>
              <div v-if="book.authors" class="text-xs text-muted-foreground truncate">
                {{ book.authors }}
              </div>
            </div>
          </button>
        </li>
      </ul>

      <button
        v-if="suggestions.length > 0"
        type="button"
        class="w-full flex items-center justify-center gap-2 border-t border-border px-3 py-2.5 text-sm font-medium text-primary hover:bg-muted transition-colors"
        @click="submitSearch"
      >
        {{ t('search.see_all_results') }}
        <ArrowRight class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>
