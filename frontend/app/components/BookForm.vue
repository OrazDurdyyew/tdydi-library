<script setup lang="ts">
interface Section {
  id: number
  parent_id: number | null
  name_tm: string
  name_ru: string
  name_en: string
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
  is_institute_book: boolean
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

const { t } = useI18n()
const loc = useLocalized()
const form = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v)
})

useAutoSlug(toRef(props.modelValue, 'title'), toRef(props.modelValue, 'slug'))

function sectionLabel(s: Section): string {
  if (s.parent_id === null) return loc(s, 'name')
  const parent = props.sections.find((p) => p.id === s.parent_id)
  return parent ? `${loc(parent, 'name')} → ${loc(s, 'name')}` : loc(s, 'name')
}
</script>

<template>
  <form @submit.prevent="emit('submit')" class="space-y-6 max-w-3xl">
    <Card>
      <CardHeader>
        <CardTitle class="text-lg">{{ t('admin.books.main') }}</CardTitle>
      </CardHeader>
      <CardContent class="space-y-5">
        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-2">
            <Label for="book-slug">{{ t('admin.common.slug') }}</Label>
            <Input id="book-slug" v-model="form.slug" readonly tabindex="-1" class="font-mono bg-muted text-muted-foreground cursor-not-allowed" />
            <p class="text-xs text-muted-foreground">{{ t('admin.common.slug_auto_hint') }}</p>
          </div>
          <div class="space-y-2">
            <Label for="book-language">{{ t('admin.books.language') }} *</Label>
            <Select id="book-language" v-model="form.language">
              <option value="tm">Türkmençe</option>
              <option value="ru">Русский</option>
              <option value="en">English</option>
            </Select>
          </div>
        </div>

        <div class="space-y-2">
          <Label for="book-title">{{ t('admin.common.title') }} *</Label>
          <Input id="book-title" v-model="form.title" required />
        </div>

        <div class="space-y-2">
          <Label for="book-authors">{{ t('admin.books.authors') }}</Label>
          <Input id="book-authors" v-model="form.authors" :placeholder="t('admin.books.authors_ph')" />
        </div>

        <div class="grid grid-cols-3 gap-4">
          <div class="space-y-2">
            <Label for="book-year">{{ t('admin.books.year') }}</Label>
            <Input id="book-year" v-model.number="form.year" type="number" />
          </div>
          <div class="space-y-2">
            <Label for="book-pages">{{ t('admin.books.pages') }}</Label>
            <Input id="book-pages" v-model.number="form.pages" type="number" />
          </div>
          <div class="space-y-2">
            <Label for="book-isbn">{{ t('admin.books.isbn') }}</Label>
            <Input id="book-isbn" v-model="form.isbn" />
          </div>
        </div>

        <div class="space-y-2">
          <Label for="book-publisher">{{ t('admin.books.publisher') }}</Label>
          <Input id="book-publisher" v-model="form.publisher" />
        </div>

        <div class="space-y-2">
          <Label for="book-annotation">{{ t('admin.books.annotation') }}</Label>
          <Textarea id="book-annotation" v-model="form.annotation" rows="4" />
        </div>
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle class="text-lg">{{ t('admin.books.sections') }} *</CardTitle>
        <CardDescription>{{ t('admin.books.sections_hint') }}</CardDescription>
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
          <span class="text-sm">{{ t('admin.books.allow_download') }}</span>
        </label>
        <label class="flex items-center gap-3 cursor-pointer">
          <Checkbox v-model="form.is_published" />
          <span class="text-sm">{{ t('admin.books.published') }}</span>
        </label>
        <label class="flex items-center gap-3 cursor-pointer">
          <Checkbox v-model="form.is_institute_book" />
          <span class="text-sm">{{ t('admin.books.is_institute_book') }}</span>
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
        {{ loading ? t('admin.common.saving') : t('admin.common.save') }}
      </Button>
      <NuxtLink to="/admin/books">
        <Button type="button" variant="outline">{{ t('admin.common.cancel') }}</Button>
      </NuxtLink>
    </div>
  </form>
</template>
