<script setup lang="ts">
interface SectionFormData {
  slug: string
  parent_id: number | null
  sort_order: number
  image_path: string
  name_tm: string
  name_ru: string
  name_en: string
  description_tm: string
  description_ru: string
  description_en: string
}

interface RootSection {
  id: number
  name_ru: string
  parent_id: number | null
}

const props = defineProps<{
  modelValue: SectionFormData
  rootSections: RootSection[]
  excludeId?: number
  loading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: SectionFormData]
  submit: []
}>()

const form = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v)
})

const availableParents = computed(() =>
  props.rootSections.filter((s) => s.parent_id === null && s.id !== props.excludeId)
)
</script>

<template>
  <form @submit.prevent="emit('submit')" class="space-y-6 max-w-3xl">
    <div class="card space-y-5">
      <div>
        <label class="block text-sm font-medium mb-2">Slug *</label>
        <input
          v-model="form.slug"
          required
          class="input font-mono text-sm"
          placeholder="economics"
        />
        <div class="text-xs text-muted-foreground mt-1.5">
          Используется в URL: /sections/{slug}. Должен быть уникальным.
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium mb-2">Родительский раздел</label>
        <select v-model="form.parent_id" class="input">
          <option :value="null">— (корневой раздел в меню)</option>
          <option v-for="s in availableParents" :key="s.id" :value="s.id">
            {{ s.name_ru }}
          </option>
        </select>
        <div class="text-xs text-muted-foreground mt-1.5">
          Если пусто — будет показан в главном меню. Иначе — как подраздел в dropdown.
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium mb-2">Порядок в меню</label>
        <input
          v-model.number="form.sort_order"
          type="number"
          class="input w-32"
        />
      </div>
    </div>

    <div class="card space-y-5">
      <h3 class="font-semibold">Названия на 3 языках</h3>
      <div class="grid grid-cols-3 gap-4">
        <div>
          <label class="block text-sm font-medium mb-2">TM *</label>
          <input v-model="form.name_tm" required class="input" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-2">RU *</label>
          <input v-model="form.name_ru" required class="input" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-2">EN *</label>
          <input v-model="form.name_en" required class="input" />
        </div>
      </div>

      <h3 class="font-semibold pt-2">Описания на 3 языках</h3>
      <div class="grid grid-cols-3 gap-4">
        <div>
          <label class="block text-sm font-medium mb-2">TM</label>
          <textarea v-model="form.description_tm" rows="3" class="input"></textarea>
        </div>
        <div>
          <label class="block text-sm font-medium mb-2">RU</label>
          <textarea v-model="form.description_ru" rows="3" class="input"></textarea>
        </div>
        <div>
          <label class="block text-sm font-medium mb-2">EN</label>
          <textarea v-model="form.description_en" rows="3" class="input"></textarea>
        </div>
      </div>
    </div>

    <div v-if="error" class="text-sm text-red-600 p-3 bg-red-50 rounded-lg">
      {{ error }}
    </div>

    <div class="flex gap-3">
      <button type="submit" :disabled="loading" class="btn-primary">
        {{ loading ? 'Сохранение…' : 'Сохранить' }}
      </button>
      <NuxtLink to="/admin/sections" class="btn-outline">Отмена</NuxtLink>
    </div>
  </form>
</template>
