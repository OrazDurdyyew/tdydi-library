<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const api = useApi()
const toast = useToast()
const localePath = useLocalePath()

const form = reactive({
  slug: '',
  name_tm: '',
  name_ru: '',
  name_en: '',
  description_tm: '',
  description_ru: '',
  description_en: '',
  show_on_home: false,
  sort_order: 0,
})

const saving = ref(false)

async function submit() {
  if (!form.name_tm) {
    toast.error('Название (TM) обязательно')
    return
  }
  saving.value = true
  try {
    const created = await api.post<{ id: number }>('/admin/collections', form)
    toast.success('Подборка создана')
    await navigateTo(localePath(`/admin/collections/${created.id}`))
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка сохранения')
  } finally {
    saving.value = false
  }
}

useHead({ title: 'Создать подборку' })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/collections" class="text-muted-foreground hover:text-foreground text-sm">← Подборки</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">Создать подборку</h1>
    </div>

    <form @submit.prevent="submit" class="space-y-6 max-w-3xl">
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Основное</CardTitle>
        </CardHeader>
        <CardContent class="space-y-5">
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <Label for="slug">Slug</Label>
              <Input id="slug" v-model="form.slug" class="font-mono text-sm" placeholder="auto-генерируется" />
            </div>
            <div class="space-y-2">
              <Label for="order">Порядок</Label>
              <Input id="order" v-model.number="form.sort_order" type="number" />
            </div>
          </div>
          <label class="flex items-center gap-3 cursor-pointer">
            <Checkbox v-model="form.show_on_home" />
            <span class="text-sm">Показывать на главной</span>
          </label>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Названия</CardTitle>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="name-tm">TM <span class="text-destructive">*</span></Label>
            <Input id="name-tm" v-model="form.name_tm" required />
          </div>
          <div class="space-y-2">
            <Label for="name-ru">RU</Label>
            <Input id="name-ru" v-model="form.name_ru" />
          </div>
          <div class="space-y-2">
            <Label for="name-en">EN</Label>
            <Input id="name-en" v-model="form.name_en" />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Описание</CardTitle>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="desc-tm">TM</Label>
            <Textarea id="desc-tm" v-model="form.description_tm" rows="3" />
          </div>
          <div class="space-y-2">
            <Label for="desc-ru">RU</Label>
            <Textarea id="desc-ru" v-model="form.description_ru" rows="3" />
          </div>
          <div class="space-y-2">
            <Label for="desc-en">EN</Label>
            <Textarea id="desc-en" v-model="form.description_en" rows="3" />
          </div>
        </CardContent>
      </Card>

      <div class="flex gap-3">
        <Button type="submit" :disabled="saving">
          {{ saving ? 'Сохранение…' : 'Создать' }}
        </Button>
        <NuxtLink to="/admin/collections">
          <Button type="button" variant="outline">Отмена</Button>
        </NuxtLink>
      </div>
    </form>
  </div>
</template>
