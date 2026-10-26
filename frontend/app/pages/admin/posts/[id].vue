<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const { t } = useI18n()
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

useAutoSlug(toRef(form, 'title_tm'), toRef(form, 'slug'))

const saving = ref(false)
const uploading = ref(false)
const coverInput = ref<HTMLInputElement | null>(null)
const activeLang = ref<'tm' | 'ru' | 'en'>('tm')

async function submit() {
  saving.value = true
  try {
    await api.put(`/admin/posts/${id}`, form)
    await refresh()
    toast.success(t('admin.common.saved'))
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error_save'))
  } finally {
    saving.value = false
  }
}

async function uploadCover() {
  const file = coverInput.value?.files?.[0]
  if (!file) return
  uploading.value = true
  try {
    await api.uploadField(`/admin/posts/${id}/cover`, 'cover', file)
    await refresh()
    toast.success(t('admin.posts.cover_uploaded'))
    if (coverInput.value) coverInput.value.value = ''
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error_upload'))
  } finally {
    uploading.value = false
  }
}

useHead({ title: () => `${t('admin.posts.edit_title')}: ${post.value?.title_ru || post.value?.title_tm}` })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/posts" class="text-muted-foreground hover:text-foreground text-sm">← {{ t('admin.posts.title') }}</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">{{ t('admin.posts.edit_title') }}</h1>
    </div>

    <Card class="mb-6 max-w-4xl">
      <CardHeader>
        <CardTitle class="text-lg">{{ t('admin.common.cover') }}</CardTitle>
        <CardDescription v-if="post?.cover_path" class="font-mono">{{ post.cover_path }}</CardDescription>
        <CardDescription v-else>{{ t('admin.posts.cover_not_uploaded') }}</CardDescription>
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
            {{ uploading ? t('admin.common.uploading') : t('admin.common.upload') }}
          </Button>
        </div>
      </CardContent>
    </Card>

    <form @submit.prevent="submit" class="space-y-6 max-w-4xl">
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">{{ t('admin.sections.main') }}</CardTitle>
        </CardHeader>
        <CardContent class="space-y-5">
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <Label for="slug">{{ t('admin.common.slug') }}</Label>
              <Input id="slug" v-model="form.slug" readonly tabindex="-1" class="font-mono text-sm bg-muted text-muted-foreground cursor-not-allowed" />
            </div>
            <div class="space-y-2">
              <Label for="status">{{ t('admin.posts.status') }}</Label>
              <Select id="status" v-model="form.status">
                <option value="draft">{{ t('admin.common.draft') }}</option>
                <option value="published">{{ t('admin.common.published') }}</option>
              </Select>
            </div>
          </div>
          <div class="space-y-2">
            <Label for="author">{{ t('admin.posts.author') }}</Label>
            <Input id="author" v-model="form.author" />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-lg">{{ t('admin.common.title') }}</CardTitle>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="title-tm">{{ t('admin.common.lang_tm') }} <span class="text-destructive">*</span></Label>
            <Input id="title-tm" v-model="form.title_tm" required />
          </div>
          <div class="space-y-2">
            <Label for="title-ru">{{ t('admin.common.lang_ru') }}</Label>
            <Input id="title-ru" v-model="form.title_ru" />
          </div>
          <div class="space-y-2">
            <Label for="title-en">{{ t('admin.common.lang_en') }}</Label>
            <Input id="title-en" v-model="form.title_en" />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <div class="flex items-center justify-between gap-3">
            <CardTitle class="text-lg">{{ t('admin.posts.content') }}</CardTitle>
            <div class="inline-flex rounded-md border border-border overflow-hidden text-sm">
              <button
                v-for="l in ['tm', 'ru', 'en'] as const"
                :key="l"
                type="button"
                class="px-3 py-1.5 transition-colors"
                :class="activeLang === l ? 'bg-primary text-primary-foreground' : 'hover:bg-muted'"
                @click="activeLang = l"
              >{{ l.toUpperCase() }}</button>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div v-show="activeLang === 'tm'">
            <RichTextEditor v-model="form.content_tm" :placeholder="t('admin.posts.content_placeholder_tm')" />
          </div>
          <div v-show="activeLang === 'ru'">
            <RichTextEditor v-model="form.content_ru" :placeholder="t('admin.posts.content_placeholder_ru')" />
          </div>
          <div v-show="activeLang === 'en'">
            <RichTextEditor v-model="form.content_en" placeholder="News content (English)…" />
          </div>
        </CardContent>
      </Card>

      <div class="flex gap-3">
        <Button type="submit" :disabled="saving">
          {{ saving ? t('admin.common.saving') : t('admin.common.save') }}
        </Button>
        <NuxtLink to="/admin/posts">
          <Button type="button" variant="outline">{{ t('admin.common.back') }}</Button>
        </NuxtLink>
      </div>
    </form>
  </div>
</template>
