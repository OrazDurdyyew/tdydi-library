<script setup lang="ts">
import { Search, ChevronDown, HelpCircle } from 'lucide-vue-next'

const api = useApi()
const { t } = useI18n()
const loc = useLocalized()

const { data: categories } = await useAsyncData('faq-list', () =>
  api.get<any[]>('/faq').catch(() => [])
)

const query = ref('')
const openItems = ref<Set<number>>(new Set())

function toggle(id: number) {
  if (openItems.value.has(id)) openItems.value.delete(id)
  else openItems.value.add(id)
  openItems.value = new Set(openItems.value)
}

const filtered = computed(() => {
  const q = query.value.trim().toLowerCase()
  return (categories.value ?? [])
    .map((cat: any) => {
      const items = (cat.items ?? []).filter((item: any) => {
        if (!q) return true
        return (
          loc(item, 'question').toLowerCase().includes(q) ||
          loc(item, 'answer').toLowerCase().includes(q)
        )
      })
      return { ...cat, items }
    })
    .filter((cat: any) => cat.items.length > 0)
})

useSeoMeta({
  title: () => t('faq.title'),
  description: () => t('faq.subtitle')
})
</script>

<template>
  <div>
    <PageHero :title="t('faq.title')" :description="t('faq.subtitle')">
      <template #breadcrumbs>
        <Breadcrumbs :items="[{ label: t('nav.home'), to: '/' }, { label: t('faq.title') }]" />
      </template>
    </PageHero>

    <div class="page-container py-10">
      <div class="max-w-3xl">
        <!-- Поиск по FAQ -->
        <div
          class="flex items-center h-11 px-4 mb-8 rounded-xl border border-border bg-card focus-within:border-ring focus-within:ring-2 focus-within:ring-ring/15 transition-colors"
        >
          <Search class="w-4 h-4 text-muted-foreground shrink-0" />
          <input
            v-model="query"
            type="search"
            :placeholder="t('faq.search_placeholder')"
            class="flex-1 h-full bg-transparent border-0 outline-none px-3 text-sm placeholder:text-muted-foreground"
          />
        </div>

        <div v-if="filtered.length" class="space-y-10">
          <section v-for="cat in filtered" :key="cat.id">
            <h2 class="text-lg font-bold tracking-tight mb-4">{{ loc(cat, 'name') }}</h2>
            <div class="space-y-2">
              <div
                v-for="item in cat.items"
                :key="item.id"
                class="card overflow-hidden"
              >
                <button
                  class="w-full flex items-center justify-between gap-4 px-5 py-4 text-left hover:bg-muted/50 transition-colors"
                  @click="toggle(item.id)"
                >
                  <span class="font-medium">{{ loc(item, 'question') }}</span>
                  <ChevronDown
                    class="w-5 h-5 text-muted-foreground shrink-0 transition-transform"
                    :class="{ 'rotate-180': openItems.has(item.id) }"
                  />
                </button>
                <div
                  v-if="openItems.has(item.id)"
                  class="px-5 pb-5 pt-1 serif text-[1.02rem] leading-[1.7] text-foreground/85 border-t border-border"
                  v-html="loc(item, 'answer')"
                />
              </div>
            </div>
          </section>
        </div>

        <div v-else class="card p-16 text-center">
          <HelpCircle class="w-10 h-10 mx-auto text-muted-foreground/50 mb-3" />
          <p class="text-muted-foreground">
            {{ query ? t('common.nothing_found') : t('faq.empty') }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
