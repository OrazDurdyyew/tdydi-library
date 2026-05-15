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
    <Card>
      <CardHeader>
        <CardTitle class="text-lg">Основное</CardTitle>
      </CardHeader>
      <CardContent class="space-y-5">
        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-2">
            <Label for="book-slug">Slug *</Label>
            <Input id="book-slug" v-model="form.slug" required class="font-mono" />
          </div>
          <div class="space-y-2">
            <Label for="book-language">Язык книги *</Label>
            <Select id="book-language" v-model="form.language">
              <option value="tm">Türkmençe</option>
              <option value="ru">Русский</option>
              <option value="en">English</option>
            </Select>
          </div>
        </div>

        <div class="space-y-2">
          <Label for="book-title">Название *</Label>
          <Input id="book-title" v-model="form.title" required />
        </div>

        <div class="space-y-2">
          <Label for="book-authors">Автор(ы)</Label>
          <Input id="book-authors" v-model="form.authors" placeholder="Иванов И.И., Петров П.П." />
        </div>

        <div class="grid grid-cols-3 gap-4">
          <div class="space-y-2">
            <Label for="book-year">Год</Label>
            <Input id="book-year" v-model.number="form.year" type="number" />
          </div>
          <div class="space-y-2">
            <Label for="book-pages">Страниц</Label>
            <Input id="book-pages" v-model.number="form.pages" type="number" />
          </div>
          <div class="space-y-2">
            <Label for="book-isbn">ISBN</Label>
            <Input id="book-isbn" v-model="form.isbn" />
          </div>
        </div>

        <div class="space-y-2">
          <Label for="book-publisher">Издательство</Label>
          <Input id="book-publisher" v-model="form.publisher" />
        </div>

        <div class="space-y-2">
          <Label for="book-annotation">Аннотация</Label>
          <Textarea id="book-annotation" v-model="form.annotation" rows="4" />
        </div>
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle class="text-lg">Разделы *</CardTitle>
        <CardDescription>Книга появится во всех выбранных разделах</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="border border-border rounded-lg p-3 max-h-56 overflow-y-auto space-y-1.5">
          <label
            v-for="s in sections"
            :key="s.id"
            class="flex items-center gap-3 cursor-pointer hover:bg-muted px-2 py-1 rounded-md"
          >
            <Checkbox :value="s.id" v-model="form.section_ids" />
            <span :class="s.parent_id === null ? 'font-medium text-sm' : 'text-sm pl-4 text-muted-foreground'">
              {{ sectionLabel(s) }}
            </span>
          </label>
        </div>
      </CardContent>
    </Card>

    <Card>
      <CardContent class="space-y-3 pt-6">
        <label class="flex items-center gap-3 cursor-pointer">
          <Checkbox v-model="form.allow_download" />
          <span class="text-sm">Разрешить скачивание</span>
        </label>
        <label class="flex items-center gap-3 cursor-pointer">
          <Checkbox v-model="form.is_published" />
          <span class="text-sm">Опубликована</span>
        </label>
      </CardContent>
    </Card>

    <div
      v-if="error"
      class="text-sm text-destructive p-3 bg-destructive/10 border border-destructive/20 rounded-lg"
    >
      {{ error }}
    </div>

    <div class="flex gap-3">
      <Button type="submit" :disabled="loading">
        {{ loading ? 'Сохранение…' : 'Сохранить' }}
      </Button>
      <NuxtLink to="/admin/books">
        <Button type="button" variant="outline">Отмена</Button>
      </NuxtLink>
    </div>
  </form>
</template>
