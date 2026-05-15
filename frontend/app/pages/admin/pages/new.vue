<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const api = useApi()
const toast = useToast()
const localePath = useLocalePath()

const form = reactive({
  slug: '',
  in_menu: false,
  menu_order: 0,
  title_tm: '',
  title_ru: '',
  title_en: '',
  content_tm: '',
  content_ru: '',
  content_en: '',
})

const saving = ref(false)

async function submit() {
  if (!form.title_tm) {
    toast.error('Заголовок (TM) обязателен')
    return
  }
  saving.value = true
  try {
    const created = await api.post<{ id: number }>('/admin/pages', form)
    toast.success('Страница создана')
    await navigateTo(localePath(`/admin/pages/${created.id}`))
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка сохранения')
  } finally {
    saving.value = false
  }
}

useHead({ title: 'Создать страницу' })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/pages" class="text-muted-foreground hover:text-foreground text-sm">← Страницы</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">Создать страницу</h1>
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
              <Input id="slug" v-model="form.slug" class="font-mono text-sm" />
            </div>
            <div class="space-y-2">
              <Label for="order">Порядок в меню</Label>
              <Input id="order" v-model.number="form.menu_order" type="number" />
            </div>
          </div>
          <label class="flex items-center gap-3 cursor-pointer">
            <Checkbox v-model="form.in_menu" />
            <span class="text-sm">Показывать в меню</span>
          </label>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Заголовки</CardTitle>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="title-tm">TM <span class="text-destructive">*</span></Label>
            <Input id="title-tm" v-model="form.title_tm" required />
          </div>
          <div class="space-y-2">
            <Label for="title-ru">RU</Label>
            <Input id="title-ru" v-model="form.title_ru" />
          </div>
          <div class="space-y-2">
            <Label for="title-en">EN</Label>
            <Input id="title-en" v-model="form.title_en" />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Контент</CardTitle>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="content-tm">TM</Label>
            <Textarea id="content-tm" v-model="form.content_tm" rows="8" />
          </div>
          <div class="space-y-2">
            <Label for="content-ru">RU</Label>
            <Textarea id="content-ru" v-model="form.content_ru" rows="8" />
          </div>
          <div class="space-y-2">
            <Label for="content-en">EN</Label>
            <Textarea id="content-en" v-model="form.content_en" rows="8" />
          </div>
        </CardContent>
      </Card>

      <div class="flex gap-3">
        <Button type="submit" :disabled="saving">
          {{ saving ? 'Сохранение…' : 'Создать' }}
        </Button>
        <NuxtLink to="/admin/pages">
          <Button type="button" variant="outline">Отмена</Button>
        </NuxtLink>
      </div>
    </form>
  </div>
</template>
