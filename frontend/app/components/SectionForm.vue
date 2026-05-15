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
  props.rootSections.filter((s) => s.parent_id == null && s.id !== props.excludeId)
)
</script>

<template>
  <form @submit.prevent="emit('submit')" class="space-y-6 max-w-3xl">
    <Card>
      <CardHeader>
        <CardTitle class="text-lg">Основное</CardTitle>
      </CardHeader>
      <CardContent class="space-y-5">
        <div class="space-y-2">
          <Label for="sec-slug">Slug *</Label>
          <Input id="sec-slug" v-model="form.slug" required class="font-mono" placeholder="economics" />
          <p class="text-xs text-muted-foreground">
            Используется в URL: /sections/{slug}. Должен быть уникальным.
          </p>
        </div>

        <div class="space-y-2">
          <Label for="sec-parent">Родительский раздел</Label>
          <Select id="sec-parent" v-model="form.parent_id">
            <option :value="null">— (корневой раздел в меню)</option>
            <option v-for="s in availableParents" :key="s.id" :value="s.id">
              {{ s.name_ru }}
            </option>
          </Select>
          <p class="text-xs text-muted-foreground">
            Если пусто — будет показан в главном меню. Иначе — как подраздел в dropdown.
          </p>
        </div>

        <div class="space-y-2">
          <Label for="sec-order">Порядок в меню</Label>
          <Input id="sec-order" v-model.number="form.sort_order" type="number" class="w-32" />
        </div>
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle class="text-lg">Названия на 3 языках</CardTitle>
      </CardHeader>
      <CardContent class="space-y-5">
        <div class="grid grid-cols-3 gap-4">
          <div class="space-y-2">
            <Label for="sec-name-tm">TM *</Label>
            <Input id="sec-name-tm" v-model="form.name_tm" required />
          </div>
          <div class="space-y-2">
            <Label for="sec-name-ru">RU *</Label>
            <Input id="sec-name-ru" v-model="form.name_ru" required />
          </div>
          <div class="space-y-2">
            <Label for="sec-name-en">EN *</Label>
            <Input id="sec-name-en" v-model="form.name_en" required />
          </div>
        </div>
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle class="text-lg">Описания на 3 языках</CardTitle>
      </CardHeader>
      <CardContent>
        <div class="grid grid-cols-3 gap-4">
          <div class="space-y-2">
            <Label for="sec-desc-tm">TM</Label>
            <Textarea id="sec-desc-tm" v-model="form.description_tm" rows="3" />
          </div>
          <div class="space-y-2">
            <Label for="sec-desc-ru">RU</Label>
            <Textarea id="sec-desc-ru" v-model="form.description_ru" rows="3" />
          </div>
          <div class="space-y-2">
            <Label for="sec-desc-en">EN</Label>
            <Textarea id="sec-desc-en" v-model="form.description_en" rows="3" />
          </div>
        </div>
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
      <NuxtLink to="/admin/sections">
        <Button type="button" variant="outline">Отмена</Button>
      </NuxtLink>
    </div>
  </form>
</template>
