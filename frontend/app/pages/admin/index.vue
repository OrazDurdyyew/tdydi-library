<script setup lang="ts">
import { BookMarked, Eye, Download, TrendingUp, Plus } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

const { admin } = useAuth()
const api = useApi()

interface Stats {
  total_books: number
  total_views: number
  total_downloads: number
  top_books: { id: number; title: string; authors: string; views_count: number; downloads_count: number }[]
  daily_activity: { date: string; views: number }[]
  top_sections: { section_id: number; section_name: string; book_count: number; views_total: number }[]
  period_days: number
}

const { data: stats, pending } = await useAsyncData<Stats>('admin-stats', () =>
  api.get<Stats>('/admin/stats')
)

useHead({ title: 'Дашборд' })
</script>

<template>
  <div>
    <div class="mb-8">
      <h1 class="text-3xl font-bold font-serif">Дашборд</h1>
      <p v-if="admin" class="text-muted-foreground mt-1">
        Добро пожаловать, <strong class="text-foreground">{{ admin.username }}</strong>
        <span class="text-sm opacity-70"> · {{ admin.role }}</span>
      </p>
    </div>

    <!-- Метрики -->
    <div v-if="!pending && stats" class="grid grid-cols-1 sm:grid-cols-3 gap-5 mb-8">
      <Card>
        <CardContent class="flex items-center gap-4 pt-6">
          <span class="p-3 rounded-md bg-primary/10">
            <BookMarked class="w-5 h-5 text-primary" />
          </span>
          <div>
            <div class="text-2xl font-bold font-serif">{{ stats.total_books }}</div>
            <div class="text-sm text-muted-foreground">Книг в каталоге</div>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent class="flex items-center gap-4 pt-6">
          <span class="p-3 rounded-md bg-primary/10">
            <Eye class="w-5 h-5 text-primary" />
          </span>
          <div>
            <div class="text-2xl font-bold font-serif">{{ stats.total_views.toLocaleString() }}</div>
            <div class="text-sm text-muted-foreground">Просмотров</div>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent class="flex items-center gap-4 pt-6">
          <span class="p-3 rounded-md bg-primary/10">
            <Download class="w-5 h-5 text-primary" />
          </span>
          <div>
            <div class="text-2xl font-bold font-serif">{{ stats.total_downloads.toLocaleString() }}</div>
            <div class="text-sm text-muted-foreground">Скачиваний</div>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Скелетон загрузки -->
    <div v-if="pending" class="grid grid-cols-1 sm:grid-cols-3 gap-5 mb-8">
      <Card v-for="i in 3" :key="i" class="h-20 animate-pulse bg-muted" />
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Топ книг -->
      <Card v-if="stats" class="p-0 overflow-hidden">
        <CardHeader class="border-b border-border">
          <div class="flex items-center gap-2">
            <TrendingUp class="w-4 h-4 text-primary" />
            <CardTitle class="text-base">Популярные книги</CardTitle>
            <span class="text-xs text-muted-foreground ml-auto">за {{ stats.period_days }} дней</span>
          </div>
        </CardHeader>
        <CardContent class="p-0">
          <div v-if="stats.top_books.length === 0" class="px-6 py-10 text-center text-muted-foreground text-sm">
            Нет данных о просмотрах
          </div>
          <ul v-else>
            <li
              v-for="(book, idx) in stats.top_books"
              :key="book.id"
              class="flex items-center gap-3 px-6 py-3 border-b border-border last:border-0 hover:bg-muted/30 transition-colors"
            >
              <span class="w-5 text-center text-sm font-mono text-muted-foreground shrink-0">{{ idx + 1 }}</span>
              <div class="min-w-0 flex-1">
                <NuxtLink :to="`/admin/books/${book.id}`" class="font-medium text-sm hover:text-primary transition-colors truncate block">
                  {{ book.title }}
                </NuxtLink>
                <div class="text-xs text-muted-foreground truncate">{{ book.authors || '—' }}</div>
              </div>
              <div class="text-right shrink-0">
                <div class="text-sm font-medium">{{ book.views_count }}</div>
                <div class="text-xs text-muted-foreground">просм.</div>
              </div>
            </li>
          </ul>
        </CardContent>
      </Card>

      <!-- Топ разделов -->
      <Card v-if="stats && stats.top_sections.length > 0" class="p-0 overflow-hidden">
        <CardHeader class="border-b border-border">
          <CardTitle class="text-base">Активные разделы</CardTitle>
        </CardHeader>
        <CardContent class="p-0">
          <ul>
            <li
              v-for="s in stats.top_sections"
              :key="s.section_id"
              class="flex items-center gap-3 px-6 py-3 border-b border-border last:border-0"
            >
              <div class="min-w-0 flex-1">
                <div class="text-sm font-medium truncate">{{ s.section_name }}</div>
                <div class="text-xs text-muted-foreground">{{ s.book_count }} книг</div>
              </div>
              <div class="text-right shrink-0 text-sm">{{ s.views_total.toLocaleString() }} просм.</div>
            </li>
          </ul>
        </CardContent>
      </Card>

      <!-- Быстрый доступ -->
      <Card v-if="!pending && (!stats || stats.top_books.length === 0)">
        <CardHeader>
          <CardTitle class="text-base">Быстрый доступ</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="grid grid-cols-2 gap-3">
            <NuxtLink to="/admin/books/new">
              <Button class="w-full"><Plus class="w-4 h-4" /> Добавить книгу</Button>
            </NuxtLink>
            <NuxtLink to="/admin/sections">
              <Button variant="outline" class="w-full">Разделы</Button>
            </NuxtLink>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
