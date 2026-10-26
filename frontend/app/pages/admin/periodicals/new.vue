<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const api = useApi()
const toast = useToast()
const localePath = useLocalePath()

const form = reactive({
  slug: '',
  type: 'magazine',
  periodicity: '',
  name_tm: '',
  name_ru: '',
  name_en: '',
  description_tm: '',
  description_ru: '',
  description_en: '',
})

useAutoSlug(toRef(form, 'name_tm'), toRef(form, 'slug'))

const saving = ref(false)

async function submit() {
  if (!form.name_tm) {
    toast.error(t('admin.common.title_required'))
    return
  }
  saving.value = true
  try {
    const created = await api.post<{ id: number }>('/admin/periodicals', form)
    toast.success(t('admin.periodicals.created'))
    await navigateTo(localePath(`/admin/periodicals/${created.id}`))
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error_save'))
  } finally {
    saving.value = false
  }
}

useHead({ title: () => t('admin.periodicals.add_new') })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/periodicals" class="text-muted-foreground hover:text-foreground text-sm">← {{ t('admin.periodicals.title') }}</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">{{ t('admin.periodicals.add_new') }}</h1>
    </div>

    <form @submit.prevent="submit" class="space-y-6 max-w-3xl">
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">{{ t('admin.sections.main') }}</CardTitle>
        </CardHeader>
        <CardContent class="space-y-5">
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <Label for="type">{{ t('admin.periodicals.type') }} <span class="text-destructive">*</span></Label>
              <Select id="type" v-model="form.type">
                <option value="magazine">{{ t('admin.periodicals.type_magazine') }}</option>
                <option value="newspaper">{{ t('admin.periodicals.type_newspaper') }}</option>
              </Select>
            </div>
            <div class="space-y-2">
              <Label for="slug">{{ t('admin.common.slug') }}</Label>
              <Input id="slug" v-model="form.slug" readonly tabindex="-1" class="font-mono text-sm bg-muted text-muted-foreground cursor-not-allowed" :placeholder="t('admin.common.slug_auto_hint')" />
            </div>
          </div>
          <div class="space-y-2">
            <Label for="period">{{ t('admin.periodicals.periodicity') }}</Label>
            <Input id="period" v-model="form.periodicity" :placeholder="t('admin.periodicals.periodicity_placeholder')" />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-lg">{{ t('admin.sections.names_3lang') }}</CardTitle>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="name-tm">{{ t('admin.common.lang_tm') }} <span class="text-destructive">*</span></Label>
            <Input id="name-tm" v-model="form.name_tm" required />
          </div>
          <div class="space-y-2">
            <Label for="name-ru">{{ t('admin.common.lang_ru') }}</Label>
            <Input id="name-ru" v-model="form.name_ru" />
          </div>
          <div class="space-y-2">
            <Label for="name-en">{{ t('admin.common.lang_en') }}</Label>
            <Input id="name-en" v-model="form.name_en" />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-lg">{{ t('admin.sections.descriptions_3lang') }}</CardTitle>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="desc-tm">{{ t('admin.common.lang_tm') }}</Label>
            <Textarea id="desc-tm" v-model="form.description_tm" rows="3" />
          </div>
          <div class="space-y-2">
            <Label for="desc-ru">{{ t('admin.common.lang_ru') }}</Label>
            <Textarea id="desc-ru" v-model="form.description_ru" rows="3" />
          </div>
          <div class="space-y-2">
            <Label for="desc-en">{{ t('admin.common.lang_en') }}</Label>
            <Textarea id="desc-en" v-model="form.description_en" rows="3" />
          </div>
        </CardContent>
      </Card>

      <div class="flex gap-3">
        <Button type="submit" :disabled="saving">
          {{ saving ? t('admin.common.saving') : t('admin.common.create') }}
        </Button>
        <NuxtLink to="/admin/periodicals">
          <Button type="button" variant="outline">{{ t('admin.common.cancel') }}</Button>
        </NuxtLink>
      </div>
    </form>
  </div>
</template>
