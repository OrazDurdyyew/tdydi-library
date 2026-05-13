<script setup lang="ts">
definePageMeta({ layout: 'admin' })

interface Book {
  id: number
  slug: string
  title: string
  authors: string
  year: number
  language: string
  is_published: boolean
  sections?: { id: number; name_ru: string }[]
}

const api = useApi()
const { data, refresh, pending } = await useAsyncData('admin-books', () =>
  api.get<{ data: Book[]; meta: { total: number } }>('/books?limit=100')
)

const books = computed(() => data.value?.data ?? [])

async function remove(id: number) {
  if (!confirm('Удалить книгу?')) return
  try {
    await api.delete(`/admin/books/${id}`)
    await refresh()
  } catch (e: any) {
    alert(e.data?.error || 'Ошибка')
  }
}

useHead({ title: 'Книги' })
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold">Книги</h1>
      <NuxtLink to="/admin/books/new" class="btn-primary">+ Добавить книгу</NuxtLink>
    </div>

    <div v-if="pending">Загрузка…</div>

    <div v-else class="card !p-0 overflow-hidden">
      <table class="w-full">
        <thead>
          <tr class="text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">
            <th class="px-6 py-4">ID</th>
            <th class="px-6 py-4">Название</th>
            <th class="px-6 py-4">Автор</th>
            <th class="px-6 py-4">Год</th>
            <th class="px-6 py-4">Язык</th>
            <th class="px-6 py-4">Статус</th>
            <th class="px-6 py-4"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="b in books" :key="b.id" class="border-t border-border hover:bg-muted/40 transition-colors">
            <td class="px-6 py-4 text-sm">{{ b.id }}</td>
            <td class="px-6 py-4 font-medium">{{ b.title }}</td>
            <td class="px-6 py-4 text-sm text-muted-foreground">{{ b.authors }}</td>
            <td class="px-6 py-4 text-sm">{{ b.year || '—' }}</td>
            <td class="px-6 py-4 text-sm uppercase">{{ b.language }}</td>
            <td class="px-6 py-4">
              <span :class="b.is_published ? 'badge-success' : 'badge-warning'">
                {{ b.is_published ? 'Опубликована' : 'Черновик' }}
              </span>
            </td>
            <td class="px-6 py-4 text-right space-x-4">
              <NuxtLink :to="`/admin/books/${b.id}`" class="text-primary hover:underline text-sm">Изменить</NuxtLink>
              <button @click="remove(b.id)" class="text-red-600 hover:underline text-sm">Удалить</button>
            </td>
          </tr>
          <tr v-if="books.length === 0">
            <td colspan="7" class="px-6 py-16 text-center text-muted-foreground">
              Книг пока нет. Добавьте первую.
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
