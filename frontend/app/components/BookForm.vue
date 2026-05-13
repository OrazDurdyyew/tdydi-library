<script setup lang="ts">
interface Section {
  id: number
  parent_id: number | null
  name_ru: string
}

interface BookFormData {
  slug: string
  isbn: string
  year: number
  pages: number
  language: string
  title: string
  annotation: string
  authors: string
  publisher: string
  allow_download: boolean
  is_published: boolean
  section_ids: number[]
}

const props = defineProps<{
  modelValue: BookFormData
  sections: Section[]
  loading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: BookFormData]
  submit: []
}>()

const form = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v)
})

function sectionLabel(s: Section): string {
  if (s.parent_id === null) return s.name_ru
  const parent = props.sections.find((p) => p.id === s.parent_id)
  return parent ? `${parent.name_ru} → ${s.name_ru}` : s.name_ru
}
</script>

<template>
  <form @submit.prevent="emit('submit')" class="space-y-6 max-w-3xl">
    <div class="card space-y-5">
      <h3 class="font-semibold">Основное</h3>

      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium mb-2">Slug *</label>
          <input v-model="form.slug" required class="input font-mono text-sm" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-2">Язык книги *</label>
          <select v-model="form.language" required class="input">
            <option value="tm">Türkmençe</option>
            <option value="ru">Русский</option>
            <option value="en">English</option>
          </select>
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium mb-2">Название *</label>
        <input v-model="form.title" required class="input" />
      </div>

      <div>
        <label class="block text-sm font-medium mb-2">Автор(ы)</label>
        <input v-model="form.authors" class="input" placeholder="Иванов И.И., Петров П.П." />
      </div>

      <div class="grid grid-cols-3 gap-4">
        <div>
          <label class="block text-sm font-medium mb-2">Год</label>
          <input v-model.number="form.year" type="number" class="input" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-2">Страниц</label>
          <input v-model.number="form.pages" type="number" class="input" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-2">ISBN</label>
          <input v-model="form.isbn" class="input" />
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium mb-2">Издательство</label>
        <input v-model="form.publisher" class="input" />
      </div>

      <div>
        <label class="block text-sm font-medium mb-2">Аннотация</label>
        <textarea
          v-model="form.annotation"
          rows="4"
          class="input"
        ></textarea>
      </div>
    </div>

    <div class="card">
      <h3 class="font-semibold mb-4">Разделы *</h3>
      <div class="border border-border rounded-lg p-3 max-h-56 overflow-y-auto space-y-1.5">
        <label v-for="s in sections" :key="s.id" class="flex items-center gap-2.5 cursor-pointer hover:bg-muted px-2 py-1 rounded-md">
          <input
            type="checkbox"
            :value="s.id"
            v-model="form.section_ids"
            class="accent-primary"
          />
          <span :class="s.parent_id === null ? 'font-medium' : 'text-sm pl-4'">
            {{ sectionLabel(s) }}
          </span>
        </label>
      </div>
      <div class="text-xs text-muted-foreground mt-2">
        Книга появится во всех выбранных разделах.
      </div>
    </div>

    <div class="card space-y-3">
      <label class="flex items-center gap-3 cursor-pointer">
        <input v-model="form.allow_download" type="checkbox" class="accent-primary" />
        <span class="text-sm">Разрешить скачивание</span>
      </label>
      <label class="flex items-center gap-3 cursor-pointer">
        <input v-model="form.is_published" type="checkbox" class="accent-primary" />
        <span class="text-sm">Опубликована</span>
      </label>
    </div>

    <div v-if="error" class="text-sm text-red-600 p-3 bg-red-50 rounded-lg">{{ error }}</div>

    <div class="flex gap-3">
      <button type="submit" :disabled="loading" class="btn-primary">
        {{ loading ? 'Сохранение…' : 'Сохранить' }}
      </button>
      <NuxtLink to="/admin/books" class="btn-outline">Отмена</NuxtLink>
    </div>
  </form>
</template>
