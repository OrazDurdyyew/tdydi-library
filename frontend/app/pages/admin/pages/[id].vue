<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const route = useRoute()
const api = useApi()
const toast = useToast()
const id = Number(route.params.id)

const { data: page, refresh } = await useAsyncData(`page-${id}`, () =>
  api.get<any>(`/pages/${id}`)
)

const form = reactive({
  slug: page.value?.slug ?? '',
  in_menu: page.value?.in_menu ?? false,
  menu_order: page.value?.menu_order ?? 0,
  title_tm: page.value?.title_tm ?? '',
  title_ru: page.value?.title_ru ?? '',
  title_en: page.value?.title_en ?? '',
  content_tm: page.value?.content_tm ?? '',
  content_ru: page.value?.content_ru ?? '',
  content_en: page.value?.content_en ?? '',
})

useAutoSlug(toRef(form, 'title_tm'), toRef(form, 'slug'))

const saving = ref(false)
const activeLang = ref<'tm' | 'ru' | 'en'>('tm')

async function submit() {
  saving.value = true
  try {
    await api.put(`/admin/pages/${id}`, form)
    await refresh()
    toast.success(t('admin.common.saved'))
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error'))
  } finally {
    saving.value = false
  }
}

useHead({ title: () => `${t('admin.pages.edit_title')}: ${page.value?.title_ru || page.value?.title_tm}` })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/pages" class="text-muted-foreground hover:text-foreground text-sm">← {{ t('admin.pages.title') }}</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">{{ page?.title_ru || page?.title_tm }}</h1>
    </div>

    <form @submit.prevent="submit" class="space-y-6 max-w-3xl">
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
              <Label for="order">{{ t('admin.pages.menu_order') }}</Label>
              <Input id="order" v-model.number="form.menu_order" type="number" />
            </div>
          </div>
          <label class="flex items-center gap-3 cursor-pointer">
            <Checkbox v-model="form.in_menu" />
            <span class="text-sm">{{ t('admin.pages.in_menu') }}</span>
          </label>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-lg">{{ t('admin.common.title') }}</CardTitle>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="title-tm">{{ t('admin.common.lang_tm') }}</Label>
            <Input id="title-tm" v-model="form.title_tm" />
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
            <CardTitle class="text-lg">{{ t('admin.pages.content') }}</CardTitle>
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
            <RichTextEditor v-model="form.content_tm" :placeholder="t('admin.pages.content_placeholder_tm')" />
          </div>
          <div v-show="activeLang === 'ru'">
            <RichTextEditor v-model="form.content_ru" :placeholder="t('admin.pages.content_placeholder_ru')" />
          </div>
          <div v-show="activeLang === 'en'">
            <RichTextEditor v-model="form.content_en" placeholder="Page content (English)…" />
          </div>
        </CardContent>
      </Card>

      <div class="flex gap-3">
        <Button type="submit" :disabled="saving">
          {{ saving ? t('admin.common.saving') : t('admin.common.save') }}
        </Button>
        <NuxtLink to="/admin/pages">
          <Button type="button" variant="outline">{{ t('admin.common.back') }}</Button>
        </NuxtLink>
      </div>
    </form>
  </div>
</template>
