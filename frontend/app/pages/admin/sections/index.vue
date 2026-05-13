<script setup lang="ts">
definePageMeta({ layout: 'admin' })

interface Section {
  id: number
  slug: string
  parent_id: number | null
  sort_order: number
  name_tm: string
  name_ru: string
  name_en: string
}

const api = useApi()
const { data, refresh, pending } = await useAsyncData('admin-sections', () => api.get<{ data: Section[] }>('/sections'))

const sections = computed(() => data.value?.data ?? [])
const parentNames = computed(() => {
  const map = new Map<number, string>()
  for (const s of sections.value) map.set(s.id, s.name_ru)
  return map
})

async function remove(id: number) {
  if (!confirm('Удалить раздел? Если в нём есть книги или подразделы — будет ошибка.')) return
  try {
    await api.delete(`/admin/sections/${id}`)
    await refresh()
  } catch (e: any) {
    alert(e.data?.error || 'Ошибка удаления')
  }
}

useHead({ title: 'Разделы' })
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold">Разделы</h1>
      <NuxtLink to="/admin/sections/new" class="btn-primary">+ Создать раздел</NuxtLink>
    </div>

    <div v-if="pending" class="text-muted-foreground">Загрузка…</div>

    <div v-else class="card !p-0 overflow-hidden">
      <table class="w-full">
        <thead>
          <tr class="text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">
            <th class="px-6 py-4">ID</th>
            <th class="px-6 py-4">Slug</th>
            <th class="px-6 py-4">Название (RU)</th>
            <th class="px-6 py-4">Родитель</th>
            <th class="px-6 py-4">Порядок</th>
            <th class="px-6 py-4"></th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="s in sections"
            :key="s.id"
            class="border-t border-border hover:bg-muted/40 transition-colors"
          >
            <td class="px-6 py-4 text-sm">{{ s.id }}</td>
            <td class="px-6 py-4 text-sm font-mono">{{ s.slug }}</td>
            <td class="px-6 py-4 font-medium">{{ s.name_ru }}</td>
            <td class="px-6 py-4 text-sm text-muted-foreground">
              {{ s.parent_id ? parentNames.get(s.parent_id) || s.parent_id : '—' }}
            </td>
            <td class="px-6 py-4 text-sm">{{ s.sort_order }}</td>
            <td class="px-6 py-4 text-right space-x-4">
              <NuxtLink :to="`/admin/sections/${s.id}`" class="text-primary hover:underline text-sm">Изменить</NuxtLink>
              <button @click="remove(s.id)" class="text-red-600 hover:underline text-sm">Удалить</button>
            </td>
          </tr>
          <tr v-if="sections.length === 0">
            <td colspan="6" class="px-6 py-16 text-center text-muted-foreground">
              Пока ни одного раздела. Создайте первый.
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
