<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const route = useRoute()
const api = useApi()
const toast = useToast()
const id = Number(route.params.id)

const { data: post, refresh } = await useAsyncData(`post-${id}`, () =>
  api.get<any>(`/posts/${id}`)
)

const form = reactive({
  slug: post.value?.slug ?? '',
  author: post.value?.author ?? '',
  status: post.value?.status ?? 'draft',
  title_tm: post.value?.title_tm ?? '',
  title_ru: post.value?.title_ru ?? '',
  title_en: post.value?.title_en ?? '',
  content_tm: post.value?.content_tm ?? '',
  content_ru: post.value?.content_ru ?? '',
  content_en: post.value?.content_en ?? '',
})

const saving = ref(false)
const uploading = ref(false)
const coverInput = ref<HTMLInputElement | null>(null)

async function submit() {
  saving.value = true
  try {
    await api.put(`/admin/posts/${id}`, form)
    await refresh()
    toast.success('Сохранено')
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка сохранения')
  } finally {
    saving.value = false
  }
}

async function uploadCover() {
  const file = coverInput.value?.files?.[0]
  if (!file) return
  uploading.value = true
  try {
    await api.upload(`/admin/posts/${id}/cover`, file)
    await refresh()
    toast.success('Обложка загружена')
    if (coverInput.value) coverInput.value.value = ''
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка загрузки')
  } finally {
    uploading.value = false
  }
}

useHead({ title: `Новость: ${post.value?.title_ru || post.value?.title_tm}` })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/posts" class="text-muted-foreground hover:text-foreground text-sm">← Новости</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">Редактировать новость</h1>
    </div>

    <Card class="mb-6 max-w-3xl">
      <CardHeader>
        <CardTitle class="text-lg">Обложка</CardTitle>
        <CardDescription v-if="post?.cover_path" class="font-mono">{{ post.cover_path }}</CardDescription>
        <CardDescription v-else>Обложка не загружена</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="flex items-center gap-3 flex-wrap">
          <input
            ref="coverInput"
            type="file"
            accept="image/*"
            class="text-sm flex-1 min-w-[200px] file:mr-3 file:rounded-md file:border-0 file:bg-muted file:px-3 file:py-1.5 file:text-sm file:font-medium hover:file:bg-muted/80"
          />
          <Button @click="uploadCover" :disabled="uploading" size="sm">
            {{ uploading ? 'Загрузка…' : 'Загрузить' }}
          </Button>
        </div>
      </CardContent>
    </Card>

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
              <Label for="status">Статус</Label>
              <Select id="status" v-model="form.status">
                <option value="draft">Черновик</option>
                <option value="published">Опубликована</option>
              </Select>
            </div>
          </div>
          <div class="space-y-2">
            <Label for="author">Автор</Label>
            <Input id="author" v-model="form.author" />
          </div>
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
            <Label for="content-tm">Текст (TM)</Label>
            <Textarea id="content-tm" v-model="form.content_tm" rows="6" />
          </div>
          <div class="space-y-2">
            <Label for="content-ru">Текст (RU)</Label>
            <Textarea id="content-ru" v-model="form.content_ru" rows="6" />
          </div>
          <div class="space-y-2">
            <Label for="content-en">Текст (EN)</Label>
            <Textarea id="content-en" v-model="form.content_en" rows="6" />
          </div>
        </CardContent>
      </Card>

      <div class="flex gap-3">
        <Button type="submit" :disabled="saving">
          {{ saving ? 'Сохранение…' : 'Сохранить' }}
        </Button>
        <NuxtLink to="/admin/posts">
          <Button type="button" variant="outline">Назад</Button>
        </NuxtLink>
      </div>
    </form>
  </div>
</template>
