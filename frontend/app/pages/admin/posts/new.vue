<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const api = useApi()
const toast = useToast()
const localePath = useLocalePath()

const form = reactive({
  slug: '',
  author: '',
  status: 'draft',
  title_tm: '',
  title_ru: '',
  title_en: '',
  content_tm: '',
  content_ru: '',
  content_en: '',
})

useAutoSlug(toRef(form, 'title_tm'), toRef(form, 'slug'))

const saving = ref(false)
const activeLang = ref<'tm' | 'ru' | 'en'>('tm')

async function submit() {
  if (!form.title_tm) {
    toast.error(t('admin.common.heading_required'))
    return
  }
  saving.value = true
  try {
    const created = await api.post<{ id: number }>('/admin/posts', form)
    toast.success(t('admin.posts.created'))
    await navigateTo(localePath(`/admin/posts/${created.id}`))
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error_save'))
  } finally {
    saving.value = false
  }
}

useHead({ title: () => t('admin.posts.add_new') })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/posts" class="text-muted-foreground hover:text-foreground text-sm">← {{ t('admin.posts.title') }}</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">{{ t('admin.posts.add_new') }}</h1>
    </div>

    <form @submit.prevent="submit" class="space-y-6 max-w-4xl">
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">{{ t('admin.sections.main') }}</CardTitle>
        </CardHeader>
        <CardContent class="space-y-5">
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <Label for="slug">{{ t('admin.common.slug') }}</Label>
              <Input id="slug" v-model="form.slug" readonly tabindex="-1" class="font-mono text-sm bg-muted text-muted-foreground cursor-not-allowed" :placeholder="t('admin.common.slug_auto_hint')" />
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
          {{ saving ? t('admin.common.saving') : t('admin.common.create') }}
        </Button>
        <NuxtLink to="/admin/posts">
          <Button type="button" variant="outline">{{ t('admin.common.cancel') }}</Button>
        </NuxtLink>
      </div>
    </form>
  </div>
</template>
